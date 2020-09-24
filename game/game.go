package game

// State is a current game state as percieved by the current turn.
type State struct {
	// Turn is a current turn number.
	// Note that turn number starts with one, not zero.
	Turn int

	// Round is a number of the current encounter.
	// Note that round number starts with one, not zero.
	Round int

	// Score is your current game score.
	Score int

	// Avatar contains information about your hero status.
	Avatar AvatarStatus

	// Creep is an information about your current opponent status.
	Creep CreepStatus

	// NextCreep is a name of the next creep.
	// Next creep is encountered after the current creep is defeated.
	// If there is no next creep, a special string "none" indicates that.
	NextCreep string

	// Deck is your cards collection.
	// It's keyed by a card type, like CardAttack.
	Deck map[CardType]Card
}

// CreepStatus is a creep status information.
type CreepStatus struct {
	HP int
}

// AvatarStatus is a hero status information.
type AvatarStatus struct {
	HP int
	MP int
}

// Card is a hero deck card information.
type Card struct {
	// Type is a card type, like "CardAttack" or "CardMagicArrow".
	Type CardType

	// Count tells how many such cards you have.
	Count int

	// MP is a card mana cost per usage.
	MP int

	// IsMagic tells whether this card effect is considered to be magical.
	IsMagic bool

	// IsOffensive tells whether this card targets enemy.
	// If it's not, it either targets you or has some special effect like "Retreat".
	IsOffensive bool
}

// CardType is an enum-like type for cards.
type CardType int

// All card types.
const (
	CardAttack CardType = iota
	CardMagicArrow
	CardRetreat
	CardRest
)
