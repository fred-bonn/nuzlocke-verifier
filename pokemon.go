package main

import (
	"fmt"
	"log"
	"slices"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type Pokemon struct {
	Base        pokeapi.BasePokemon
	Level       int
	IVs         map[string]int
	Nature      []string
	Moves       []*pokeapi.BaseMove
	LockedMove  *pokeapi.BaseMove
	Stats       map[string]int
	Stages      map[string]int
	Hp          int
	Fainted     bool
	Ailments    map[string]*Ailment
	Item        *item
	Ability     string
	Unnerved    bool
	FlashFire   bool
	Unburden    bool
	Trace       bool
	FocusEnergy bool
	LaserFocus  bool
}

var ivMap = map[string]string{
	"hp":  "hp",
	"atk": "attack",
	"def": "defense",
	"spa": "special-attack",
	"spd": "special-defense",
	"spe": "speed",
}

func initPokemon(base pokeapi.BasePokemon, level int, ivs map[string]int, nature string, moves []*pokeapi.BaseMove, hp int, status string) (Pokemon, error) {
	if level < 1 || level > 100 {
		return Pokemon{}, fmt.Errorf("invalid level: %d", level)
	}

	nat, err := pokemon.GetNature(nature)
	if err != nil {
		return Pokemon{}, err
	}

	res := Pokemon{
		Base:  base,
		Level: level,
		IVs: map[string]int{
			"hp":              31,
			"attack":          31,
			"defense":         31,
			"speed":           31,
			"special-attack":  31,
			"special-defense": 31,
		},
		Nature: nat,
		Moves:  moves,
		Stats:  make(map[string]int, 6),
		Stages: map[string]int{
			"attack":          0,
			"defense":         0,
			"speed":           0,
			"special-attack":  0,
			"special-defense": 0,
			"accuracy":        0,
			"evasion":         0,
		},
		Hp:       0,
		Fainted:  false,
		Ailments: make(map[string]*Ailment),
	}

	setIVs(&res, ivs)

	calculateStats(&res)

	if hp == -1 {
		hp = res.maxHP()
	}

	res.Hp = max(1, min(res.maxHP(), hp))

	if _, ok := nonVolatileStatuses[status]; ok {
		res.Ailments[status] = generateAilment(status, nil)
	}

	return res, nil
}

func setIVs(pokemon *Pokemon, ivs map[string]int) {
	for key, val := range ivs {
		key = ivMap[key]
		pokemon.IVs[key] = max(0, min(31, val))
	}
}

func calculateStats(pokemon *Pokemon) {
	for key, val := range pokemon.Base.Stats {
		pokemon.Stats[key] = ((val*2+pokemon.IVs[key])*pokemon.Level)/100 + 5
	}
	// Shedinja case: if HP is 1, it stays 1 regardless of level or IVs
	if pokemon.Stats["hp"] == 1 {
		pokemon.Stats["hp"] = 1
	} else {
		pokemon.Stats["hp"] += pokemon.Level + 5
	}

	// Apply nature modifiers
	posNat := pokemon.Nature[0]
	negNat := pokemon.Nature[1]

	if posNat != negNat {
		pokemon.Stats[posNat] = (pokemon.Stats[posNat] * 110) / 100
		pokemon.Stats[negNat] = (pokemon.Stats[negNat] * 90) / 100
	}
}

func (p *Pokemon) switchReset() {
	for a := range volatileStatuses {
		delete(p.Ailments, a)
	}

	for stat := range p.Stages {
		p.Stages[stat] = 0
	}

	if toxic, ok := p.Ailments["toxic"]; ok {
		toxic.Turns = 0
	}

	if p.Trace {
		p.Trace = false
		p.Ability = "trace"
	}

	p.LockedMove = nil
	p.FlashFire = false
	p.Unburden = false
	p.FocusEnergy = false
	p.LaserFocus = false
}

func (p *Pokemon) effectiveStat(stat string, crit bool) int {
	if _, ok := p.Stages[stat]; !ok {
		panic("invalid stat")
	}

	stage := p.Stages[stat]
	base := p.Stats[stat]

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

func (p *Pokemon) effectiveSpeed(bs battleState) int {
	stage := p.Stages["speed"]
	base := p.Stats["speed"]
	numerator := 1
	denominator := 1

	if p.Item.name == "iron-ball" {
		denominator *= 2
	} else if p.Unburden && p.Ability == "unburden" {
		numerator *= 2
	}
	if _, ok := p.Ailments["paralysis"]; ok {
		denominator *= 4
	}

	base = base * numerator / denominator

	if stage >= 0 {
		return base * (2 + stage) / 2
	}
	return base * 2 / (2 - stage)
}

func (p *Pokemon) isFasterThan(bs battleState, mon *Pokemon) bool {
	return p.effectiveSpeed(bs) >= mon.effectiveSpeed(bs)
}

func (p *Pokemon) evasionFraction(keenEye bool) (int, int) {
	if keenEye {
		return 1, 1
	}

	stage := p.Stages["evasion"]
	if stage == 0 {
		return 3, 3
	} else if stage > 0 {
		return 3, 3 + stage
	}
	return 3 - stage, 3
}

func (p *Pokemon) accuracyFraction() (int, int) {
	stage := p.Stages["accuracy"]
	if stage == 0 {
		return 3, 3
	} else if stage > 0 {
		return 3 + stage, 3
	}
	return 3, 3 - stage
}

func (p *Pokemon) hasType(typeName string) bool {
	for _, t := range p.Base.Types {
		if typeName == t {
			return true
		}
	}
	return false
}

func (p *Pokemon) applyAilment(ailment string, move *pokeapi.BaseMove, afflictedBy *slot) {
	if _, ok := volatileStatuses[ailment]; !ok {
		if _, ok := nonVolatileStatuses[ailment]; !ok {
			panic("invalid ailment")
		}
	}

	if _, ok := p.Ailments[ailment]; ok {
		return
	}
	if _, ok := nonVolatileStatuses[ailment]; ok && p.hasNonVolatileAilment() {
		return
	}
	if ailment == "burn" && (p.hasType("fire") || p.Ability == "water-veil") {
		return
	}
	if ailment == "paralysis" && (p.hasType("electric") || p.Ability == "limber") {
		return
	}
	if ailment == "poison" || ailment == "toxic" {
		if p.Ability == "immunity" {
			return
		}
		if (p.hasType("poison") || p.hasType("steel")) && (afflictedBy == nil || afflictedBy.mon.Ability != "corrosion") {
			return
		}
	}
	if ailment == "freeze" && p.hasType("ice") {
		return
	}
	if ailment == "sleep" {
		if _, ok := sleepBlockingAbilities[p.Ability]; ok {
			return
		}
	}
	if ailment == "yawn" {
		if _, ok := sleepBlockingAbilities[p.Ability]; ok || p.hasNonVolatileAilment() {
			return
		}
	}

	switch ailment {
	case "trap":
		p.Ailments[ailment] = generateTrap(move.MinTurns, move.MaxTurns, afflictedBy)
		return
	case "poison":
		if move != nil && (move.Name == "toxic" || move.Name == "poison-fang") {
			ailment = "toxic"
		}
	case "infatuation":
		if afflictedBy.mon.Ability == "oblivious" {
			return
		}
	}

	p.Ailments[ailment] = generateAilment(ailment, afflictedBy)
	log.Printf("%s became afflicted with %s", p.Base.Name, ailment)
	if _, ok := nonVolatileStatuses[ailment]; ok && p.Ability == "synchronize" {
		afflictedBy.mon.applyAilment(ailment, nil, nil)
	}
	p.checkItemTrigger(true, nil)
}

func (p *Pokemon) hasAilment(ailment string) *Ailment {
	if a, ok := p.Ailments[ailment]; ok {
		return a
	}
	return nil
}

func (p *Pokemon) hasNonVolatileAilment() bool {
	for ailment := range nonVolatileStatuses {
		if p.hasAilment(ailment) != nil {
			return true
		}
	}
	return false
}

func (p *Pokemon) isGrounded() bool {
	if p.Item.name == "iron-ball" {
		return true
	}
	if p.hasType("flying") || p.Ability == "levitate" {
		return false
	}
	return true
}

func (p *Pokemon) changeHpBy(change int) {
	p.Hp = min(p.Hp+change, p.maxHP())
	p.checkItemTrigger(true, nil)
}

func (p *Pokemon) hasMovePredicate(f func(*pokeapi.BaseMove) bool) bool {
	return slices.ContainsFunc(p.Moves, f)
}

func (p *Pokemon) changeStatStageBy(stat string, change int, offensive bool) {
	if offensive && (p.Ability == "clear-smoke" || p.Ability == "clear-body") {
		log.Printf("blocked by clear body")
		return
	}
	if p.Ability == "keen-eye" && stat == "accuracy" && change < 0 {
		return
	}

	p.Stages[stat] = max(-6, min(6, p.Stages[stat]+change))
	log.Printf("%s's %s changed by %d stages (%d)", p.Base.Name, stat, change, p.Stages[stat])
}

func (p *Pokemon) maxHP() int {
	return p.Stats["hp"]
}

func (p *Pokemon) serenceGraceBonus() int {
	if p.Ability == "serence-grace" {
		return 2
	}
	return 1
}

func (p *Pokemon) applyMoveType(num, dem int, moveType string) (int, int) {
	for _, t := range p.Base.Types {
		if t == "flying" && moveType == "ground" && p.isGrounded() {
			continue
		}
		switch pokemon.GetEffectiveness(moveType, t) {
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
