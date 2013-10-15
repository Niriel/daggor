package world

import (
	"fmt"
)

// This module deals with the behavior of creatures in the game.
// Because most creatures are controlled by the computer, we can refer to this
// module as AI (artificial intelligence).  However, the character controled by the
// player uses the same concepts.

// Actions are the result of the AI or the player commands.
// "Move north", "Cast fireball" and "Wear helmet" are actions.
// Actions are not authoritative.  If the AI decides on the "Move north" action but
// there is an obstacle in the way, the creature will not manage to walk north.
// Actions may fail.

// Note that the player commands are not actions, and the player keypresses/clicks
// are not commands.  The player press `w`.  It is translated into the "Forward"
// command, which is translated into the "Move east" actions.  Actions take place in
// the game world, while commands exist in the real world.  There are commands to
// save the game, close a window, but there are no actions to do so.  Actions are
// role-play, commands are out-of-character.  As to keypresses, they are not
// commands for the obvious reason that you can remap your keys.

// Actions are executed by actors.  Actors are the subject of the action.  If the
// player wears a helmet, the player is the subject.  If a trap fires an arrow, the
// trap is the subject.  Note that a trap is not a creature, but it is an actor.
// Every actor has a unique identifier.
type ActorId uint64

type Actor struct{}

func MakeActor() Actor {
	return Actor{}
}

type Actors struct {
	Next_id ActorId
	Content map[ActorId]Actor
}

func MakeActors() Actors {
	return Actors{
		Content: make(map[ActorId]Actor),
	}
}

func (self Actors) Copy() Actors {
	result := Actors{
		Next_id: self.Next_id,
		Content: make(map[ActorId]Actor),
	}
	for key, value := range self.Content {
		result.Content[key] = value
	}
	return result
}

func (self Actors) Spawn() (Actors, ActorId, Actor) {
	actor_id := self.Next_id
	actor := MakeActor()
	actors := self.Set(actor_id, actor)
	actors.Next_id += 1
	return actors, actor_id, actor
}

func (self Actors) Add(actor Actor) (Actors, ActorId) {
	actor_id := self.Next_id
	actors := self.Set(actor_id, actor)
	actors.Next_id += 1
	return actors, actor_id
}

func (actors Actors) Set(actor_id ActorId, actor Actor) Actors {
	new_actors := actors.Copy()
	new_actors.Content[actor_id] = actor
	return new_actors
}

// Actions modify the world.  They can potentially modify it entirely, I am not
// limiting them here.  Most actions will have very little effect on the world
// though.
//
// Creating an action does not execute it.  It must be executed to have an effect.
//
// Note that the Execute method returns a World.  This is probably overkill.  It
// may be more efficient to just return some 'deltas' instead.  These deltas would
// then be combined and apply all at once to the world, instead of creating a new
// world each time.  But this is easier, and I go for correctness before trying
// to be smart.
type Action interface {
	Execute(world World) (World, error)
}

// Wait: That action does nothing.
// How is that different from a nil action?  Not much.  Except that nil could
// indicate AI failure while Wait is a deliberate choice, for example.  Must think.
type ActionWait struct{}

// The Wait action does nothing at all, it does not even increment a time variable.
func (self ActionWait) Execute(world World) (World, error) {
	return world, nil
}

// Move: That action moves one actor to a neighboring tile.
type ActionMoveAbsolute struct {
	Subject_id ActorId
	Direction  AbsoluteDirection
	Steps      int
}

func (action ActionMoveAbsolute) Execute(world World) (World, error) {
	// Trivial case: no movement.
	if action.Steps <= 0 {
		return world, nil
	}
	// Only creatures can move.
	creature_id, ok := world.Level.Creature_actor.GetCreature(action.Subject_id)
	if !ok {
		return world, fmt.Errorf(
			"Actor %v does not have a corresponding creature.",
			action.Subject_id,
		)
	}
	// We start computing the new location from the current one.
	new_loc, ok := world.Level.Creature_location.GetLocation(creature_id)
	if !ok {
		return world, fmt.Errorf(
			"Actor %v creature %v does not have a corresponding position.",
			action.Subject_id,
			creature_id,
		)
	}
	for step_id := 0; step_id < action.Steps; step_id++ {
		if world.Level.IsPassable(new_loc, action.Direction) {
			new_loc = new_loc.MoveAbsolute(action.Direction, 1)
		} else {
			return world, fmt.Errorf(
				"Actor %v creature %v cannot pass %v.",
				action.Subject_id,
				creature_id,
				new_loc,
			)
		}
	}
	// Move the creature.
	locations, err := world.Level.Creature_location.Move(creature_id, new_loc)
	if err != nil {
		return world, err
	}
	// World was passed by value, we can modify it.
	world.Level.Creature_location = locations
	return world, nil
}

// Move: That action moves one actor to a neighboring tile.
type ActionMoveRelative struct {
	Subject_id ActorId
	Direction  RelativeDirection
	Steps      int
}

func (action ActionMoveRelative) Execute(world World) (World, error) {
	if action.Steps <= 0 {
		return world, nil
	}
	creature_id, ok := world.Level.Creature_actor.GetCreature(action.Subject_id)
	if !ok {
		return world, fmt.Errorf(
			"Actor %v does not have a corresponding creature.",
			action.Subject_id,
		)
	}
	creature, ok := world.Level.Creatures.Get(creature_id)
	if !ok {
		return world, fmt.Errorf(
			"Actor %v creature %v does not have a corresponding creature.",
			action.Subject_id,
			creature_id,
		)
	}
	direction := creature.F.Add(action.Direction)
	new_loc, ok := world.Level.Creature_location.GetLocation(creature_id)
	if !ok {
		return world, fmt.Errorf(
			"Actor %v creature %v does not have a corresponding position.",
			action.Subject_id,
			creature_id,
		)
	}
	for step_id := 0; step_id < action.Steps; step_id++ {
		if world.Level.IsPassable(new_loc, direction) {
			new_loc = new_loc.MoveAbsolute(direction, 1)
		} else {
			return world, fmt.Errorf(
				"Actor %v creature %v cannot pass %v.",
				action.Subject_id,
				creature_id,
				new_loc,
			)
		}
	}
	// Move the creature.
	locations, err := world.Level.Creature_location.Move(creature_id, new_loc)
	if err != nil {
		return world, err
	}
	// World was passed by value, we can modify it.
	world.Level.Creature_location = locations
	return world, nil
}

// Turn: That action rotates an actor.
// Should it be an action?  I mean, should it take one turn?  That's left to the
// user to choose.
type ActionTurn struct {
	Subject_id ActorId
	Direction  RelativeDirection
	Steps      int
}

func (action ActionTurn) Execute(world World) (World, error) {
	if action.Steps <= 0 {
		return world, nil
	}
	creature_id, ok := world.Level.Creature_actor.GetCreature(action.Subject_id)
	if !ok {
		return world, fmt.Errorf(
			"Actor %v does not have a corresponding creature.",
			action.Subject_id,
		)
	}
	creature, ok := world.Level.Creatures.Get(creature_id)
	if !ok {
		return world, fmt.Errorf(
			"Actor %v creature %v does not have a corresponding creature.",
			action.Subject_id,
			creature_id,
		)
	}
	// Payload.
	facing := creature.F
	for step_id := 0; step_id < action.Steps; step_id++ {
		facing = facing.Add(action.Direction)
	}
	creature.F = facing
	// /Payload.
	world.Level.Creatures = world.Level.Creatures.Set(creature_id, creature)
	return world, nil
}

func DecideAction(subject_id ActorId) Action {
	return ActionTurn{
		Subject_id: subject_id,
		Direction:  LEFT(),
		Steps:      1,
	}
}

//

type ActorTime struct {
	Time            uint64
	Actor_id        ActorId
	Stability_index uint64 // To ensure stable sorting.
}

type ActorSchedule struct {
	Actor_times          []ActorTime
	Next_stability_index uint64 // To ensure stable sorting.
}

func MakeActorSchedule() ActorSchedule {
	return ActorSchedule{
		Actor_times: make([]ActorTime, 0, 8),
	}
}

// Implement sort.Interface.
func (self ActorSchedule) Len() int {
	return len(self.Actor_times)
}
func (self ActorSchedule) Less(i, j int) bool {
	// Order by time has priority.
	if self.Actor_times[i].Time < self.Actor_times[j].Time {
		return true
	}
	// But in case of tie (two actors have the same time), the first actor
	// scheduled wins.
	return self.Actor_times[i].Stability_index < self.Actor_times[j].Stability_index
}
func (self ActorSchedule) Swap(i, j int) {
	self.Actor_times[i], self.Actor_times[j] = self.Actor_times[i], self.Actor_times[j]
}

func (self ActorSchedule) Copy() ActorSchedule {
	// `self` is already a copy.  We just need to copy the content of the slice.
	new_slice := make([]ActorTime, len(self.Actor_times))
	copy(new_slice, self.Actor_times)
	self.Actor_times = new_slice
	return self
}

// Find an actor that should be acting at the provided time.
// Brute force search that does ot assume that the actors are sorted by time.
func (self ActorSchedule) Next(time uint64) (ActorTime, bool) {
	for _, actor_time := range self.Actor_times {
		if actor_time.Time <= time {
			return actor_time, true
		}
	}
	return ActorTime{}, false
}

func (self ActorSchedule) PosActorId(actor_id ActorId) int {
	for index, actor_time := range self.Actor_times {
		if actor_time.Actor_id == actor_id {
			return index
		}
	}
	return -1
}

func (self ActorSchedule) Pos(actor_time ActorTime) int {
	for index, actor_time0 := range self.Actor_times {
		if actor_time0 == actor_time {
			return index
		}
	}
	return -1
}

func (self ActorSchedule) Remove(actor_time ActorTime) (ActorSchedule, bool) {
	index := self.Pos(actor_time)
	if index == -1 {
		return self, false
	}
	// A cheap remove consists in swapping the last and the current, then reduce
	// the length.  I would avoid that for two reasons.
	// 1: my slice is a priori shared, so I need an expensive  deep copy anyway.
	// 2: I would like to keep the order intact in order to speed up search-by
	//    time.

	// I call Copy in order to deep-copy the Actor_times slice.
	// If I do not perform a deep copy of the slice, then the call to `append`
	// later will affect the content of the original slice, introducing side
	// effects.
	self = self.Copy()
	self.Actor_times = append(
		self.Actor_times[:index],
		self.Actor_times[index+1:]...)
	return self, true
}

func (self ActorSchedule) Add(actor_id ActorId, time uint64) ActorSchedule {
	new_entry := ActorTime{
		Actor_id:        actor_id,
		Time:            time,
		Stability_index: self.Next_stability_index,
	}
	len_slice := len(self.Actor_times)
	result := ActorSchedule{
		Next_stability_index: self.Next_stability_index + 1,
		Actor_times:          make([]ActorTime, len_slice, len_slice+1),
	}
	copy(result.Actor_times, self.Actor_times)
	result.Actor_times = append(result.Actor_times, new_entry)
	return result
}
