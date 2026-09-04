package main

func cloneFieldEffects(src map[fieldEffect]int) map[fieldEffect]int {
	if src == nil {
		return make(map[fieldEffect]int)
	}
	clone := make(map[fieldEffect]int, len(src))
	for k, v := range src {
		clone[k] = v
	}
	return clone
}

func clonePokemonParty(src []*pokemon) []*pokemon {
	clone := make([]*pokemon, len(src))
	for i, mon := range src {
		clone[i] = clonePokemon(mon)
	}
	return clone
}

func clonePokemon(p *pokemon) *pokemon {
	if p == nil {
		return nil
	}

	copyP := *p
	copyP.base.Types = append([]pokemonType(nil), p.base.Types...)
	copyP.base.Stats = make(map[string]int, len(p.base.Stats))
	for stat, value := range p.base.Stats {
		copyP.base.Stats[stat] = value
	}
	copyP.ivs = append([]int(nil), p.ivs...)
	copyP.moves = make([]*Move, len(p.moves))
	moveCopies := make(map[*Move]*Move, len(p.moves))
	for i, move := range p.moves {
		copyP.moves[i] = cloneMove(move)
		if move != nil {
			moveCopies[move] = copyP.moves[i]
		}
	}
	copyP.stats = append([]int(nil), p.stats...)
	copyP.stages = append([]int(nil), p.stages...)
	if p.lockedMove != nil {
		if move, ok := moveCopies[p.lockedMove]; ok {
			copyP.lockedMove = move
		} else {
			copyP.lockedMove = cloneMove(p.lockedMove)
		}
	}

	copyP.ailments = make(map[ailmentState]*ailment, len(p.ailments))
	for state, ailment := range p.ailments {
		if ailment == nil {
			continue
		}
		copyA := *ailment
		copyP.ailments[state] = &copyA
	}

	if p.item != nil {
		copyP.item, _ = registerItem(p.item.state, &copyP)
		copyP.item.consumed = p.item.consumed
	}

	return &copyP
}

func cloneMove(move *Move) *Move {
	if move == nil {
		return nil
	}

	copyMove := *move
	copyMove.StatChanges = make(map[string]int, len(move.StatChanges))
	for stat, change := range move.StatChanges {
		copyMove.StatChanges[stat] = change
	}
	return &copyMove
}
