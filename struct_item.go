package main

import (
	"fmt"
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
	"liechi-berry":   makeStatBoostBerryMiddleware("liechi-berry", attack),
	"ganlon-berry":   makeStatBoostBerryMiddleware("ganlon-berry", defense),
	"salac-berry":    makeStatBoostBerryMiddleware("salac-berry", speed),
	"petaya-berry":   makeStatBoostBerryMiddleware("petaya-berry", specialAttack),
	"apicot-berry":   makeStatBoostBerryMiddleware("apicot-berry", specialDefense),
	"babiri-berry":   makeResistBerryMiddleware("babiri-berry", steelType),
	"chilan-berry":   makeResistBerryMiddleware("chilan-berry", normalType),
	"charti-berry":   makeResistBerryMiddleware("charti-berry", rockType),
	"chople-berry":   makeResistBerryMiddleware("chople-berry", fightingType),
	"coba-berry":     makeResistBerryMiddleware("coba-berry", flyingType),
	"colbur-berry":   makeResistBerryMiddleware("colbur-berry", darkType),
	"haban-berry":    makeResistBerryMiddleware("haban-berry", dragonType),
	"kasib-berry":    makeResistBerryMiddleware("kasib-berry", ghostType),
	"kebia-berry":    makeResistBerryMiddleware("kebia-berry", poisonType),
	"occa-berry":     makeResistBerryMiddleware("occa-berry", fireType),
	"passho-berry":   makeResistBerryMiddleware("passho-berry", waterType),
	"payapa-berry":   makeResistBerryMiddleware("payapa-berry", psychicType),
	"rindo-berry":    makeResistBerryMiddleware("rindo-berry", grassType),
	"roseli-berry":   makeResistBerryMiddleware("roseli-berry", fairyType),
	"shuca-berry":    makeResistBerryMiddleware("shuca-berry", groundType),
	"tanga-berry":    makeResistBerryMiddleware("tanga-berry", bugType),
	"wacan-berry":    makeResistBerryMiddleware("wacan-berry", electricType),
	"yache-berry":    makeResistBerryMiddleware("yache-berry", iceType),
	"iron-ball":      makePassiveItemMiddleware("iron-ball"),
	"scope-lens":     makePassiveItemMiddleware("scope-lens"),
	"leftovers":      makePassiveItemMiddleware("leftovers"),
	"safety-goggles": makePassiveItemMiddleware("safety-goggles"),
	"mystic-water":   makeTypeBoostingItemMiddleware("mystic-water", waterType),
	"dragon-fang":    makeTypeBoostingItemMiddleware("dragon-fang", dragonType),
	"silver-powder":  makeTypeBoostingItemMiddleware("silver-powder", bugType),
	"magnet":         makeTypeBoostingItemMiddleware("magnet", electricType),
	"black-belt":     makeTypeBoostingItemMiddleware("black-belt", fightingType),
	"charcoal":       makeTypeBoostingItemMiddleware("charcoal", fireType),
	"sharp-beak":     makeTypeBoostingItemMiddleware("sharp-beak", flyingType),
	"spell-tag":      makeTypeBoostingItemMiddleware("spell-tag", ghostType),
	"miracle-seed":   makeTypeBoostingItemMiddleware("miracle-seed", grassType),
	"soft-sand":      makeTypeBoostingItemMiddleware("soft-sand", groundType),
	"never-melt-ice": makeTypeBoostingItemMiddleware("never-melt-ice", iceType),
	"silk-scarf":     makeTypeBoostingItemMiddleware("silk-scarf", normalType),
	"poison-barb":    makeTypeBoostingItemMiddleware("poison-barb", poisonType),
	"twisted-spoon":  makeTypeBoostingItemMiddleware("twisted-spoon", psychicType),
	"hard-stone":     makeTypeBoostingItemMiddleware("hard-stone", rockType),
	"metal-coat":     makeTypeBoostingItemMiddleware("metal-coat", steelType),
	"pixie-plate":    makeTypeBoostingItemMiddleware("pixir-plate", fairyType),
	"dark-glasses":   makeTypeBoostingItemMiddleware("dark-glasses", darkType),
	"normal-gem":     makeGemMiddleware(normalType),
	"fire-gem":       makeGemMiddleware(fireType),
	"fighting-gem":   makeGemMiddleware(fightingType),
	"water-gem":      makeGemMiddleware(waterType),
	"flying-gem":     makeGemMiddleware(flyingType),
	"grass-gem":      makeGemMiddleware(grassType),
	"poison-gem":     makeGemMiddleware(poisonType),
	"electric-gem":   makeGemMiddleware(electricType),
	"ground-gem":     makeGemMiddleware(groundType),
	"psychic-gem":    makeGemMiddleware(psychicType),
	"rock-gem":       makeGemMiddleware(rockType),
	"ice-gem":        makeGemMiddleware(iceType),
	"bug-gem":        makeGemMiddleware(bugType),
	"dragon-gem":     makeGemMiddleware(dragonType),
	"ghost-gem":      makeGemMiddleware(ghostType),
	"dark-gem":       makeGemMiddleware(darkType),
	"steel-gem":      makeGemMiddleware(steelType),
	"fairy-gem":      makeGemMiddleware(fairyType),
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

func makeTypeBoostingItemMiddleware(itemName string, pokemonType pokemonType) func(mon *pokemon) *item {
	return func(mon *pokemon) *item {
		var p *int
		return &item{
			name: itemName,
			trigger: func(e any) bool {
				event, ok := e.(moveBoostingEvent)
				if !ok {
					return false
				}
				if pokemonType != event.pokemonType {
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
			vlogItem("%s drank its berry juice and restored 20 hp", mon.base.Name)
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
			vlogItem("%s ate its oran berry and restored 10 hp", mon.base.Name)
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
			vlogItem("%s ate its sitrus berry and restored %d hp", mon.base.Name, restore)
			cheekPouch(mon)
		},
	}
}

func makeCheriBerry(mon *pokemon) *item {
	return &item{
		name: "cheri-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hasAilment(paralysisAilment) != nil
		},
		activate: func() {
			vlogItem("%s ate its cheri berry", mon.base.Name)
			delete(mon.ailments, paralysisAilment)
			cheekPouch(mon)
		},
	}
}

func makeChestoBerry(mon *pokemon) *item {
	return &item{
		name: "chesto-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hasAilment(sleepAilment) != nil
		},
		activate: func() {
			vlogItem("%s ate its chesto berry", mon.base.Name)
			delete(mon.ailments, sleepAilment)
			cheekPouch(mon)
		},
	}
}

func makePechaBerry(mon *pokemon) *item {
	return &item{
		name: "pecha-berry",
		trigger: func(any) bool {
			return !mon.unnerved && (mon.hasAilment(poisonAilment) != nil || mon.hasAilment(toxicAilment) != nil)
		},
		activate: func() {
			vlogItem("%s ate its pecha berry", mon.base.Name)
			delete(mon.ailments, poisonAilment)
			delete(mon.ailments, toxicAilment)
			cheekPouch(mon)
		},
	}
}

func makeRawstBerry(mon *pokemon) *item {
	return &item{
		name: "rawst-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hasAilment(burnAilment) != nil
		},
		activate: func() {
			vlogItem("%s ate its rawst berry", mon.base.Name)
			delete(mon.ailments, burnAilment)
			cheekPouch(mon)
		},
	}
}

func makeAspearBerry(mon *pokemon) *item {
	return &item{
		name: "aspear-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hasAilment(freezeAilment) != nil
		},
		activate: func() {
			vlogItem("%s ate its aspear berry", mon.base.Name)
			delete(mon.ailments, freezeAilment)
			cheekPouch(mon)
		},
	}
}

func makePersimBerry(mon *pokemon) *item {
	return &item{
		name: "persim-berry",
		trigger: func(any) bool {
			return !mon.unnerved && mon.hasAilment(confusionAilment) != nil
		},
		activate: func() {
			vlogItem("%s ate its persim berry", mon.base.Name)
			delete(mon.ailments, confusionAilment)
			cheekPouch(mon)
		},
	}
}

func makeLumBerry(mon *pokemon) *item {
	return &item{
		name: "lum-berry",
		trigger: func(any) bool {
			return !mon.unnerved && (mon.hasNonVolatileAilment() || mon.hasAilment(confusionAilment) != nil)
		},
		activate: func() {
			vlogItem("%s ate its lum berry", mon.base.Name)
			for ailment := range nonVolatileStatuses {
				if mon.hasAilment(ailment) != nil {
					delete(mon.ailments, ailment)
					vlogItem("%s had its %s removed", mon.base.Name, ailment.String())
				}
			}
			if mon.hasAilment(confusionAilment) != nil {
				delete(mon.ailments, confusionAilment)
				vlogItem("%s had its confusion removed", mon.base.Name)
				cheekPouch(mon)
			}
		},
	}
}

func makeLeppaBerry(mon *pokemon) *item {
	var m *Move
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

func makeStatBoostBerryMiddleware(name string, stat stats) func(mon *pokemon) *item {
	return func(mon *pokemon) *item {
		return &item{
			name: name,
			trigger: func(any) bool {
				if mon.ability == gluttonyAbility {
					return !mon.unnerved && mon.hp > 0 && mon.hp*2 <= mon.maxHP()
				}
				return !mon.unnerved && mon.hp > 0 && mon.hp*4 <= mon.maxHP()
			},
			activate: func() {
				vlogItem("%s ate its %s", mon.base.Name, name)
				mon.changeStatStageBy(stat, 1, false)
				cheekPouch(mon)
			},
		}
	}
}

func makeResistBerryMiddleware(name string, pokemonType pokemonType) func(mon *pokemon) *item {
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
				return event.pokemonType == pokemonType
			},
			activate: func() {
				if d == nil {
					vlogItem("%s ate its %s and reduced the damage", mon.base.Name, name)
					cheekPouch(mon)
				} else {
					*d /= 2
				}
			},
		}
	}
}

func makeGemMiddleware(pokemonType pokemonType) func(mon *pokemon) *item {
	var p *int
	return func(mon *pokemon) *item {
		return &item{
			name: fmt.Sprintf("%s-gem", pokemonType.String()),
			trigger: func(e any) bool {
				event, ok := e.(gemEvent)
				if !ok {
					return false
				}
				p = event.power
				return event.pokemonType == pokemonType
			},
			activate: func() {
				if p == nil {
					vlogItem("%s consumed its %s gem and boosted the damage", mon.base.Name, pokemonType.String())
				} else {
					*p = *p * 3 / 2
				}
			},
		}
	}
}

func makeChoiceScarf(mon *pokemon) *item {
	mon.stats[speed] = mon.stats[speed] * 3 / 2
	return &item{
		name: "choice-scarf",
	}
}

func makeAssaultVest(mon *pokemon) *item {
	mon.stats[specialDefense] = mon.stats[specialDefense] * 3 / 2
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
			if event.move.Class != physicalClass {
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
			if event.move.Class != specialClass {
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

type resistBerryEvent struct {
	pokemonType pokemonType
	damage      *int
}

type gemEvent struct {
	pokemonType pokemonType
	power       *int
}

type leppaBerryEvent struct {
	move *Move
}

type choiceItemEvent struct {
	move *Move
	stat *int
}

type focusSashEvent struct {
	damage *int
}

type moveBoostingEvent struct {
	power       *int
	pokemonType pokemonType
}
