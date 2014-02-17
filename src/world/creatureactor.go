package world

// Each creature has at most one actor.
// Each actor has at most one creature.

type CreatureActorError int

const (
	CA_INSANE_LENGTH = CreatureActorError(iota)
	CA_INSANE_BIJECTION
	CA_CREATURE_ALREADY_IN
	CA_ACTOR_ALREADY_IN
)

var creature_actor_error_text = map[CreatureActorError]string{
	CA_INSANE_LENGTH:       "insane because of length mismatch",
	CA_INSANE_BIJECTION:    "insane because of bijection mismatch",
	CA_CREATURE_ALREADY_IN: "creature already registered",
	CA_ACTOR_ALREADY_IN:    "actor already registered",
}

func (self CreatureActorError) Error() string {
	return creature_actor_error_text[self]
}

type CreatureActor struct {
	// Upper case for debugging and for Gob, not so that you use it.
	Ca map[CreatureId]ActorID
	Ac map[ActorID]CreatureId
}

func MakeCreatureActor() CreatureActor {
	const CAPACITY = 0
	return CreatureActor{
		Ca: make(map[CreatureId]ActorID, CAPACITY),
		Ac: make(map[ActorID]CreatureId, CAPACITY),
	}
}

func (self CreatureActor) IsSane() error {
	if len(self.Ca) != len(self.Ac) {
		return CA_INSANE_LENGTH
	}
	for Ca_c, Ca_l := range self.Ca {
		Ac_c, ok := self.Ac[Ca_l]
		if !ok {
			return CA_INSANE_BIJECTION
		}
		if Ac_c != Ca_c {
			return CA_INSANE_BIJECTION
		}
	}
	for Ac_l, Ac_c := range self.Ac {
		Ca_l, ok := self.Ca[Ac_c]
		if !ok {
			return CA_INSANE_BIJECTION
		}
		if Ca_l != Ac_l {
			return CA_INSANE_BIJECTION
		}
	}
	return nil
}

func (self CreatureActor) GetCreature(actor_id ActorID) (CreatureId, bool) {
	creature, ok := self.Ac[actor_id]
	return creature, ok
}

func (self CreatureActor) GetActor(creature_id CreatureId) (ActorID, bool) {
	actor_id, ok := self.Ca[creature_id]
	return actor_id, ok
}

func (self CreatureActor) Copy() CreatureActor {
	n_items := len(self.Ca)
	result := CreatureActor{
		Ca: make(map[CreatureId]ActorID, n_items),
		Ac: make(map[ActorID]CreatureId, n_items),
	}
	for c, a := range self.Ca {
		result.Ca[c] = a
		result.Ac[a] = c
	}
	return result
}

func (self CreatureActor) Add(creature_id CreatureId, actor_id ActorID) (CreatureActor, error) {
	// First make sure that the creature or actor aren't already taken.
	_, ok := self.Ca[creature_id]
	if ok {
		return self, CA_CREATURE_ALREADY_IN
	}
	_, ok = self.Ac[actor_id]
	if ok {
		return self, CA_ACTOR_ALREADY_IN
	}
	result := self.Copy()
	result.Ca[creature_id] = actor_id
	result.Ac[actor_id] = creature_id
	return result, nil
}

func (self CreatureActor) RemoveCreature(creature_id CreatureId) (CreatureActor, bool) {
	result := self.Copy()
	actor_id, ok := result.Ca[creature_id]
	if !ok {
		return self, false
	}
	delete(result.Ca, creature_id)
	delete(result.Ac, actor_id)
	return result, true
}

func (self CreatureActor) RemoveActor(actor_id ActorID) (CreatureActor, bool) {
	result := self.Copy()
	creature_id, ok := result.Ac[actor_id]
	if !ok {
		return self, false
	}
	delete(result.Ca, creature_id)
	delete(result.Ac, actor_id)
	return result, true
}
