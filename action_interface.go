package main

type action interface {
	invoke(bs battleState)
	prio() int
	speed() int
}

func rollInt(numerator int, denominator int) int {
	if roll(numerator, denominator) {
		return 1
	}
	return 0
}
