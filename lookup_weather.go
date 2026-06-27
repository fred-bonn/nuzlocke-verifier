package main

type weatherState int

const (
	rainWeather weatherState = iota
	sunWeather
	sandstormWeather
	hailWeather
	noneWeather
)

var weatherFuncs = map[weatherState]func(*int, *int, pokemonType){
	rainWeather: func(num, den *int, t pokemonType) {
		switch t {
		case waterType:
			*num = *num * 3
			*den = *den * 2
		case fireType:
			*den = *den * 2
		}
	},
	sunWeather: func(num, den *int, t pokemonType) {
		switch t {
		case waterType:
			*den = *den * 2
		case fireType:
			*num = *num * 3
			*den = *den * 2
		}
	},
}

func (ws weatherState) affectsMon(mon *pokemon) bool {
	if mon.ability == "overcoat" || mon.ability == "magic-guard" || mon.item.name == "safety-goggles" {
		return false
	}

	switch ws {
	case sandstormWeather:
		if !mon.hasType(rockType) && !mon.hasType(steelType) && !mon.hasType(groundType) && mon.ability != "sand-veil" && mon.ability != "sand-rush" && mon.ability != "sand-force" {
			return true
		}
	case hailWeather:
		if mon.ability == "ice-body" {
			return false
		} else if !mon.hasType(iceType) && mon.ability != "snow-cloak" {
			return true
		}
	}

	return false
}

func (ws weatherState) activateMonAbility(bs battleState, slot *slot) {
	mon := slot.mon

	switch ws {
	case hailWeather:
		if mon.ability == "ice-body" {
			vlogf("%s healed due to ice body", mon.base.Name)
			mon.changeHpBy(mon.maxHP() / 16)
		}
	case rainWeather:
		switch mon.ability {
		case "rain-dish":
			vlogf("%s healed due to rain dish", mon.base.Name)
			mon.changeHpBy(mon.maxHP() / 16)
		case "dry-skin":
			takeResidualDamage(bs, slot, "dry skin", 1, 8)
			vlogf("%s healed due to dry skin", mon.base.Name)
			mon.changeHpBy(mon.maxHP() / 8)
		case "hydration":
			if mon.hasNonVolatileAilment() {
				for ailment := range nonVolatileStatuses {
					if mon.hasAilment(ailment) != nil {
						delete(mon.ailments, ailment)
						vlogf("%s had its %s removed", mon.base.Name, ailment.String())
						return
					}
				}
			}
		}
	case sunWeather:
		switch mon.ability {
		case "dry-skin":
			takeResidualDamage(bs, slot, "dry skin", 1, 8)
		case "solar-power":
			takeResidualDamage(bs, slot, "solar power", 1, 8)
		}
	}
}

func (ws weatherState) String() string {
	switch ws {
	case rainWeather:
		return "rain"
	case sunWeather:
		return "sun"
	case hailWeather:
		return "hail"
	case sandstormWeather:
		return "sandstorm"
	default:
		return ""
	}
}

func (ws weatherState) onset() {
	switch ws {
	case rainWeather:
		vlogln("it started to rain")
	case sunWeather:
		vlogln("the sunlight turned harsh")
	case sandstormWeather:
		vlogln("a sandstorm brewed")
	case hailWeather:
		vlogln("it started to hail")
	}
}
