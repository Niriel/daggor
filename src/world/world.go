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
	Player_id ActorId  // Later, there will also be a LevelId too in here.
	Level     Level    // Later, there will be many.
	Actions   []Action // IA for the entire world.
	Time      uint64   // Nanoseconds.
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

func (world World) ExecuteActions() World {
	for _, action := range world.Actions {
		new_world, err := action.Execute(world)
		if err == nil {
			world = new_world
		} else {
			fmt.Println(action, err)
		}
	}
	world.Actions = make([]Action, 0, 16)
	return world
}
