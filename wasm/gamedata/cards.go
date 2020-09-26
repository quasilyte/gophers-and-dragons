package gamedata

import (
	"github.com/quasilyte/gophers-and-dragons/game"
)

var Cards = map[game.CardType]game.CardStats{
	game.CardAttack: {
		MP:          0,
		IsMagic:     false,
		IsOffensive: true,
		Power:       game.IntRange{2, 4},
		Effect:      "damage",
	},

	game.CardPowerAttack: {
		MP:          0,
		IsMagic:     false,
		IsOffensive: true,
		Power:       game.IntRange{4, 5},
		Effect:      "damage",
	},

	game.CardStun: {
		MP:          0,
		IsMagic:     false,
		IsOffensive: true,
		Power:       game.IntRange{2, 2},
		Effect:      "turns skipped",
	},

	game.CardMagicArrow: {
		MP:          1,
		IsMagic:     true,
		IsOffensive: true,
		Power:       game.IntRange{3, 3},
		Effect:      "magical damage",
	},

	game.CardFirebolt: {
		MP:          3,
		IsMagic:     true,
		IsOffensive: true,
		Power:       game.IntRange{4, 6},
		Effect:      "magical damage",
	},

	game.CardRetreat: {
		MP:          0,
		IsMagic:     false,
		IsOffensive: false,
	},

	game.CardRest: {
		MP:          2,
		IsMagic:     false,
		IsOffensive: false,
		Power:       game.IntRange{3, 3},
		Effect:      "HP recovered",
	},

	game.CardHeal: {
		MP:          4,
		IsMagic:     true,
		IsOffensive: false,
		Power:       game.IntRange{10, 15},
		Effect:      "HP recovered",
	},

	game.CardParry: {
		MP:          0,
		IsMagic:     false,
		IsOffensive: false,
	},
}

func GetCardStats(typ game.CardType) game.CardStats {
	return Cards[typ]
}
