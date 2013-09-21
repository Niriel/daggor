package world

import (
	"encoding/gob"
)

func init() {
	// For Gob to be able to deal with interfaces, we must
	// teach it the concrete types that implement these
	// interfaces.
	gob.Register(MakeBaseBuilding(0))
	gob.Register(MakeOrientedBuilding(0, 0))
}

type Building interface {
	Model() ModelId
}

type BaseBuilding struct {
	Model_ ModelId
}

type OrientedBuilding struct {
	BaseBuilding
	Facing int
}

func MakeBaseBuilding(model ModelId) BaseBuilding {
	return BaseBuilding{Model_: model}
}

func (self BaseBuilding) Model() ModelId {
	return self.Model_
}

func MakeOrientedBuilding(model ModelId, facing int) OrientedBuilding {
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
