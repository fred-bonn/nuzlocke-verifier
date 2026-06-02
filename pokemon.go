package main

import (
	"fmt"
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type Pokemon struct {
	Base     pokeapi.BasePokemon
	Level    int
	IVs      map[string]int
	Nature   []string
	Moves    []*pokeapi.BaseMove
	Stats    map[string]int
	Stages   map[string]int
	Hp       int
	Fainted  bool
	Ailments map[string]int
	Item     *item
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
		Ailments: make(map[string]int),
	}

	setIVs(&res, ivs)

	calculateStats(&res)

	if hp == -1 {
		hp = res.Stats["hp"]
	}

	res.Hp = max(1, min(res.Stats["hp"], hp))

	if _, ok := pokemon.NonVolatileStatuses[status]; ok {
		res.Ailments[status] = pokemon.GenerateAilment(status)
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
	if pokemon.Base.Stats["hp"] == 1 {
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

func (p *Pokemon) changeStatStage(stat string, change int) {
	if _, ok := p.Stages[stat]; !ok {
		panic("invalid stat")
	}
	p.Stages[stat] = max(-6, min(6, p.Stages[stat]+change))
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

func (p *Pokemon) isFasterThan(mon *Pokemon) bool {
	return p.effectiveSpeed() >= mon.effectiveSpeed()
}

func (p *Pokemon) effectiveSpeed() int {
	stage := p.Stages["speed"]
	base := p.Stats["speed"]
	if p.Item.name == "iron-ball" {
		base /= 2
	}
	if _, ok := p.Ailments["paralysis"]; ok {
		base /= 4
	}

	if stage >= 0 {
		return base * (2 + stage) / 2
	}
	return base * 2 / (2 - stage)
}

func (p *Pokemon) evasionFraction() (int, int) {
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

func (p *Pokemon) switchReset() {
	delete(p.Ailments, "confusion")
	for stat := range p.Stages {
		p.Stages[stat] = 0
	}
	if _, ok := p.Ailments["toxic"]; ok {
		p.Ailments["toxic"] = 0
	}
}

func (p *Pokemon) hasType(typeName string) bool {
	for _, t := range p.Base.Types {
		if typeName == t {
			return true
		}
	}
	return false
}

func (p *Pokemon) applyAilment(ailment string, move *pokeapi.BaseMove) {
	if _, ok := pokemon.ValidAilments[ailment]; !ok {
		panic("invalid ailment")
	}
	if _, ok := p.Ailments[ailment]; ok {
		return
	}
	if _, ok := pokemon.NonVolatileStatuses[ailment]; ok {
		if ailment == "burn" && p.hasType("fire") {
			return
		}
		if ailment == "paralysis" && p.hasType("electric") {
			return
		}
		if ailment == "poison" && p.hasType("poison") {
			return
		}
		for a := range p.Ailments {
			if _, ok := pokemon.NonVolatileStatuses[a]; ok {
				return
			}
		}
	}

	if ailment == "trap" {
		p.Ailments[ailment] = pokemon.GenerateTrap(move.MinTurns, move.MaxTurns)
		return
	}
	if ailment == "poison" && (move.Name == "toxic" || move.Name == "poison-fang") {
		ailment = "toxic"
	}
	p.Ailments[ailment] = pokemon.GenerateAilment(ailment)
	log.Printf("%s became afflicted with %s", p.Base.Name, ailment)
	p.Item.checkTrigger(true, nil)
}

func (p *Pokemon) hasAilment(ailment string) bool {
	if _, ok := p.Ailments[ailment]; ok {
		return true
	}
	return false
}

func (p *Pokemon) hasNonVolatileAilment() bool {
	for ailment := range pokemon.NonVolatileStatuses {
		if p.hasAilment(ailment) {
			return true
		}
	}
	return false
}

func (p *Pokemon) hasMoveClass(c string) bool {
	for _, m := range p.Moves {
		if m.Class == c {
			return true
		}
	}
	return false
}

func (p *Pokemon) isGrounded() bool {
	if p.hasType("flying") && p.Item.name != "iron-ball" {
		return false
	}
	return true
}

func (p *Pokemon) changeHp(change int) {
	p.Hp = min(p.Hp+change, p.Stats["hp"])
	if p.Item != nil {
		p.Item.checkTrigger(true, nil)
	}
}
