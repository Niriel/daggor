// world project world.go
package world

import (
	"encoding/gob"
	"fmt"
	"os"
)

// This should go into a package that knows about models.
type ModelId uint16

type World struct {
	Player_id ActorId // Later, there will also be a LevelId too in here.
	Level     Level   // Later, there will be many.
	Time      uint64  // Nanoseconds.
}

func Load() (*World, error) {
	f, err := os.Open("quicksave.sav")
	defer func(f *os.File) {
		if err_close := f.Close(); err_close != nil {
			fmt.Printf("File %v closed with error %v.", f, err_close.Error())
		}
	}(f)
	if err != nil {
		return nil, err
	}
	var world World
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&world)
	return &world, err
}

func (world *World) Save() error {
	f, err := os.Create("quicksave.sav")
	defer func(f *os.File) {
		if err_close := f.Close(); err_close != nil {
			fmt.Printf("File %v closed with error %v.", f, err_close.Error())
		}
	}(f)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(f)
	return encoder.Encode(world)
}

func MakeWorld() World {
	var world World
	world.Level = MakeLevel()

	// Place the player in the world.
	level := world.Level
	creature := Creature{F: EAST()}
	actors, actor_id, _ := level.Actors.Spawn()
	creatures, creature_id := level.Creatures.Add(creature)
	level.Actors = actors
	level.Creatures = creatures
	level.Creature_actor, _ = level.Creature_actor.Add(creature_id, actor_id)
	level.Creature_location, _ = level.Creature_location.Add(creature_id, Location{})
	level.Actor_schedule = MakeActorSchedule()
	world.Level = level
	world.Player_id = actor_id
	return world
}

func (world World) SetTime(time uint64) World {
	world.Time = time
	return world
}

func (world World) SetActorSchedule(actor_schedule ActorSchedule) World {
	world.Level = world.Level.SetActorSchedule(actor_schedule)
	return world
}
