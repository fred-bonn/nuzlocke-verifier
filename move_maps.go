package main

var pivotMoves = map[string]struct{}{
	"u-turn":           {},
	"volt-switch":      {},
	"flip-turn":        {},
	"parting-shot":     {},
	"teleport":         {},
	"chilly-reception": {},
	"baton-pass":       {},
	"shed-tail":        {},
}

var speedControlMoves = map[string]struct{}{
	"electroweb": {},
	"icy-wind":   {},
	"low-sweep":  {},
	"mud-shot":   {},
	"rock-tomb":  {},
	"bulldoze":   {},
	"glaciate":   {},
}

var offenseControlMoves = map[string]string{
	"mystical-fire":  "special",
	"skitter-smack":  "special",
	"breaking-swipe": "physical",
	"snarl":          "special",
	"struggle-bug":   "special",
	"trop-kick":      "special",
	"chilling-water": "physical",
	"lunge":          "physical",
}

var protectMoves = map[string]struct{}{
	"protect": {},
	"detect":  {},
}

var multipleTurnMoves = map[string]struct{}{
	"bounce": {},
}

var paralysisMoves = map[string]struct{}{
	"thunder-wave": {},
	"glare":        {},
	"stun-spore":   {},
	"nuzzle":       {},
}
