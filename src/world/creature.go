package world

// Each creature in the world is identified with a unique ID.
// A creature does not know its own ID.  This is to avoid inconsistencies.
// Indeed, creatures are contained in a map[CreatureId]Creature.
type CreatureId uint64

type Creature struct {
	F AbsoluteDirection
}

func MakeCreature() Creature {
	return Creature{F: EAST()}
}

// Implement Facer interface.
func (self Creature) Facing() AbsoluteDirection {
	return self.F
}

type Creatures struct {
	Next_id CreatureId
	Content map[CreatureId]Creature
}

func MakeCreatures() Creatures {
	return Creatures{Content: make(map[CreatureId]Creature)}
}

func (self Creatures) Copy() Creatures {
	result := Creatures{
		Next_id: self.Next_id,
		Content: make(map[CreatureId]Creature),
	}
	for key, value := range self.Content {
		result.Content[key] = value
	}
	return result
}

func (self Creatures) Spawn() (Creatures, CreatureId, Creature) {
	creature_id := self.Next_id
	creature := MakeCreature()
	creatures := self.Set(creature_id, creature)
	creatures.Next_id += 1
	return creatures, creature_id, creature
}

func (self Creatures) Add(creature Creature) (Creatures, CreatureId) {
	creature_id := self.Next_id
	creatures := self.Set(creature_id, creature)
	creatures.Next_id += 1
	return creatures, creature_id
}

func (self Creatures) Get(creature_id CreatureId) (Creature, bool) {
	creature, ok := self.Content[creature_id]
	return creature, ok
}

func (self Creatures) Set(creature_id CreatureId, creature Creature) Creatures {
	result := self.Copy()
	result.Content[creature_id] = creature
	return result
}
