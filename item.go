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

func (i *item) checkTrigger(consume bool, event any) {
	if i.trigger == nil || i.consumed {
		return
	}

	if i.trigger(event) {
		if consume {
			i.consumed = true
		}
		i.activate()
	}
}

type ItemFactoryBuilder func(*Pokemon) *item

var itemBuilders = map[string]ItemFactoryBuilder{
	"oran-berry":   makeOranBerry,
	"sitrus-berry": makeSitrusBerry,
	"lum-berry":    makeLumBerry,
	"leppa-berry":  makeLeppaBerry,
	"babiri-berry": makeResistBerryMiddleWare("babiri-berry", "steel"),
	"chilan-berry": makeResistBerryMiddleWare("chilan-berry", "normal"),
	"charti-berry": makeResistBerryMiddleWare("charti-berry", "rock"),
	"chople-berry": makeResistBerryMiddleWare("chople-berry", "fighting"),
	"coba-berry":   makeResistBerryMiddleWare("coba-berry", "coba"),
	"colbur-berry": makeResistBerryMiddleWare("colbur-berry", "dark"),
	"haban-berry":  makeResistBerryMiddleWare("haban-berry", "dragon"),
	"kasib-berry":  makeResistBerryMiddleWare("kasib-berry", "ghost"),
	"kebia-berry":  makeResistBerryMiddleWare("kebia-berry", "poison"),
	"occa-berry":   makeResistBerryMiddleWare("occa-berry", "fire"),
	"passho-berry": makeResistBerryMiddleWare("passho-berry", "water"),
	"payapa-berry": makeResistBerryMiddleWare("payapa-berry", "psychic"),
	"rindo-berry":  makeResistBerryMiddleWare("rindo-berry", "grass"),
	"roseli-berry": makeResistBerryMiddleWare("roseli-berry", "fairy"),
	"shuca-berry":  makeResistBerryMiddleWare("shuca-berry", "ground"),
	"tanga-berry":  makeResistBerryMiddleWare("tanga-berry", "bug"),
	"wacan-berry":  makeResistBerryMiddleWare("wacan-berry", "electric"),
	"yache-berry":  makeResistBerryMiddleWare("yache-berry", "ice"),
	"iron-ball":    makePassiveItemMiddleWare("iron-ball"),
	"normal-gem":   makeGemMiddleWare("normal"),
	"fire-gem":     makeGemMiddleWare("fire"),
	"fighting-gem": makeGemMiddleWare("fighting"),
	"water-gem":    makeGemMiddleWare("water"),
	"flying-gem":   makeGemMiddleWare("flying"),
	"grass-gem":    makeGemMiddleWare("grass"),
	"poison-gem":   makeGemMiddleWare("poison"),
	"electric-gem": makeGemMiddleWare("electric"),
	"ground-gem":   makeGemMiddleWare("ground"),
	"psychic-gem":  makeGemMiddleWare("psychic"),
	"rock-gem":     makeGemMiddleWare("rock"),
	"ice-gem":      makeGemMiddleWare("ice"),
	"bug-gem":      makeGemMiddleWare("bug"),
	"dragon-gem":   makeGemMiddleWare("dragon"),
	"ghost-gem":    makeGemMiddleWare("ghost"),
	"dark-gem":     makeGemMiddleWare("dark"),
	"steel-gem":    makeGemMiddleWare("steel"),
	"fairy-gem":    makeGemMiddleWare("fairy"),
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
	slots := bs.getAllSlots()
	for _, slot := range slots {
		item := slot.mon.Item
		if slot.mon.Item == nil {
			continue
		}

		item.checkTrigger(true, e)
	}
}

func makeOranBerry(mon *Pokemon) *item {
	return &item{
		name: "oran-berry",
		trigger: func(any) bool {
			return mon.Hp > 0 && mon.Hp <= mon.Stats["hp"]/2
		},
		activate: func() {
			mon.changeHp(10)
			log.Printf("%s ate its berry and restored 10 hp", mon.Base.Name)
		},
	}
}

func makeSitrusBerry(mon *Pokemon) *item {
	return &item{
		name: "sitrus-berry",
		trigger: func(any) bool {
			return mon.Hp > 0 && mon.Hp <= mon.Stats["hp"]/2
		},
		activate: func() {
			restore := mon.Stats["hp"] / 4
			mon.changeHp(restore)
			log.Printf("%s ate its berry and restored %d hp", mon.Base.Name, restore)
		},
	}
}

func makeResistBerryMiddleWare(name, typeName string) func(mon *Pokemon) *item {
	var d *int
	return func(mon *Pokemon) *item {
		return &item{
			name: name,
			trigger: func(e any) bool {
				event, ok := e.(resistBerryEvent)
				if !ok {
					return false
				}
				d = event.denominator
				return event.typeName == typeName
			},
			activate: func() {
				if d == nil {
					log.Printf("%s ate its %s and reduced the damage", mon.Base.Name, name)
				} else {
					*d *= 2
				}
			},
		}
	}
}

func makeGemMiddleWare(typeName string) func(mon *Pokemon) *item {
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

func makeLumBerry(mon *Pokemon) *item {
	return &item{
		trigger: func(any) bool {
			return mon.hasNonVolatileAilment() || mon.hasAilment("confusion") != nil
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
			}
		},
	}
}

func makeLeppaBerry(mon *Pokemon) *item {
	var m *pokeapi.BaseMove
	return &item{
		trigger: func(e any) bool {
			event, ok := e.(leppaBerryEvent)
			if !ok {
				return false
			}
			m = event.move
			return event.move.PP <= 0
		},
		activate: func() {
			m.PP += min(10, m.MaxPP)
		},
	}
}

func makePassiveItemMiddleWare(itemName string) func(mon *Pokemon) *item {
	return func(mon *Pokemon) *item {
		return &item{
			name: itemName,
			trigger: func(e any) bool {
				return false
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
	var c bool
	return &item{
		trigger: func(e any) bool {
			event, ok := e.(focusSashEvent)
			if !ok {
				return false
			}
			dmg = event.damage
			c = event.consume
			return mon.Hp == mon.Stats["hp"] && *event.damage >= mon.Hp
		},
		activate: func() {
			if c {
				log.Printf("%s held on with its focus sash", mon.Base.Name)
			}
			*dmg = mon.Hp - 1
			dmg = nil
		},
	}
}
