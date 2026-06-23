package main

type stats int

const (
	HitPoints stats = iota
	Attack
	Defense
	SpecialAttack
	SpecialDefense
	Speed
	Accuracy
	Evasion
)

func (s stats) String() string {
	switch s {
	case HitPoints:
		return "hp"
	case Attack:
		return "attack"
	case Defense:
		return "defense"
	case SpecialAttack:
		return "special attack"
	case SpecialDefense:
		return "special defense"
	case Speed:
		return "speed"
	case Accuracy:
		return "accuracy"
	case Evasion:
		return "evasion"
	default:
		return ""
	}
}

func stringToStat(stat string) stats {
	switch stat {
	case "hp":
		return HitPoints
	case "attack":
		return Attack
	case "defense":
		return Defense
	case "special-attack":
		return SpecialAttack
	case "special-defense":
		return SpecialDefense
	case "speed":
		return Speed
	case "accuracy":
		return Accuracy
	case "evasion":
		return Evasion
	default:
		elogFatalf("%s is not a valid stat", stat)
		return 0
	}
}

var natureChart = map[string][]stats{
	"hardy":   {Attack, Attack},
	"lonely":  {Attack, Defense},
	"adamant": {Attack, SpecialAttack},
	"naughty": {Attack, SpecialDefense},
	"brave":   {Attack, Speed},
	"bold":    {Defense, Attack},
	"docile":  {Defense, Defense},
	"impish":  {Defense, SpecialAttack},
	"lax":     {Defense, SpecialDefense},
	"relaxed": {Defense, Speed},
	"modest":  {SpecialAttack, Attack},
	"mild":    {SpecialAttack, Defense},
	"bashful": {SpecialAttack, SpecialAttack},
	"rash":    {SpecialAttack, SpecialDefense},
	"quiet":   {SpecialAttack, Speed},
	"calm":    {SpecialDefense, Attack},
	"gentle":  {SpecialDefense, Defense},
	"careful": {SpecialDefense, Speed},
	"quirky":  {SpecialDefense, SpecialDefense},
	"sassy":   {SpecialDefense, Speed},
	"timid":   {Speed, Attack},
	"hasty":   {Speed, Defense},
	"jolly":   {Speed, SpecialAttack},
	"naive":   {Speed, SpecialDefense},
	"serious": {Speed, Speed},
}
