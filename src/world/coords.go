package world

import (
	"encoding/gob"
	"fmt"
)

func init() {
	// For Gob to be able to deal with interfaces, we must
	// teach it the concrete types that implement these
	// interfaces.
	gob.Register(absoluteDirection{0})
	gob.Register(relativeDirection{0})
}

// I am not exporting these types because I want to limit their instances.
// Directions MUST be bound to 0..3.  There is some modulo artihmetics involved
// in rotation, which is easy to get wrong with Go's modulo operator, especially
// when using negative values.
type absoluteDirection struct {
	value int
}
type relativeDirection struct {
	value int
}

// Although the following values are not constant, please do not change them.
// The user could do something like EAST, NORTH = NORTH, EAST and mess up the whole
// thing.  Just don't.

// Absolute directions.
var EAST = absoluteDirection{0}
var NORTH = absoluteDirection{1}
var WEST = absoluteDirection{2}
var SOUTH = absoluteDirection{3}

// Relative directions.
var FRONT = relativeDirection{0}
var LEFT = relativeDirection{1}
var BACK = relativeDirection{2}
var RIGHT = relativeDirection{3}

// The rotation arithmetics is abstracted behind these two interfaces.  You can
// combine relative directions together.  For example, left of left is back.
// Back of right is left. Note that this is commutative: right of back is left too.
type RelativeDirection interface {
	// This `concrete` unexported field does two things.
	// First, it prevents anybody else from creating a new type that satisfies
	// this RelativeDirection interface, since the relativeDirection type is
	// private.  Then, it is handy in the Add method.  The name "concrete" refers
	// to the fact that this function returns a variable of the concrete type
	// implementing the interface.
	concrete() relativeDirection
	Value() int
	Add(rel RelativeDirection) RelativeDirection
}

// Relative directions can be applied to absolute directions to get a new absolute
// direction.
type AbsoluteDirection interface {
	concrete() absoluteDirection
	Value() int
	Add(b RelativeDirection) AbsoluteDirection
	DxDy() (Coord, Coord)
}

type Coord int

var COS = [...]Coord{1, 0, -1, 0}
var SIN = [...]Coord{0, 1, 0, -1}

func (rel relativeDirection) concrete() relativeDirection {
	return rel
}
func (rel relativeDirection) Value() int {
	return rel.value
}
func (rel0 relativeDirection) Add(rel1 RelativeDirection) RelativeDirection {
	return relativeDirection{(rel0.value + rel1.Value()) % 4}
}

func (dir absoluteDirection) concrete() absoluteDirection {
	return dir
}
func (dir absoluteDirection) Value() int {
	return dir.value
}
func (dir absoluteDirection) Add(rel RelativeDirection) AbsoluteDirection {
	return absoluteDirection{(dir.value + rel.concrete().value) % 4}
}

func (dir absoluteDirection) GobEncode() ([]byte, error) {
	if dir.value >= 0 && dir.value <= 3 {
		slice := []byte{byte(dir.value)}
		return slice, nil
	}
	return nil, fmt.Errorf("Internal value of an absolute direction should be in [0..4], not %v.", dir.value)
}

func (dir *absoluteDirection) GobDecode(bytes []byte) error {
	if len(bytes) != 1 {
		return fmt.Errorf("absoluteDirection needs exactly one byte of data, not %v.", len(bytes))
	}
	value := bytes[0]
	if value < 0 || value > 3 {
		return fmt.Errorf("absoluteDirection needs to contain a value between 0 and 3, not %v.", value)
	}
	dir.value = int(value)
	return nil
}

func (dir absoluteDirection) DxDy() (Coord, Coord) {
	return COS[dir.value], SIN[dir.value]
}

type Location struct {
	X, Y Coord
}

type Position struct {
	Location
	F AbsoluteDirection
}

func (self Location) ToPosition(facing AbsoluteDirection) Position {
	var position Position
	position.X = self.X
	position.Y = self.Y
	position.F = facing
	return position
}

func (self Position) ToLocation() Location {
	var location Location
	location.X = self.X
	location.Y = self.Y
	return location
}

func (self Location) SetX(x Coord) Location {
	self.X = x
	return self
}
func (self Location) SetY(y Coord) Location {
	self.Y = y
	return self
}

func (self Position) SetX(x Coord) Position {
	self.X = x
	return self
}
func (self Position) SetY(y Coord) Position {
	self.Y = y
	return self
}
func (self Position) SetF(f AbsoluteDirection) Position {
	self.F = f
	return self
}

func (self Position) TurnLeft() Position {
	self.F = self.F.Add(LEFT)
	return self
}
func (self Position) TurnRight() Position {
	self.F = self.F.Add(RIGHT)
	return self
}

func (self Position) MoveAbsolute(absdir AbsoluteDirection, steps int) Position {
	dx, dy := absdir.DxDy()
	// Ugly casting.  dx and dy are distances, which I multiply by a pure number
	// `steps`.  The result are also distances, which I add to a position to get
	// a position.  It's a mess of types, units and dimensions.
	self.X += dx * Coord(steps)
	self.Y += dy * Coord(steps)
	return self
}

func (self Position) MoveRelative(reldir RelativeDirection, steps int) Position {
	return self.MoveAbsolute(self.F.Add(reldir), steps)
}

func (self Position) MoveForward(steps int) Position {
	return self.MoveRelative(FRONT, steps)
}
func (self Position) MoveLeft(steps int) Position {
	return self.MoveRelative(LEFT, steps)
}
func (self Position) MoveBackward(steps int) Position {
	return self.MoveRelative(BACK, steps)
}
func (self Position) MoveRight(steps int) Position {
	return self.MoveRelative(RIGHT, steps)
}
