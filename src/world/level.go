package world

type Level struct {
	Floors           Buildings
	Ceilings         Buildings
	Walls            [4]Buildings // Sorted by facing.
	Columns          Buildings
	Dynamic          Dynamic
	Actors           Actors
	Creatures        Creatures
	CreatureLocation CreatureLocation
	CreatureActor    CreatureActor
	ActorSchedule    ActorSchedule
}

func MakeLevel() Level {
	return Level{
		Floors:           MakeBuildings(),
		Ceilings:         MakeBuildings(),
		Columns:          MakeBuildings(),
		Actors:           MakeActors(),
		Creatures:        MakeCreatures(),
		CreatureLocation: MakeCreatureLocation(),
		CreatureActor:    MakeCreatureActor(),
	}
}

func (level *Level) IsPassable(location Location, direction AbsoluteDirection) bool {
	wall_passable := true   // By default, no wall is good.
	floor_passable := false // By default, no floor is bad.
	// Western walls face East.
	wall_facing := direction.Add(BACK())
	wall_index := wall_facing.Value()

	building, ok := level.Walls[wall_index].Get(location.X, location.Y)
	if ok {
		wall_passable = building.(Wall).IsPassable()
	}

	new_loc := location.MoveAbsolute(direction, 1)
	building, ok = level.Floors.Get(new_loc.X, new_loc.Y)
	if ok {
		floor_passable = building.(Floor).IsPassable()
	}

	return wall_passable && floor_passable
}

func (self Level) ActorLocation(actor_id ActorID) (Location, bool) {
	creature_id, ok := self.CreatureActor.GetCreature(actor_id)
	if !ok {
		return Location{}, false
	}
	location, ok := self.CreatureLocation.GetLocation(creature_id)
	return location, ok
}

func (self Level) ActorPosition(actor_id ActorID) (Position, bool) {
	creature_id, ok := self.CreatureActor.GetCreature(actor_id)
	if !ok {
		return Position{}, false
	}
	location, ok := self.CreatureLocation.GetLocation(creature_id)
	if !ok {
		return Position{}, false
	}
	creature, ok := self.Creatures.Get(creature_id)
	if !ok {
		return Position{}, false
	}
	return location.ToPosition(creature.F), true
}

func (self Level) SetActorSchedule(actor_schedule ActorSchedule) Level {
	self.ActorSchedule = actor_schedule
	return self
}
