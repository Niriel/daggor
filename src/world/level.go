package world

type Level struct {
	Floors   Buildings
	Ceilings Buildings
	Walls    [4]Buildings // Sorted by facing.
	Columns  Buildings
	Dynamic  Dynamic
}
