package main

import (
	"strings"
)

type trainer struct {
	pokemonParty []*Pokemon
	player       bool
	ai           ai
	field        map[string]struct{}
	lost         bool
}

func (t *trainer) nextAction(bs battleState, slot *slot) action {
	if slot.invulnerableAction != nil {
		return slot.invulnerableAction
	}

	opponentSlots := bs.getOtherSlots(slot)

	possibleActions := make([]*moveAction, 0)
	if slot.mon.LockedMove != nil && slot.mon.LockedMove.PP > 0 {
		for _, oppSlot := range opponentSlots {
			possibleActions = append(possibleActions, &moveAction{
				userSlot:   slot,
				targetSlot: oppSlot,
				move:       slot.mon.LockedMove,
			})
		}
	} else {
		for _, move := range slot.mon.Moves {
			if move.PP <= 0 {
				continue
			}
			if slot.mon.Item.name == "assault-vest" && move.Class != "status" {
				continue
			}
			for _, oppSlot := range opponentSlots {
				possibleActions = append(possibleActions, &moveAction{
					userSlot:   slot,
					targetSlot: oppSlot,
					move:       move,
				})
			}
		}
	}

	// if there are no possible moves, struggle
	if len(possibleActions) == 0 {
		for _, oppSlot := range opponentSlots {
			if oppSlot.trainer != slot.trainer {
				possibleActions = append(possibleActions, &moveAction{
					userSlot:   slot,
					targetSlot: oppSlot,
					move:       &struggleMove,
				})
			}
		}
	}

	action, score := t.ai.evaluateActions(bs, possibleActions)
	if strings.HasPrefix(slot.mon.Item.name, "choice") {
		slot.mon.LockedMove = action.move
	}
	if score > 0 {
		return action
	}
	if roll(1, 2) || slot.mon.Hp <= slot.mon.maxHP()/2 || !slot.trainer.canReplace(bs) || slot.isTrapped() {
		return action
	}

	var possibleMons []*Pokemon
	for _, mon := range t.pokemonParty {
		if mon == slot.mon || mon.Fainted || bs.getActions().containstSwitchTo(mon) {
			continue
		}
		possibleMons = append(possibleMons, mon)
	}
	return &switchAction{
		oldSlot: slot,
		new:     t.ai.evaluteSwitchIns(bs, possibleMons, bs.getOpponentSlot(slot)),
	}
}

func (t *trainer) selectSwitchIn(bs battleState, slot *slot) *Pokemon {
	var possibleMons []*Pokemon
	for _, mon := range t.pokemonParty {
		if mon == slot.mon || mon.Fainted {
			continue
		}
		possibleMons = append(possibleMons, mon)
	}

	if len(possibleMons) == 0 {
		t.lost = true
		return nil
	}

	return t.ai.evaluteSwitchIns(bs, possibleMons, bs.getOpponentSlot(slot))
}

func (t *trainer) canReplace(bs battleState) bool {
	count := 0
	for _, mon := range t.pokemonParty {
		if !mon.Fainted {
			count++
		}
		if count > 1 {
			return true
		}
	}
	return false
}
