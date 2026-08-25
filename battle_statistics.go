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
