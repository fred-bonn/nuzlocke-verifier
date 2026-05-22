package main

import (
	"fmt"
	"log"
)

type item struct {
	trigger  func(...any) bool
	activate func(...any)
	consumed bool
}

func (i *item) checkTrigger(consume bool, params ...any) {
	if i.trigger == nil || i.consumed {
		return
	}

	if i.trigger(params...) {
		if consume {
			i.consumed = true
		}
		i.activate(params...)
	}
}

type ItemFactoryBuilder func(battleState, *Pokemon) *item

var itemBuilders = map[string]ItemFactoryBuilder{
	"oran-berry":   makeOranBerry,
	"sitrus-berry": makeSitrusBerry,
	"chilan-berry": makeResistBerryMiddleWare("normal"),
}

func createItemFactory(builder ItemFactoryBuilder, mon *Pokemon) func(battleState) *item {
	return func(bs battleState) *item {
		return builder(bs, mon)
	}
}

func registerItem(bs battleState, itemName string, mon *Pokemon, params ...any) (*item, error) {
	if itemName == "" {
		return nil, nil
	}

	builder, ok := itemBuilders[itemName]
	if !ok {
		return nil, fmt.Errorf("invalid item: %s", itemName)
	}

	factory := createItemFactory(builder, mon)
	return factory(bs), nil
}

func checkItemTriggers(bs battleState, params ...any) {
	slots := bs.getAllSlots()
	for _, slot := range slots {
		item := slot.mon.Item
		if slot.mon.Item == nil {
			continue
		}

		item.checkTrigger(true, params...)
	}
}

func makeOranBerry(bs battleState, mon *Pokemon) *item {
	return &item{
		trigger: func(...any) bool {
			return mon.Hp > 0 && mon.Hp <= mon.Stats["hp"]/2
		},
		activate: func(...any) {
			mon.ChangeHp(10)
			log.Printf("%s ate its berry and restored 10 hp", mon.Base.Name)
		},
	}
}

func makeSitrusBerry(bs battleState, mon *Pokemon) *item {
	return &item{
		trigger: func(...any) bool {
			return mon.Hp > 0 && mon.Hp <= mon.Stats["hp"]/2
		},
		activate: func(...any) {
			restore := mon.Stats["hp"] / 4
			mon.ChangeHp(restore)
			log.Printf("%s ate its berry and restored %d hp", mon.Base.Name, restore)
		},
	}
}

func makeResistBerryMiddleWare(typeName string) func(bs battleState, mon *Pokemon) *item {
	return func(bs battleState, mon *Pokemon) *item {
		return &item{
			trigger: func(p ...any) bool {
				if len(p) != 2 {
					return false
				}
				return p[0].(string) == typeName
			},
			activate: func(p ...any) {
				damage := p[1].(*int)
				*damage = *damage / 2
				log.Printf("%s ate its berry and reduced the damage", mon.Base.Name)
			},
		}
	}
}
