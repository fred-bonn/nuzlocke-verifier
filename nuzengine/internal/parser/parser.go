package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type parser struct {
	tokens  []token
	current token
	peek    token
	pos     int
}

func (p *parser) generateError() parseError {
	return parseError{
		previous: p.tokens[p.pos-2],
		current:  p.current,
	}
}

func newParser(tokens []token) *parser {
	p := &parser{
		tokens:  tokens,
		peek:    tokens[0],
		current: token{Type: t_START},
	}

	return p
}

func (p *parser) parsePokemon() (ParsedPokemon, error) {
	res := ParsedPokemon{}
	res.IVs = make(map[string]int)

	p.consumeNewlines()

	// parse name @ item line
	err := p.parseNameLine(&res)
	if err != nil {
		return ParsedPokemon{}, err
	}

	// parse level line
	err = p.parseLevelLine(&res)
	if err != nil {
		return ParsedPokemon{}, err
	}

	// parse nature line
	err = p.parseNatureLine(&res)
	if err != nil {
		return ParsedPokemon{}, err
	}

	// parse ability line
	err = p.parseAbilityLine(&res)
	if err != nil {
		return ParsedPokemon{}, err
	}

	// parse status line
	if p.expectedToken(t_STATUS) {
		err = p.parseStatusLine(&res)
		if err != nil {
			return ParsedPokemon{}, err
		}
	}

	// parse predamage line
	if p.expectedToken(t_HP) {
		err = p.parsePredamageLine(&res)
		if err != nil {
			return ParsedPokemon{}, err
		}
	} else {
		res.HP = -1
	}

	// parse IVs line
	if p.expectedToken((t_IVS)) {
		err = p.parseIVsLine(&res)
		if err != nil {
			return ParsedPokemon{}, err
		}
	}

	// parse move lines, max 4
	var moves int
	for p.expectedToken(t_MOVE) && moves < 4 {
		err = p.parseMoveLine(&res)
		if err != nil {
			return ParsedPokemon{}, err
		}
		moves++
	}

	return res, nil
}

func (p *parser) parseNameLine(pokemon *ParsedPokemon) error {
	if !p.expectedToken(t_IDENT) {
		return p.generateError()
	}

	pokemon.Name = p.current.Literal

	p.nextToken()

	if p.expectedToken(t_ITEM) {
		pokemon.Item = p.current.Literal

		p.nextToken()
	}

	p.consumeNewlines()

	return nil
}

func (p *parser) parseLevelLine(pokemon *ParsedPokemon) error {
	if !p.expectedToken(t_LEVEL) {
		return p.generateError()
	}

	p.nextToken()

	if !p.expectedToken(t_IDENT) {
		return p.generateError()
	}

	val, err := strconv.Atoi(p.current.Literal)
	if err != nil {
		return fmt.Errorf("error: %s could not be converted to a number", p.current.Literal)
	}

	pokemon.Level = val

	p.nextToken()

	p.consumeNewlines()

	return nil
}

func (p *parser) parseNatureLine(pokemon *ParsedPokemon) error {
	if !p.expectedToken(t_IDENT) {
		return p.generateError()
	}

	parts := strings.Fields(p.current.Literal)

	if len(parts) < 2 {
		return fmt.Errorf("error: nature line not enough fields: '%s'", p.current.Literal)
	}
	if parts[1] != "Nature" {
		return fmt.Errorf("error: incorrect nature line format: '%s'", p.current.Literal)
	}
	if _, ok := validNature[parts[0]]; !ok {
		return fmt.Errorf("error: unrecognized nature: '%s'", p.current.Literal)
	}

	pokemon.Nature = strings.ToLower(parts[0])

	p.nextToken()

	p.consumeNewlines()

	return nil
}

func (p *parser) parseAbilityLine(pokemon *ParsedPokemon) error {
	if !p.expectedToken(t_ABILITY) {
		return p.generateError()
	}

	p.nextToken()

	if !p.expectedToken(t_IDENT) {
		return p.generateError()
	}

	pokemon.Ability = p.current.Literal

	p.nextToken()

	p.consumeNewlines()

	return nil
}

func (p *parser) parseStatusLine(pokemon *ParsedPokemon) error {
	p.nextToken()

	if !p.expectedToken(t_IDENT) {
		return p.generateError()
	}
	if _, ok := validAilments[p.current.Literal]; !ok {
		return fmt.Errorf("error: unrecognized status: '%s'", p.current.Literal)
	}

	pokemon.Status = strings.ToLower(p.current.Literal)

	p.nextToken()

	p.consumeNewlines()

	return nil
}

func (p *parser) parsePredamageLine(pokemon *ParsedPokemon) error {
	p.nextToken()

	if !p.expectedToken(t_IDENT) {
		return p.generateError()
	}

	val, err := strconv.Atoi(p.current.Literal)
	if err != nil {
		return fmt.Errorf("error: %s could not be converted to a number", p.current.Literal)
	}

	pokemon.HP = val

	p.nextToken()

	p.consumeNewlines()

	return nil
}

func (p *parser) parseIVsLine(pokemon *ParsedPokemon) error {
	p.nextToken()

	for !p.expectedToken(t_NEWLINE) {
		if !p.expectedToken(t_IDENT) {
			return p.generateError()
		}
		parts := strings.Fields(p.current.Literal)

		if len(parts) < 2 {
			return fmt.Errorf("error: iv not enough fields: '%s'", p.current.Literal)
		}
		iv := parts[1]

		if _, ok := validStat[iv]; !ok {
			return fmt.Errorf("error: unrecognized iv: '%s'", iv)
		}

		val, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("error: %s could not be converted to a number", parts[0])
		}

		pokemon.IVs[strings.ToLower(iv)] = val

		p.nextToken()
	}

	p.consumeNewlines()
	return nil
}

func (p *parser) parseMoveLine(pokemon *ParsedPokemon) error {
	p.nextToken()

	if !p.expectedToken(t_IDENT) {
		return p.generateError()
	}

	pokemon.Moves = append(pokemon.Moves, p.current.Literal)

	p.nextToken()

	p.consumeNewlines()

	return nil
}

func (p *parser) consumeNewlines() {
	for p.current.Type == t_NEWLINE {
		p.nextToken()
	}
}

func (p *parser) nextToken() {
	p.current = p.peek

	if p.current.Type == t_EOF {
		return
	}

	p.pos++
	p.peek = p.tokens[p.pos]
}

func (p *parser) expectedToken(t tokenType) bool {
	return p.current.Type == t
}
