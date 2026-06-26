package main

import (
	"fmt"
	"slices"
)

type pokemon struct {
	base        basePokemon
	level       int
	ivs         []int
	nature      []stats
	moves       []*move
	lockedMove  *move
	stats       []int
	stages      []int
	hp          int
	fainted     bool
	ailments    map[ailmentState]*ailment
	item        *item
	ability     string
	unnerved    bool
	flashFire   bool
	unburden    bool
	trace       bool
	focusEnergy bool
	laserFocus  bool
}

func getNature(nature string) ([]stats, error) {
	res, ok := natureChart[nature]
	if !ok {
		return nil, fmt.Errorf("invalid nature: %s", nature)
	}

	return res, nil

}

func initPokemon(base basePokemon, level int, ivs map[string]int, nature string, moves []*move, hp int, status ailmentState) (pokemon, error) {
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

	setIVs(&res, ivs)

	calculateStats(&res)

	if hp == -1 {
		hp = res.maxHP()
	}

	res.hp = max(1, min(res.maxHP(), hp))

	if _, ok := nonVolatileStatuses[status]; ok {
		res.ailments[status] = generateAilment(status, nil)
	}

	return res, nil
}

func setIVs(pokemon *pokemon, ivs map[string]int) {
	for key, val := range ivs {
		stat := stringToStat(key)
		pokemon.ivs[stat] = max(0, min(31, val))
	}
}

func calculateStats(pokemon *pokemon) {
	for key, val := range pokemon.base.Stats {
		stat := stringToStat(key)
		pokemon.stats[stat] = ((val*2+pokemon.ivs[stat])*pokemon.level)/100 + 5
	}
	// Shedinja case: if HP is 1, it stays 1 regardless of level or IVs
	if pokemon.stats[HitPoints] == 1 {
		pokemon.stats[HitPoints] = 1
	} else {
		pokemon.stats[HitPoints] += pokemon.level + 5
	}

	// Apply nature modifiers
	posNat := pokemon.nature[0]
	negNat := pokemon.nature[1]

	if posNat != negNat {
		pokemon.stats[posNat] = (pokemon.stats[posNat] * 110) / 100
		pokemon.stats[negNat] = (pokemon.stats[negNat] * 90) / 100
	}
}

func (p *pokemon) switchReset() {
	for a := range volatileStatuses {
		delete(p.ailments, a)
	}

	for stat := range p.stages {
		p.stages[stat] = 0
	}

	if toxic, ok := p.ailments[Toxic]; ok {
		toxic.turns = 0
	}

	if p.trace {
		p.trace = false
		p.ability = "trace"
	}

	p.lockedMove = nil
	p.flashFire = false
	p.unburden = false
	p.focusEnergy = false
	p.laserFocus = false
}

func (p *pokemon) effectiveStat(stat stats, crit bool) int {
	stage := p.stages[stat]
	base := p.stats[stat]

	if crit {
		switch stat {
		case Defense, SpecialDefense:
			stage = min(0, stage)
		case Attack, SpecialAttack:
			stage = max(0, stage)
		}
	}

	if stage >= 0 {
		return base * (2 + stage) / 2
	}
	return base * 2 / (2 - stage)
}

func (p *pokemon) effectiveSpeed(bs battleState) int {
	stage := p.stages[Speed]
	base := p.stats[Speed]
	numerator := 1
	denominator := 1

	if p.item.name == "iron-ball" {
		denominator *= 2
	} else if p.unburden && p.ability == "unburden" {
		numerator *= 2
	}
	if _, ok := p.ailments[Paralysis]; ok {
		denominator *= 4
	}
	if bs.getWeather() != NoneWeather {
		switch bs.getWeather() {
		case Rain:
			if p.ability == "swift-swim" {
				numerator *= 2
			}
		case Sun:
			if p.ability == "chlorophyll" {
				numerator *= 2
			}
		case Hail:
			if p.ability == "slush-rush" {
				numerator *= 2
			}
		case Sandstorm:
			if p.ability == "sand-rush" {
				numerator *= 2
			}
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

	stage := p.stages[Evasion]
	if stage == 0 {
		return 3, 3
	} else if stage > 0 {
		return 3, 3 + stage
	}
	return 3 - stage, 3
}

func (p *pokemon) accuracyFraction() (int, int) {
	stage := p.stages[Accuracy]
	if stage == 0 {
		return 3, 3
	} else if stage > 0 {
		return 3 + stage, 3
	}
	return 3, 3 - stage
}

func (p *pokemon) hasType(typeName string) bool {
	for _, t := range p.base.Types {
		if typeName == t {
			return true
		}
	}
	return false
}

func (p *pokemon) applyAilment(ailment ailmentState, move *move, afflictedBy *slot) {
	if ailment == NoneAilment {
		elogf("%s applies an ailment but is none", ailment.String())
		return
	}

	if _, ok := p.ailments[ailment]; ok {
		return
	}
	if _, ok := nonVolatileStatuses[ailment]; ok && p.hasNonVolatileAilment() {
		return
	}
	if ailment == Burn && (p.hasType("fire") || p.ability == "water-veil") {
		return
	}
	if ailment == Paralysis && (p.hasType("electric") || p.ability == "limber") {
		return
	}
	if ailment == Poison || ailment == Toxic {
		if p.ability == "immunity" {
			return
		}
		if (p.hasType("poison") || p.hasType("steel")) && (afflictedBy == nil || afflictedBy.mon.ability != "corrosion") {
			return
		}
	}
	if ailment == Freeze && p.hasType("ice") {
		return
	}
	if ailment == Sleep {
		if _, ok := sleepBlockingAbilities[p.ability]; ok {
			return
		}
	}
	if ailment == Yawn {
		if _, ok := sleepBlockingAbilities[p.ability]; ok || p.hasNonVolatileAilment() {
			return
		}
	}

	switch ailment {
	case Trap:
		p.ailments[ailment] = generateTrap(move.MinTurns, move.MaxTurns, afflictedBy)
		return
	case Poison:
		if move != nil && (move.Name == "toxic" || move.Name == "poison-fang") {
			ailment = Toxic
		}
	case Infatuation:
		if afflictedBy.mon.ability == "oblivious" {
			return
		}
	}

	p.ailments[ailment] = generateAilment(ailment, afflictedBy)
	vlogf("%s became afflicted with %s", p.base.Name, ailment.String())
	if _, ok := nonVolatileStatuses[ailment]; ok && p.ability == "synchronize" {
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
		if ailment <= Sleep {
			return true
		}
	}
	return false
}

func (p *pokemon) isGrounded() bool {
	if p.item.name == "iron-ball" {
		return true
	}
	if p.hasType("flying") || p.ability == "levitate" {
		return false
	}
	return true
}

func (p *pokemon) changeHpBy(change int) {
	p.hp = min(p.hp+change, p.maxHP())
	p.checkItemTrigger(true, nil)
}

func (p *pokemon) hasMovePredicate(f func(*move) bool) bool {
	return slices.ContainsFunc(p.moves, f)
}

func (p *pokemon) changeStatStageBy(stat stats, change int, offensive bool) {
	if offensive && (p.ability == "clear-smoke" || p.ability == "clear-body") {
		vlogf("blocked by clear body")
		return
	}
	if p.ability == "keen-eye" && stat == Accuracy && change < 0 {
		return
	}

	p.stages[stat] = max(-6, min(6, p.stages[stat]+change))
	vlogf("%s's %s changed by %d stages (%d)", p.base.Name, stat, change, p.stages[stat])
}

func (p *pokemon) maxHP() int {
	return p.stats[HitPoints]
}

func (p *pokemon) serenceGraceBonus() int {
	if p.ability == "serence-grace" {
		return 2
	}
	return 1
}

func (p *pokemon) applyMoveType(num, dem int, moveType string) (int, int) {
	for _, t := range p.base.Types {
		if t == "flying" && moveType == "ground" && p.isGrounded() {
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
