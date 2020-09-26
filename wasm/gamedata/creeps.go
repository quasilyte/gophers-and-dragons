package gamedata

import (
	"github.com/quasilyte/gophers-and-dragons/game"
)

var creeps = map[game.CreepType]game.CreepStats{
	game.CreepCheepy: {
		MaxHP:       4,
		Damage:      game.IntRange{1, 4},
		ScoreReward: 3,
		CardsReward: 1,
		Traits: []game.CreepTrait{
			game.TraitCoward,
		},
	},

	game.CreepImp: {
		MaxHP:       5,
		Damage:      game.IntRange{3, 4},
		ScoreReward: 5,
		CardsReward: 1,
	},

	game.CreepLion: {
		MaxHP:       10,
		Damage:      game.IntRange{2, 3},
		ScoreReward: 6,
		CardsReward: 2,
	},

	game.CreepFairy: {
		MaxHP:       9,
		Damage:      game.IntRange{4, 5},
		ScoreReward: 11,
		CardsReward: 2,
		Traits: []game.CreepTrait{
			game.TraitRanged,
		},
	},

	game.CreepMummy: {
		MaxHP:       18,
		Damage:      game.IntRange{3, 4},
		ScoreReward: 15,
		CardsReward: 3,
		Traits: []game.CreepTrait{
			game.TraitWeakToFire,
			game.TraitSlow,
		},
	},

	game.CreepDragon: {
		MaxHP:       30,
		Damage:      game.IntRange{5, 6},
		ScoreReward: 35,
		CardsReward: 0,
		Traits: []game.CreepTrait{
			game.TraitMagicImmunity,
		},
	},
}

func GetCreepStats(typ game.CreepType) game.CreepStats {
	return creeps[typ]
}
