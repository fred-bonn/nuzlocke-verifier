package main

import (
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type item struct {
	trigger  func() bool
	activate func()
	consumed bool
}

func registerOranBerry(mon *pokemon.Pokemon) {
	oranTrigger := func() bool {
		return mon.Hp < mon.Stats["hp"]/2 && mon.Hp > 0 && !mon.Fainted
	}
	oranActivate := func() {
		mon.ChangeHp(10)
		log.Printf("%s ate its berry and restored 10 hp", mon.Base.Name)
	}

	mon.Item = &item{
		trigger:  oranTrigger,
		activate: oranActivate,
		consumed: false,
	}
}

func checkItemTriggers(bs battleState) {
	slots := bs.getAllSlots()
	for _, slot := range slots {
		if slot.mon.Item == nil {
			continue
		}

		item, ok := slot.mon.Item.(*item)
		if !ok || item.trigger == nil {
			continue
		}

		if !item.consumed && item.trigger() {
			item.activate()
			item.consumed = true
		}
	}
}
