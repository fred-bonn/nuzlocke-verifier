package main

import "log"

type slot struct {
	mon                *pokemon
	trainer            *trainer
	firstTurn          bool
	suckerPunch        bool
	protected          bool
	protectTurns       int
	invulnerableAction *moveAction
	unnerved           bool
}

func (s *slot) setMon(bs battleState, new *pokemon) {
	s.mon.switchReset()
	s.firstTurn = true
	s.suckerPunch = false

	new.unnerved = s.mon.unnerved
	s.mon = new
	for effect := range s.trainer.field {
		if effect == "stealth-rock" {
			log.Printf("%s took damage from stealth rock", s.mon.base.Name)
			num, dem := s.mon.applyMoveType(1, 1, "rock")
			s.mon.changeHpBy(-s.mon.maxHP() / 8 * num / dem)
			if s.mon.hp <= 0 {
				monFainted(bs, s, false)
			}
			continue
		}
		if !s.mon.isGrounded() {
			continue
		}
		switch effect {
		case "spikes":
			log.Printf("%s took damage from spikes", s.mon.base.Name)
			s.mon.changeHpBy(-s.mon.maxHP() / 8)
		case "toxic-spikes":
			s.mon.applyAilment("poison", nil, nil)
		case "sticky-web":
			s.mon.changeStatStageBy("speed", -1, true)
		}
	}
}

func (s *slot) isTrapped() bool {
	return s.mon.hasAilment("trap") != nil || s.mon.hasAilment("bound") != nil
}

func (s *slot) resolveProtect() {
	denominator := 1
	for i := 0; i < s.protectTurns; i++ {
		denominator *= 3
	}
	if roll(1, denominator) {
		s.protected = true
		s.protectTurns++
	} else {
		log.Println("but it failed")
	}
}

func (s *slot) hasFieldEffect(effect string) bool {
	if _, ok := s.trainer.field[effect]; ok {
		return true
	}
	return false
}

func (s *slot) applyFieldEffect(effect string) {
	s.trainer.field[effect] = struct{}{}
}
