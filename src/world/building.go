package world

import (
	"encoding/gob"
)

func init() {
	// For Gob to be able to deal with interfaces, we must
	// teach it the concrete types that implement these
	// interfaces.
	gob.Register(MakeBaseBuilding(0))
	gob.Register(MakeOrientedBuilding(0, EAST()))
	gob.Register(MakeFloor(0, EAST(), false))
	gob.Register(MakeWall(0, false))
}

type Building interface {
	Model() ModelId
}

type BaseBuilding struct {
	// Underscore means that this field is public but should not be used.
	// Please use the Model() function to read it, and do NOT change it.
	// It is public only because Gob requires it.  I may spend time writing
	// a GobEncoder later to hide this field again, but I do not want to
	// spend time maintaining code that is still so likely to change.
	Model_ ModelId
}

type OrientedBuilding struct {
	BaseBuilding
	Facing AbsoluteDirection
}

type MaybePassable struct {
	Passable_ bool
}

func (self MaybePassable) IsPassable() bool {
	return self.Passable_
}

type Floor struct {
	OrientedBuilding
	MaybePassable
}

type Wall struct {
	BaseBuilding
	MaybePassable
}

func MakeBaseBuilding(model ModelId) BaseBuilding {
	return BaseBuilding{Model_: model}
}

func (self BaseBuilding) Model() ModelId {
	return self.Model_
}

func MakeFloor(model ModelId, facing AbsoluteDirection, passable bool) Floor {
	var floor Floor
	floor.Model_ = model
	floor.Facing = facing
	floor.Passable_ = passable
	return floor
}

func MakeWall(model ModelId, passable bool) Wall {
	var wall Wall
	wall.Model_ = model
	wall.Passable_ = passable
	return wall
}

func MakeOrientedBuilding(model ModelId, facing AbsoluteDirection) OrientedBuilding {
	var result OrientedBuilding
	result.Model_ = model
	result.Facing = facing
	return result
}

type Buildings map[Location]Building

// Making copies is required to produce updated version of maps.
func (src Buildings) Copy() Buildings {
	dst := make(Buildings, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func (src Buildings) Set(x, y Coord, value Building) Buildings {
	dst := src.Copy()
	dst[Location{x, y}] = value
	return dst
}

func (src Buildings) Get(x, y Coord) (Building, bool) {
	building, ok := src[Location{x, y}]
	return building, ok
}

func (src Buildings) Delete(x, y Coord) Buildings {
	dst := src.Copy()
	delete(dst, Location{x, y})
	return dst
}
