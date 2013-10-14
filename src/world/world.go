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
	Player_id      ActorId // Later, there will also be a LevelId too in here.
	Level          Level   // Later, there will be many.
	Time           uint64  // Nanoseconds.
	Actor_schedule ActorSchedule
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
	player_id := ActorId(0)
	var player_position Position
	player_position.X = 0
	player_position.Y = 0
	player_position.F = EAST()
	player := Actor{Pos: player_position}
	actors := make(Actors, 1)
	actors[player_id] = player
	world.Level.Actors = actors
	world.Actor_schedule = MakeActorSchedule()
	return world
}

func (world World) SetTime(time uint64) World {
	world.Time = time
	return world
}

func (world World) SetActorSchedule(actor_schedule ActorSchedule) World {
	world.Actor_schedule = actor_schedule
	return world
}
