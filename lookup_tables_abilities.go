package main

import "log"

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
	"rough-skin": roughSkin,
	"iron-barbs": roughSkin,
	"cute-charm": cuteCharm,
	"flame-body": flameBody,
}

var contactOffensiveAbilities = map[string]func(u, t *slot){
	"poison-touch": poisonTouch,
}

func roughSkin(u, t *slot) {
	change := u.mon.Stats["hp"] * 1 / 8
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

func poisonTouch(u, t *slot) {
	if roll(30, 100) {
		t.mon.applyAilment("poison", nil, u)
	}
}
