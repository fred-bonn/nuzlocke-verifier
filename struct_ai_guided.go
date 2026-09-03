package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type guidedAi struct {
	input  *bufio.Reader
	output io.Writer
	err    error
}

type guidedActionChooser interface {
	chooseAction(bs battleState, slot *slot, party []*pokemon, actions []*moveAction) action
}

func newGuidedAi(input io.Reader, output io.Writer) *guidedAi {
	if input == nil {
		input = os.Stdin
	}
	return &guidedAi{input: bufio.NewReader(input), output: output}
}

func (ga *guidedAi) evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int) {
	ga.printBattleState(bs, actions[0].userSlot)
	ga.print("Choose an action:\n")
	for i, action := range actions {
		ga.print("%d. Use %s against %s\n", i+1, action.move.Name, action.targetSlot.mon.base.Name)
	}

	choice := ga.readChoice(len(actions))
	if choice < 0 {
		return nil, 0
	}
	return actions[choice], 1
}

func (ga *guidedAi) chooseAction(bs battleState, slot *slot, party []*pokemon, actions []*moveAction) action {
	ga.printBattleState(bs, slot)

	type menuAction struct {
		move *moveAction
		mon  *pokemon
	}
	menu := make([]menuAction, 0, len(actions)+len(party))
	for _, action := range actions {
		menu = append(menu, menuAction{move: action})
	}
	if canReplace(party) && !slot.isTrapped() {
		for _, mon := range party {
			if mon != slot.mon && !mon.fainted && !bs.getActions().containstSwitchTo(mon) {
				menu = append(menu, menuAction{mon: mon})
			}
		}
	}

	ga.print("Choose an action:\n")
	for i, option := range menu {
		if option.move != nil {
			ga.print("%d. Use %s against %s\n", i+1, option.move.move.Name, option.move.targetSlot.mon.base.Name)
		} else {
			ga.print("%d. Switch to %s (%d/%d HP)\n", i+1, option.mon.base.Name, option.mon.hp, option.mon.maxHP())
		}
	}

	choice := ga.readChoice(len(menu))
	if choice < 0 {
		return &dummyAction{}
	}
	if menu[choice].move != nil {
		return menu[choice].move
	}
	return &switchAction{oldSlot: slot, new: menu[choice].mon}
}

func (ga *guidedAi) evaluteSwitchIns(bs battleState, mons []*pokemon, opponentSlot *slot) *pokemon {
	ga.printBattleState(bs, bs.getOpponentSlot(opponentSlot))
	ga.print("Choose a Pokemon to switch in:\n")
	for i, mon := range mons {
		ga.print("%d. %s (%d/%d HP)\n", i+1, mon.base.Name, mon.hp, mon.maxHP())
	}

	choice := ga.readChoice(len(mons))
	if choice < 0 {
		return nil
	}
	return mons[choice]
}

func (ga *guidedAi) printBattleState(bs battleState, slot *slot) {
	opponent := bs.getOpponentSlot(slot)
	ga.print("Battle state: %s %d/%d - %s %d/%d\n", slot.mon.base.Name, slot.mon.hp, slot.mon.maxHP(), opponent.mon.base.Name, opponent.mon.hp, opponent.mon.maxHP())
}

func (ga *guidedAi) shouldSwitch(bs battleState, slot *slot, score int, party []*pokemon) bool {
	return false
}

func (ga *guidedAi) print(format string, args ...any) {
	output := ga.output
	if output == nil {
		output = os.Stdout
	}
	fmt.Fprintf(output, format, args...)
}

func (ga *guidedAi) readChoice(numberOfChoices int) int {
	if ga.input == nil {
		ga.input = bufio.NewReader(os.Stdin)
	}

	for {
		ga.print("> ")
		line, err := ga.input.ReadString('\n')
		if err != nil && len(line) == 0 {
			ga.err = fmt.Errorf("guided AI input ended before a valid choice was entered")
			return -1
		}

		choice, parseErr := strconv.Atoi(strings.TrimSpace(line))
		if parseErr == nil && choice >= 1 && choice <= numberOfChoices {
			return choice - 1
		}
		ga.print("Please enter a number from 1 to %d.\n", numberOfChoices)
		if err != nil {
			return -1
		}
	}
}
