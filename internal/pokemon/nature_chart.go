package pokemon

import "fmt"

var natureChart = map[string][]string{
	"hardy":   {"attack", "attack"},
	"lonely":  {"attack", "defense"},
	"adamant": {"attack", "special-attack"},
	"naughty": {"attack", "special-defense"},
	"brave":   {"attack", "speed"},
	"bold":    {"defense", "attack"},
	"docile":  {"defense", "defense"},
	"impish":  {"defense", "special-attack"},
	"lax":     {"defense", "special-defense"},
	"relaxed": {"defense", "speed"},
	"modest":  {"special-attack", "attack"},
	"mild":    {"special-attack", "defense"},
	"bashful": {"special-attack", "special-attack"},
	"rash":    {"special-attack", "special-defense"},
	"quiet":   {"special-attack", "speed"},
	"calm":    {"special-defense", "attack"},
	"gentle":  {"special-defense", "defense"},
	"careful": {"special-defense", "speed"},
	"quirky":  {"special-defense", "special-defense"},
	"sassy":   {"special-defense", "speed"},
	"timid":   {"speed", "attack"},
	"hasty":   {"speed", "defense"},
	"jolly":   {"speed", "special-attack"},
	"naive":   {"speed", "special-defense"},
	"serious": {"speed", "speed"},
}

func getNature(nature string) ([]string, error) {
	res, ok := natureChart[nature]
	if !ok {
		return nil, fmt.Errorf("invalid nature: %s", nature)
	}

	return res, nil

}
