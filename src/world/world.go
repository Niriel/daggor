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
	Player Player
	Level  Level
	Time   uint64 // Nanoseconds.
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
