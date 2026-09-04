package main

type singleBattleState struct {
	activePlayerSlot   *slot
	activeOpponentSlot *slot
	player             *trainer
	opponent           *trainer
	actions            actionQueue
	weather            weatherState
	fieldEffects       map[fieldEffect]int
	err                error
	initialPlayer      trainer
	initialOpponent    trainer
	initialWeather     weatherState
	statistics         battleStatistics
}

func (sbs *singleBattleState) execute() error {
	vprintln("\nStarting battle...")

	for k := 0; !sbs.player.lost && !sbs.opponent.lost; k++ {
		vprintln("=====")
		vprintf("Turn %d:", k+1)
		vprintf("%s %d/%d - %s %d/%d", sbs.activePlayerSlot.mon.base.Name, sbs.activePlayerSlot.mon.hp, sbs.activePlayerSlot.mon.maxHP(), sbs.activeOpponentSlot.mon.base.Name, sbs.activeOpponentSlot.mon.hp, sbs.activeOpponentSlot.mon.maxHP())

		sbs.gatherActions()
		sbs.actions.sort(sbs)
		for len(sbs.actions.queue) > 0 {
			action, _ := sbs.actions.queue.pop()
			action.invoke(sbs)
			if sbs.err != nil {
				return sbs.err
			}
		}
		resolveEndOfTurn(sbs)
		// if the end of turn causes mons to faint, empty the queue for replace actions
		for len(sbs.actions.queue) > 0 {
			action, _ := sbs.actions.queue.pop()
			action.invoke(sbs)
			if sbs.err != nil {
				return sbs.err
			}
		}
	}
	vprintln("=====")
	vprintln("Ending battle...")

	return nil
}

func (sbs *singleBattleState) setError(err error) {
	sbs.err = err
}

func (sbs *singleBattleState) gatherActions() {
	sbs.actions.queue.push(chooseNextAction(sbs, sbs.activePlayerSlot, sbs.player.pokemonParty, sbs.player.ai))
	sbs.actions.queue.push(chooseNextAction(sbs, sbs.activeOpponentSlot, sbs.opponent.pokemonParty, sbs.opponent.ai))
}

func (sbs *singleBattleState) getAllSlots() []*slot {
	return []*slot{
		sbs.activePlayerSlot,
		sbs.activeOpponentSlot,
	}
}

func (sbs *singleBattleState) getOtherSlots(s *slot) []*slot {
	if s == sbs.activePlayerSlot {
		return []*slot{sbs.activeOpponentSlot}
	}
	return []*slot{sbs.activePlayerSlot}
}

func (sbs *singleBattleState) getOpponentSlot(s *slot) *slot {
	if s == sbs.activePlayerSlot {
		return sbs.activeOpponentSlot
	}
	return sbs.activePlayerSlot
}

func (sbs *singleBattleState) getActions() *actionQueue {
	return &sbs.actions
}

func (sbs *singleBattleState) getWeather() weatherState {
	return sbs.weather
}

func (sbs *singleBattleState) setWeather(w weatherState) {
	sbs.weather = w
	w.onset()
}

func (sbs *singleBattleState) getFieldEffects() map[fieldEffect]int {
	return sbs.fieldEffects
}

func (sbs *singleBattleState) getStatistics() *battleStatistics {
	return &sbs.statistics
}

func (sbs *singleBattleState) recordStatistics() {
	sbs.statistics.record(sbs.player)
}

func (sbs *singleBattleState) printStatistics() {
	sbs.statistics.print(sbs.initialPlayer.pokemonParty)
}

func (sbs *singleBattleState) reset() error {
	playerParty := clonePokemonParty(sbs.initialPlayer.pokemonParty)
	opponentParty := clonePokemonParty(sbs.initialOpponent.pokemonParty)
	resetPokemonPartyPPs(playerParty)
	resetPokemonPartyPPs(opponentParty)

	player := sbs.initialPlayer
	player.pokemonParty = playerParty
	player.lost = false
	player.fieldEffects = cloneFieldEffects(sbs.initialPlayer.fieldEffects)

	opponent := sbs.initialOpponent
	opponent.pokemonParty = opponentParty
	opponent.lost = false
	opponent.fieldEffects = cloneFieldEffects(sbs.initialOpponent.fieldEffects)

	sbs.activePlayerSlot = &slot{
		mon:       playerParty[0],
		trainer:   &player,
		firstTurn: true,
	}
	sbs.activeOpponentSlot = &slot{
		mon:       opponentParty[0],
		trainer:   &opponent,
		firstTurn: true,
	}
	sbs.player = &player
	sbs.opponent = &opponent
	sbs.actions = actionQueue{queue: make(priorityQueue[action], 0, 3)}
	sbs.weather = sbs.initialWeather
	sbs.err = nil

	sbs.setWeather(sbs.initialWeather)
	resolveOnEntry(sbs)
	return nil
}

func initSingleBattleState(player, opponent trainer, playerParty, opponentParty []*pokemon, weather weatherState) *singleBattleState {
	player.pokemonParty = playerParty
	opponent.pokemonParty = opponentParty

	res := singleBattleState{
		activePlayerSlot: &slot{
			mon:       playerParty[0],
			trainer:   &player,
			firstTurn: true,
		},
		activeOpponentSlot: &slot{
			mon:       opponentParty[0],
			trainer:   &opponent,
			firstTurn: true,
		},
		player:   &player,
		opponent: &opponent,
		actions: actionQueue{
			queue: make(priorityQueue[action], 0, 3),
		},
		initialPlayer:   player,
		initialOpponent: opponent,
		initialWeather:  weather,
		statistics:      newBattleStatistics(playerParty),
	}

	res.initialPlayer.pokemonParty = clonePokemonParty(playerParty)
	res.initialOpponent.pokemonParty = clonePokemonParty(opponentParty)
	res.initialPlayer.fieldEffects = cloneFieldEffects(player.fieldEffects)
	res.initialOpponent.fieldEffects = cloneFieldEffects(opponent.fieldEffects)
	res.player.pokemonParty = playerParty
	res.opponent.pokemonParty = opponentParty

	res.setWeather(weather)
	resolveOnEntry(&res)

	return &res
}
