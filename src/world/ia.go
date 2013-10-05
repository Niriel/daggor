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

type Actor struct {
	Pos   Position
	Brain Brain
}

type Actors map[ActorId]Actor

func (src Actors) Copy() Actors {
	dst := make(Actors, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func (actors Actors) Set(actor_id ActorId, actor Actor) Actors {
	new_actors := actors.Copy()
	new_actors[actor_id] = actor
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
	actors := world.Level.Actors
	subject, ok := actors[action.Subject_id]
	if !ok {
		return world, fmt.Errorf("Actor %v not found.", action.Subject_id)
	}
	new_loc := subject.Pos.Location
	for step_id := 0; step_id < action.Steps; step_id++ {
		if world.Level.IsPassable(new_loc, action.Direction) {
			new_loc = new_loc.MoveAbsolute(action.Direction, 1)
		} else {
			return world, fmt.Errorf("Unpassable.")
		}
	}
	// Subject is a value, not a pointer, we are free to modify it.
	subject.Pos = subject.Pos.SetLocation(new_loc) // Preserve facing.
	// However, actors is a map and therefore contains a pointer, we need a new
	// map.
	actors = actors.Set(action.Subject_id, subject)
	// World is a value, we can modify it.
	world.Level.Actors = actors
	return world, nil
}

// Move: That action moves one actor to a neighboring tile.
type ActionMoveRelative struct {
	Subject_id ActorId
	Direction  RelativeDirection
	Steps      int
}

func (action ActionMoveRelative) Execute(world World) (World, error) {
	actors := world.Level.Actors
	subject, ok := actors[action.Subject_id]
	if !ok {
		return world, fmt.Errorf("Actor %v not found.", action.Subject_id)
	}
	new_loc := subject.Pos.Location
	direction := subject.Pos.F.Add(action.Direction)
	for step_id := 0; step_id < action.Steps; step_id++ {
		if world.Level.IsPassable(new_loc, direction) {
			new_loc = new_loc.MoveAbsolute(direction, 1)
		} else {
			return world, fmt.Errorf("Unpassable.")
		}
	}
	// Subject is a value, not a pointer, we are free to modify it.
	subject.Pos = subject.Pos.SetLocation(new_loc) // Preserve facing.
	// However, actors is a map and therefore contains a pointer, we need a new
	// map.
	actors = actors.Set(action.Subject_id, subject)
	// World is a value, we can modify it.
	world.Level.Actors = actors
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
	actors := world.Level.Actors
	subject, ok := actors[action.Subject_id]
	if !ok {
		return world, fmt.Errorf("Actor %v not found.", action.Subject_id)
	}
	subject.Pos = subject.Pos.Turn(action.Direction, action.Steps)
	actors = actors.Set(action.Subject_id, subject)
	world.Level.Actors = actors
	return world, nil
}

// The AI should not fire an action every frame.  Unless that is actually what we
// want.  There are very neat ways of doing smart things, like Finite State
// Machines, or better, Behavior Trees.  But before I do that, I need a dirt simple
// cooldown mechanism.

type Brain struct {
	Cooldown_timer uint64
}

func (brain Brain) WarmUp(dt uint64) Brain {
	if brain.Cooldown_timer != 0 {
		panic("Brain is already warm.")
	}
	brain.Cooldown_timer = dt
	return brain
}

func (brain Brain) CoolDown(dt uint64) Brain {
	if dt > brain.Cooldown_timer {
		brain.Cooldown_timer = 0
	} else {
		brain.Cooldown_timer -= dt
	}
	return brain
}

func (brain Brain) IsCold() bool {
	return brain.Cooldown_timer == 0
}

func DecideAction(subject_id ActorId) Action {
	return ActionTurn{
		Subject_id: subject_id,
		Direction:  LEFT(),
		Steps:      1,
	}
}
