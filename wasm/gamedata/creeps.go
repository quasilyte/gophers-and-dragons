package gamedata

import (
	"github.com/quasilyte/gophers-and-dragons/game"
)

func CreepProto(typ game.CreepType) game.CreepStatus {
	creep := game.CreepStatus{Type: typ}

	switch typ {
	case game.CreepCheepy:
		creep.HP = 4

	case game.CreepLion:
		creep.HP = 7

	case game.CreepFairy:
		creep.HP = 6

	case game.CreepMummy:
		creep.HP = 20

	case game.CreepDragon:
		creep.HP = 25

	default:
		panic("unknown creep type")
	}

	return creep
}
