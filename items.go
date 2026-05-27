package main

import (
	"fmt"
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type item struct {
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
	"babiri-berry": makeResistBerryMiddleWare("steel"),
	"chilan-berry": makeResistBerryMiddleWare("normal"),
	"charti-berry": makeResistBerryMiddleWare("rock"),
	"chople-berry": makeResistBerryMiddleWare("fighting"),
	"coba-berry":   makeResistBerryMiddleWare("coba"),
	"colbur-berry": makeResistBerryMiddleWare("dark"),
	"haban-berry":  makeResistBerryMiddleWare("dragon"),
	"kasib-berry":  makeResistBerryMiddleWare("ghost"),
	"kebia-berry":  makeResistBerryMiddleWare("poison"),
	"occa-berry":   makeResistBerryMiddleWare("fire"),
	"passho-berry": makeResistBerryMiddleWare("water"),
	"payapa-berry": makeResistBerryMiddleWare("psychic"),
	"rindo-berry":  makeResistBerryMiddleWare("grass"),
	"roseli-berry": makeResistBerryMiddleWare("fairy"),
	"shuca-berry":  makeResistBerryMiddleWare("ground"),
	"tanga-berry":  makeResistBerryMiddleWare("bug"),
	"wacan-berry":  makeResistBerryMiddleWare("electric"),
	"yache-berry":  makeResistBerryMiddleWare("ice"),
}

func createItemFactory(builder ItemFactoryBuilder, mon *Pokemon) func() *item {
	return func() *item {
		return builder(mon)
	}
}

func registerItem(itemName string, mon *Pokemon) (*item, error) {
	if itemName == "" {
		return &item{}, nil
	}

	builder, ok := itemBuilders[itemName]
	if !ok {
		return nil, fmt.Errorf("invalid item: %s", itemName)
	}

	factory := createItemFactory(builder, mon)
	return factory(), nil
}

func checkItemTriggers(bs battleState, event any) {
	slots := bs.getAllSlots()
	for _, slot := range slots {
		item := slot.mon.Item
		if slot.mon.Item == nil {
			continue
		}

		item.checkTrigger(true, event)
	}
}

func makeOranBerry(mon *Pokemon) *item {
	return &item{
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

func makeResistBerryMiddleWare(typeName string) func(mon *Pokemon) *item {
	var d *int
	return func(mon *Pokemon) *item {
		return &item{
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
					log.Printf("%s ate its berry and reduced the damage", mon.Base.Name)
				} else {
					*d *= 2
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
