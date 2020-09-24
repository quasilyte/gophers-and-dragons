package gamedata

import (
	"github.com/quasilyte/gophers-and-dragons/game"
)

func CardProto(typ game.CardType) game.Card {
	card := game.Card{Type: typ}

	switch typ {
	case game.CardAttack:
		card.MP = 0
		card.IsMagic = false
		card.IsOffensive = true

	case game.CardMagicArrow:
		card.MP = 1
		card.IsMagic = true
		card.IsOffensive = true

	case game.CardRetreat:
		card.MP = 0
		card.IsMagic = false
		card.IsOffensive = false

	case game.CardRest:
		card.MP = 2
		card.IsMagic = false
		card.IsOffensive = false

	default:
		panic("unknown card type")
	}

	return card
}
