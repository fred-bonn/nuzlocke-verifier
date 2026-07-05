package main

var sleepMoves = map[string]struct{}{
	"dark void":     {},
	"grass whistle": {},
	"hypnosis":      {},
	"lovely-kiss":   {},
	"sing":          {},
	"sleep-powder":  {},
	"spore":         {},
	"yawn":          {},
}

var pivotMoves = map[string]struct{}{
	"u turn":           {},
	"volt switch":      {},
	"flip turn":        {},
	"parting shot":     {},
	"teleport":         {},
	"chilly reception": {},
	"baton pass":       {},
	"shed tail":        {},
}

var speedControlMoves = map[string]struct{}{
	"electroweb": {},
	"icy wind":   {},
	"low sweep":  {},
	"mud shot":   {},
	"rock tomb":  {},
	"bulldoze":   {},
	"glaciate":   {},
}

var offenseControlMoves = map[string]moveClass{
	"mystical fire":  specialClass,
	"skitter smack":  specialClass,
	"breaking swipe": physicalClass,
	"snarl":          specialClass,
	"struggle bug":   specialClass,
	"trop kick":      specialClass,
	"chilling water": physicalClass,
	"lunge":          physicalClass,
}

var protectMoves = map[string]struct{}{
	"protect":      {},
	"detect":       {},
	"kings shield": {},
}

var multipleTurnMoves = map[string]struct{}{
	"bounce": {},
}

var paralysisMoves = map[string]struct{}{
	"thunder wave": {},
	"glare":        {},
	"stun spore":   {},
	"nuzzle":       {},
}

var powderMoves = map[string]struct{}{
	"powder":        {},
	"spore":         {},
	"sleep powder":  {},
	"stun spore":    {},
	"poison powder": {},
	"rage powder":   {},
	"cotten spore":  {},
}

var selfThawingMoves = map[string]struct{}{
	"flame wheel":     {},
	"sacred fire":     {},
	"flare blitz":     {},
	"fusion flare":    {},
	"scald":           {},
	"steam eruption":  {},
	"burn up":         {},
	"pyro ball":       {},
	"scorching-sands": {},
}
