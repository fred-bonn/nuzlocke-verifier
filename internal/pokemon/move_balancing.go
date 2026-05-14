package pokemon

import (
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

type MoveBalance struct {
	Power    *int
	Accuracy *int
	PP       *int
	Type     *string
}

func (mb MoveBalance) Apply(m *pokeapi.BaseMove) {
	if mb.Power != nil {
		m.Power = *mb.Power
	}
	if mb.Accuracy != nil {
		m.Accuracy = *mb.Accuracy
	}
	if mb.PP != nil {
		m.PP = *mb.PP
	}
	if mb.Type != nil {
		m.Type = *mb.Type
	}
}

var MoveBalanceMap = map[string]*MoveBalance{
	"absorb": {
		Power: new(40),
	},
	"aeroblast": {
		Accuracy: new(100),
	},
	"air-cutter": {
		Accuracy: new(100),
	},
	"air-slash": {
		Accuracy: new(100),
	},
	"aqua-tail": {
		Accuracy: new(95),
	},
	"astonish": {
		Power: new(40),
	},
	"baby-doll-eyes": {
		PP: new(10),
	},
	"barrage": {
		Accuracy: new(100),
	},
	"belch": {
		Accuracy: new(100),
	},
	"bind": {
		Accuracy: new(100),
	},
	"blaze-kick": {
		Accuracy: new(100),
	},
	"blizzard": {
		Accuracy: new(80),
	},
	"blue-flare": {
		Accuracy: new(90),
	},
	"bolt-strike": {
		Accuracy: new(90),
	},
	"bone-club": {
		Accuracy: new(100),
	},
	"bonemerang": {
		Accuracy: new(100),
	},
	"bounce": {
		Accuracy: new(95),
	},
	"captivate": {
		PP: new(5),
	},
	"charge-beam": {
		Power:    new(40),
		Accuracy: new(100),
	},
	"charm": {
		PP: new(5),
	},
	"circle-throw": {
		Accuracy: new(95),
	},
	"clamp": {
		Accuracy: new(100),
	},
	"comet-punch": {
		Accuracy: new(90),
	},
	"confide": {
		PP: new(10),
	},
	"covet": {
		Type: new("fairy"),
	},
	"crabhammer": {
		Accuracy: new(100),
	},
	"cross-chop": {
		Accuracy: new(90),
	},
	"cut": {
		Accuracy: new(100),
	},
	"dark-void": {
		Accuracy: new(80),
	},
	"diamond-storm": {
		Accuracy: new(100),
	},
	"double-hit": {
		Accuracy: new(100),
	},
	"double-slap": {
		Accuracy: new(100),
	},
	"draco-meteor": {
		Accuracy: new(100),
	},
	"dragon-rush": {
		Accuracy: new(85),
	},
	"dragon-tail": {
		Accuracy: new(95),
	},
	"drill-run": {
		Accuracy: new(100),
	},
	"dual-chop": {
		Accuracy: new(100),
	},
	"dual-wingbeat": {
		Accuracy: new(100),
	},
	"eerie-impulse": {
		PP: new(5),
	},
	"electroweb": {
		Accuracy: new(100),
	},
	"fake-out": {
		PP: new(5),
	},
	"fake-tears": {
		PP: new(5),
	},
	"feather-dance": {
		PP: new(5),
	},
	"fire-fang": {
		Accuracy: new(100),
	},
	"fire-spin": {
		Accuracy: new(100),
	},
	"flash": {
		Accuracy: new(70),
	},
	"fleur-cannon": {
		Accuracy: new(100),
	},
	"fly": {
		Accuracy: new(100),
	},
	"flying-press": {
		Accuracy: new(100),
	},
	"focus-blast": {
		Accuracy: new(80),
	},
	"freeze-shock": {
		Accuracy: new(100),
	},
	"frenzy-plant": {
		Accuracy: new(100),
	},
	"frost-breath": {
		Accuracy: new(100),
	},
	"frustration": {
		Power: new(102),
	},
	"fury-attack": {
		Accuracy: new(100),
	},
	"fury-swipes": {
		Accuracy: new(90),
	},
	"gear-grind": {
		Accuracy: new(100),
	},
	"giga-impact": {
		Accuracy: new(100),
	},
	"glaciate": {
		Accuracy: new(100),
	},
	"grass-whistle": {
		Accuracy: new(70),
	},
	"growl": {
		PP: new(10),
	},
	"gunk-shot": {
		Accuracy: new(85),
	},
	"hammer-arm": {
		Accuracy: new(100),
	},
	"harden": {
		PP: new(5),
	},
	"head-smash": {
		Accuracy: new(85),
	},
	"heat-wave": {
		Accuracy: new(100),
	},
	"high-horsepower": {
		Accuracy: new(100),
	},
	"hurricane": {
		Accuracy: new(80),
	},
	"hydro-cannon": {
		Accuracy: new(100),
	},
	"hydro-pump": {
		Accuracy: new(85),
	},
	"hyper-beam": {
		Accuracy: new(100),
	},
	"hyper-fang": {
		Accuracy: new(100),
	},
	"hypnosis": {
		Accuracy: new(70),
	},
	"ice-burn": {
		Accuracy: new(100),
	},
	"ice-fang": {
		Accuracy: new(100),
	},
	"ice-hammer": {
		Accuracy: new(100),
	},
	"icicle-crash": {
		Accuracy: new(100),
	},
	"icy-wind": {
		Accuracy: new(100),
	},
	"iron-tail": {
		Accuracy: new(85),
	},
	"kinesis": {
		Accuracy: new(100),
	},
	"leaf-storm": {
		Accuracy: new(100),
	},
	"leaf-tornado": {
		Accuracy: new(100),
	},
	"leech-seed": {
		Accuracy: new(100),
	},
	"leer": {
		PP: new(10),
	},
	"lick": {
		Power: new(40),
	},
	"light-of-ruin": {
		Accuracy: new(100),
	},
	"lovely-kiss": {
		Accuracy: new(80),
	},
	"magma-storm": {
		Accuracy: new(90),
	},
	"mega-drain": {
		Power: new(60),
	},
	"mega-kick": {
		Accuracy: new(85),
	},
	"mega-punch": {
		Accuracy: new(100),
	},
	"megahorn": {
		Accuracy: new(90),
	},
	"metal-claw": {
		Accuracy: new(100),
	},
	"metal-sound": {
		PP:       new(5),
		Accuracy: new(100),
	},
	"meteor-beam": {
		Accuracy: new(100),
	},
	"meteor-mash": {
		Accuracy: new(100),
	},
	"mirror-shot": {
		Accuracy: new(100),
	},
	"misty-explosion": {
		Power: new(200),
	},
	"mud-bomb": {
		Accuracy: new(100),
	},
	"muddy-water": {
		Accuracy: new(95),
	},
	"natures-madness": {
		Accuracy: new(100),
	},
	"night-daze": {
		Accuracy: new(100),
	},
	"noble-roar": {
		PP: new(10),
	},
	"octazooka": {
		Power:    new(80),
		Accuracy: new(100),
	},
	"origin-pulse": {
		Accuracy: new(100),
	},
	"overheat": {
		Accuracy: new(100),
	},
	"pin-missile": {
		Accuracy: new(100),
	},
	"play-nice": {
		PP: new(10),
	},
	"play-rough": {
		Accuracy: new(100),
	},
	"poison-powder": {
		Accuracy: new(90),
	},
	"power-whip": {
		Accuracy: new(90),
	},
	"precipice-blades": {
		Accuracy: new(100),
	},
	"psycho-boost": {
		Accuracy: new(100),
	},
	"razor-leaf": {
		Accuracy: new(100),
	},
	"razor-shell": {
		Accuracy: new(100),
	},
	"return": {
		Power: new(102),
	},
	"roar-of-time": {
		Accuracy: new(100),
	},
	"rock-blast": {
		Accuracy: new(100),
	},
	"rock-climb": {
		Accuracy: new(95),
	},
	"rock-slide": {
		Accuracy: new(100),
	},
	"rock-throw": {
		Accuracy: new(100),
	},
	"rock-wrecker": {
		Accuracy: new(100),
	},
	"rolling-kick": {
		Accuracy: new(100),
	},
	"roost": {
		PP: new(5),
	},
	"sacred-fire": {
		Accuracy: new(100),
	},
	"sand-tomb": {
		Accuracy: new(100),
	},
	"sand-attack": {
		PP: new(5),
	},
	"scale-shot": {
		Accuracy: new(100),
	},
	"screech": {
		PP:       new(5),
		Accuracy: new(100),
	},
	"seed-flare": {
		Accuracy: new(90),
	},
	"sing": {
		Accuracy: new(70),
	},
	"skitter-smack": {
		Accuracy: new(100),
	},
	"sky-attack": {
		Accuracy: new(100),
	},
	"sky-uppercut": {
		Accuracy: new(100),
	},
	"slam": {
		Accuracy: new(90),
	},
	"sleep-powder": {
		Accuracy: new(80),
	},
	"smog": {
		Accuracy: new(90),
	},
	"snarl": {
		PP:       new(10),
		Accuracy: new(100),
	},
	"sonic-boom": {
		Accuracy: new(100),
	},
	"spacial-rend": {
		Accuracy: new(100),
	},
	"steam-eruption": {
		Accuracy: new(100),
	},
	"steel-beam": {
		Accuracy: new(100),
	},
	"steel-wing": {
		Accuracy: new(100),
	},
	"stone-edge": {
		Accuracy: new(85),
	},
	"strange-steam": {
		Accuracy: new(100),
	},
	"struggle-bug": {
		PP: new(10),
	},
	"stun-spore": {
		Accuracy: new(90),
	},
	"submission": {
		Accuracy: new(100),
	},
	"super-fang": {
		Accuracy: new(100),
		Type:     new("dark"),
	},
	"supersonic": {
		Accuracy: new(70),
	},
	"swagger": {
		Accuracy: new(90),
	},
	"sweet-kiss": {
		Accuracy: new(80),
	},
	"tail-slap": {
		Accuracy: new(100),
	},
	"take-down": {
		Accuracy: new(100),
	},
	"tearful-look": {
		PP: new(10),
	},
	"thunder": {
		Accuracy: new(80),
	},
	"thunder-cage": {
		Accuracy: new(100),
	},
	"thunder-fang": {
		Accuracy: new(100),
	},
	"tickle": {
		PP: new(10),
	},
	"whirlpool": {
		Accuracy: new(100),
	},
	"wrap": {
		Accuracy: new(100),
	},
	"zen-headbutt": {
		Accuracy: new(100),
	},
}
