package main

import (
	"log"
	"math/rand"
)

var onSwitchAbilities = map[string]func(s *slot, bs battleState, switchIn bool){
	"unnerve":     unnerve,
	"intimidate":  intimidate,
	"regenerator": regenerator,
}

func unnerve(s *slot, bs battleState, switchIn bool) {
	for _, otherSlot := range bs.getOtherSlots(s) {
		if s.trainer != otherSlot.trainer {
			otherSlot.mon.Unnerved = switchIn
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
		if slot.mon.Ability == "inner-focus" {
			continue
		}
		slot.mon.changeStatStageBy("attack", -1, true)
	}
}

func regenerator(s *slot, bs battleState, switchIn bool) {
	if switchIn {
		return
	}

	s.mon.changeHpBy(s.mon.maxHP() / 3)
}

var typeConvertingAbilities = map[string]func(t *string, n, d *int){
	"aerilate":    typeConvertingAbilitiesMiddleware("flying"),
	"pixilate":    typeConvertingAbilitiesMiddleware("fairy"),
	"galvanize":   typeConvertingAbilitiesMiddleware("electric"),
	"refrigerate": typeConvertingAbilitiesMiddleware("ice"),
	"normalize":   normalize,
}

func normalize(t *string, n, d *int) {
	if *t != "normal" {
		*t = "normal"
		*n *= 6
		*d *= 5
	}
}

func typeConvertingAbilitiesMiddleware(t1 string) func(t *string, n, d *int) {
	return func(t2 *string, n, d *int) {
		if *t2 == "normal" {
			*t2 = t1
			*n *= 6
			*d *= 5
		}
	}
}

var typeImmunityAbilities = map[string]func(u *Pokemon, t string, s bool) bool{
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

func flashFire(p *Pokemon, t string, s bool) bool {
	if t != "fire" {
		return false
	}
	p.FlashFire = true
	return true
}

// still need to implement sunlight penalty
func drySkin(p *Pokemon, t string, s bool) bool {
	if t != "water" {
		return false
	}
	if s {
		return true
	}
	log.Printf("%s restored health with %s", p.Base.Name, p.Ability)
	p.changeHpBy(p.maxHP() / 4)
	return true
}

func stormDrain(p *Pokemon, t string, s bool) bool {
	if t != "water" {
		return false
	}
	if s {
		return true
	}
	p.changeStatStageBy("special-attack", 1, false)
	return true
}

func voltAbsorb(p *Pokemon, t string, s bool) bool {
	if t != "electric" {
		return false
	}
	if s {
		return true
	}
	log.Printf("%s restored health with %s", p.Base.Name, p.Ability)
	p.changeHpBy(p.maxHP() / 4)
	return true
}

func lightningRod(p *Pokemon, t string, s bool) bool {
	if t != "electric" {
		return false
	}
	if s {
		return true
	}
	p.changeStatStageBy("special-attack", 1, false)
	return true
}

func motorDrive(p *Pokemon, t string, s bool) bool {
	if t != "electric" {
		return false
	}
	if s {
		return true
	}
	p.changeStatStageBy("speed", 1, false)
	return true
}

func sapSipper(p *Pokemon, t string, s bool) bool {
	if t != "grass" {
		return false
	}
	if s {
		return true
	}
	p.changeStatStageBy("attack", 1, false)
	return true
}

func levitate(p *Pokemon, t string, s bool) bool {
	return t == "ground"
}

var pinchAbilities = map[string]string{
	"overgrow": "grass",
	"blaze":    "fire",
	"torrent":  "water",
	"swarm":    "bug",
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
	log.Printf("%s was hurt by %s", u.mon.Base.Name, t.mon.Ability)
}

func cuteCharm(u, t *slot) {
	if roll(30, 100) {
		u.mon.applyAilment("infatuation", nil, t)
	}
}

func flameBody(u, t *slot) {
	if roll(30, 100) {
		u.mon.applyAilment("burn", nil, t)
	}
}

func poisonPoint(u, t *slot) {
	if roll(30, 100) {
		u.mon.applyAilment("poison", nil, t)
	}
}

func effectSpore(u, t *slot) {
	if u.mon.hasType("grass") || u.mon.Ability == "overcoat" || u.mon.Item.name == "safety-goggles" {
		return
	}
	if roll(30, 100) {
		ailmentRoll := rand.Intn(30)
		if ailmentRoll <= 8 {
			u.mon.applyAilment("poison", nil, t)
		} else if ailmentRoll <= 18 {
			u.mon.applyAilment("paralysis", nil, t)
		} else {
			u.mon.applyAilment("sleep", nil, t)
		}
	}
}

var contactOffensiveAbilities = map[string]func(u, t *slot){
	"poison-touch": poisonTouch,
}

func poisonTouch(u, t *slot) {
	if roll(30, 100) {
		t.mon.applyAilment("poison", nil, u)
	}
}

func cheekPouch(mon *Pokemon) {
	if mon.Ability == "cheek-pouch" {
		restore := mon.maxHP() / 3
		mon.changeHpBy(restore)
		log.Printf("%s ate its cheek pouch and restored %d hp", mon.Base.Name, restore)
	}
}
