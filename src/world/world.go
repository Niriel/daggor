// world project world.go
package world

type Coord int

type Coords struct {
	X, Y Coord
}

type Building uint16

type Buildings map[Coords]Building

type Level struct {
	Floors   Buildings
	Ceilings Buildings
	WallsE   Buildings
	WallsN   Buildings
	WallsW   Buildings
	WallsS   Buildings
}

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
	dst[Coords{x, y}] = value
	return dst
}

func (src Buildings) Get(x, y Coord) (Building, bool) {
	building, ok := src[Coords{x, y}]
	return building, ok
}

func (src Buildings) Delete(x, y Coord) Buildings {
	dst := src.Copy()
	delete(dst, Coords{x, y})
	return dst
}
