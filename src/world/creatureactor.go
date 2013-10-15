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
	// This is why I add an underscore.  You don't want to use fields that end
	// with an underscore.
	Ca_ map[CreatureId]ActorId
	Ac_ map[ActorId]CreatureId
}

func MakeCreatureActor() CreatureActor {
	const CAPACITY = 0
	return CreatureActor{
		Ca_: make(map[CreatureId]ActorId, CAPACITY),
		Ac_: make(map[ActorId]CreatureId, CAPACITY),
	}
}

func (self CreatureActor) IsSane() error {
	if len(self.Ca_) != len(self.Ac_) {
		return CA_INSANE_LENGTH
	}
	for Ca__c, Ca__l := range self.Ca_ {
		Ac__c, ok := self.Ac_[Ca__l]
		if !ok {
			return CA_INSANE_BIJECTION
		}
		if Ac__c != Ca__c {
			return CA_INSANE_BIJECTION
		}
	}
	for Ac__l, Ac__c := range self.Ac_ {
		Ca__l, ok := self.Ca_[Ac__c]
		if !ok {
			return CA_INSANE_BIJECTION
		}
		if Ca__l != Ac__l {
			return CA_INSANE_BIJECTION
		}
	}
	return nil
}

func (self CreatureActor) GetCreature(actor_id ActorId) (CreatureId, bool) {
	creature, ok := self.Ac_[actor_id]
	return creature, ok
}

func (self CreatureActor) GetActor(creature_id CreatureId) (ActorId, bool) {
	actor_id, ok := self.Ca_[creature_id]
	return actor_id, ok
}

func (self CreatureActor) Copy() CreatureActor {
	n_items := len(self.Ca_)
	result := CreatureActor{
		Ca_: make(map[CreatureId]ActorId, n_items),
		Ac_: make(map[ActorId]CreatureId, n_items),
	}
	for c, a := range self.Ca_ {
		result.Ca_[c] = a
		result.Ac_[a] = c
	}
	return result
}

func (self CreatureActor) Add(creature_id CreatureId, actor_id ActorId) (CreatureActor, error) {
	// First make sure that the creature or actor aren't already taken.
	_, ok := self.Ca_[creature_id]
	if ok {
		return self, CA_CREATURE_ALREADY_IN
	}
	_, ok = self.Ac_[actor_id]
	if ok {
		return self, CA_ACTOR_ALREADY_IN
	}
	result := self.Copy()
	result.Ca_[creature_id] = actor_id
	result.Ac_[actor_id] = creature_id
	return result, nil
}

func (self CreatureActor) RemoveCreature(creature_id CreatureId) (CreatureActor, bool) {
	result := self.Copy()
	actor_id, ok := result.Ca_[creature_id]
	if !ok {
		return self, false
	}
	delete(result.Ca_, creature_id)
	delete(result.Ac_, actor_id)
	return result, true
}

func (self CreatureActor) RemoveActor(actor_id ActorId) (CreatureActor, bool) {
	result := self.Copy()
	creature_id, ok := result.Ac_[actor_id]
	if !ok {
		return self, false
	}
	delete(result.Ca_, creature_id)
	delete(result.Ac_, actor_id)
	return result, true
}
