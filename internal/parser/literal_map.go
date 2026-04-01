package parser

var literalMap = map[string]Token{
	"Level": Token{
		Type:    LEVEL,
		Literal: "Level",
	},
	"Ability": Token{
		Type:    ABILITY,
		Literal: "Ability",
	},
	"IVs": Token{
		Type:    IVS,
		Literal: "IVs",
	},
	"Status": Token{
		Type:    STATUS,
		Literal: "Status",
	},
	"HP": Token{
		Type:    HP,
		Literal: "HP",
	},
}
