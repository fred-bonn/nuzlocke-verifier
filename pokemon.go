package main

import (
	"fmt"
	"log"
	"slices"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

type pokemon struct {
	base        pokeapi.BasePokemon
	level       int
	ivs         map[string]int
	nature      []string
	moves       []*pokeapi.BaseMove
	lockedMove  *pokeapi.BaseMove
	stats       map[string]int
	stages      map[string]int
	hp          int
	fainted     bool
	ailments    map[string]*ailment
	item        *item
	ability     string
	unnerved    bool
	flashFire   bool
	unburden    bool
	trace       bool
	focusEnergy bool
	laserFocus  bool
}

var ivMap = map[string]string{
	"hp":  "hp",
	"atk": "attack",
	"def": "defense",
	"spa": "special-attack",
	"spd": "special-defense",
	"spe": "speed",
}

var natureChart = map[string][]string{
	"hardy":   {"attack", "attack"},
	"lonely":  {"attack", "defense"},
	"adamant": {"attack", "special-attack"},
	"naughty": {"attack", "special-defense"},
	"brave":   {"attack", "speed"},
	"bold":    {"defense", "attack"},
	"docile":  {"defense", "defense"},
	"impish":  {"defense", "special-attack"},
	"lax":     {"defense", "special-defense"},
	"relaxed": {"defense", "speed"},
	"modest":  {"special-attack", "attack"},
	"mild":    {"special-attack", "defense"},
	"bashful": {"special-attack", "special-attack"},
	"rash":    {"special-attack", "special-defense"},
	"quiet":   {"special-attack", "speed"},
	"calm":    {"special-defense", "attack"},
	"gentle":  {"special-defense", "defense"},
	"careful": {"special-defense", "speed"},
	"quirky":  {"special-defense", "special-defense"},
	"sassy":   {"special-defense", "speed"},
	"timid":   {"speed", "attack"},
	"hasty":   {"speed", "defense"},
	"jolly":   {"speed", "special-attack"},
	"naive":   {"speed", "special-defense"},
	"serious": {"speed", "speed"},
}

func getNature(nature string) ([]string, error) {
	res, ok := natureChart[nature]
	if !ok {
		return nil, fmt.Errorf("invalid nature: %s", nature)
	}

	return res, nil

}

func initPokemon(base pokeapi.BasePokemon, level int, ivs map[string]int, nature string, moves []*pokeapi.BaseMove, hp int, status string) (pokemon, error) {
	if level < 1 || level > 100 {
		return pokemon{}, fmt.Errorf("invalid level: %d", level)
	}

	nat, err := getNature(nature)
	if err != nil {
		return pokemon{}, err
	}

	res := pokemon{
		base:  base,
		level: level,
		ivs: map[string]int{
			"hp":              31,
			"attack":          31,
			"defense":         31,
			"speed":           31,
			"special-attack":  31,
			"special-defense": 31,
		},
		nature: nat,
		moves:  moves,
		stats:  make(map[string]int, 6),
		stages: map[string]int{
			"attack":          0,
			"defense":         0,
			"speed":           0,
			"special-attack":  0,
			"special-defense": 0,
			"accuracy":        0,
			"evasion":         0,
		},
		hp:       0,
		fainted:  false,
		ailments: make(map[string]*ailment),
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
		key = ivMap[key]
		pokemon.ivs[key] = max(0, min(31, val))
	}
}

func calculateStats(pokemon *pokemon) {
	for key, val := range pokemon.base.Stats {
		pokemon.stats[key] = ((val*2+pokemon.ivs[key])*pokemon.level)/100 + 5
	}
	// Shedinja case: if HP is 1, it stays 1 regardless of level or IVs
	if pokemon.stats["hp"] == 1 {
		pokemon.stats["hp"] = 1
	} else {
		pokemon.stats["hp"] += pokemon.level + 5
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

	if toxic, ok := p.ailments["toxic"]; ok {
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

func (p *pokemon) effectiveStat(stat string, crit bool) int {
	if _, ok := p.stages[stat]; !ok {
		panic("invalid stat")
	}

	stage := p.stages[stat]
	base := p.stats[stat]

	if crit {
		switch stat {
		case "defense", "special-defense":
			stage = min(0, stage)
		case "attack", "special-attack":
			stage = max(0, stage)
		}
	}

	if stage >= 0 {
		return base * (2 + stage) / 2
	}
	return base * 2 / (2 - stage)
}

func (p *pokemon) effectiveSpeed(bs battleState) int {
	stage := p.stages["speed"]
	base := p.stats["speed"]
	numerator := 1
	denominator := 1

	if p.item.name == "iron-ball" {
		denominator *= 2
	} else if p.unburden && p.ability == "unburden" {
		numerator *= 2
	}
	if _, ok := p.ailments["paralysis"]; ok {
		denominator *= 4
	}
	if bs.getWeather() != None {
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

	stage := p.stages["evasion"]
	if stage == 0 {
		return 3, 3
	} else if stage > 0 {
		return 3, 3 + stage
	}
	return 3 - stage, 3
}

func (p *pokemon) accuracyFraction() (int, int) {
	stage := p.stages["accuracy"]
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

func (p *pokemon) applyAilment(ailment string, move *pokeapi.BaseMove, afflictedBy *slot) {
	if _, ok := volatileStatuses[ailment]; !ok {
		if _, ok := nonVolatileStatuses[ailment]; !ok {
			panic("invalid ailment")
		}
	}

	if _, ok := p.ailments[ailment]; ok {
		return
	}
	if _, ok := nonVolatileStatuses[ailment]; ok && p.hasNonVolatileAilment() {
		return
	}
	if ailment == "burn" && (p.hasType("fire") || p.ability == "water-veil") {
		return
	}
	if ailment == "paralysis" && (p.hasType("electric") || p.ability == "limber") {
		return
	}
	if ailment == "poison" || ailment == "toxic" {
		if p.ability == "immunity" {
			return
		}
		if (p.hasType("poison") || p.hasType("steel")) && (afflictedBy == nil || afflictedBy.mon.ability != "corrosion") {
			return
		}
	}
	if ailment == "freeze" && p.hasType("ice") {
		return
	}
	if ailment == "sleep" {
		if _, ok := sleepBlockingAbilities[p.ability]; ok {
			return
		}
	}
	if ailment == "yawn" {
		if _, ok := sleepBlockingAbilities[p.ability]; ok || p.hasNonVolatileAilment() {
			return
		}
	}

	switch ailment {
	case "trap":
		p.ailments[ailment] = generateTrap(move.MinTurns, move.MaxTurns, afflictedBy)
		return
	case "poison":
		if move != nil && (move.Name == "toxic" || move.Name == "poison-fang") {
			ailment = "toxic"
		}
	case "infatuation":
		if afflictedBy.mon.ability == "oblivious" {
			return
		}
	}

	p.ailments[ailment] = generateAilment(ailment, afflictedBy)
	log.Printf("%s became afflicted with %s", p.base.Name, ailment)
	if _, ok := nonVolatileStatuses[ailment]; ok && p.ability == "synchronize" {
		afflictedBy.mon.applyAilment(ailment, nil, nil)
	}
	p.checkItemTrigger(true, nil)
}

func (p *pokemon) hasAilment(ailment string) *ailment {
	if a, ok := p.ailments[ailment]; ok {
		return a
	}
	return nil
}

func (p *pokemon) hasNonVolatileAilment() bool {
	for ailment := range nonVolatileStatuses {
		if p.hasAilment(ailment) != nil {
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

func (p *pokemon) hasMovePredicate(f func(*pokeapi.BaseMove) bool) bool {
	return slices.ContainsFunc(p.moves, f)
}

func (p *pokemon) changeStatStageBy(stat string, change int, offensive bool) {
	if offensive && (p.ability == "clear-smoke" || p.ability == "clear-body") {
		log.Printf("blocked by clear body")
		return
	}
	if p.ability == "keen-eye" && stat == "accuracy" && change < 0 {
		return
	}

	p.stages[stat] = max(-6, min(6, p.stages[stat]+change))
	log.Printf("%s's %s changed by %d stages (%d)", p.base.Name, stat, change, p.stages[stat])
}

func (p *pokemon) maxHP() int {
	return p.stats["hp"]
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
