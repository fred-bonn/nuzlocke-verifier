package main

import (
	"encoding/json"
	"sort"
)

type discreteBattleState struct {
	PlayerPokemon             string   `json:"player_pokemon"`
	PlayerAilments            []string `json:"player_ailments"`
	OpponentPokemon           string   `json:"opponent_pokemon"`
	OpponentAilments          []string `json:"opponent_ailments"`
	SpeedHierarchy            string   `json:"speed_hierarchy"`
	PlayerMonIsFaster         bool     `json:"player_mon_is_faster"`
	PlayerMonHasMoveThatKills bool     `json:"player_mon_has_move_that_kills"`
	OpponentMovesThatKill     []string `json:"opponent_moves_that_kill"`
}

// discretizeBattleState returns a canonical key suitable for a tabular policy.
func discretizeBattleState(bs battleState) string {
	state := discreteBattleState{
		PlayerAilments:        []string{},
		OpponentAilments:      []string{},
		OpponentMovesThatKill: []string{},
	}
	if bs == nil {
		encoded, _ := json.Marshal(state)
		return string(encoded)
	}

	var playerSlot *slot
	for _, slot := range bs.getAllSlots() {
		if slot != nil && slot.trainer != nil && slot.trainer.player {
			playerSlot = slot
			break
		}
	}
	if playerSlot == nil {
		return string(mustMarshalDiscreteState(state))
	}
	opponentSlot := bs.getOpponentSlot(playerSlot)
	if playerSlot.mon == nil || opponentSlot == nil || opponentSlot.mon == nil {
		return string(mustMarshalDiscreteState(state))
	}

	state.PlayerPokemon = playerSlot.mon.base.Name
	state.PlayerAilments = discreteAilments(playerSlot.mon)
	state.OpponentPokemon = opponentSlot.mon.base.Name
	state.OpponentAilments = discreteAilments(opponentSlot.mon)
	playerSpeed := playerSlot.mon.effectiveSpeed(bs)
	opponentSpeed := opponentSlot.mon.effectiveSpeed(bs)
	state.PlayerMonIsFaster = playerSpeed > opponentSpeed
	state.SpeedHierarchy = "slower"
	if playerSpeed == opponentSpeed {
		state.SpeedHierarchy = "tie"
	} else if state.PlayerMonIsFaster {
		state.SpeedHierarchy = "faster"
	}

	state.PlayerMonHasMoveThatKills = monHasMoveThatKills(bs, playerSlot.mon, opponentSlot.mon)
	for _, move := range opponentSlot.mon.moves {
		if move.PP <= 0 || move.Class == statusClass {
			continue
		}
		if moveCanKill(bs, opponentSlot.mon, playerSlot.mon, move) {
			state.OpponentMovesThatKill = append(state.OpponentMovesThatKill, move.Name)
		}
	}
	sort.Strings(state.OpponentMovesThatKill)

	return string(mustMarshalDiscreteState(state))
}

func mustMarshalDiscreteState(state discreteBattleState) []byte {
	encoded, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	return encoded
}

func discreteAilments(mon *pokemon) []string {
	result := make([]string, 0, len(mon.ailments))
	for ailment := range mon.ailments {
		result = append(result, ailment.String())
	}
	sort.Strings(result)
	return result
}

func monHasMoveThatKills(bs battleState, user, target *pokemon) bool {
	for _, move := range user.moves {
		if move.PP > 0 && move.Class != statusClass && moveCanKill(bs, user, target, move) {
			return true
		}
	}
	return false
}

func moveCanKill(bs battleState, user, target *pokemon, move *Move) bool {
	rolls := 1
	if move.MaxHits == 5 {
		rolls = 3
	} else if move.MaxHits > 0 {
		rolls = move.MaxHits
	}
	damage := 0
	critRate := determineCritRate(user, move)
	for i := 0; i < rolls; i++ {
		damage += calculateDamage(user, target, move, new(critRate >= 3), bs.getWeather(), true, true, false)
	}
	target.checkItemTrigger(false, makeFocusSashEvent(&damage))
	if target.ability == sturdyAbility && target.hp == target.maxHP() {
		damage = min(damage, target.hp-1)
	}
	return damage >= target.hp
}

type learningAi struct {
	policy map[string][]string
}

func newLearningAi() *learningAi {
	return &learningAi{
		policy: make(map[string][]string),
	}
}

func (la *learningAi) setPolicyAction(bs battleState, action string) {
	if la == nil {
		return
	}
	if la.policy == nil {
		la.policy = make(map[string][]string)
	}
	la.policy[discretizeBattleState(bs)] = []string{action}
}

func (la *learningAi) evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int) {
	if len(actions) == 0 {
		return nil, 0
	}
	if la == nil {
		return rnbAi{}.evaluateActions(bs, actions)
	}
	if la.policy == nil {
		la.policy = make(map[string][]string)
	}
	firstAction := actions[0]
	la.policy[discretizeBattleState(bs)] = moveActionKeys(actions)
	return firstAction, 0
}

func (la *learningAi) evaluteSwitchIns(bs battleState, mons []*pokemon, opponentSlot *slot) *pokemon {
	if len(mons) == 0 {
		return nil
	}
	if la == nil {
		return rnbAi{}.evaluteSwitchIns(bs, mons, opponentSlot)
	}
	if la.policy == nil {
		la.policy = make(map[string][]string)
	}
	firstMon := mons[0]
	stateKey := discretizeBattleState(bs)
	for _, mon := range mons {
		la.policy[stateKey] = appendUnique(la.policy[stateKey], actionKeyForSwitch(mon))
	}
	return firstMon
}

func (la *learningAi) shouldSwitch(bs battleState, slot *slot, score int, party []*pokemon) bool {
	return false
}

func actionKey(action *moveAction) string {
	return "move:" + action.move.Name
}

func moveActionKeys(actions []*moveAction) []string {
	keys := make([]string, 0, len(actions))
	for _, action := range actions {
		keys = appendUnique(keys, actionKey(action))
	}
	return keys
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func actionKeyForSwitch(mon *pokemon) string {
	return "switch:" + mon.base.Name
}
