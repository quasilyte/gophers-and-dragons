package sim

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/quasilyte/gophers-and-dragons/game"
	"github.com/quasilyte/gophers-and-dragons/wasm/gamedata"
	"github.com/quasilyte/gophers-and-dragons/wasm/simstep"
)

type Config struct {
	AvatarHP int
	AvatarMP int
	Rounds   int
}

func Run(config *Config, chooseCard func(game.State) game.CardType) []simstep.Action {
	runner := newRunner(config, chooseCard)
	return runner.Run()
}

func newGameState(config *Config) *game.State {
	avatarStats := game.AvatarStats{
		MaxHP: config.AvatarHP,
		MaxMP: config.AvatarMP,
	}
	return &game.State{
		Round: 1,
		Turn:  1,
		Avatar: game.Avatar{
			HP:          avatarStats.MaxHP,
			MP:          avatarStats.MaxMP,
			AvatarStats: avatarStats,
		},
		Deck: make(map[game.CardType]game.Card),
	}
}

type runner struct {
	state         *game.State
	config        *Config
	out           []simstep.Action
	rand          *rand.Rand
	chooseCard    func(game.State) game.CardType
	peekableCards []game.CardType
	badMoves      int
}

func newRunner(config *Config, chooseCard func(game.State) game.CardType) *runner {
	return &runner{
		config:     config,
		state:      newGameState(config),
		chooseCard: chooseCard,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *runner) Run() []simstep.Action {
	r.initWorld()
	r.out = append(r.out, simstep.NextRound{})
	for {
		if r.badMoves >= 10 {
			r.emitRedLogf("Game over: too many illegal moves!")
			break
		}
		if r.state.RoundTurn >= 50 {
			r.emitRedLogf("Game over: round lasted for too long!")
			break
		}
		if r.state.Round > r.config.Rounds {
			r.victory()
			break
		}
		stop := r.runTurn()
		if stop {
			break
		}
	}
	return r.out
}

func (r *runner) victory() {
	r.out = append(r.out, simstep.Victory{})

	bonus := r.state.Avatar.HP
	r.out = append(r.out, simstep.UpdateScore{Delta: bonus})
	r.emitGreenLogf("Got %d survival bonus points", bonus)
}

func (r *runner) initWorld() {
	r.state.Creep = newCreep(r.peekCreep(1))
	r.state.NextCreep = r.peekCreep(2)
	r.initDeck()
}

func (r *runner) initDeck() {
	deck := r.state.Deck
	for typ, cardStats := range gamedata.Cards {
		card := game.Card{
			Type:      typ,
			CardStats: cardStats,
		}
		switch typ {
		case game.CardAttack, game.CardMagicArrow, game.CardRest, game.CardRetreat:
			card.Count = -1
		default:
			r.peekableCards = append(r.peekableCards, typ)
		}
		deck[typ] = card
	}
}

func (r *runner) peekCard() game.CardType {
	return r.peekableCards[r.rand.Intn(len(r.peekableCards))]
}

func (r *runner) peekCreep(round int) game.CreepType {
	// This handles the next creep for the last round.
	if round > r.config.Rounds {
		return game.CreepNone
	}

	// Dragon is always encountered at the last round.
	if round == r.config.Rounds {
		return game.CreepDragon
	}
	// Cheepy is always encountered at the first round.
	if round == 1 {
		return game.CreepCheepy
	}
	// Imp is always encountered at the second
	if round == 2 {
		return game.CreepImp
	}

	roll := r.rand.Intn(99)

	// First 5 rounds can't have high-tier enemies.
	if r.state.Round <= 5 {
		switch {
		case roll >= 90: // 10%
			return game.CreepFairy
		case roll >= 50: // 40%
			return game.CreepLion
		case roll >= 30: // 20%
			return game.CreepImp
		default: // 30%
			return game.CreepCheepy
		}
	}

	switch {
	case roll >= 70: // 30%
		return game.CreepMummy
	case roll >= 50: // 20%
		return game.CreepFairy
	case roll >= 30: // 20%
		return game.CreepLion
	case roll >= 10: // 20%
		return game.CreepImp
	default: // 10%
		return game.CreepCheepy
	}
}

func (r *runner) runCreepAction(parried bool) {
	creep := &r.state.Creep
	avatar := &r.state.Avatar

	damageRoll := r.rangeRand(creep.Damage)
	if parried {
		if !creep.Traits.Has(game.TraitRanged) {
			creep.HP -= damageRoll
			r.out = append(r.out, simstep.UpdateCreepHP{Delta: -damageRoll})
			r.emitLogf("%d damage is reflected back to %s", damageRoll, creep.Type.String())
			return
		}
		r.emitRedLogf("Failed to parry a ranged attack")
	}

	avatar.HP -= damageRoll
	r.out = append(r.out, simstep.UpdateHP{Delta: -damageRoll})
	r.emitRedLogf("%s deals %d damage", creep.Type.String(), damageRoll)
}

func (r *runner) runAvatarAction(cardType game.CardType, card game.CardStats) {
	creep := &r.state.Creep
	avatar := &r.state.Avatar

	cardCount := r.state.Deck[cardType].Count
	if cardCount == 0 {
		r.emitRedLogf("Tried to use unavailable card %s", cardType.String())
		r.badMoves++
		return
	}
	if cardCount != -1 {
		r.out = append(r.out, simstep.ChangeCardCount{
			Name:  cardType.String(),
			Delta: -1,
		})
		changeDeckCardCount(r.state.Deck, cardType, -1)
	}

	if card.MP != 0 {
		if avatar.MP < card.MP {
			r.emitRedLogf("Not enough mana to use %s", cardType.String())
			r.badMoves++
			return
		}
		avatar.MP -= card.MP
		r.out = append(r.out, simstep.UpdateMP{Delta: -card.MP})
	}

	switch cardType {
	case game.CardAttack, game.CardPowerAttack:
		damageRoll := r.rangeRand(card.Power)
		creep.HP -= damageRoll
		r.out = append(r.out, simstep.UpdateCreepHP{Delta: -damageRoll})
		r.emitLogf("Your %s deals %d damage", cardType.String(), damageRoll)

	case game.CardStun:
		stunRoll := r.rangeRand(card.Power)
		creep.Stun = stunRoll
		r.emitLogf("%s is stunned for %d turns", creep.Type.String(), stunRoll)

	case game.CardMagicArrow:
		if creep.Traits.Has(game.TraitMagicImmunity) {
			r.emitRedLogf("%s failed: is immune to magic", cardType.String())
			break
		}
		damageRoll := r.rangeRand(card.Power)
		creep.HP -= damageRoll
		r.out = append(r.out, simstep.UpdateCreepHP{Delta: -damageRoll})
		r.emitLogf("Your %s deals %d damage", cardType.String(), damageRoll)

	case game.CardFirebolt:
		if creep.Traits.Has(game.TraitMagicImmunity) {
			r.emitRedLogf("%s failed: is immune to magic", cardType.String())
			break
		}
		damageRoll := r.rangeRand(card.Power)
		if creep.Traits.Has(game.TraitWeakToFire) {
			damageRoll *= 2
		}
		creep.HP -= damageRoll
		r.out = append(r.out, simstep.UpdateCreepHP{Delta: -damageRoll})
		r.emitLogf("Your %s deals %d damage", cardType.String(), damageRoll)

	case game.CardRest, game.CardHeal:
		r.avatarHeal(cardType, card)
	}
}

func (r *runner) avatarHeal(cardType game.CardType, card game.CardStats) {
	avatar := &r.state.Avatar

	roll := r.rangeRand(card.Power)
	healed := calculateHealed(roll, avatar.HP, r.config.AvatarHP)
	avatar.HP += healed
	r.out = append(r.out, simstep.UpdateHP{Delta: healed})
	r.emitGreenLogf("Got %d HP from %s", healed, cardType.String())
}

func (r *runner) creepDefeated() {
	creep := &r.state.Creep

	r.state.Score += creep.ScoreReward
	r.emitGreenLogf("%s is defeated! %d score points received",
		creep.Type.String(), creep.ScoreReward)
	r.out = append(r.out, simstep.UpdateScore{Delta: creep.ScoreReward})

	for i := 0; i < creep.CardsReward; i++ {
		rewardCardType := r.peekCard()
		r.emitGreenLogf("Collected %s card", rewardCardType.String())
		r.out = append(r.out, simstep.ChangeCardCount{
			Name:  rewardCardType.String(),
			Delta: 1,
		})
		changeDeckCardCount(r.state.Deck, rewardCardType, 1)
	}

	r.nextRound()
}

func (r *runner) runTurn() bool {
	r.beginTurn()
	defer r.endTurn()

	creep := &r.state.Creep
	avatar := &r.state.Avatar

	cardType := r.chooseCard(*r.state)
	card := gamedata.GetCardStats(cardType)
	r.runAvatarAction(cardType, card)

	if creep.HP <= 0 {
		r.creepDefeated()
		return false
	}

	retreatedBeforeAttacked := cardType == game.CardRetreat &&
		creep.Traits.Has(game.TraitSlow)
	skipsAttack := creep.IsFull() &&
		creep.Traits.Has(game.TraitCoward)
	stunned := creep.IsStunned()
	if !stunned && !retreatedBeforeAttacked && !skipsAttack {
		parried := cardType == game.CardParry
		r.runCreepAction(parried)
		if parried && creep.HP <= 0 {
			r.creepDefeated()
			return false
		}
	}
	if skipsAttack && cardType == game.CardParry {
		r.emitRedLogf("Tried to parry, but the enemy was not attacking")
	}

	if creep.Stun > 0 {
		creep.Stun--
	}

	if avatar.HP <= 0 {
		r.out = append(r.out, simstep.Defeat{})
		r.emitRedLogf("Game over: avatar has been defeated!")
		return true
	}

	if cardType == game.CardRetreat {
		r.emitLogf("Retreated from %s!", creep.Type.String())
		r.nextRound()
	}

	return false
}

func (r *runner) nextRound() {
	r.state.Round++
	r.state.RoundTurn = 0

	r.state.Creep = newCreep(r.state.NextCreep)
	r.state.NextCreep = r.peekCreep(r.state.Round + 1)
	r.out = append(r.out, simstep.SetCreep{
		Name: r.state.Creep.Type.String(),
		HP:   r.state.Creep.HP,
	})
	r.out = append(r.out, simstep.SetNextCreep{
		Name: r.state.NextCreep.String(),
		HP:   gamedata.GetCreepStats(r.state.NextCreep).MaxHP,
	})
	r.out = append(r.out, simstep.NextRound{})
}

func (r *runner) rangeRand(rng game.IntRange) int {
	if rng.IsZero() {
		return 0
	}
	v := r.rand.Intn(rng.High() - rng.Low() + 1)
	return v + rng.Low()
}

func (r *runner) beginTurn() {
	r.emitLogf("--- Turn %d ---", r.state.Turn)
}

func (r *runner) endTurn() {
	r.state.Turn++
	r.state.RoundTurn++
	r.out = append(r.out, simstep.Wait{})
}

func (r *runner) emitLogf(format string, args ...interface{}) {
	r.out = append(r.out, simstep.Log{Message: fmt.Sprintf(format, args...)})
}

func (r *runner) emitRedLogf(format string, args ...interface{}) {
	r.out = append(r.out, simstep.RedLog{Message: fmt.Sprintf(format, args...)})
}

func (r *runner) emitGreenLogf(format string, args ...interface{}) {
	r.out = append(r.out, simstep.GreenLog{Message: fmt.Sprintf(format, args...)})
}
