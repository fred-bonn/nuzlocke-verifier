package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBattleStatisticsRecord(t *testing.T) {
	survivor := &pokemon{base: BasePokemon{Name: "survivor"}}
	fainted := &pokemon{base: BasePokemon{Name: "fainted"}, fainted: true}
	player := &trainer{pokemonParty: []*pokemon{survivor, fainted}}
	statistics := newBattleStatistics(player.pokemonParty)

	statistics.record(player)
	player.lost = true
	statistics.record(player)

	if statistics.battleCount != 2 || statistics.winCount != 1 {
		t.Fatalf("unexpected battle totals: battles=%d wins=%d", statistics.battleCount, statistics.winCount)
	}
	if got := statistics.pokemonSurvivors[0]; got != 2 {
		t.Fatalf("survivor count = %d, want 2", got)
	}
	if got := statistics.pokemonSurvivors[1]; got != 0 {
		t.Fatalf("fainted Pokemon survivor count = %d, want 0", got)
	}
}

func TestBattleStatisticsOutcomeScorePrefersSafetyOverSurvival(t *testing.T) {
	stats := newBattleStatistics([]*pokemon{{base: BasePokemon{Name: "a"}}, {base: BasePokemon{Name: "b"}}})
	stats.battleCount = 1
	stats.winCount = 0
	if got := stats.outcomeScore(); got >= 0 {
		t.Fatalf("lost battle should be scored negatively: got %.2f", got)
	}

	stats = newBattleStatistics([]*pokemon{{base: BasePokemon{Name: "a"}}, {base: BasePokemon{Name: "b"}}})
	stats.battleCount = 1
	stats.winCount = 1
	stats.pokemonSurvivors = []int{1, 1}
	if got := stats.outcomeScore(); got <= 0 {
		t.Fatalf("won battle with surviving party should be scored positively: got %.2f", got)
	}
}

func TestBattleStatisticsOutcomeScoreHeavilyRewardsWinsAndPunishesLosses(t *testing.T) {
	winStats := newBattleStatistics([]*pokemon{{base: BasePokemon{Name: "a"}}, {base: BasePokemon{Name: "b"}}})
	winStats.battleCount = 1
	winStats.winCount = 1
	winStats.pokemonSurvivors = []int{1, 1}
	if got := winStats.outcomeScore(); got <= 500000.0 {
		t.Fatalf("winning should be rewarded heavily: got %.2f", got)
	}

	lossStats := newBattleStatistics([]*pokemon{{base: BasePokemon{Name: "a"}}, {base: BasePokemon{Name: "b"}}})
	lossStats.battleCount = 1
	lossStats.winCount = 0
	lossStats.pokemonSurvivors = []int{1, 0}
	if got := lossStats.outcomeScore(); got >= -500000.0 {
		t.Fatalf("losing should be catastrophic: got %.2f", got)
	}
}

func TestLearningAiChoosesSafeStateActionByScore(t *testing.T) {
	player := testSwitchPokemon("player", 100, 100, 100, 100, nil)
	opponent := testSwitchPokemon("opponent", 100, 50, 100, 100, nil)
	bs := testSwitchBattleState(player, opponent, opponent)
	bs.player.player = true

	safeMove := &Move{Name: "safe", Power: 1, PP: 1, Class: physicalClass}
	riskyMove := &Move{Name: "risky", Power: 100, PP: 1, Class: physicalClass}
	la := newLearningAi()
	stateKey := discretizeBattleState(bs)
	la.policy[stateKey] = []string{"move:safe", "move:risky"}
	la.scores[stateKey] = map[string]float64{"move:safe": 200, "move:risky": -2000}

	got, _ := la.evaluateActions(bs, []*moveAction{{userSlot: bs.activePlayerSlot, targetSlot: bs.activeOpponentSlot, move: safeMove}, {userSlot: bs.activePlayerSlot, targetSlot: bs.activeOpponentSlot, move: riskyMove}})
	if got.move != safeMove {
		t.Fatalf("learning AI chose risky action despite negative score; got %+v", got.move)
	}
}

func TestLearningAiTracksCountsAndDecaysStaleScores(t *testing.T) {
	la := newLearningAi()
	stateKey := "state"
	actionKey := "move:test"
	la.ensureState(stateKey)
	la.scores[stateKey][actionKey] = 100
	la.counts[stateKey][actionKey] = 5

	la.decayScores(0.8)
	if la.counts[stateKey][actionKey] != 5 {
		t.Fatalf("counts should remain visible: got %d", la.counts[stateKey][actionKey])
	}
	if la.scores[stateKey][actionKey] >= 100 {
		t.Fatalf("scores should decay over time: got %f", la.scores[stateKey][actionKey])
	}
}
func TestLearningAiAllowsReasonableSwitchDecisionsAndTracksReplacement(t *testing.T) {
	current := testSwitchPokemon("current", 100, 100, 100, 100, nil)
	replacement := testSwitchPokemon("replacement", 100, 100, 100, 100, nil)
	opponent := testSwitchPokemon("opponent", 100, 50, 100, 100, nil)
	bs := testSwitchBattleState(current, replacement, opponent)
	bs.player.player = true
	current.hp = 20

	la := newLearningAi()
	if !la.shouldSwitch(bs, bs.activePlayerSlot, 10, []*pokemon{current, replacement}) {
		t.Fatal("learning AI should consider switching when the current mon is low HP and no score exists yet")
	}

	stateKey := discretizeBattleState(bs)
	if _, ok := la.policy[stateKey]; ok && len(la.policy[stateKey]) > 0 {
		t.Fatal("learning AI should only record the switch decision after an actual action is selected")
	}

	action := chooseNextAction(bs, bs.activePlayerSlot, []*pokemon{current, replacement}, la)
	if _, ok := action.(*switchAction); !ok {
		t.Fatalf("expected a switch action to be selected, got %T", action)
	}
	stateKey = discretizeBattleState(bs)
	if _, ok := la.policy[stateKey]; !ok {
		t.Fatal("learning AI did not record the selected replacement in policy history")
	}
	if !strings.Contains(strings.Join(la.policy[stateKey], ","), "switch:replacement") {
		t.Fatalf("replacement action was not recorded for learning: %#v", la.policy[stateKey])
	}
}
func TestBattleStatisticsRewardsFullPartySurvival(t *testing.T) {
	stats := newBattleStatistics([]*pokemon{{base: BasePokemon{Name: "a"}}, {base: BasePokemon{Name: "b"}}})
	stats.battleCount = 2
	stats.winCount = 2
	stats.pokemonSurvivors = []int{2, 2}

	if !stats.allPartySurvived() {
		t.Fatal("all party survival should be detected when every member survives every battle")
	}
	if got := stats.outcomeScore(); got < 50000 {
		t.Fatalf("full-party survival should carry a very large bonus: got %.2f", got)
	}
}

func TestBattleStatisticsPrioritizesWholePartySurvivalOverPartialWins(t *testing.T) {
	fullParty := newBattleStatistics([]*pokemon{{base: BasePokemon{Name: "a"}}, {base: BasePokemon{Name: "b"}}, {base: BasePokemon{Name: "c"}}})
	fullParty.battleCount = 4
	fullParty.winCount = 3
	fullParty.pokemonSurvivors = []int{4, 4, 4}

	partialWin := newBattleStatistics([]*pokemon{{base: BasePokemon{Name: "a"}}, {base: BasePokemon{Name: "b"}}, {base: BasePokemon{Name: "c"}}})
	partialWin.battleCount = 4
	partialWin.winCount = 4
	partialWin.pokemonSurvivors = []int{4, 4, 0}

	if got := fullParty.outcomeScore(); got <= partialWin.outcomeScore() {
		t.Fatalf("whole-party survival should score higher than partial survival after a win: full=%.2f partial=%.2f", got, partialWin.outcomeScore())
	}
}

func TestDiscretizeBattleStateFlagsOpponentCritKillRisk(t *testing.T) {
	player := testSwitchPokemon("player", 100, 100, 100, 100, nil)
	opponent := testSwitchPokemon("opponent", 100, 50, 100, 100, &Move{Name: "critical-slap", Power: 80, PP: 1, Class: physicalClass, CritRate: 1})
	bs := testSwitchBattleState(player, opponent, opponent)
	bs.player.player = true
	player.hp = 50

	state := discretizeBattleState(bs)
	if !contains(state, `"opponent_crit_kill_risk":true`) {
		t.Fatalf("state did not detect opponent lethal crit risk: %s", state)
	}
	if !contains(state, `"player_hp_percent":50`) {
		t.Fatalf("state did not encode player hp percent: %s", state)
	}
	if !contains(state, `"opponent_hp_percent":100`) {
		t.Fatalf("state did not encode opponent hp percent: %s", state)
	}
}

func TestPolicyRoundTripPreservesLearnedRewards(t *testing.T) {
	la := newLearningAi()
	la.policy["state"] = []string{"move:Bubble Beam", "switch:Ponyta"}
	la.scores["state"] = map[string]float64{"move:Bubble Beam": 42.5, "switch:Ponyta": -12.25}
	la.counts["state"] = map[string]int{"move:Bubble Beam": 9, "switch:Ponyta": 3}

	loaded := loadLearningAIFromPolicy(&savedPolicy{
		Policy: cloneActionMap(la.policy),
		Scores: cloneScoreMap(la.scores),
		Counts: cloneCountMap(la.counts),
	})

	if got := loaded.scores["state"]["move:Bubble Beam"]; got != 42.5 {
		t.Fatalf("loaded score was reset: got %.2f want %.2f", got, 42.5)
	}
	if got := loaded.counts["state"]["switch:Ponyta"]; got != 3 {
		t.Fatalf("loaded count was reset: got %d want %d", got, 3)
	}
	if got := len(loaded.policy["state"]); got != 2 {
		t.Fatalf("loaded policy actions were lost: got %d want 2", got)
	}
}

func TestPolicySaveLoadSmokeTest(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir temp dir failed: %v", err)
	}
	defer func() { _ = os.Chdir(cwd) }()

	playerParty := []*pokemon{{base: BasePokemon{Name: "Horsea"}}, {base: BasePokemon{Name: "Ponyta"}}}
	opponentParty := []*pokemon{{base: BasePokemon{Name: "Dwebble"}}}
	la := newLearningAi()
	la.policy["state"] = []string{"move:Bubble Beam"}
	la.scores["state"] = map[string]float64{"move:Bubble Beam": 128.75}
	la.counts["state"] = map[string]int{"move:Bubble Beam": 11}

	if err := savePolicyToDisk(la, "player.txt", "rnb_trainer_1.txt", playerParty, opponentParty); err != nil {
		t.Fatalf("save policy failed: %v", err)
	}

	path := filepath.Join("policies", "player__vs__rnb_trainer_1.json")
	policy, err := loadPolicyFromDisk(path)
	if err != nil {
		t.Fatalf("load policy failed: %v", err)
	}
	loaded := loadLearningAIFromPolicy(policy)
	if got := loaded.scores["state"]["move:Bubble Beam"]; got != 128.75 {
		t.Fatalf("score value did not survive disk round trip: got %.2f want %.2f", got, 128.75)
	}
	if got := loaded.counts["state"]["move:Bubble Beam"]; got != 11 {
		t.Fatalf("count value did not survive disk round trip: got %d want %d", got, 11)
	}
	if err := validatePolicyCompatibility(policy, playerParty, opponentParty); err != nil {
		t.Fatalf("saved policy was incompatible after reload: %v", err)
	}
}

func TestPolicyCliEndToEndSmoke(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	defer func() { _ = os.Chdir(cwd) }()

	cmdSave := exec.Command("go", "run", ".", "--player-learning-ai", "--save-policy", "--iterations", "10", "data/player.txt", "data/rnb_trainer_1.txt")
	cmdSave.Dir = cwd
	output, err := cmdSave.CombinedOutput()
	if err != nil {
		t.Fatalf("save policy CLI failed: %v\n%s", err, output)
	}
	if !bytes.Contains(output, []byte("policy saved to")) {
		t.Fatalf("save CLI did not report policy saved:\n%s", output)
	}

	policyPath := filepath.Join(cwd, "policies", "player__vs__rnb_trainer_1.json")
	if _, err := os.Stat(policyPath); err != nil {
		t.Fatalf("saved policy file was not created at %s: %v", policyPath, err)
	}

	cmdLoad := exec.Command("go", "run", ".", "--policy-file", policyPath, "data/player.txt", "data/rnb_trainer_1.txt", "--iterations", "1")
	cmdLoad.Dir = cwd
	output, err = cmdLoad.CombinedOutput()
	if err != nil {
		t.Fatalf("load policy CLI failed: %v\n%s", err, output)
	}
	if !bytes.Contains(output, []byte("loaded policy from")) {
		t.Fatalf("load CLI did not report the loaded policy summary:\n%s", output)
	}
}

func TestLoadedPolicyUsesStaticScores(t *testing.T) {
	policy := &savedPolicy{
		PlayerParty:   []string{"Horsea"},
		OpponentParty: []string{"Dwebble"},
		Policy: map[string][]string{
			"state": {"move:Bubble Beam", "move:Twister"},
		},
		Scores: map[string]map[string]float64{
			"state": {"move:Bubble Beam": 100, "move:Twister": 5},
		},
		Counts: map[string]map[string]int{
			"state": {"move:Bubble Beam": 12, "move:Twister": 2},
		},
	}

	ai := newStaticPolicyAiFromPolicy(policy)
	if ai == nil {
		t.Fatal("expected static policy AI to be created")
	}
	if got, _ := ai.evaluateActions(nil, []*moveAction{{move: &Move{Name: "Bubble Beam"}}, {move: &Move{Name: "Twister"}}}); got.move.Name != "Bubble Beam" {
		t.Fatalf("loaded policy should choose the highest scored action, got %s", got.move.Name)
	}
}

func TestPolicyPathAndCompatibilityValidation(t *testing.T) {
	path := policyPathForInputs("data/player.txt", "data/rnb_trainer_1.txt")
	if path != "policies/player__vs__rnb_trainer_1.json" {
		t.Fatalf("unexpected policy path: %s", path)
	}

	playerParty := []*pokemon{{base: BasePokemon{Name: "Horsea"}}, {base: BasePokemon{Name: "Ponyta"}}}
	opponentParty := []*pokemon{{base: BasePokemon{Name: "Dwebble"}}}
	policy := &savedPolicy{
		PlayerParty:   []string{"Horsea", "Ponyta"},
		OpponentParty: []string{"Dwebble"},
		Policy:        map[string][]string{"state": {"move:Bubble Beam", "switch:Ponyta"}},
		Scores:        map[string]map[string]float64{"state": {"move:Bubble Beam": 1}},
		Counts:        map[string]map[string]int{"state": {"move:Bubble Beam": 1}},
	}
	if err := validatePolicyCompatibility(policy, playerParty, opponentParty); err != nil {
		t.Fatalf("policy compatibility failed unexpectedly: %v", err)
	}

	policy.PlayerParty = []string{"Wrong"}
	if err := validatePolicyCompatibility(policy, playerParty, opponentParty); err == nil {
		t.Fatal("policy compatibility should reject mismatched party")
	}
	if err := validatePolicyCompatibility(nil, playerParty, opponentParty); err == nil {
		t.Fatal("nil policy should be rejected")
	}
}

func TestPolicyLoaderInitializesMissingMaps(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "policy.json")
	payload := []byte(`{"player_party":["Horsea"],"opponent_party":["Dwebble"]}`)
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write policy fixture: %v", err)
	}

	policy, err := loadPolicyFromDisk(path)
	if err != nil {
		t.Fatalf("loadPolicyFromDisk returned unexpected error: %v", err)
	}
	if policy.Policy == nil || policy.Scores == nil || policy.Counts == nil {
		t.Fatal("loadPolicyFromDisk should initialize missing maps")
	}
	if len(policy.Policy) != 0 || len(policy.Scores) != 0 || len(policy.Counts) != 0 {
		t.Fatalf("unexpected non-empty maps after load: policy=%d scores=%d counts=%d", len(policy.Policy), len(policy.Scores), len(policy.Counts))
	}
}

func TestStaticPolicyAiUsesPolicyFallbackAndSwitchGuard(t *testing.T) {
	policy := &savedPolicy{
		PlayerParty:   []string{"Horsea"},
		OpponentParty: []string{"Dwebble"},
		Policy: map[string][]string{
			"state": {"move:Bubble Beam", "switch:Ponyta"},
		},
		Scores: map[string]map[string]float64{
			"state": {"move:Bubble Beam": 15, "switch:Ponyta": 22},
		},
		Counts: map[string]map[string]int{
			"state": {"move:Bubble Beam": 2},
		},
	}
	ai := newStaticPolicyAiFromPolicy(policy)
	if got := ai.scoreFor("state", "move:Bubble Beam"); got != 15 {
		t.Fatalf("scoreFor returned unexpected value: got %.2f want 15", got)
	}
	if got := ai.scoreFor("state", "move:Missing"); got != -1e18 {
		t.Fatalf("unknown action should return the failure score: got %.2f", got)
	}
	if got := ai.scoreFor("state", "switch:Ponyta"); got != 22 {
		t.Fatalf("switch score was not loaded: got %.2f", got)
	}
	if got := ai.scoreFor("missing", "switch:Ponyta"); got != -1e18 {
		t.Fatalf("missing state should reject actions: got %.2f", got)
	}

	if ai.scoreFor("state", "switch:Ponyta") <= ai.scoreFor("state", "move:Bubble Beam") {
		t.Fatal("switch action should outrank the backup move in the saved policy scores")
	}
	if ai.scoreFor("missing", "move:Bubble Beam") != -1e18 {
		t.Fatal("missing state should not accept values outside the stored policy map")
	}
}

func TestLearningAiRecordBattleOutcomeAppliesCritPenaltyAndClearsHistory(t *testing.T) {
	makeState := func(risk bool) string {
		if risk {
			return "{\"opponent_crit_kill_risk\":true}"
		}
		return "{\"opponent_crit_kill_risk\":false}"
	}

	makeAi := func(stateKey string) *learningAi {
		la := newLearningAi()
		la.ensureState(stateKey)
		la.scores[stateKey]["move:critical-slap"] = 120
		la.counts[stateKey]["move:critical-slap"] = 2
		la.history = []stateActionEntry{{stateKey: stateKey, action: "move:critical-slap"}}
		la.seen[stateKey] = map[string]bool{"move:critical-slap": true}
		return la
	}

	safeState := makeState(false)
	riskyState := makeState(true)
	safeAI := makeAi(safeState)
	riskyAI := makeAi(riskyState)

	stats := newBattleStatistics([]*pokemon{{base: BasePokemon{Name: "a"}}, {base: BasePokemon{Name: "b"}}})
	stats.battleCount = 2
	stats.winCount = 1
	stats.pokemonSurvivors = []int{2, 1}

	safeAI.recordBattleOutcome(&stats)
	riskyAI.recordBattleOutcome(&stats)

	if len(safeAI.history) != 0 || len(riskyAI.history) != 0 {
		t.Fatal("history should be cleared after recording the outcome")
	}
	if safeAI.scores[safeState]["move:critical-slap"] <= riskyAI.scores[riskyState]["move:critical-slap"] {
		t.Fatalf("crit-risk states should lose more reward than safe states; safe=%.2f risky=%.2f", safeAI.scores[safeState]["move:critical-slap"], riskyAI.scores[riskyState]["move:critical-slap"])
	}
}
