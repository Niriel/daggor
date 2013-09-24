package world

type Coord int

type Location struct {
	X, Y Coord
}

type Position struct {
	Location
	F int // int sucks, I should do a Facing arithmetic.
}

func (self Location) ToPosition(facing int) Position {
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

var SIN = [...]Coord{0, 1, 0, -1}
var COS = [...]Coord{1, 0, -1, 0}

func (self Position) TurnLeft() Position {
	self.F = (self.F + 1) % 4
	return self
}
func (self Position) TurnRight() Position {
	// Here I add 3, because if I subtract 1 I get the stupid
	// go result: -1 % 4 = -1 (go) instead of -1 % 4 = 3 (python).
	// See this discussion:
	// https://code.google.com/p/go/issues/detail?id=448
	self.F = (self.F + 3) % 4
	return self
}
func (self Position) Forward() Position {
	self.X += COS[self.F]
	self.Y += SIN[self.F]
	return self
}
func (self Position) Backward() Position {
	self.X -= COS[self.F]
	self.Y -= SIN[self.F]
	return self
}
func (self Position) StrafeLeft() Position {
	self.X -= SIN[self.F]
	self.Y += COS[self.F]
	return self
}
func (self Position) StrafeRight() Position {
	self.X += SIN[self.F]
	self.Y -= COS[self.F]
	return self
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
func (self Position) SetF(f int) Position {
	self.F = f
	return self
}
