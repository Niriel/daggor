package world

import (
	"fmt"
)

// ActorID is used as a unique identifier for an Actor.
type ActorID uint64

// An Actor is anything that can be at the origin of an Action (see the IA
// package for Actions).  It can be a creature, a trap, a mechanism, etc.
// An Actor is not aware of its unique identifier, you have to keep track of it
// yourself.
type Actor struct{}

// MakeActor creates, initializes and returns an Actor.
func MakeActor() Actor {
	return Actor{}
}

// An Actors object holds a collection of Actor objects s identified by an
// ActorID.
// It contains a map of Actors indexed by their unique identifier.
// It also manages unique identifiers.  Use the Add method to get a new unique
// identifier for an Actor.
type Actors struct {
	// NextIDprivate is exported only so that it can be encoded by gob.
	// Do not use.
	// It starts at 0 and increases each time the Add method is called.
	NextIDprivate ActorID
	// ContentPrivate is exported only so that it can be encoded by gob.
	// Do not use.  Use the Actors.Content() method instead.
	ContentPrivate map[ActorID]Actor
}

// MakeActors creates and returns an empty collection of Actors.
func MakeActors() Actors {
	return Actors{
		ContentPrivate: make(map[ActorID]Actor),
	}
}

// Content returns a copy of the private map[ActorID]Actor.
func (actors Actors) Content() map[ActorID]Actor {
	contentCopy := make(map[ActorID]Actor)
	for key, value := range actors.ContentPrivate {
		contentCopy[key] = value
	}
	return contentCopy
}

// Copy returns a deep-copy of the Actors receiver.
func (actors Actors) Copy() Actors {
	newActors := Actors{
		NextIDprivate:  actors.NextIDprivate,
		ContentPrivate: make(map[ActorID]Actor),
	}
	for key, value := range actors.ContentPrivate {
		newActors.ContentPrivate[key] = value
	}
	return newActors
}

// Replace returns a new Actors object in which the Actor actor is given the
// ActorID actorID.  There must be an existing Actor with this ID.  If not,
// this method panics.
func (actors Actors) Replace(actorID ActorID, actor Actor) Actors {
	_, ok := actors.ContentPrivate[actorID]
	if !ok {
		panic(fmt.Sprintf("cannot replace inexistent Actor ID=%v", actorID))
	}
	newActors := actors.Copy()
	newActors.ContentPrivate[actorID] = actor
	return newActors
}

// Add returns a new Actors object to which the given Actor actor is added.
// The method also returns the ActorID that was given to actor.
func (actors Actors) Add(actor Actor) (Actors, ActorID) {
	actorID := actors.NextIDprivate
	newActors := actors.Copy()
	newActors.ContentPrivate[actorID] = actor
	newActors.NextIDprivate++
	return newActors, actorID
}

// Delete returns a new Actors object from which the entry corresponding to
// the provided actorID is deleted.  If actorID does not exist, this method
// panics.
func (actors Actors) Delete(actorID ActorID) Actors {
	_, ok := actors.ContentPrivate[actorID]
	if !ok {
		panic(fmt.Sprintf("cannot remove inexistent Actor ID %v", actorID))
	}
	newActors := actors.Copy()
	delete(newActors.ContentPrivate, actorID)
	return newActors
}
