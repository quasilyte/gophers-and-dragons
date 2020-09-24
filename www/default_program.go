package tactic

import "github.com/quasilyte/gophers-and-dragons/game"

func ChooseCard(s *game.State) game.CardType {
	return tactic1(s)
}

// tactic1 is a trivial tactic that always retreats.
// You can't win like this, but it's a minimal working
// example that manages to pass the entire game without dying.
func tactic1(s *game.State) game.CardType {
	return game.CardRetreat
}

// tactic2 will only fight with the easiest kind of
// monsters and will route if wounded.
func tactic2(s *game.State) game.CardType {
	if s.Avatar.HP < 10 {
		return game.CardRetreat
	}
	if s.Creep.Type == game.CreepCheepy {
		return game.CardAttack
	}
	return game.CardRetreat
}
