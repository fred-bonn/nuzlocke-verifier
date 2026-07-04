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
	for effect := range s.trainer.fieldEffects {
		if effect == stealthRockEffect {
			num, dem := s.mon.applyMoveType(1, 8, rockType)
			takeResidualDamage(bs, s, "stealth rock", num, dem)
			continue
		}
		if !s.mon.isGrounded() {
			continue
		}
		switch effect {
		case spikesEffect:
			takeResidualDamage(bs, s, "spikes", 1, 8)
		case toxicSpikesEffect:
			s.mon.applyAilment(poisonAilment, nil, nil)
		case stickyWebEffect:
			s.mon.changeStatStageBy(speed, -1, true)
		}
	}
}

func (s *slot) isTrapped() bool {
	return s.mon.hasAilment(trapAilment) != nil || s.mon.hasAilment(boundAilment) != nil
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
		vprintln("but it failed")
	}
}

func (s *slot) hasFieldEffect(effect string) bool {
	if _, ok := s.trainer.fieldEffects[stringToFieldEffect(effect)]; ok {
		return true
	}
	return false
}

func (s *slot) applyFieldEffect(effect string) {
	e := stringToFieldEffect(effect)
	turns := 0
	// code here to assign turn number based on which field effect it is
	s.trainer.fieldEffects[e] = turns
}
