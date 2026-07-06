package main

import (
	"fmt"
	"slices"
)

type pokemon struct {
	base        BasePokemon
	level       int
	ivs         []int
	nature      []stat
	moves       []*Move
	lockedMove  *Move
	stats       []int
	stages      []int
	hp          int
	fainted     bool
	ailments    map[ailmentState]*ailment
	item        *item
	ability     ability
	unnerved    bool
	flashFire   bool
	unburden    bool
	trace       bool
	focusEnergy bool
	laserFocus  bool
}

func getNature(nature string) ([]stat, error) {
	res, ok := natureChart[nature]
	if !ok {
		return nil, fmt.Errorf("invalid nature: %s", nature)
	}

	return res, nil

}

func initPokemon(base BasePokemon, level int, ivs map[string]int, nature string, moves []*Move, hp int, status ailmentState) (pokemon, error) {
	if level < 1 || level > 100 {
		return pokemon{}, fmt.Errorf("invalid level: %d", level)
	}

	nat, err := getNature(nature)
	if err != nil {
		return pokemon{}, err
	}

	res := pokemon{
		base:     base,
		level:    level,
		ivs:      []int{31, 31, 31, 31, 31, 31},
		nature:   nat,
		moves:    moves,
		stats:    []int{0, 0, 0, 0, 0, 0},
		stages:   []int{0, 0, 0, 0, 0, 0, 0, 0},
		hp:       0,
		fainted:  false,
		ailments: make(map[ailmentState]*ailment),
	}

	err = setIVs(&res, ivs)
	if err != nil {
		return pokemon{}, err
	}

	err = calculateStats(&res)
	if err != nil {
		return pokemon{}, err
	}

	if hp == -1 {
		hp = res.maxHP()
	}

	res.hp = max(1, min(res.maxHP(), hp))

	if status.isNonVolatileStatus() {
		res.ailments[status] = generateAilment(status, nil)
	}

	return res, nil
}

func setIVs(pokemon *pokemon, ivs map[string]int) error {
	for key, val := range ivs {
		stat := stringToStat(key)
		if stat == noStat {
			return fmt.Errorf("%s is not a valid stat for %s", key, pokemon.base.Name)
		}
		pokemon.ivs[stat] = max(0, min(31, val))
	}

	return nil
}

func calculateStats(pokemon *pokemon) error {
	for key, val := range pokemon.base.Stats {
		stat := stringToStat(key)
		if stat == noStat {
			return fmt.Errorf("%s is not a valid stat for %s", key, pokemon.base.Name)
		}
		pokemon.stats[stat] = ((val*2+pokemon.ivs[stat])*pokemon.level)/100 + 5
	}
	// Shedinja case: if HP is 1, it stays 1 regardless of level or IVs
	if pokemon.stats[hitPoints] == 1 {
		pokemon.stats[hitPoints] = 1
	} else {
		pokemon.stats[hitPoints] += pokemon.level + 5
	}

	// Apply nature modifiers
	posNat := pokemon.nature[0]
	negNat := pokemon.nature[1]

	if posNat != negNat {
		pokemon.stats[posNat] = (pokemon.stats[posNat] * 110) / 100
		pokemon.stats[negNat] = (pokemon.stats[negNat] * 90) / 100
	}

	return nil
}

func (p *pokemon) switchReset() {
	for a := range volatileStatuses {
		delete(p.ailments, a)
	}

	for stat := range p.stages {
		p.stages[stat] = 0
	}

	if toxic, ok := p.ailments[toxicAilment]; ok {
		toxic.turns = 0
	}

	if p.trace {
		p.trace = false
		p.ability = traceAbility
	}

	p.lockedMove = nil
	p.flashFire = false
	p.unburden = false
	p.focusEnergy = false
	p.laserFocus = false
}

func (p *pokemon) effectiveStat(stat stat, crit bool) int {
	stage := p.stages[stat]
	base := p.stats[stat]

	if crit {
		switch stat {
		case defense, specialDefense:
			stage = min(0, stage)
		case attack, specialAttack:
			stage = max(0, stage)
		}
	}

	if stage >= 0 {
		return base * (2 + stage) / 2
	}
	return base * 2 / (2 - stage)
}

func (p *pokemon) effectiveSpeed(bs battleState) int {
	stage := p.stages[speed]
	base := p.stats[speed]
	numerator := 1
	denominator := 1

	if p.item.state == ironBall {
		denominator *= 2
	} else if p.unburden && p.ability == unburdenAbility {
		numerator *= 2
	}
	if _, ok := p.ailments[paralysisAilment]; ok {
		denominator *= 4
	}
	switch bs.getWeather() {
	case rainWeather:
		if p.ability == swiftSwimAbility {
			numerator *= 2
		}
	case sunWeather:
		if p.ability == chlorophyllAbility {
			numerator *= 2
		}
	case hailWeather:
		if p.ability == slushRushAbility {
			numerator *= 2
		}
	case sandstormWeather:
		if p.ability == sandRushAbility {
			numerator *= 2
		}
	}

	base = base * numerator / denominator

	if stage >= 0 {
		return base * (2 + stage) / 2
	}
	return base * 2 / (2 - stage)
}

func (p *pokemon) isFasterThan(bs battleState, mon *pokemon) bool {
	return p.effectiveSpeed(bs) >= mon.effectiveSpeed(bs)
}

func (p *pokemon) evasionFraction(keenEye bool) (int, int) {
	if keenEye {
		return 1, 1
	}

	stage := p.stages[evasion]
	if stage == 0 {
		return 3, 3
	} else if stage > 0 {
		return 3, 3 + stage
	}
	return 3 - stage, 3
}

func (p *pokemon) accuracyFraction() (int, int) {
	stage := p.stages[accuracy]
	if stage == 0 {
		return 3, 3
	} else if stage > 0 {
		return 3 + stage, 3
	}
	return 3, 3 - stage
}

func (p *pokemon) hasType(typeName pokemonType) bool {
	for _, t := range p.base.Types {
		if typeName == t {
			return true
		}
	}
	return false
}

func (p *pokemon) applyAilment(ailment ailmentState, move *Move, afflictedBy *slot) {
	if ailment == noneAilment {
		elogf("warning: %s applies an ailment but is none", ailment.String())
		return
	}

	if _, ok := p.ailments[ailment]; ok {
		return
	}
	if ailment.isNonVolatileStatus() && p.hasNonVolatileAilment() {
		return
	}

	switch ailment {
	case burnAilment:
		if p.hasType(fireType) || p.ability == waterVeilAbility {
			return
		}
	case paralysisAilment:
		if p.hasType(electricType) || p.ability == limberAbility {
			return
		}
	case poisonAilment, toxicAilment:
		if p.ability == immunityAbility {
			return
		}
		if (p.hasType(poisonType) || p.hasType(steelType)) && (afflictedBy == nil || afflictedBy.mon.ability != corrosionAbility) {
			return
		}
	case freezeAilment:
		if p.hasType(iceType) || p.ability == magmaArmorAbility {
			return
		}
	case sleepAilment, yawnAilment:
		if p.ability.blocksSleep() || p.hasNonVolatileAilment() {
			return
		}
	case trapAilment:
		p.ailments[ailment] = generateTrap(move.MinTurns, move.MaxTurns, afflictedBy)
		return
	case infatuationAilment:
		if afflictedBy.mon.ability == obliviousAbility {
			return
		}
	}

	if ailment == poisonAilment {
		if move != nil && (move.Name == "toxic" || move.Name == "poison-fang") {
			ailment = toxicAilment
		}
	}

	p.ailments[ailment] = generateAilment(ailment, afflictedBy)
	vprintf("%s became afflicted with %s", p.base.Name, ailment.String())
	if ailment.isNonVolatileStatus() && p.ability == synchronizeAbility {
		afflictedBy.mon.applyAilment(ailment, nil, nil)
	}
	p.checkItemTrigger(true, nil)
}

func (p *pokemon) hasAilment(ailment ailmentState) *ailment {
	if a, ok := p.ailments[ailment]; ok {
		return a
	}
	return nil
}

func (p *pokemon) hasNonVolatileAilment() bool {
	for ailment := range p.ailments {
		if ailment <= sleepAilment {
			return true
		}
	}
	return false
}

func (p *pokemon) isGrounded() bool {
	if p.item.state == ironBall {
		return true
	}
	if p.hasType(flyingType) || p.ability == levitateAbility {
		return false
	}
	return true
}

func (p *pokemon) changeHpBy(change int) {
	p.hp = min(p.hp+change, p.maxHP())
	p.checkItemTrigger(true, nil)
}

func (p *pokemon) hasMovePredicate(f func(*Move) bool) bool {
	return slices.ContainsFunc(p.moves, f)
}

func (p *pokemon) changeStatStageBy(stat stat, change int, offensive bool) {
	if offensive && (p.ability == clearBodyAbility || p.ability == clearSmokeAbility) {
		vprintf("blocked by clear body")
		return
	}
	if p.ability == keenEyeAbility && stat == accuracy && change < 0 {
		return
	}

	p.stages[stat] = max(-6, min(6, p.stages[stat]+change))
	vprintf("%s's %s changed by %d stages (%d)", p.base.Name, stat, change, p.stages[stat])
}

func (p *pokemon) maxHP() int {
	return p.stats[hitPoints]
}

func (p *pokemon) serenceGraceBonus() int {
	if p.ability == serenceGraceAbility {
		return 2
	}
	return 1
}

func (p *pokemon) applyMoveType(num, dem int, moveType pokemonType) (int, int) {
	for _, t := range p.base.Types {
		if t == flyingType && moveType == groundType && p.isGrounded() {
			continue
		}
		switch getEffectiveness(moveType, t) {
		case 0:
			num = 0
		case 0.5:
			dem *= 2
		case 1:
		case 2:
			num *= 2
		}
	}

	return num, dem
}

func (p *pokemon) isImmuneToPowderMoves() bool {
	return p.hasType(grassType) || p.ability == overcoatAbility
}
