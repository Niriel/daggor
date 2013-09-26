package world

// Each creature in the world is identified with a unique ID.
// A creature does not know its own ID.  This is to avoid inconsistencies.
// Indeed, creatures are contained in a map[CreatureUid]Creature.
type CreatureUid uint64

type Creature interface {
	Pos() Position
	Appearance() ModelId
}

type Creatures map[CreatureUid]Creature

func (src Creatures) Copy() Creatures {
	dst := make(Creatures, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}
