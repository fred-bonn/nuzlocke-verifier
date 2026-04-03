package pokemon

import (
	"fmt"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
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
}

func InitializePokemon(base pokeapi.BasePokemon, level int, ivs map[string]int, nature string, moves []pokeapi.BaseMove, hp int, status string) (Pokemon, error) {
	if level < 1 || level > 100 {
		return Pokemon{}, fmt.Errorf("invalid level: %d", level)
	}

	nat, err := getNature(nature)
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

	_, ok := validAilments[status]
	if ok && status != "" {
		res.Ailments[status] = GenerateAilment(status)
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
		pokemon.Stats[posNat] = int(float64(pokemon.Stats[posNat]) * 1.1)
		pokemon.Stats[negNat] = int(float64(pokemon.Stats[negNat]) * 0.9)
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
