package main

import "log"

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
	p.changeStatStageBy("special-attack", 1)
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
	p.changeStatStageBy("special-attack", 1)
	return true
}

func motorDrive(p *Pokemon, t string, s bool) bool {
	if t != "electric" {
		return false
	}
	if s {
		return true
	}
	p.changeStatStageBy("speed", 1)
	return true
}

func sapSipper(p *Pokemon, t string, s bool) bool {
	if t != "grass" {
		return false
	}
	if s {
		return true
	}
	p.changeStatStageBy("attack", 1)
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

var contactOffensiveAbilities = map[string]func(u, t *slot){
	"poison-touch": poisonTouch,
}

func poisonTouch(u, t *slot) {
	if roll(30, 100) {
		t.mon.applyAilment("poison", nil, u)
	}
}
