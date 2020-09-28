package game

// State is a current game state as percieved by the current turn.
type State struct {
	// Turn is a current turn number.
	// Note that turn number starts with one, not zero.
	Turn int

	// Round is a number of the current encounter.
	// Note that round number starts with one, not zero.
	Round int

	// RoundTurn is a round-local turn number.
	// If it's 1, then it's a first turn in the round.
	// Note that round turn number starts with one, not zero.
	RoundTurn int

	// Score is your current game score.
	Score int

	// Avatar contains information about your hero status.
	Avatar Avatar

	// Creep is an information about your current opponent.
	Creep Creep

	// NextCreep is a type of the next creep.
	// Next creep is encountered after the current creep is defeated.
	// If there is no next creep, a special type CreepNone indicates that.
	NextCreep CreepType

	// Deck is your cards collection.
	// It's keyed by a card type, like CardAttack.
	Deck map[CardType]Card
}

// Can reports whether it's legal to do a cardType move.
func (st *State) Can(cardType CardType) bool {
	if st.Deck[cardType].Count == 0 {
		return false // Card is unavailable
	}
	if st.Avatar.MP < st.Deck[cardType].MP {
		return false // Not enougn mana
	}
	return true
}

// Creep is a particular creep information.
type Creep struct {
	Type CreepType

	HP int

	// Stun is a number of turns this creep is going to skip.
	// You probably want to use Creep.IsStunned() instead of this.
	Stun int

	CreepStats
}

// IsFull reports whether creep health is full.
func (c *Creep) IsFull() bool { return c.HP == c.MaxHP }

// IsStunned reports whether creep is currently stunned.
func (c *Creep) IsStunned() bool { return c.Stun > 0 }

// CreepStats is a set of creep statistics.
type CreepStats struct {
	MaxHP       int
	Damage      IntRange
	ScoreReward int
	CardsReward int
	Traits      CreepTraitList
}

// Avatar is a hero status information.
type Avatar struct {
	HP int
	MP int
	AvatarStats
}

// AvatarStats is a set of avatar statistics.
type AvatarStats struct {
	MaxHP int
	MaxMP int
}

// Card is a hero deck card information.
type Card struct {
	// Type is a card type, like "CardAttack" or "CardMagicArrow".
	Type CardType

	// Count tells how many such cards you have.
	// -1 means "unlimited".
	Count int

	CardStats
}

// CardStats is a set of card statistics.
type CardStats struct {
	// MP is a card mana cost per usage.
	MP int

	// IsMagic tells whether this card effect is considered to be magical.
	IsMagic bool

	// Effect is a description-like string that explains the Power field meaning.
	Effect string

	// Power is a spell effectiveness.
	// For offensive spells, it's the damage they deal.
	// For other spells it can mean different things (see Effect field).
	Power IntRange

	// IsOffensive tells whether this card targets enemy.
	// If it's not, it either targets you or has some special effect like "Retreat".
	IsOffensive bool
}

// IntRange is an inclusive integer range from Low() to High().
type IntRange [2]int

func (rng IntRange) Low() int     { return rng[0] }
func (rng IntRange) High() int    { return rng[1] }
func (rng IntRange) IsZero() bool { return rng.Low() == 0 && rng.High() == 0 }

// CardType is an enum-like type for cards.
type CardType int

// All card types.
//go:generate stringer -type=CardType -trimprefix=Card
const (
	// Infinite cards.

	CardAttack CardType = iota
	CardMagicArrow
	CardRetreat
	CardRest

	// Cards that need to be obtained during the gameplay.

	CardPowerAttack
	CardFirebolt
	CardStun
	CardHeal
	CardParry
)

// CreepType is an enum-like type for creeps.
type CreepType int

// All creep types.
//go:generate stringer -type=CreepType -trimprefix=Creep
const (
	CreepNone CreepType = iota
	CreepCheepy
	CreepImp
	CreepLion
	CreepFairy
	CreepMummy
	CreepDragon
)

// CreepTraitList is convenience wrapper over a slice of creep traits.
type CreepTraitList []CreepTrait

// Has reports whether a creep trait list contains the specified trait.
func (list CreepTraitList) Has(x CreepTrait) bool {
	for _, trait := range list {
		if trait == x {
			return true
		}
	}
	return false
}

// CreepTrait is an enum-like type for creep special traits.
type CreepTrait int

// All creep traits.
//go:generate stringer -type=CreepTrait -trimprefix=Trait
const (
	TraitCoward CreepTrait = iota
	TraitMagicImmunity
	TraitWeakToFire
	TraitSlow
	TraitRanged
)
