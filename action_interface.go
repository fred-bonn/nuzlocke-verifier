package main

type action interface {
	invoke(bs battleState)
	prio() int
	speed() int
}
