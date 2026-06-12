package main

import (
	"fmt"
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

type item struct {
	name     string
	trigger  func(any) bool
	activate func()
	consumed bool
}

func (p *Pokemon) checkItemTrigger(consume bool, event any) {
	if p.Item.trigger == nil || p.Item.consumed {
		return
	}

	if p.Item.trigger(event) {
		if consume {
			p.Item.consumed = true
			p.Unburden = true
		}
		p.Item.activate()
	}
}

type ItemFactoryBuilder func(*Pokemon) *item

var itemBuilders = map[string]ItemFactoryBuilder{
	"oran-berry":   makeOranBerry,
	"sitrus-berry": makeSitrusBerry,
	"lum-berry":    makeLumBerry,
	"leppa-berry":  makeLeppaBerry,
	"cheri-berry":  makeCheriBerry,
	"chesto-berry": makeChestoBerry,
	"pecha-berry":  makePechaBerry,
	"rawst-berry":  makeRawstBerry,
	"aspear-berry": makeAspearBerry,
	"liechi-berry": makeStatBoostBerryMiddleware("liechi-berry", "attack"),
	"ganlon-berry": makeStatBoostBerryMiddleware("ganlon-berry", "defense"),
	"salac-berry":  makeStatBoostBerryMiddleware("salac-berry", "speed"),
	"petaya-berry": makeStatBoostBerryMiddleware("petaya-berry", "special-attack"),
	"apicot-berry": makeStatBoostBerryMiddleware("apicot-berry", "special-defense"),
	"babiri-berry": makeResistBerryMiddleware("babiri-berry", "steel"),
	"chilan-berry": makeResistBerryMiddleware("chilan-berry", "normal"),
	"charti-berry": makeResistBerryMiddleware("charti-berry", "rock"),
	"chople-berry": makeResistBerryMiddleware("chople-berry", "fighting"),
	"coba-berry":   makeResistBerryMiddleware("coba-berry", "coba"),
	"colbur-berry": makeResistBerryMiddleware("colbur-berry", "dark"),
	"haban-berry":  makeResistBerryMiddleware("haban-berry", "dragon"),
	"kasib-berry":  makeResistBerryMiddleware("kasib-berry", "ghost"),
	"kebia-berry":  makeResistBerryMiddleware("kebia-berry", "poison"),
	"occa-berry":   makeResistBerryMiddleware("occa-berry", "fire"),
	"passho-berry": makeResistBerryMiddleware("passho-berry", "water"),
	"payapa-berry": makeResistBerryMiddleware("payapa-berry", "psychic"),
	"rindo-berry":  makeResistBerryMiddleware("rindo-berry", "grass"),
	"roseli-berry": makeResistBerryMiddleware("roseli-berry", "fairy"),
	"shuca-berry":  makeResistBerryMiddleware("shuca-berry", "ground"),
	"tanga-berry":  makeResistBerryMiddleware("tanga-berry", "bug"),
	"wacan-berry":  makeResistBerryMiddleware("wacan-berry", "electric"),
	"yache-berry":  makeResistBerryMiddleware("yache-berry", "ice"),
	"iron-ball":    makePassiveItemMiddleware("iron-ball"),
	"scope-lens":   makePassiveItemMiddleware("scope-lens"),
	"normal-gem":   makeGemMiddleware("normal"),
	"fire-gem":     makeGemMiddleware("fire"),
	"fighting-gem": makeGemMiddleware("fighting"),
	"water-gem":    makeGemMiddleware("water"),
	"flying-gem":   makeGemMiddleware("flying"),
	"grass-gem":    makeGemMiddleware("grass"),
	"poison-gem":   makeGemMiddleware("poison"),
	"electric-gem": makeGemMiddleware("electric"),
	"ground-gem":   makeGemMiddleware("ground"),
	"psychic-gem":  makeGemMiddleware("psychic"),
	"rock-gem":     makeGemMiddleware("rock"),
	"ice-gem":      makeGemMiddleware("ice"),
	"bug-gem":      makeGemMiddleware("bug"),
	"dragon-gem":   makeGemMiddleware("dragon"),
	"ghost-gem":    makeGemMiddleware("ghost"),
	"dark-gem":     makeGemMiddleware("dark"),
	"steel-gem":    makeGemMiddleware("steel"),
	"fairy-gem":    makeGemMiddleware("fairy"),
	"choice-scarf": makeChoiceScarf,
	"assault-vest": makeAssaultVest,
	"choice-band":  makeChoiceBand,
	"choice-specs": makeChoiceSpecs,
	"focus-sash":   makeFocusSash,
}

func createItemFactory(builder ItemFactoryBuilder, mon *Pokemon) func() *item {
	return func() *item {
		return builder(mon)
	}
}

func registerItem(itemName string, mon *Pokemon) (*item, error) {
	if itemName == "" {
		return &item{
			consumed: true,
		}, nil
	}

	builder, ok := itemBuilders[itemName]
	if !ok {
		return nil, fmt.Errorf("invalid item: %s", itemName)
	}

	factory := createItemFactory(builder, mon)
	return factory(), nil
}

func checkItemTriggers(bs battleState, e any) {
	for _, slot := range bs.getAllSlots() {
		slot.mon.checkItemTrigger(true, e)
	}
}

func makePassiveItemMiddleware(itemName string) func(mon *Pokemon) *item {
	return func(mon *Pokemon) *item {
		return &item{
			name: itemName,
			trigger: func(e any) bool {
				return false
			},
		}
	}
}

func makeOranBerry(mon *Pokemon) *item {
	return &item{
		name: "oran-berry",
		trigger: func(any) bool {
			return !mon.Unnerved && mon.Hp > 0 && mon.Hp*2 <= mon.maxHP()
		},
		activate: func() {
			mon.changeHpBy(10)
			log.Printf("%s ate its berry and restored 10 hp", mon.Base.Name)
			cheekPouch(mon)
		},
	}
}

func makeSitrusBerry(mon *Pokemon) *item {
	return &item{
		name: "sitrus-berry",
		trigger: func(any) bool {
			return !mon.Unnerved && mon.Hp > 0 && mon.Hp*2 <= mon.maxHP()
		},
		activate: func() {
			restore := mon.maxHP() / 4
			mon.changeHpBy(restore)
			log.Printf("%s ate its berry and restored %d hp", mon.Base.Name, restore)
			cheekPouch(mon)
		},
	}
}

func makeCheriBerry(mon *Pokemon) *item {
	return &item{
		name: "cheri-berry",
		trigger: func(any) bool {
			return !mon.Unnerved && mon.hasAilment("paralysis") != nil
		},
		activate: func() {
			log.Printf("%s ate its cheri berry", mon.Base.Name)
			delete(mon.Ailments, "paralysis")
			cheekPouch(mon)
		},
	}
}

func makeChestoBerry(mon *Pokemon) *item {
	return &item{
		name: "chesto-berry",
		trigger: func(any) bool {
			return !mon.Unnerved && mon.hasAilment("sleep") != nil
		},
		activate: func() {
			log.Printf("%s ate its chesto berry", mon.Base.Name)
			delete(mon.Ailments, "sleep")
			cheekPouch(mon)
		},
	}
}

func makePechaBerry(mon *Pokemon) *item {
	return &item{
		name: "pecha-berry",
		trigger: func(any) bool {
			return !mon.Unnerved && (mon.hasAilment("poison") != nil || mon.hasAilment("toxic") != nil)
		},
		activate: func() {
			log.Printf("%s ate its pecha berry", mon.Base.Name)
			delete(mon.Ailments, "poison")
			delete(mon.Ailments, "toxic")
			cheekPouch(mon)
		},
	}
}

func makeRawstBerry(mon *Pokemon) *item {
	return &item{
		name: "rawst-berry",
		trigger: func(any) bool {
			return !mon.Unnerved && mon.hasAilment("burn") != nil
		},
		activate: func() {
			log.Printf("%s ate its rawst berry", mon.Base.Name)
			delete(mon.Ailments, "burn")
			cheekPouch(mon)
		},
	}
}

func makeAspearBerry(mon *Pokemon) *item {
	return &item{
		name: "aspear-berry",
		trigger: func(any) bool {
			return !mon.Unnerved && mon.hasAilment("bufreezern") != nil
		},
		activate: func() {
			log.Printf("%s ate its aspear berry", mon.Base.Name)
			delete(mon.Ailments, "freeze")
			cheekPouch(mon)
		},
	}
}

func makeLumBerry(mon *Pokemon) *item {
	return &item{
		name: "lum-berry",
		trigger: func(any) bool {
			return !mon.Unnerved && (mon.hasNonVolatileAilment() || mon.hasAilment("confusion") != nil)
		},
		activate: func() {
			log.Printf("%s ate its lum berry", mon.Base.Name)
			for ailment := range nonVolatileStatuses {
				if mon.hasAilment(ailment) != nil {
					delete(mon.Ailments, ailment)
					log.Printf("%s had its %s removed", mon.Base.Name, ailment)
				}
			}
			if mon.hasAilment("confusion") != nil {
				delete(mon.Ailments, "confusion")
				log.Printf("%s had its confusion removed", mon.Base.Name)
				cheekPouch(mon)
			}
		},
	}
}

func makeLeppaBerry(mon *Pokemon) *item {
	var m *pokeapi.BaseMove
	return &item{
		trigger: func(e any) bool {
			event, ok := e.(leppaBerryEvent)
			if !ok || mon.Unnerved {
				return false
			}
			m = event.move
			return event.move.PP <= 0
		},
		activate: func() {
			m.PP += min(10, m.MaxPP)
			cheekPouch(mon)
		},
	}
}

func makeStatBoostBerryMiddleware(name, stat string) func(mon *Pokemon) *item {
	return func(mon *Pokemon) *item {
		d := 4
		if mon.Ability == "gluttony" {
			d = 2
		}
		return &item{
			name: name,
			trigger: func(any) bool {
				if mon.Ability == "gluttony" {

				}
				return !mon.Unnerved && mon.Hp > 0 && mon.Hp*d <= mon.maxHP()
			},
			activate: func() {
				log.Printf("%s ate its %s", mon.Base.Name, name)
				mon.changeStatStageBy(stat, 1, false)
				cheekPouch(mon)
			},
		}
	}
}

func makeResistBerryMiddleware(name, typeName string) func(mon *Pokemon) *item {
	var d *int
	return func(mon *Pokemon) *item {
		return &item{
			name: name,
			trigger: func(e any) bool {
				event, ok := e.(resistBerryEvent)
				if !ok || mon.Unnerved {
					return false
				}
				d = event.denominator
				return event.typeName == typeName
			},
			activate: func() {
				if d == nil {
					log.Printf("%s ate its %s and reduced the damage", mon.Base.Name, name)
					cheekPouch(mon)
				} else {
					*d *= 2
				}
			},
		}
	}
}

func makeGemMiddleware(typeName string) func(mon *Pokemon) *item {
	var d *int
	var n *int
	return func(mon *Pokemon) *item {
		return &item{
			name: fmt.Sprintf("%s-gem", typeName),
			trigger: func(e any) bool {
				event, ok := e.(gemEvent)
				if !ok {
					return false
				}
				d = event.denominator
				n = event.numerator
				return event.typeName == typeName
			},
			activate: func() {
				if n == nil {
					log.Printf("%s consumed its %s gem and boosted the damage", mon.Base.Name, typeName)
				} else {
					*n *= 3
					*d *= 2
				}
			},
		}
	}
}

func makeChoiceScarf(mon *Pokemon) *item {
	mon.Stats["speed"] = mon.Stats["speed"] * 3 / 2
	return &item{
		name: "choice-scarf",
	}
}

func makeAssaultVest(mon *Pokemon) *item {
	mon.Stats["special-defense"] = mon.Stats["special-defense"] * 3 / 2
	return &item{
		name: "assault-vest",
	}
}

func makeChoiceBand(mon *Pokemon) *item {
	var d *int
	var n *int
	return &item{
		name: "choice-band",
		trigger: func(e any) bool {
			event, ok := e.(choiceItemEvent)
			if !ok {
				return false
			}
			if event.move.Class != "physical" {
				return false
			}
			d = event.denominator
			n = event.numerator
			return true
		},
		activate: func() {
			*n *= 3
			*d *= 2
		},
	}
}

func makeChoiceSpecs(mon *Pokemon) *item {
	var d *int
	var n *int
	return &item{
		name: "choice-specs",
		trigger: func(e any) bool {
			event, ok := e.(choiceItemEvent)
			if !ok {
				return false
			}
			if event.move.Class != "special" {
				return false
			}
			d = event.denominator
			n = event.numerator
			return true
		},
		activate: func() {
			*n *= 3
			*d *= 2
		},
	}
}

func makeFocusSash(mon *Pokemon) *item {
	var dmg *int
	return &item{
		trigger: func(e any) bool {
			event, ok := e.(focusSashEvent)
			if !ok {
				return false
			}
			dmg = event.damage
			return mon.Hp == mon.maxHP() && *event.damage >= mon.Hp
		},
		activate: func() {
			*dmg = mon.Hp - 1
		},
	}
}
