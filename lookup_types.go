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

type effectiveness float64

const (
	immune  effectiveness = 0.0
	notVery effectiveness = 0.5
	normal  effectiveness = 1.0
	super   effectiveness = 2.0
)

var typeChart = map[pokemonType]map[pokemonType]effectiveness{
	normalType: {
		rockType: notVery, ghostType: immune, steelType: notVery,
	},
	fireType: {
		fireType: notVery, waterType: notVery, grassType: super, iceType: super,
		bugType: super, rockType: notVery, dragonType: notVery, steelType: super,
	},
	waterType: {
		fireType: super, waterType: notVery, grassType: notVery, groundType: super,
		rockType: super, dragonType: notVery,
	},
	electricType: {
		waterType: super, electricType: notVery, grassType: notVery, groundType: immune,
		flyingType: super, dragonType: notVery,
	},
	grassType: {
		fireType: notVery, waterType: super, grassType: notVery, poisonType: notVery,
		groundType: super, flyingType: notVery, bugType: notVery, rockType: super,
		dragonType: notVery, steelType: notVery,
	},
	iceType: {
		fireType: notVery, waterType: notVery, grassType: super, groundType: super,
		flyingType: super, dragonType: super, steelType: notVery, iceType: notVery,
	},
	fightingType: {
		normalType: super, iceType: super, rockType: super, darkType: super,
		steelType: super, poisonType: notVery, flyingType: notVery, psychicType: notVery,
		bugType: notVery, fairyType: notVery, ghostType: immune,
	},
	poisonType: {
		grassType: super, fairyType: super, poisonType: notVery, groundType: notVery,
		rockType: notVery, ghostType: notVery, steelType: immune,
	},
	groundType: {
		fireType: super, electricType: super, grassType: notVery, poisonType: super,
		flyingType: immune, bugType: notVery, rockType: super, steelType: super,
	},
	flyingType: {
		electricType: notVery, grassType: super, fightingType: super,
		bugType: super, rockType: notVery, steelType: notVery,
	},
	psychicType: {
		fightingType: super, poisonType: super, psychicType: notVery,
		steelType: notVery, darkType: immune,
	},
	bugType: {
		fireType: notVery, grassType: super, fightingType: notVery, poisonType: notVery,
		flyingType: notVery, psychicType: super, ghostType: notVery,
		darkType: super, steelType: notVery, fairyType: notVery,
	},
	rockType: {
		fireType: super, iceType: super, flyingType: super, bugType: super,
		fightingType: notVery, groundType: notVery, steelType: notVery,
	},
	ghostType: {
		normalType: immune, psychicType: super, ghostType: super, darkType: notVery,
	},
	dragonType: {
		dragonType: super, steelType: notVery, fairyType: immune,
	},
	darkType: {
		fightingType: notVery, psychicType: super, ghostType: super,
		darkType: notVery, fairyType: notVery,
	},
	steelType: {
		fireType: notVery, waterType: notVery, electricType: notVery, iceType: super,
		rockType: super, fairyType: super, steelType: notVery,
	},
	fairyType: {
		fireType: notVery, fightingType: super, poisonType: notVery,
		dragonType: super, darkType: super, steelType: notVery,
	},
}

func getEffectiveness(attacking, defending pokemonType) effectiveness {
	if row, ok := typeChart[attacking]; ok {
		if eff, ok := row[defending]; ok {
			return eff
		}
	}
	return normal
}
