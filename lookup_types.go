package main

type pokemonType int

const (
	normalType pokemonType = iota
	fireType
	waterType
	electricType
	grassType
	iceType
	fightingType
	poisonType
	groundType
	flyingType
	psychicType
	bugType
	rockType
	ghostType
	dragonType
	darkType
	steelType
	fairyType
	noType
)

func stringToPokemonType(s string) pokemonType {
	switch s {
	case "normal":
		return normalType
	case "fire":
		return fireType
	case "water":
		return waterType
	case "electric":
		return electricType
	case "grass":
		return grassType
	case "ice":
		return iceType
	case "fighting":
		return fightingType
	case "poison":
		return poisonType
	case "ground":
		return groundType
	case "flying":
		return flyingType
	case "psychic":
		return psychicType
	case "bug":
		return bugType
	case "rock":
		return rockType
	case "ghost":
		return ghostType
	case "dragon":
		return dragonType
	case "dark":
		return darkType
	case "steel":
		return steelType
	case "fairy":
		return fairyType
	default:
		elogFatalf("%s is not a valid type", s)
		return noType
	}
}

func (pt pokemonType) String() string {
	switch pt {
	case normalType:
		return "normal"
	case fireType:
		return "fire"
	case waterType:
		return "water"
	case electricType:
		return "electric"
	case grassType:
		return "grass"
	case iceType:
		return "ice"
	case fightingType:
		return "fighting"
	case poisonType:
		return "poison"
	case groundType:
		return "ground"
	case flyingType:
		return "flying"
	case psychicType:
		return "psychic"
	case bugType:
		return "bug"
	case rockType:
		return "rock"
	case ghostType:
		return "ghost"
	case dragonType:
		return "dragon"
	case darkType:
		return "dark"
	case steelType:
		return "steel"
	case fairyType:
		return "fairy"
	default:
		return "no type"
	}
}

type Effectiveness float64

const (
	Immune  Effectiveness = 0.0
	NotVery Effectiveness = 0.5
	Normal  Effectiveness = 1.0
	Super   Effectiveness = 2.0
)

var typeChart = map[pokemonType]map[pokemonType]Effectiveness{
	normalType: {
		rockType: NotVery, ghostType: Immune, steelType: NotVery,
	},
	fireType: {
		fireType: NotVery, waterType: NotVery, grassType: Super, iceType: Super,
		bugType: Super, rockType: NotVery, dragonType: NotVery, steelType: Super,
	},
	waterType: {
		fireType: Super, waterType: NotVery, grassType: NotVery, groundType: Super,
		rockType: Super, dragonType: NotVery,
	},
	electricType: {
		waterType: Super, electricType: NotVery, grassType: NotVery, groundType: Immune,
		flyingType: Super, dragonType: NotVery,
	},
	grassType: {
		fireType: NotVery, waterType: Super, grassType: NotVery, poisonType: NotVery,
		groundType: Super, flyingType: NotVery, bugType: NotVery, rockType: Super,
		dragonType: NotVery, steelType: NotVery,
	},
	iceType: {
		fireType: NotVery, waterType: NotVery, grassType: Super, groundType: Super,
		flyingType: Super, dragonType: Super, steelType: NotVery, iceType: NotVery,
	},
	fightingType: {
		normalType: Super, iceType: Super, rockType: Super, darkType: Super,
		steelType: Super, poisonType: NotVery, flyingType: NotVery, psychicType: NotVery,
		bugType: NotVery, fairyType: NotVery, ghostType: Immune,
	},
	poisonType: {
		grassType: Super, fairyType: Super, poisonType: NotVery, groundType: NotVery,
		rockType: NotVery, ghostType: NotVery, steelType: Immune,
	},
	groundType: {
		fireType: Super, electricType: Super, grassType: NotVery, poisonType: Super,
		flyingType: Immune, bugType: NotVery, rockType: Super, steelType: Super,
	},
	flyingType: {
		electricType: NotVery, grassType: Super, fightingType: Super,
		bugType: Super, rockType: NotVery, steelType: NotVery,
	},
	psychicType: {
		fightingType: Super, poisonType: Super, psychicType: NotVery,
		steelType: NotVery, darkType: Immune,
	},
	bugType: {
		fireType: NotVery, grassType: Super, fightingType: NotVery, poisonType: NotVery,
		flyingType: NotVery, psychicType: Super, ghostType: NotVery,
		darkType: Super, steelType: NotVery, fairyType: NotVery,
	},
	rockType: {
		fireType: Super, iceType: Super, flyingType: Super, bugType: Super,
		fightingType: NotVery, groundType: NotVery, steelType: NotVery,
	},
	ghostType: {
		normalType: Immune, psychicType: Super, ghostType: Super, darkType: NotVery,
	},
	dragonType: {
		dragonType: Super, steelType: NotVery, fairyType: Immune,
	},
	darkType: {
		fightingType: NotVery, psychicType: Super, ghostType: Super,
		darkType: NotVery, fairyType: NotVery,
	},
	steelType: {
		fireType: NotVery, waterType: NotVery, electricType: NotVery, iceType: Super,
		rockType: Super, fairyType: Super, steelType: NotVery,
	},
	fairyType: {
		fireType: NotVery, fightingType: Super, poisonType: NotVery,
		dragonType: Super, darkType: Super, steelType: NotVery,
	},
}

func getEffectiveness(attacking, defending pokemonType) Effectiveness {
	if row, ok := typeChart[attacking]; ok {
		if eff, ok := row[defending]; ok {
			return eff
		}
	}
	return Normal
}
