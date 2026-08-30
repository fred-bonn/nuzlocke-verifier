package main

import "log"

type battleStatistics struct {
	winCount         int
	battleCount      int
	pokemonSurvivors []int
}

func newBattleStatistics(party []*pokemon) battleStatistics {
	return battleStatistics{
		pokemonSurvivors: make([]int, len(party)),
	}
}

func (bs *battleStatistics) record(player *trainer) {
	bs.battleCount++
	if !player.lost {
		bs.winCount++
	}

	for index, mon := range player.pokemonParty {
		if !mon.fainted {
			bs.pokemonSurvivors[index]++
		}
	}
}

func (bs *battleStatistics) partySurvivalRatio() float64 {
	if bs == nil || bs.battleCount == 0 || len(bs.pokemonSurvivors) == 0 {
		return 0
	}

	alive := 0.0
	for _, survivors := range bs.pokemonSurvivors {
		alive += float64(survivors) / float64(bs.battleCount)
	}
	return alive / float64(len(bs.pokemonSurvivors))
}

func (bs *battleStatistics) allPartySurvived() bool {
	if bs == nil || bs.battleCount == 0 || len(bs.pokemonSurvivors) == 0 {
		return false
	}
	for _, survivors := range bs.pokemonSurvivors {
		if survivors != bs.battleCount {
			return false
		}
	}
	return true
}

func (bs *battleStatistics) outcomeScore() float64 {
	if bs == nil || bs.battleCount == 0 {
		return 0
	}

	winRate := float64(bs.winCount) / float64(bs.battleCount)
	survivalRatio := bs.partySurvivalRatio()
	lossRate := 1.0 - winRate
	missingMembers := 0
	aliveMembers := 0
	for _, survivors := range bs.pokemonSurvivors {
		if survivors == bs.battleCount {
			aliveMembers++
		}
		if survivors != bs.battleCount {
			missingMembers++
		}
	}
	aliveRatio := float64(aliveMembers) / float64(len(bs.pokemonSurvivors))

	if bs.allPartySurvived() {
		return 75000.0 + 20000.0*winRate + 15000.0*survivalRatio
	}

	return (2.0*winRate-1.0)*20000.0 + aliveRatio*35000.0 + survivalRatio*12000.0 - lossRate*22000.0 - float64(missingMembers)*14000.0
}

func (bs *battleStatistics) print(party []*pokemon) {
	if bs.battleCount == 0 {
		return
	}

	winRate := float64(bs.winCount) * 100.0 / float64(bs.battleCount)
	log.Printf("player win rate: %.2f%% (%d/%d)", winRate, bs.winCount, bs.battleCount)
	for index, mon := range party {
		survivalRate := float64(bs.pokemonSurvivors[index]) * 100.0 / float64(bs.battleCount)
		log.Printf("%s survival rate: %.2f%% (%d/%d)", mon.base.Name, survivalRate, bs.pokemonSurvivors[index], bs.battleCount)
	}
}
