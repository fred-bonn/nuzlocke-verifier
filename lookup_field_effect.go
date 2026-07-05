package main

type fieldEffect int

const (
	noneEffect fieldEffect = iota
	spikesEffect
	stealthRockEffect
	stickyWebEffect
	toxicSpikesEffect
	reflectEffect
	lightScreenEffect
	auroraVeilEffect
	tailwindEffect
	safeguardEffect
	luckyChantEffect
	gravityEffect
	trickRoomEffect
	magicRoomEffect
	wonderRoomEffect
)

var fieldEffectMap = map[string]fieldEffect{
	"spikes":       spikesEffect,
	"stealth rock": stealthRockEffect,
	"sticky web":   stickyWebEffect,
	"toxic spikes": toxicSpikesEffect,
	"reflect":      reflectEffect,
	"light screen": lightScreenEffect,
	"aurora veil":  auroraVeilEffect,
	"tailwind":     tailwindEffect,
	"safeguard":    safeguardEffect,
	"lucky chant":  luckyChantEffect,
	"gravity":      gravityEffect,
	"trick room":   trickRoomEffect,
	"magic room":   magicRoomEffect,
	"wonder room":  wonderRoomEffect,
}

func stringToFieldEffect(s string) fieldEffect {
	if e, ok := fieldEffectMap[s]; ok {
		return e
	}
	return noneEffect
}
