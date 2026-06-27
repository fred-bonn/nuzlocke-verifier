package main

import (
	"math/rand"
)

var sleepBlockingAbilities = map[string]struct{}{
	"insomnia":     {},
	"vital-spirit": {},
	"sweet-veil":   {},
}

var onSwitchAbilities = map[string]func(s *slot, bs battleState, switchIn bool){
	"trace":        trace,
	"unnerve":      unnerve,
	"intimidate":   intimidate,
	"regenerator":  regenerator,
	"natural-cure": naturalCure,
	"drizzle":      drizzle,
	"drought":      drought,
	"snow-warning": snowWarning,
	"sand-stream":  sandStream,
}

func trace(s *slot, bs battleState, switchIn bool) {
	if !switchIn {
		return
	}

	opponentMons := make([]*pokemon, 0)
	for _, slot := range bs.getOtherSlots(s) {
		if slot.trainer == s.trainer {
			continue
		}
		opponentMons = append(opponentMons, slot.mon)
	}

	s.mon.ability = opponentMons[rand.Int()%len(opponentMons)].ability
	s.mon.trace = true
	vlogf("%s traced %s", s.mon.base.Name, s.mon.ability)
}

func unnerve(s *slot, bs battleState, switchIn bool) {
	for _, otherSlot := range bs.getOtherSlots(s) {
		if s.trainer != otherSlot.trainer {
			otherSlot.mon.unnerved = switchIn
			otherSlot.mon.checkItemTrigger(true, nil)
		}
	}
}

func intimidate(s *slot, bs battleState, switchIn bool) {
	if !switchIn {
		return
	}

	for _, slot := range bs.getAllSlots() {
		if slot.trainer == s.trainer {
			continue
		}
		if slot.mon.ability == "inner-focus" {
			continue
		}
		slot.mon.changeStatStageBy(attack, -1, true)
	}
}

func regenerator(s *slot, bs battleState, switchIn bool) {
	if switchIn || s.mon.fainted {
		return
	}

	s.mon.changeHpBy(s.mon.maxHP() / 3)
}

func naturalCure(s *slot, bs battleState, switchIn bool) {
	if switchIn {
		return
	}

	for ailment := range nonVolatileStatuses {
		delete(s.mon.ailments, ailment)
	}
}

func drizzle(s *slot, bs battleState, switchIn bool) {
	if !switchIn {
		return
	}
	vlogln("it started to rain")
	bs.setWeather(rainWeather)
}

func drought(s *slot, bs battleState, switchIn bool) {
	if !switchIn {
		return
	}
	vlogln("the sunlight turned harsh")
	bs.setWeather(sunWeather)
}

func snowWarning(s *slot, bs battleState, switchIn bool) {
	if !switchIn {
		return
	}
	vlogln("it started to hail")
	bs.setWeather(hailWeather)
}

func sandStream(s *slot, bs battleState, switchIn bool) {
	if !switchIn {
		return
	}
	vlogln("a sandstorm brewed")
	bs.setWeather(sandstormWeather)
}

var typeConvertingAbilities = map[string]func(t *pokemonType, p *int){
	"aerilate":    typeConvertingAbilitiesMiddleware(flyingType),
	"pixilate":    typeConvertingAbilitiesMiddleware(fairyType),
	"galvanize":   typeConvertingAbilitiesMiddleware(electricType),
	"refrigerate": typeConvertingAbilitiesMiddleware(iceType),
	"normalize":   normalize,
}

func normalize(t *pokemonType, p *int) {
	if *t != normalType {
		*t = normalType
		*p = *p * 6 / 5
	}
}

func typeConvertingAbilitiesMiddleware(t1 pokemonType) func(t *pokemonType, p *int) {
	return func(t2 *pokemonType, p *int) {
		if *t2 == normalType {
			*t2 = t1
			*p = *p * 6 / 5
		}
	}
}

var typeImmunityAbilities = map[string]func(u *pokemon, t pokemonType, s bool) bool{
	"flash-fire":    flashFire,
	"dry-skin":      drySkin,
	"water-absorb":  drySkin,
	"storm-drain":   stormDrain,
	"volt-absorb":   voltAbsorb,
	"lightning-rod": lightningRod,
	"motor-drive":   motorDrive,
	"sap-sipper":    sapSipper,
	"levitate":      levitate,
}

func flashFire(p *pokemon, t pokemonType, s bool) bool {
	if t != fireType {
		return false
	}
	p.flashFire = true
	return true
}

// still need to implement sunlight penalty
func drySkin(p *pokemon, t pokemonType, s bool) bool {
	if t != waterType {
		return false
	}
	if s {
		return true
	}
	vlogf("%s restored health with %s", p.base.Name, p.ability)
	p.changeHpBy(p.maxHP() / 4)
	return true
}

func stormDrain(p *pokemon, t pokemonType, s bool) bool {
	if t != waterType {
		return false
	}
	if s {
		return true
	}
	p.changeStatStageBy(specialAttack, 1, false)
	return true
}

func voltAbsorb(p *pokemon, t pokemonType, s bool) bool {
	if t != electricType {
		return false
	}
	if s {
		return true
	}
	vlogf("%s restored health with %s", p.base.Name, p.ability)
	p.changeHpBy(p.maxHP() / 4)
	return true
}

func lightningRod(p *pokemon, t pokemonType, s bool) bool {
	if t != electricType {
		return false
	}
	if s {
		return true
	}
	p.changeStatStageBy(specialAttack, 1, false)
	return true
}

func motorDrive(p *pokemon, t pokemonType, s bool) bool {
	if t != electricType {
		return false
	}
	if s {
		return true
	}
	p.changeStatStageBy(speed, 1, false)
	return true
}

func sapSipper(p *pokemon, t pokemonType, s bool) bool {
	if t != grassType {
		return false
	}
	if s {
		return true
	}
	p.changeStatStageBy(attack, 1, false)
	return true
}

func levitate(p *pokemon, t pokemonType, s bool) bool {
	return t == groundType
}

var pinchAbilities = map[string]pokemonType{
	"overgrow": grassType,
	"blaze":    fireType,
	"torrent":  waterType,
	"swarm":    bugType,
}

var critBlockingAbilities = map[string]struct{}{
	"battle-armor": {},
	"shell-armor":  {},
	"magma'armor":  {},
}

var contactDefensiveAbilities = map[string]func(u, t *slot){
	"rough-skin":   roughSkin,
	"iron-barbs":   roughSkin,
	"cute-charm":   cuteCharm,
	"flame-body":   flameBody,
	"poison-point": poisonPoint,
	"effect-spore": effectSpore,
}

func roughSkin(u, t *slot) {
	change := u.mon.maxHP() * 1 / 8
	u.mon.changeHpBy(-change)
	vlogf("%s was hurt by %s", u.mon.base.Name, t.mon.ability)
}

func cuteCharm(u, t *slot) {
	if roll(30, 100) {
		u.mon.applyAilment(infatuationAilment, nil, t)
	}
}

func flameBody(u, t *slot) {
	if roll(30, 100) {
		u.mon.applyAilment(burnAilment, nil, t)
	}
}

func poisonPoint(u, t *slot) {
	if roll(30, 100) {
		u.mon.applyAilment(poisonAilment, nil, t)
	}
}

func effectSpore(u, t *slot) {
	if u.mon.hasType(grassType) || u.mon.ability == "overcoat" || u.mon.item.name == "safety-goggles" {
		return
	}
	if roll(30, 100) {
		ailmentRoll := rand.Intn(30)
		if ailmentRoll <= 8 {
			u.mon.applyAilment(poisonAilment, nil, t)
		} else if ailmentRoll <= 18 {
			u.mon.applyAilment(paralysisAilment, nil, t)
		} else {
			u.mon.applyAilment(sleepAilment, nil, t)
		}
	}
}

var contactOffensiveAbilities = map[string]func(u, t *slot){
	"poison-touch": poisonTouch,
}

func poisonTouch(u, t *slot) {
	if roll(30, 100) {
		t.mon.applyAilment(poisonAilment, nil, u)
	}
}

func cheekPouch(mon *pokemon) {
	if mon.ability == "cheek-pouch" {
		restore := mon.maxHP() / 3
		mon.changeHpBy(restore)
		vlogf("%s ate its cheek pouch and restored %d hp", mon.base.Name, restore)
	}
}
