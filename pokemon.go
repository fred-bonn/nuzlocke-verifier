package main

import (
	"fmt"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type Pokemon struct {
	Base     pokeapi.BasePokemon
	Level    int
	IVs      map[string]int
	Nature   []string
	Moves    []pokeapi.BaseMove
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

func InitializePokemon(base pokeapi.BasePokemon, level int, ivs map[string]int, nature string, moves []pokeapi.BaseMove, hp int, status string) (Pokemon, error) {
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

func (p *Pokemon) String() string {
	status := "Status: "
	if len(p.Ailments) == 0 {
		status += "None"
	} else {
		for ailment, duration := range p.Ailments {
			status += fmt.Sprintf("%s (%d turns) ", ailment, duration)
		}
	}
	ivs := fmt.Sprintf("IVS:    (%d, %d, %d, %d, %d, %d)", p.IVs["hp"], p.IVs["attack"], p.IVs["defense"], p.IVs["special-attack"], p.IVs["special-defense"], p.IVs["speed"])
	stats := fmt.Sprintf("Stats:  (%d, %d, %d, %d, %d, %d)", p.Stats["hp"], p.Stats["attack"], p.Stats["defense"], p.Stats["special-attack"], p.Stats["special-defense"], p.Stats["speed"])
	stages := fmt.Sprintf("Stages: (%d, %d, %d, %d, %d, %d, %d)", p.Stages["attack"], p.Stages["defense"], p.Stages["special-attack"], p.Stages["special-defense"], p.Stages["speed"], p.Stages["accuracy"], p.Stages["evasion"])

	return fmt.Sprintf("%s (Level %d) - HP: %d/%d\nNature: (%s,%s)\n%s\n%s\n%s\n%s", p.Base.Name, p.Level, p.Hp, p.Stats["hp"], p.Nature[0], p.Nature[1], status, ivs, stats, stages)
}

func (p *Pokemon) ChangeStatStage(stat string, change int) {
	if _, ok := p.Stages[stat]; !ok {
		panic("invalid stat")
	}
	p.Stages[stat] = max(-6, min(6, p.Stages[stat]+change))
}

func (p *Pokemon) EffectiveStat(stat string, crit bool) int {
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

func (p *Pokemon) IsFasterThan(mon *Pokemon) bool {
	return p.EffectiveSpeed() >= mon.EffectiveSpeed()
}

func (p *Pokemon) EffectiveSpeed() int {
	stage := p.Stages["speed"]
	base := p.Stats["speed"]
	if _, ok := p.Ailments["paralysis"]; ok {
		base = base / 4
	}

	if stage >= 0 {
		return base * (2 + stage) / 2
	}
	return base * 2 / (2 - stage)
}

func (p *Pokemon) EvasionFraction() (int, int) {
	stage := p.Stages["evasion"]
	if stage == 0 {
		return 3, 3
	} else if stage > 0 {
		return 3, 3 + stage
	}
	return 3 - stage, 3
}

func (p *Pokemon) AccuracyFraction() (int, int) {
	stage := p.Stages["accuracy"]
	if stage == 0 {
		return 3, 3
	} else if stage > 0 {
		return 3 + stage, 3
	}
	return 3, 3 - stage
}

func (p *Pokemon) SwitchReset() {
	delete(p.Ailments, "confusion")
	for stat := range p.Stages {
		p.Stages[stat] = 0
	}
	if _, ok := p.Ailments["toxic"]; ok {
		p.Ailments["toxic"] = 0
	}
}

func (p *Pokemon) HasType(typeName string) bool {
	for _, t := range p.Base.Types {
		if typeName == t {
			return true
		}
	}
	return false
}

func (p *Pokemon) ApplyAilment(ailment string, move *pokeapi.BaseMove) bool {
	if _, ok := pokemon.ValidAilments[ailment]; !ok {
		panic("invalid ailment")
	}
	if _, ok := p.Ailments[ailment]; ok {
		return false
	}
	if _, ok := pokemon.NonVolatileStatuses[ailment]; ok {
		if ailment == "burn" && p.HasType("fire") {
			return false
		}
		if ailment == "paralysis" && p.HasType("electric") {
			return false
		}
		if (ailment == "toxic" || ailment == "poison") && p.HasType("poison") {
			return false
		}
		for a := range p.Ailments {
			if _, ok := pokemon.NonVolatileStatuses[a]; ok {
				return false
			}
		}
	}

	if ailment == "trap" {
		p.Ailments[ailment] = pokemon.GenerateTrap(move.MinTurns, move.MaxTurns)
		return true
	}
	if ailment == "poison" && (move.Name == "toxic" || move.Name == "poison-fang") {
		ailment = "toxic"
	}
	p.Ailments[ailment] = pokemon.GenerateAilment(ailment)
	return true
}

func (p *Pokemon) HasAilment(ailment string) bool {
	if _, ok := p.Ailments[ailment]; ok {
		return true
	}
	return false
}

func (p *Pokemon) ChangeHp(change int) {
	p.Hp = min(p.Hp+change, p.Stats["hp"])
	if p.Item != nil {
		p.Item.checkTrigger(true, nil)
	}
}
