package main

type itemEvent interface {
	getPayload() any
}

type resistBerryEvent struct {
	typeName    string
	denominator *int
}

func (e *resistBerryEvent) getPayload() any {
	return []any{e.typeName, e.denominator}
}
