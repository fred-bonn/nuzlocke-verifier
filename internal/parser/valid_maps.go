package parser

var validNature = map[string]struct{}{
	"Hardy":   {},
	"Lonely":  {},
	"Adamant": {},
	"Naughty": {},
	"Brave":   {},
	"Bold":    {},
	"Docile":  {},
	"Impish":  {},
	"Lax":     {},
	"Relaxed": {},
	"Modest":  {},
	"Mild":    {},
	"Bashful": {},
	"Rash":    {},
	"Quiet":   {},
	"Calm":    {},
	"Gentle":  {},
	"Careful": {},
	"Quirky":  {},
	"Sassy":   {},
	"Timid":   {},
	"Hasty":   {},
	"Jolly":   {},
	"Naive":   {},
	"Serious": {},
}

var validAilments = map[string]struct{}{
	"Paralysis": {},
	"Poison":    {},
	"Toxic":     {},
	"Burn":      {},
	"Freeze":    {},
	"Sleep":     {},
}

var validStat = map[string]struct{}{
	"HP":  {},
	"Atk": {},
	"Def": {},
	"SpA": {},
	"SpD": {},
	"Spe": {},
}
