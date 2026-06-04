package main

import "log"

var critBlockingAbilities = map[string]struct{}{
	"battle-armor": {},
	"shell-armor":  {},
	"magma'armor":  {},
}

var contactAbilities = map[string]func(u, t *slot){
	"rough-skin": roughSkin,
	"iron-barbs": roughSkin,
	"cute-charm": cuteCharm,
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

func poisonTouch(u, t *slot) {
	if roll(1, 1) {
		t.mon.applyAilment("poison", nil, u)
	}
}
