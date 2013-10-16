package world

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
