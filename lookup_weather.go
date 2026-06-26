package main

type weatherState int

const (
	NoneWeather weatherState = iota
	Rain
	Sun
	Sandstorm
	Hail
)

var weatherFuncs = map[weatherState]func(*int, *int, string){
	Rain: func(num, den *int, moveType string) {
		switch moveType {
		case "water":
			*num = *num * 3
			*den = *den * 2
		case "fire":
			*den = *den * 2
		}
	},
	Sun: func(num, den *int, typeName string) {
		switch typeName {
		case "water":
			*den = *den * 2
		case "fire":
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
	case Sandstorm:
		if !mon.hasType("rock") && !mon.hasType("steel") && !mon.hasType("ground") && mon.ability != "sand-veil" && mon.ability != "sand-rush" && mon.ability != "sand-force" {
			return true
		}
	case Hail:
		if mon.ability == "ice-body" {
			return false
		} else if !mon.hasType("ice") && mon.ability != "snow-cloak" {
			return true
		}
	}

	return false
}

func (ws weatherState) activateMonAbility(bs battleState, slot *slot) {
	mon := slot.mon

	switch ws {
	case Hail:
		if mon.ability == "ice-body" {
			vlogf("%s healed due to ice body", mon.base.Name)
			mon.changeHpBy(mon.maxHP() / 16)
		}
	case Rain:
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
	case Sun:
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
	case Rain:
		return "rain"
	case Sun:
		return "sun"
	case Hail:
		return "hail"
	case Sandstorm:
		return "sandstorm"
	default:
		return ""
	}
}

func (ws weatherState) onset() {
	switch ws {
	case Rain:
		vlogln("it started to rain")
	case Sun:
		vlogln("the sunlight turned harsh")
	case Sandstorm:
		vlogln("a sandstorm brewed")
	case Hail:
		vlogln("it started to hail")
	}
}
