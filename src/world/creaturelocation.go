package world

// Each creature has at most one location.
// Each location has at most one creature.

type CreatureLocationError int

const (
	CL_NOOP = CreatureLocationError(iota)
	CL_INSANE_LENGTH
	CL_INSANE_BIJECTION
	CL_CREATURE_ALREADY_IN
	CL_LOCATION_ALREADY_IN
	CL_NOT_FOUND
	CL_OCCUPIED
)

var creature_location_error_text = map[CreatureLocationError]string{
	CL_NOOP:                "no operation",
	CL_INSANE_LENGTH:       "insane because of length mismatch",
	CL_INSANE_BIJECTION:    "insane because of bijection mismatch",
	CL_CREATURE_ALREADY_IN: "creature already registered",
	CL_LOCATION_ALREADY_IN: "location already registered",
	CL_NOT_FOUND:           "not found",
	CL_OCCUPIED:            "destination occupied",
}

func (self CreatureLocationError) Error() string {
	return creature_location_error_text[self]
}

type CreatureLocation struct {
	// Upper case for debugging and for Gob, not so that you use it.
	// This is why I add an underscore.  You don't want to use fields that end
	// with an underscore.
	Cl_ map[CreatureId]Location
	Lc_ map[Location]CreatureId
}

func MakeCreatureLocation() CreatureLocation {
	const CAPACITY = 0
	return CreatureLocation{
		Cl_: make(map[CreatureId]Location, CAPACITY),
		Lc_: make(map[Location]CreatureId, CAPACITY),
	}
}

func (self CreatureLocation) IsSane() error {
	if len(self.Cl_) != len(self.Lc_) {
		return CL_INSANE_LENGTH
	}
	for Cl__c, Cl__l := range self.Cl_ {
		Lc__c, ok := self.Lc_[Cl__l]
		if !ok {
			return CL_INSANE_BIJECTION
		}
		if Lc__c != Cl__c {
			return CL_INSANE_BIJECTION
		}
	}
	for Lc__l, Lc__c := range self.Lc_ {
		Cl__l, ok := self.Cl_[Lc__c]
		if !ok {
			return CL_INSANE_BIJECTION
		}
		if Cl__l != Lc__l {
			return CL_INSANE_BIJECTION
		}
	}
	return nil
}

func (self CreatureLocation) GetCreature(loc Location) (CreatureId, bool) {
	creature, ok := self.Lc_[loc]
	return creature, ok
}

func (self CreatureLocation) GetLocation(creature_id CreatureId) (Location, bool) {
	location, ok := self.Cl_[creature_id]
	return location, ok
}

func (self CreatureLocation) Copy() CreatureLocation {
	n_items := len(self.Cl_)
	result := CreatureLocation{
		Cl_: make(map[CreatureId]Location, n_items),
		Lc_: make(map[Location]CreatureId, n_items),
	}
	for c, l := range self.Cl_ {
		result.Cl_[c] = l
		result.Lc_[l] = c
	}
	return result
}

func (self CreatureLocation) Add(creature_id CreatureId, location Location) (CreatureLocation, error) {
	// First make sure that the creature or location aren't already taken.
	_, ok := self.Cl_[creature_id]
	if ok {
		return self, CL_CREATURE_ALREADY_IN
	}
	_, ok = self.Lc_[location]
	if ok {
		return self, CL_LOCATION_ALREADY_IN
	}
	result := self.Copy()
	result.Cl_[creature_id] = location
	result.Lc_[location] = creature_id
	return result, nil
}

func (self CreatureLocation) RemoveCreature(creature_id CreatureId) (CreatureLocation, bool) {
	result := self.Copy()
	location, ok := result.Cl_[creature_id]
	if !ok {
		return self, false
	}
	delete(result.Cl_, creature_id)
	delete(result.Lc_, location)
	return result, true
}

func (self CreatureLocation) RemoveLocation(location Location) (CreatureLocation, bool) {
	result := self.Copy()
	creature_id, ok := result.Lc_[location]
	if !ok {
		return self, false
	}
	delete(result.Cl_, creature_id)
	delete(result.Lc_, location)
	return result, true
}

func (self CreatureLocation) Move(creature_id CreatureId, location Location) (CreatureLocation, error) {
	crt_location, ok := self.Cl_[creature_id]
	if !ok {
		// We do not know that creature so we cannot move it.  Add it first.
		return self, CL_NOT_FOUND
	}
	if crt_location == location {
		// No-op.
		// It returns false because you should NOT be spending time doing
		// no-ops.  There is a good chance that this situation results from a
		// logical error from the programmer.
		return self, CL_NOOP
	}
	_, ok = self.Lc_[location]
	if ok {
		// There already is a creature there, cannot move.
		return self, CL_OCCUPIED
	}
	result := self.Copy()
	result.Cl_[creature_id] = location
	result.Lc_[location] = creature_id
	delete(result.Lc_, crt_location)
	return result, nil
}
