package main

type fieldEffect int

const (
	noneEffect fieldEffect = iota
	spikesEffect
	stealthRockEffect
	stickyWebEffect
	toxicSpikesEffect
)

var fieldEffectMap = map[string]fieldEffect{
	"spikes":       spikesEffect,
	"stealth rock": stealthRockEffect,
	"sticky web":   stickyWebEffect,
	"toxic spikes": toxicSpikesEffect,
}

func stringToFieldEffect(s string) fieldEffect {
	if e, ok := fieldEffectMap[s]; ok {
		return e
	}
	elogFatalf("error: %s is not a valid field effect", s)
	return noneEffect
}
