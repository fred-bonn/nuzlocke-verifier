package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
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
	OpponentCritKillRisk      bool     `json:"opponent_crit_kill_risk"`
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
		if moveCanCritKill(bs, opponentSlot.mon, playerSlot.mon, move) {
			state.OpponentCritKillRisk = true
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

func moveCanCritKill(bs battleState, user, target *pokemon, move *Move) bool {
	if move == nil || move.PP <= 0 || move.Class == statusClass {
		return false
	}
	critRate := determineCritRate(user, move)
	if critRate <= 0 {
		return false
	}
	critMult := 1
	if critRate >= 3 {
		critMult = 2
	}
	damage := 0
	rolls := 1
	if move.MaxHits == 5 {
		rolls = 3
	} else if move.MaxHits > 0 {
		rolls = move.MaxHits
	}
	for i := 0; i < rolls; i++ {
		damage += calculateDamage(user, target, move, new(true), bs.getWeather(), true, true, false)
	}
	if critMult > 1 {
		damage *= critMult
	}
	target.checkItemTrigger(false, makeFocusSashEvent(&damage))
	if target.ability == sturdyAbility && target.hp == target.maxHP() {
		damage = min(damage, target.hp-1)
	}
	return damage >= target.hp
}

type learningAi struct {
	policy      map[string][]string
	scores      map[string]map[string]float64
	counts      map[string]map[string]int
	history     []stateActionEntry
	seen        map[string]map[string]bool
	decayFactor float64
}

type stateActionEntry struct {
	stateKey string
	action   string
}

type savedPolicy struct {
	PlayerParty   []string                      `json:"player_party"`
	OpponentParty []string                      `json:"opponent_party"`
	Policy        map[string][]string           `json:"policy"`
	Scores        map[string]map[string]float64 `json:"scores"`
	Counts        map[string]map[string]int     `json:"counts"`
	CreatedAt     string                        `json:"created_at"`
	Version       int                           `json:"version"`
	Metadata      map[string]string             `json:"metadata,omitempty"`
}

func newLearningAi() *learningAi {
	return &learningAi{
		policy:      make(map[string][]string),
		scores:      make(map[string]map[string]float64),
		counts:      make(map[string]map[string]int),
		seen:        make(map[string]map[string]bool),
		decayFactor: 0.98,
	}
}

func policyPathForInputs(playerInput, opponentInput string) string {
	playerName := strings.TrimSuffix(filepath.Base(playerInput), filepath.Ext(playerInput))
	opponentName := strings.TrimSuffix(filepath.Base(opponentInput), filepath.Ext(opponentInput))
	return filepath.Join("policies", fmt.Sprintf("%s__vs__%s.json", playerName, opponentName))
}

func loadPolicyFromDisk(path string) (*savedPolicy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	policy := &savedPolicy{}
	if err := json.Unmarshal(data, policy); err != nil {
		return nil, err
	}
	if policy.Policy == nil {
		policy.Policy = make(map[string][]string)
	}
	if policy.Scores == nil {
		policy.Scores = make(map[string]map[string]float64)
	}
	if policy.Counts == nil {
		policy.Counts = make(map[string]map[string]int)
	}
	return policy, nil
}

func savePolicyToDisk(la *learningAi, playerInput, opponentInput string, playerParty, opponentParty []*pokemon) error {
	if la == nil {
		return fmt.Errorf("nil learning AI")
	}
	path := policyPathForInputs(playerInput, opponentInput)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	policy := &savedPolicy{
		PlayerParty:   partyNames(playerParty),
		OpponentParty: partyNames(opponentParty),
		Policy:        cloneActionMap(la.policy),
		Scores:        cloneScoreMap(la.scores),
		Counts:        cloneCountMap(la.counts),
		CreatedAt:     nowISO(),
		Version:       1,
		Metadata:      map[string]string{"player_input": playerInput, "opponent_input": opponentInput},
	}
	data, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func cloneActionMap(source map[string][]string) map[string][]string {
	if source == nil {
		return make(map[string][]string)
	}
	clone := make(map[string][]string, len(source))
	for key, values := range source {
		copied := make([]string, len(values))
		copy(copied, values)
		clone[key] = copied
	}
	return clone
}

func cloneScoreMap(source map[string]map[string]float64) map[string]map[string]float64 {
	if source == nil {
		return make(map[string]map[string]float64)
	}
	clone := make(map[string]map[string]float64, len(source))
	for key, values := range source {
		copied := make(map[string]float64, len(values))
		for action, score := range values {
			copied[action] = score
		}
		clone[key] = copied
	}
	return clone
}

func cloneCountMap(source map[string]map[string]int) map[string]map[string]int {
	if source == nil {
		return make(map[string]map[string]int)
	}
	clone := make(map[string]map[string]int, len(source))
	for key, values := range source {
		copied := make(map[string]int, len(values))
		for action, count := range values {
			copied[action] = count
		}
		clone[key] = copied
	}
	return clone
}

func partyNames(party []*pokemon) []string {
	result := make([]string, 0, len(party))
	for _, mon := range party {
		if mon == nil || mon.base.Name == "" {
			continue
		}
		result = append(result, mon.base.Name)
	}
	return result
}

func countScoreEntries(scores map[string]map[string]float64) int {
	count := 0
	for _, actions := range scores {
		count += len(actions)
	}
	return count
}

func countCountEntries(counts map[string]map[string]int) int {
	count := 0
	for _, actions := range counts {
		count += len(actions)
	}
	return count
}

func nowISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func loadLearningAIFromPolicy(policy *savedPolicy) *learningAi {
	if policy == nil {
		return newLearningAi()
	}
	la := &learningAi{
		policy:      cloneActionMap(policy.Policy),
		scores:      cloneScoreMap(policy.Scores),
		counts:      cloneCountMap(policy.Counts),
		seen:        make(map[string]map[string]bool),
		decayFactor: 0.98,
	}
	if la.policy == nil {
		la.policy = make(map[string][]string)
	}
	if la.scores == nil {
		la.scores = make(map[string]map[string]float64)
	}
	if la.counts == nil {
		la.counts = make(map[string]map[string]int)
	}
	for stateKey, actions := range la.policy {
		la.seen[stateKey] = make(map[string]bool)
		for _, action := range actions {
			la.seen[stateKey][action] = true
		}
	}
	return la
}

type staticPolicyAi struct {
	policy map[string][]string
	scores map[string]map[string]float64
	counts map[string]map[string]int
}

func newStaticPolicyAiFromPolicy(policy *savedPolicy) *staticPolicyAi {
	if policy == nil {
		return &staticPolicyAi{policy: map[string][]string{}, scores: map[string]map[string]float64{}, counts: map[string]map[string]int{}}
	}
	return &staticPolicyAi{
		policy: cloneActionMap(policy.Policy),
		scores: cloneScoreMap(policy.Scores),
		counts: cloneCountMap(policy.Counts),
	}
}

func (spa *staticPolicyAi) evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int) {
	if len(actions) == 0 {
		return nil, 0
	}
	stateKey := discretizeBattleState(bs)
	bestAction := actions[0]
	bestScore := -1e18
	for _, action := range actions {
		score := spa.scoreFor(stateKey, actionKey(action))
		if score > bestScore {
			bestScore = score
			bestAction = action
		}
	}
	return bestAction, int(bestScore)
}

func (spa *staticPolicyAi) evaluteSwitchIns(bs battleState, mons []*pokemon, opponentSlot *slot) *pokemon {
	if len(mons) == 0 {
		return nil
	}
	stateKey := discretizeBattleState(bs)
	bestMon := mons[0]
	bestScore := -1e18
	for _, mon := range mons {
		score := spa.scoreFor(stateKey, actionKeyForSwitch(mon))
		if score > bestScore {
			bestScore = score
			bestMon = mon
		}
	}
	return bestMon
}

func (spa *staticPolicyAi) shouldSwitch(bs battleState, slot *slot, score int, party []*pokemon) bool {
	if spa == nil || len(party) <= 1 {
		return false
	}
	stateKey := discretizeBattleState(bs)
	for _, mon := range party {
		if mon == nil || mon == slot.mon || mon.fainted {
			continue
		}
		if spa.scoreFor(stateKey, actionKeyForSwitch(mon)) > float64(score) {
			return true
		}
	}
	return false
}

func (spa *staticPolicyAi) scoreFor(stateKey, actionKey string) float64 {
	if spa == nil {
		return 0
	}
	if spa.scores == nil {
		return 0
	}
	if stateScores, ok := spa.scores[stateKey]; ok {
		if score, ok := stateScores[actionKey]; ok {
			return score
		}
	}
	if spa.policy != nil {
		for _, allowedKey := range spa.policy[stateKey] {
			if allowedKey == actionKey {
				return 0
			}
		}
	}
	return -1e18
}

func validatePolicyCompatibility(policy *savedPolicy, playerParty, opponentParty []*pokemon) error {
	if policy == nil {
		return fmt.Errorf("policy is nil")
	}
	toNames := func(p []*pokemon) []string {
		result := make([]string, 0, len(p))
		for _, mon := range p {
			if mon == nil || mon.base.Name == "" {
				continue
			}
			result = append(result, mon.base.Name)
		}
		return result
	}
	if !stringSliceEquivalent(policy.PlayerParty, toNames(playerParty)) {
		return fmt.Errorf("policy player party does not match input player party")
	}
	if !stringSliceEquivalent(policy.OpponentParty, toNames(opponentParty)) {
		return fmt.Errorf("policy opponent party does not match input opponent party")
	}
	return nil
}

func stringSliceEquivalent(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	lookup := make(map[string]int)
	for _, value := range b {
		lookup[strings.ToLower(value)]++
	}
	for _, value := range a {
		key := strings.ToLower(value)
		if lookup[key] <= 0 {
			return false
		}
		lookup[key]--
	}
	return true
}

func (la *learningAi) ensureState(stateKey string) {
	if la == nil {
		return
	}
	if la.policy == nil {
		la.policy = make(map[string][]string)
	}
	if la.scores == nil {
		la.scores = make(map[string]map[string]float64)
	}
	if la.counts == nil {
		la.counts = make(map[string]map[string]int)
	}
	if la.seen == nil {
		la.seen = make(map[string]map[string]bool)
	}
	if _, ok := la.policy[stateKey]; !ok {
		la.policy[stateKey] = []string{}
	}
	if _, ok := la.scores[stateKey]; !ok {
		la.scores[stateKey] = make(map[string]float64)
	}
	if _, ok := la.counts[stateKey]; !ok {
		la.counts[stateKey] = make(map[string]int)
	}
	if _, ok := la.seen[stateKey]; !ok {
		la.seen[stateKey] = make(map[string]bool)
	}
}

func (la *learningAi) decayScores(factor float64) {
	if la == nil {
		return
	}
	if factor <= 0 {
		factor = la.decayFactor
	}
	if factor == 0 {
		factor = 0.98
	}
	for stateKey, actionScores := range la.scores {
		for action, score := range actionScores {
			actionScores[action] = score * factor
			if la.counts[stateKey] != nil && la.counts[stateKey][action] > 0 {
				actionScores[action] += float64(la.counts[stateKey][action]) * 0.5
			}
		}
	}
}

func (la *learningAi) recordStateAction(stateKey string, action string) {
	if la == nil || stateKey == "" || action == "" {
		return
	}
	la.ensureState(stateKey)
	la.policy[stateKey] = appendUnique(la.policy[stateKey], action)
	la.counts[stateKey][action]++
	if la.seen[stateKey][action] {
		return
	}
	la.seen[stateKey][action] = true
	la.history = append(la.history, stateActionEntry{stateKey: stateKey, action: action})
}

func (la *learningAi) recordBattleOutcome(stats *battleStatistics) {
	if la == nil || stats == nil {
		return
	}

	la.decayScores(la.decayFactor)
	reward := stats.outcomeScore()
	for _, entry := range la.history {
		la.ensureState(entry.stateKey)
		count := float64(la.counts[entry.stateKey][entry.action])
		la.scores[entry.stateKey][entry.action] += reward * (0.25 + count*0.1)
		if stateHasCritKillRisk(entry.stateKey) {
			la.scores[entry.stateKey][entry.action] -= 1500.0
		}
	}
	la.history = nil
	la.seen = make(map[string]map[string]bool)
}

func (la *learningAi) setPolicyAction(bs battleState, action string) {
	if la == nil {
		return
	}
	stateKey := discretizeBattleState(bs)
	la.ensureState(stateKey)
	la.policy[stateKey] = []string{action}
}

func (la *learningAi) evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int) {
	if len(actions) == 0 {
		return nil, 0
	}
	if la == nil {
		return rnbAi{}.evaluateActions(bs, actions)
	}
	stateKey := discretizeBattleState(bs)
	la.ensureState(stateKey)
	firstAction := actions[0]
	for _, action := range actions {
		key := actionKey(action)
		la.policy[stateKey] = appendUnique(la.policy[stateKey], key)
		if _, ok := la.scores[stateKey][key]; !ok {
			la.scores[stateKey][key] = 0
		}
		if _, ok := la.counts[stateKey][key]; !ok {
			la.counts[stateKey][key] = 0
		}
	}
	bestAction := firstAction
	bestScore := -1e18
	riskPenalty := 0.0
	if stateHasCritKillRisk(stateKey) {
		riskPenalty = 1500.0
	}
	for _, action := range actions {
		key := actionKey(action)
		score := la.scores[stateKey][key]
		if count := la.counts[stateKey][key]; count > 0 {
			score += float64(count) * 2.5
		}
		score -= riskPenalty
		if score > bestScore {
			bestAction = action
			bestScore = score
		}
	}
	if bestScore <= -1e18 {
		la.recordStateAction(stateKey, actionKey(firstAction))
		return firstAction, 0
	}
	la.recordStateAction(stateKey, actionKey(bestAction))
	return bestAction, int(bestScore)
}

func (la *learningAi) evaluteSwitchIns(bs battleState, mons []*pokemon, opponentSlot *slot) *pokemon {
	if len(mons) == 0 {
		return nil
	}
	if la == nil {
		return rnbAi{}.evaluteSwitchIns(bs, mons, opponentSlot)
	}
	stateKey := discretizeBattleState(bs)
	la.ensureState(stateKey)
	bestMon := mons[0]
	bestScore := -1e18
	riskPenalty := 0.0
	if stateHasCritKillRisk(stateKey) {
		riskPenalty = 1500.0
	}
	for _, mon := range mons {
		key := actionKeyForSwitch(mon)
		la.policy[stateKey] = appendUnique(la.policy[stateKey], key)
		if _, ok := la.scores[stateKey][key]; !ok {
			la.scores[stateKey][key] = 0
		}
		if _, ok := la.counts[stateKey][key]; !ok {
			la.counts[stateKey][key] = 0
		}
		score := la.scores[stateKey][key]
		if count := la.counts[stateKey][key]; count > 0 {
			score += float64(count) * 2.5
		}
		score -= riskPenalty
		if score > bestScore {
			bestMon = mon
			bestScore = score
		}
	}
	la.recordStateAction(stateKey, actionKeyForSwitch(bestMon))
	return bestMon
}

func (la *learningAi) shouldSwitch(bs battleState, slot *slot, score int, party []*pokemon) bool {
	if la == nil || len(party) <= 1 {
		return false
	}
	stateKey := discretizeBattleState(bs)
	la.ensureState(stateKey)
	for _, mon := range party {
		if mon == slot.mon || mon.fainted {
			continue
		}
		key := actionKeyForSwitch(mon)
		if learned, ok := la.scores[stateKey][key]; ok && learned > float64(score) {
			return true
		}
	}
	return false
}

func stateHasCritKillRisk(stateKey string) bool {
	if stateKey == "" {
		return false
	}
	return strings.Contains(stateKey, "\"opponent_crit_kill_risk\":true")
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
