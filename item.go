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

func (p *pokemon) checkItemTrigger(consume bool, event any) {
	if p.item.trigger == nil || p.item.consumed {
		return
	}

	if p.item.trigger(event) {
		if consume {
			p.item.consumed = true
			p.unburden = true
		}
		p.item.activate()
	}
}

type ItemFactoryBuilder func(*pokemon) *item

var itemBuilders = map[string]ItemFactoryBuilder{
	"berry-juice":    makeBerryJuice,
	"oran-berry":     makeOranBerry,
	"sitrus-berry":   makeSitrusBerry,
	"lum-berry":      makeLumBerry,
	"leppa-berry":    makeLeppaBerry,
	"cheri-berry":    makeCheriBerry,
	"chesto-berry":   makeChestoBerry,
	"pecha-berry":    makePechaBerry,
	"rawst-berry":    makeRawstBerry,
	"aspear-berry":   makeAspearBerry,
	"persim-berry":   makePersimBerry,
	"liechi-berry":   makeStatBoostBerryMiddleware("liechi-berry", "attack"),
	"ganlon-berry":   makeStatBoostBerryMiddleware("ganlon-berry", "defense"),
	"salac-berry":    makeStatBoostBerryMiddleware("salac-berry", "speed"),
	"petaya-berry":   makeStatBoostBerryMiddleware("petaya-berry", "special-attack"),
	"apicot-berry":   makeStatBoostBerryMiddleware("apicot-berry", "special-defense"),
	"babiri-berry":   makeResistBerryMiddleware("babiri-berry", "steel"),
	"chilan-berry":   makeResistBerryMiddleware("chilan-berry", "normal"),
	"charti-berry":   makeResistBerryMiddleware("charti-berry", "rock"),
	"chople-berry":   makeResistBerryMiddleware("chople-berry", "fighting"),
	"coba-berry":     makeResistBerryMiddleware("coba-berry", "coba"),
	"colbur-berry":   makeResistBerryMiddleware("colbur-berry", "dark"),
	"haban-berry":    makeResistBerryMiddleware("haban-berry", "dragon"),
	"kasib-berry":    makeResistBerryMiddleware("kasib-berry", "ghost"),
	"kebia-berry":    makeResistBerryMiddleware("kebia-berry", "poison"),
	"occa-berry":     makeResistBerryMiddleware("occa-berry", "fire"),
	"passho-berry":   makeResistBerryMiddleware("passho-berry", "water"),
	"payapa-berry":   makeResistBerryMiddleware("payapa-berry", "psychic"),
	"rindo-berry":    makeResistBerryMiddleware("rindo-berry", "grass"),
	"roseli-berry":   makeResistBerryMiddleware("roseli-berry", "fairy"),
	"shuca-berry":    makeResistBerryMiddleware("shuca-berry", "ground"),
	"tanga-berry":    makeResistBerryMiddleware("tanga-berry", "bug"),
	"wacan-berry":    makeResistBerryMiddleware("wacan-berry", "electric"),
	"yache-berry":    makeResistBerryMiddleware("yache-berry", "ice"),
	"iron-ball":      makePassiveItemMiddleware("iron-ball"),
	"scope-lens":     makePassiveItemMiddleware("scope-lens"),
	"leftovers":      makePassiveItemMiddleware("leftovers"),
	"safety-goggles": makePassiveItemMiddleware("safety-goggles"),
	"mystic-water":   makeTypeBoostingItemMiddleware("mystic-water", "water"),
	"dragon-fang":    makeTypeBoostingItemMiddleware("dragon-fang", "dragon"),
	"silver-powder":  makeTypeBoostingItemMiddleware("silver-powder", "bug"),
	"magnet":         makeTypeBoostingItemMiddleware("magnet", "electric"),
	"black-belt":     makeTypeBoostingItemMiddleware("black-belt", "fighting"),
	"charcoal":       makeTypeBoostingItemMiddleware("charcoal", "fire"),
	"sharp-beak":     makeTypeBoostingItemMiddleware("sharp-beak", "flying"),
	"spell-tag":      makeTypeBoostingItemMiddleware("spell-tag", "ghost"),
	"miracle-seed":   makeTypeBoostingItemMiddleware("miracle-seed", "grass"),
	"soft-sand":      makeTypeBoostingItemMiddleware("soft-sand", "ground"),
	"never-melt-ice": makeTypeBoostingItemMiddleware("never-melt-ice", "ice"),
	"silk-scarf":     makeTypeBoostingItemMiddleware("silk-scarf", "normal"),
	"poison-barb":    makeTypeBoostingItemMiddleware("poison-barb", "poison"),
	"twisted-spoon":  makeTypeBoostingItemMiddleware("twisted-spoon", "psychic"),
	"hard-stone":     makeTypeBoostingItemMiddleware("hard-stone", "rock"),
	"metal-coat":     makeTypeBoostingItemMiddleware("metal-coat", "steel"),
	"pixie-plate":    makeTypeBoostingItemMiddleware("pixir-plate", "fairy"),
	"dark-glasses":   makeTypeBoostingItemMiddleware("dark-glasses", "dark"),
	"normal-gem":     makeGemMiddleware("normal"),
	"fire-gem":       makeGemMiddleware("fire"),
	"fighting-gem":   makeGemMiddleware("fighting"),
	"water-gem":      makeGemMiddleware("water"),
	"flying-gem":     makeGemMiddleware("flying"),
	"grass-gem":      makeGemMiddleware("grass"),
	"poison-gem":     makeGemMiddleware("poison"),
	"electric-gem":   makeGemMiddleware("electric"),
	"ground-gem":     makeGemMiddleware("ground"),
	"psychic-gem":    makeGemMiddleware("psychic"),
	"rock-gem":       makeGemMiddleware("rock"),
	"ice-gem":        makeGemMiddleware("ice"),
	"bug-gem":        makeGemMiddleware("bug"),
	"dragon-gem":     makeGemMiddleware("dragon"),
	"ghost-gem":      makeGemMiddleware("ghost"),
	"dark-gem":       makeGemMiddleware("dark"),
	"steel-gem":      makeGemMiddleware("steel"),
	"fairy-gem":      makeGemMiddleware("fairy"),
	"choice-scarf":   makeChoiceScarf,
	"assault-vest":   makeAssaultVest,
	"choice-band":    makeChoiceBand,
	"choice-specs":   makeChoiceSpecs,
	"focus-sash":     makeFocusSash,
}

func createItemFactory(builder ItemFactoryBuilder, mon *pokemon) func() *item {
	return func() *item {
		return builder(mon)
	}
}

func registerItem(itemName string, mon *pokemon) (*item, error) {
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

func makePassiveItemMiddleware(itemName string) func(mon *pokemon) *item {
	return func(mon *pokemon) *item {
		return &item{
			name: itemName,
			trigger: func(e any) bool {
				return false
			},
		}
	}
}

func makeTypeBoostingItemMiddleware(itemName, typeName string) func(mon *pokemon) *item {
	return func(mon *pokemon) *item {
		var p *int
		return &item{
			name: itemName,
			trigger: func(e any) bool {
				event, ok := e.(moveBoostingEvent)
				if !ok {
					return false
				}
				if typeName != event.typeName {
					return false
				}
				p = event.power
				return true
			},
			activate: func() {
				*p = *p * 6 / 5
			},
		}
	}
}

func makeBerryJuice(mon *pokemon) *item {
	return &item{
		name: "berry-juice",
		trigger: func(any) bool {
			return mon.hp > 0 && mon.hp*2 <= mon.maxHP()
		},
		activate: func() {
			mon.changeHpBy(20)
			log.Printf("%s drank its berry juice and restored 20 hp", mon.base.Name)
			cheekPouch(mon)
		},
	}
}

func makeOranBerry(mon *pokemon) *item {
	return &item{
		name: "oran-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hp > 0 && mon.hp*2 <= mon.maxHP()
		},
		activate: func() {
			mon.changeHpBy(10)
			log.Printf("%s ate its oran berry and restored 10 hp", mon.base.Name)
			cheekPouch(mon)
		},
	}
}

func makeSitrusBerry(mon *pokemon) *item {
	return &item{
		name: "sitrus-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hp > 0 && mon.hp*2 <= mon.maxHP()
		},
		activate: func() {
			restore := mon.maxHP() / 4
			mon.changeHpBy(restore)
			log.Printf("%s ate its sitrus berry and restored %d hp", mon.base.Name, restore)
			cheekPouch(mon)
		},
	}
}

func makeCheriBerry(mon *pokemon) *item {
	return &item{
		name: "cheri-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hasAilment("paralysis") != nil
		},
		activate: func() {
			log.Printf("%s ate its cheri berry", mon.base.Name)
			delete(mon.ailments, "paralysis")
			cheekPouch(mon)
		},
	}
}

func makeChestoBerry(mon *pokemon) *item {
	return &item{
		name: "chesto-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hasAilment("sleep") != nil
		},
		activate: func() {
			log.Printf("%s ate its chesto berry", mon.base.Name)
			delete(mon.ailments, "sleep")
			cheekPouch(mon)
		},
	}
}

func makePechaBerry(mon *pokemon) *item {
	return &item{
		name: "pecha-berry",
		trigger: func(any) bool {
			return !mon.unnerved && (mon.hasAilment("poison") != nil || mon.hasAilment("toxic") != nil)
		},
		activate: func() {
			log.Printf("%s ate its pecha berry", mon.base.Name)
			delete(mon.ailments, "poison")
			delete(mon.ailments, "toxic")
			cheekPouch(mon)
		},
	}
}

func makeRawstBerry(mon *pokemon) *item {
	return &item{
		name: "rawst-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hasAilment("burn") != nil
		},
		activate: func() {
			log.Printf("%s ate its rawst berry", mon.base.Name)
			delete(mon.ailments, "burn")
			cheekPouch(mon)
		},
	}
}

func makeAspearBerry(mon *pokemon) *item {
	return &item{
		name: "aspear-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hasAilment("freeze") != nil
		},
		activate: func() {
			log.Printf("%s ate its aspear berry", mon.base.Name)
			delete(mon.ailments, "freeze")
			cheekPouch(mon)
		},
	}
}

func makePersimBerry(mon *pokemon) *item {
	return &item{
		name: "persim-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hasAilment("confusion") != nil
		},
		activate: func() {
			log.Printf("%s ate its persim berry", mon.base.Name)
			delete(mon.ailments, "confusion")
			cheekPouch(mon)
		},
	}
}

func makeLumBerry(mon *pokemon) *item {
	return &item{
		name: "lum-berry",
		trigger: func(any) bool {
			return !mon.unnerved && (mon.hasNonVolatileAilment() || mon.hasAilment("confusion") != nil)
		},
		activate: func() {
			log.Printf("%s ate its lum berry", mon.base.Name)
			for ailment := range nonVolatileStatuses {
				if mon.hasAilment(ailment) != nil {
					delete(mon.ailments, ailment)
					log.Printf("%s had its %s removed", mon.base.Name, ailment)
				}
			}
			if mon.hasAilment("confusion") != nil {
				delete(mon.ailments, "confusion")
				log.Printf("%s had its confusion removed", mon.base.Name)
				cheekPouch(mon)
			}
		},
	}
}

func makeLeppaBerry(mon *pokemon) *item {
	var m *pokeapi.BaseMove
	return &item{
		trigger: func(e any) bool {
			event, ok := e.(leppaBerryEvent)
			if !ok || mon.unnerved {
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

func makeStatBoostBerryMiddleware(name, stat string) func(mon *pokemon) *item {
	return func(mon *pokemon) *item {
		return &item{
			name: name,
			trigger: func(any) bool {
				if mon.ability == "gluttony" {
					return !mon.unnerved && mon.hp > 0 && mon.hp*2 <= mon.maxHP()
				}
				return !mon.unnerved && mon.hp > 0 && mon.hp*4 <= mon.maxHP()
			},
			activate: func() {
				log.Printf("%s ate its %s", mon.base.Name, name)
				mon.changeStatStageBy(stat, 1, false)
				cheekPouch(mon)
			},
		}
	}
}

func makeResistBerryMiddleware(name, typeName string) func(mon *pokemon) *item {
	var d *int
	return func(mon *pokemon) *item {
		return &item{
			name: name,
			trigger: func(e any) bool {
				event, ok := e.(resistBerryEvent)
				if !ok || mon.unnerved {
					return false
				}
				d = event.damage
				return event.typeName == typeName
			},
			activate: func() {
				if d == nil {
					log.Printf("%s ate its %s and reduced the damage", mon.base.Name, name)
					cheekPouch(mon)
				} else {
					*d /= 2
				}
			},
		}
	}
}

func makeGemMiddleware(typeName string) func(mon *pokemon) *item {
	var p *int
	return func(mon *pokemon) *item {
		return &item{
			name: fmt.Sprintf("%s-gem", typeName),
			trigger: func(e any) bool {
				event, ok := e.(gemEvent)
				if !ok {
					return false
				}
				p = event.power
				return event.typeName == typeName
			},
			activate: func() {
				if p == nil {
					log.Printf("%s consumed its %s gem and boosted the damage", mon.base.Name, typeName)
				} else {
					*p = *p * 3 / 2
				}
			},
		}
	}
}

func makeChoiceScarf(mon *pokemon) *item {
	mon.stats["speed"] = mon.stats["speed"] * 3 / 2
	return &item{
		name: "choice-scarf",
	}
}

func makeAssaultVest(mon *pokemon) *item {
	mon.stats["special-defense"] = mon.stats["special-defense"] * 3 / 2
	return &item{
		name: "assault-vest",
	}
}

func makeChoiceBand(mon *pokemon) *item {
	var s *int
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
			s = event.stat
			return true
		},
		activate: func() {
			*s = *s * 3 / 2
		},
	}
}

func makeChoiceSpecs(mon *pokemon) *item {
	var s *int
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
			s = event.stat
			return true
		},
		activate: func() {
			*s = *s * 3 / 2
		},
	}
}

func makeFocusSash(mon *pokemon) *item {
	var dmg *int
	return &item{
		trigger: func(e any) bool {
			event, ok := e.(focusSashEvent)
			if !ok {
				return false
			}
			dmg = event.damage
			return mon.hp == mon.maxHP() && *event.damage >= mon.hp
		},
		activate: func() {
			*dmg = mon.hp - 1
		},
	}
}
