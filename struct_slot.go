package main

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
			num, dem := s.mon.applyMoveType(1, 8, "rock")
			takeResidualDamage(bs, s, "stealth rock", num, dem)
			continue
		}
		if !s.mon.isGrounded() {
			continue
		}
		switch effect {
		case "spikes":
			takeResidualDamage(bs, s, "spikes", 1, 8)
		case "toxic-spikes":
			s.mon.applyAilment("poison", nil, nil)
		case "sticky-web":
			s.mon.changeStatStageBy(Speed, -1, true)
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
		vlogln("but it failed")
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
