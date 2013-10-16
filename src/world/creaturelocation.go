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
	Cl map[CreatureId]Location
	Lc map[Location]CreatureId
}

func MakeCreatureLocation() CreatureLocation {
	const CAPACITY = 0
	return CreatureLocation{
		Cl: make(map[CreatureId]Location, CAPACITY),
		Lc: make(map[Location]CreatureId, CAPACITY),
	}
}

func (self CreatureLocation) IsSane() error {
	if len(self.Cl) != len(self.Lc) {
		return CL_INSANE_LENGTH
	}
	for Cl_c, Cl_l := range self.Cl {
		Lc_c, ok := self.Lc[Cl_l]
		if !ok {
			return CL_INSANE_BIJECTION
		}
		if Lc_c != Cl_c {
			return CL_INSANE_BIJECTION
		}
	}
	for Lc_l, Lc_c := range self.Lc {
		Cl_l, ok := self.Cl[Lc_c]
		if !ok {
			return CL_INSANE_BIJECTION
		}
		if Cl_l != Lc_l {
			return CL_INSANE_BIJECTION
		}
	}
	return nil
}

func (self CreatureLocation) GetCreature(loc Location) (CreatureId, bool) {
	creature, ok := self.Lc[loc]
	return creature, ok
}

func (self CreatureLocation) GetLocation(creature_id CreatureId) (Location, bool) {
	location, ok := self.Cl[creature_id]
	return location, ok
}

func (self CreatureLocation) Copy() CreatureLocation {
	n_items := len(self.Cl)
	result := CreatureLocation{
		Cl: make(map[CreatureId]Location, n_items),
		Lc: make(map[Location]CreatureId, n_items),
	}
	for c, l := range self.Cl {
		result.Cl[c] = l
		result.Lc[l] = c
	}
	return result
}

func (self CreatureLocation) Add(creature_id CreatureId, location Location) (CreatureLocation, error) {
	// First make sure that the creature or location aren't already taken.
	_, ok := self.Cl[creature_id]
	if ok {
		return self, CL_CREATURE_ALREADY_IN
	}
	_, ok = self.Lc[location]
	if ok {
		return self, CL_LOCATION_ALREADY_IN
	}
	result := self.Copy()
	result.Cl[creature_id] = location
	result.Lc[location] = creature_id
	return result, nil
}

func (self CreatureLocation) RemoveCreature(creature_id CreatureId) (CreatureLocation, bool) {
	result := self.Copy()
	location, ok := result.Cl[creature_id]
	if !ok {
		return self, false
	}
	delete(result.Cl, creature_id)
	delete(result.Lc, location)
	return result, true
}

func (self CreatureLocation) RemoveLocation(location Location) (CreatureLocation, bool) {
	result := self.Copy()
	creature_id, ok := result.Lc[location]
	if !ok {
		return self, false
	}
	delete(result.Cl, creature_id)
	delete(result.Lc, location)
	return result, true
}

func (self CreatureLocation) Move(creature_id CreatureId, location Location) (CreatureLocation, error) {
	crt_location, ok := self.Cl[creature_id]
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
	_, ok = self.Lc[location]
	if ok {
		// There already is a creature there, cannot move.
		return self, CL_OCCUPIED
	}
	result := self.Copy()
	result.Cl[creature_id] = location
	result.Lc[location] = creature_id
	delete(result.Lc, crt_location)
	return result, nil
}
