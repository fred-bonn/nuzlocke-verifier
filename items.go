package main

import (
	"fmt"
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
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
			mon.ChangeHp(10)
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
			mon.ChangeHp(restore)
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

func makeGemMiddleWare(name, typeName string) func(mon *Pokemon) *item {
	var n *int
	return func(mon *Pokemon) *item {
		return &item{
			name: name,
			trigger: func(e any) bool {
				event, ok := e.(gemEvent)
				if !ok {
					return false
				}
				n = event.numerator
				return event.typeName == typeName
			},
			activate: func() {
				if n == nil {
					log.Printf("%s consumed its %s and boosted the damage", mon.Base.Name, name)
				} else {
					*n *= 2
				}
			},
		}
	}
}

func makeLumBerry(mon *Pokemon) *item {
	return &item{
		trigger: func(any) bool {
			return mon.HasNonVolatileAilment() || mon.HasAilment("confusion")
		},
		activate: func() {
			log.Printf("%s ate its lum berry", mon.Base.Name)
			for ailment := range pokemon.NonVolatileStatuses {
				if mon.HasAilment(ailment) {
					delete(mon.Ailments, ailment)
					log.Printf("%s had its %s removed", mon.Base.Name, ailment)
				}
			}
			if mon.HasAilment("confusion") {
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
