package ia

import (
	"fmt"
	"world"
)

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
	Execute(world.World) (world.World, error)
}

// This module deals with the behavior of creatures in the game.
// Because most creatures are controlled by the computer, we can refer to this
// module as AI (artificial intelligence).  However, the character controled by the
// player uses the same concepts.

// Wait: That action does nothing.
// How is that different from a nil action?  Not much.  Except that nil could
// indicate AI failure while Wait is a deliberate choice, for example.  Must think.
type ActionWait struct{}

// The Wait action does nothing at all, it does not even increment a time variable.
func (self ActionWait) Execute(w world.World) (world.World, error) {
	return w, nil
}

// Move: That action moves one actor to a neighboring tile.
type ActionMoveAbsolute struct {
	Subject_id world.ActorId
	Direction  world.AbsoluteDirection
	Steps      int
}

func (action ActionMoveAbsolute) Execute(w world.World) (world.World, error) {
	// Trivial case: no movement.
	if action.Steps <= 0 {
		return w, nil
	}
	// Only creatures can move.
	creature_id, ok := w.Level.Creature_actor.GetCreature(action.Subject_id)
	if !ok {
		return w, fmt.Errorf(
			"Actor %v does not have a corresponding creature.",
			action.Subject_id,
		)
	}
	// We start computing the new location from the current one.
	new_loc, ok := w.Level.Creature_location.GetLocation(creature_id)
	if !ok {
		return w, fmt.Errorf(
			"Actor %v creature %v does not have a corresponding position.",
			action.Subject_id,
			creature_id,
		)
	}
	for step_id := 0; step_id < action.Steps; step_id++ {
		if w.Level.IsPassable(new_loc, action.Direction) {
			new_loc = new_loc.MoveAbsolute(action.Direction, 1)
		} else {
			return w, fmt.Errorf(
				"Actor %v creature %v cannot pass %v.",
				action.Subject_id,
				creature_id,
				new_loc,
			)
		}
	}
	// Move the creature.
	locations, err := w.Level.Creature_location.Move(creature_id, new_loc)
	if err != nil {
		return w, err
	}
	// World was passed by value, we can modify it.
	w.Level.Creature_location = locations
	return w, nil
}

// Move: That action moves one actor to a neighboring tile.
type ActionMoveRelative struct {
	Subject_id world.ActorId
	Direction  world.RelativeDirection
	Steps      int
}

func (action ActionMoveRelative) Execute(w world.World) (world.World, error) {
	if action.Steps <= 0 {
		return w, nil
	}
	creature_id, ok := w.Level.Creature_actor.GetCreature(action.Subject_id)
	if !ok {
		return w, fmt.Errorf(
			"Actor %v does not have a corresponding creature.",
			action.Subject_id,
		)
	}
	creature, ok := w.Level.Creatures.Get(creature_id)
	if !ok {
		return w, fmt.Errorf(
			"Actor %v creature %v does not have a corresponding creature.",
			action.Subject_id,
			creature_id,
		)
	}
	direction := creature.F.Add(action.Direction)
	new_loc, ok := w.Level.Creature_location.GetLocation(creature_id)
	if !ok {
		return w, fmt.Errorf(
			"Actor %v creature %v does not have a corresponding position.",
			action.Subject_id,
			creature_id,
		)
	}
	for step_id := 0; step_id < action.Steps; step_id++ {
		if w.Level.IsPassable(new_loc, direction) {
			new_loc = new_loc.MoveAbsolute(direction, 1)
		} else {
			return w, fmt.Errorf(
				"Actor %v creature %v cannot pass %v.",
				action.Subject_id,
				creature_id,
				new_loc,
			)
		}
	}
	// Move the creature.
	locations, err := w.Level.Creature_location.Move(creature_id, new_loc)
	if err != nil {
		return w, err
	}
	// World was passed by value, we can modify it.
	w.Level.Creature_location = locations
	return w, nil
}

// Turn: That action rotates an actor.
// Should it be an action?  I mean, should it take one turn?  That's left to the
// user to choose.
type ActionTurn struct {
	Subject_id world.ActorId
	Direction  world.RelativeDirection
	Steps      int
}

func (action ActionTurn) Execute(w world.World) (world.World, error) {
	if action.Steps <= 0 {
		return w, nil
	}
	creature_id, ok := w.Level.Creature_actor.GetCreature(action.Subject_id)
	if !ok {
		return w, fmt.Errorf(
			"Actor %v does not have a corresponding creature.",
			action.Subject_id,
		)
	}
	creature, ok := w.Level.Creatures.Get(creature_id)
	if !ok {
		return w, fmt.Errorf(
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
	w.Level.Creatures = w.Level.Creatures.Set(creature_id, creature)
	return w, nil
}

func DecideAction(subject_id world.ActorId) Action {
	return ActionTurn{
		Subject_id: subject_id,
		Direction:  world.LEFT(),
		Steps:      1,
	}
}
