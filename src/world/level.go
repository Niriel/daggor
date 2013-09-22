package world

type Level struct {
	Floors   Buildings
	Ceilings Buildings
	Walls    [4]Buildings // Sorted by facing.
	// 0: facing East, therefore western wall.
	// 1: facing North, therefore southern wall.
	// 2: facing West, therefore eastern wall.
	// 3: facing South, therefore northern wall.
}
