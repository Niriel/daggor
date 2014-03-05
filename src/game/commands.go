package main

import (
	"fmt"
	glfw "github.com/go-gl/glfw3"
	"ia"
	"world"
)

type glfwKeyEvent struct {
	key      glfw.Key
	scancode int
	action   glfw.Action
	mods     glfw.ModifierKey
}

const eventListCap = 4

type glfwKeyEventList struct {
	list []glfwKeyEvent
}

func makeGlfwKeyEventList() *glfwKeyEventList {
	return &glfwKeyEventList{
		make([]glfwKeyEvent, 0, eventListCap),
	}
}

func (keyEventList *glfwKeyEventList) Freeze() []glfwKeyEvent {
	// The list of key events is double buffered.  This allows the application
	// to process events during a frame without having to worry about new
	// events arriving and growing the list.
	frozen := keyEventList.list
	keyEventList.list = make([]glfwKeyEvent, 0, eventListCap)
	return frozen
}

func (keyEventList *glfwKeyEventList) Callback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	event := glfwKeyEvent{key, scancode, action, mods}
	keyEventList.list = append(keyEventList.list, event)
}

type command int

const (
	commandForward = command(iota)
	commandBackward
	commandStrafeLeft
	commandStrafeRight
	commandTurnLeft
	commandTurnRight
	commandPlaceFloor
	commandPlaceCeiling
	commandPlaceWall
	commandPlaceColumn
	commandRemoveFloor
	commandRemoveCeiling
	commandRemoveWall
	commandRemoveColumn
	commandRotateFloorDirect
	commandRotateFloorRetrograde
	commandRotateCeilingDirect
	commandRotateCeilingRetrograde
	commandRotateColumnDirect
	commandRotateColumnRetrograde
	commandPlaceMonster
	commandRemoveMonster
	commandSave
	commandLoad
)

func commands(events []glfwKeyEvent) []command {
	if len(events) == 0 {
		return nil
	}
	result := make([]command, 0, 1)
	for _, event := range events {
		if event.action == glfw.Press {
			switch event.key {
			case glfw.KeyW:
				result = append(result, commandForward)
			case glfw.KeyS:
				result = append(result, commandBackward)
			case glfw.KeyA:
				result = append(result, commandStrafeLeft)
			case glfw.KeyD:
				result = append(result, commandStrafeRight)
			case glfw.KeyQ:
				result = append(result, commandTurnLeft)
			case glfw.KeyE:
				result = append(result, commandTurnRight)
			case glfw.KeyF:
				if event.mods&glfw.ModShift == 0 {
					result = append(result, commandPlaceFloor)
				} else {
					result = append(result, commandRemoveFloor)
				}
			case glfw.KeyC:
				if event.mods&glfw.ModShift == 0 {
					result = append(result, commandPlaceCeiling)
				} else {
					result = append(result, commandRemoveCeiling)
				}
			case glfw.KeyR:
				if event.mods&glfw.ModShift == 0 {
					result = append(result, commandPlaceWall)
				} else {
					result = append(result, commandRemoveWall)
				}
			case glfw.KeyK:
				if event.mods&glfw.ModShift == 0 {
					result = append(result, commandPlaceColumn)
				} else {
					result = append(result, commandRemoveColumn)
				}
			case glfw.KeyM:
				if event.mods&glfw.ModShift == 0 {
					result = append(result, commandPlaceMonster)
				} else {
					result = append(result, commandRemoveMonster)
				}
			case glfw.KeyLeftBracket:
				result = append(result, commandRotateCeilingDirect)
			case glfw.KeyRightBracket:
				result = append(result, commandRotateCeilingRetrograde)
			case glfw.KeySemicolon:
				result = append(result, commandRotateColumnDirect)
			case glfw.KeyApostrophe:
				result = append(result, commandRotateColumnRetrograde)
			case glfw.KeyPeriod:
				result = append(result, commandRotateFloorDirect)
			case glfw.KeySlash:
				result = append(result, commandRotateFloorRetrograde)
			case glfw.KeyF4:
				result = append(result, commandSave)
			case glfw.KeyF5:
				result = append(result, commandLoad)
			}
		}
	}
	return result
}

func commandToAction(command command, subjectID world.ActorID) ia.Action {
	var action ia.Action
	switch command {
	// If an action is what an actor does when it's its turn to play, then
	// maybe we don't want turning to be one.  Moving yes, turning no.  We'll
	// see.
	case commandTurnLeft:
		action = ia.ActionTurn{
			SubjectID: subjectID,
			Direction: world.LEFT(),
			Steps:     1,
		}
	case commandTurnRight:
		action = ia.ActionTurn{
			SubjectID: subjectID,
			Direction: world.RIGHT(),
			Steps:     1,
		}
	case commandForward:
		action = ia.ActionMoveRelative{
			SubjectID: subjectID,
			Direction: world.FRONT(),
			Steps:     1,
		}
	case commandStrafeLeft:
		action = ia.ActionMoveRelative{
			SubjectID: subjectID,
			Direction: world.LEFT(),
			Steps:     1,
		}
	case commandBackward:
		action = ia.ActionMoveRelative{
			SubjectID: subjectID,
			Direction: world.BACK(),
			Steps:     1,
		}
	case commandStrafeRight:
		action = ia.ActionMoveRelative{
			SubjectID: subjectID,
			Direction: world.RIGHT(),
			Steps:     1,
		}
	}
	return action
}

func commandsToAction(commands []command, subjectID world.ActorID) (ia.Action, []command) {
	var actionResult ia.Action
	commandsResult := make([]command, 0, cap(commands))
	for _, command := range commands {
		action := commandToAction(command, subjectID)
		if action == nil {
			commandsResult = append(commandsResult, command)
		} else {
			if actionResult == nil {
				// Keep the first action only, the other are discarded.  It should
				// not be a big loss anyway as this function is called every frame.
				// How many keys can you hope to press in 15 milliseconds?
				actionResult = action
			} else {
				fmt.Println("Discarded action ", action)
			}
		}
	}
	return actionResult, commandsResult
}

func levelCommand(level world.Level, position world.Position, command command) world.Level {
	hereX, hereY := position.X, position.Y
	there := position.MoveForward(1)
	thereX, thereY := there.X, there.Y
	switch command {
	case commandPlaceFloor:
		floor := world.MakeFloor(floorID, world.EAST(), true)
		level.Floors = level.Floors.Set(thereX, thereY, floor)
	case commandPlaceCeiling:
		ceiling := world.MakeOrientedBuilding(ceilingID, world.EAST())
		level.Ceilings = level.Ceilings.Set(thereX, thereY, ceiling)
	case commandPlaceWall:
		{
			// If the player faces North, then the wall must face South in order
			// to face the player.
			facing := position.F.Add(world.BACK())
			index := facing.Value()
			wall := world.MakeWall(wallID, false)
			level.Walls[index] = level.Walls[index].Set(hereX, hereY, wall)
		}
	case commandPlaceColumn:
		level.Columns = level.Columns.Set(thereX, thereY, world.MakeOrientedBuilding(columnID, world.EAST()))
	case commandRotateFloorDirect, commandRotateFloorRetrograde:
		{
			building := level.Floors[world.Location{X: thereX, Y: thereY}]
			floor, ok := building.(world.Floor)
			if ok {
				var relDir world.RelativeDirection
				if command == commandRotateFloorDirect {
					relDir = world.LEFT()
				} else {
					relDir = world.RIGHT()
				}
				floor.F = floor.F.Add(relDir)
				level.Floors = level.Floors.Set(thereX, thereY, floor)
			} else {
				fmt.Println("You cannot rotate that.")
			}
		}
	case commandRotateColumnDirect, commandRotateColumnRetrograde:
		{
			var relDir world.RelativeDirection
			if command == commandRotateColumnDirect {
				relDir = world.LEFT()
			} else {
				relDir = world.RIGHT()
			}
			column := level.Columns[world.Location{X: thereX, Y: thereY}]
			orientable, ok := column.(world.OrientedBuilding)
			if ok {
				orientable.F = orientable.F.Add(relDir)
				level.Columns = level.Columns.Set(thereX, thereY, orientable)
			} else {
				fmt.Println("You cannot rotate that.")
			}
		}
	case commandRotateCeilingDirect, commandRotateCeilingRetrograde:
		{
			var relDir world.RelativeDirection
			if command == commandRotateCeilingDirect {
				relDir = world.LEFT()
			} else {
				relDir = world.RIGHT()
			}
			ceiling := level.Ceilings[world.Location{X: thereX, Y: thereY}]
			orientable, ok := ceiling.(world.OrientedBuilding)
			if ok {
				orientable.F = orientable.F.Add(relDir)
				level.Ceilings = level.Ceilings.Set(thereX, thereY, orientable)
			} else {
				fmt.Println("You cannot rotate that.")
			}
		}
	case commandRemoveFloor:
		level.Floors = level.Floors.Delete(thereX, thereY)
	case commandRemoveColumn:
		level.Columns = level.Columns.Delete(thereX, thereY)
	case commandRemoveCeiling:
		level.Ceilings = level.Ceilings.Delete(thereX, thereY)
	case commandRemoveWall:
		{
			// If the player faces North, then the wall must face South in order
			// to face the player.
			facing := position.F.Add(world.BACK())
			index := facing.Value()
			level.Walls[index] = level.Walls[index].Delete(hereX, hereY)
		}

	case commandPlaceMonster:
		{
			creature := world.MakeCreature()
			creature.F = position.F
			creatures, creatureID := level.Creatures.Add(creature)
			actors, actorID := level.Actors.Add(world.MakeActor())
			creatureActors, err := level.CreatureActor.Add(creatureID, actorID)
			if err != nil {
				fmt.Println(err)
				break
			}
			creatureLocations, err := level.CreatureLocation.Add(creatureID, there.ToLocation())
			if err != nil {
				fmt.Println(err)
				break
			}
			level.Creatures = creatures
			level.Actors = actors
			level.CreatureActor = creatureActors
			level.CreatureLocation = creatureLocations
		}
	case commandRemoveMonster:
		{
			fmt.Printf("Not implemented.")
		}
	}
	return level
}

func executeCommands(programState programState, commands []command) programState {
	for _, command := range commands {
		switch {
		case command < commandSave:
			position, ok := programState.World.Level.ActorPosition(programState.World.Player_id)
			if !ok {
				break
			}
			// Modify the world around the player character.
			programState.World.Level = levelCommand(programState.World.Level, position, command)
		case command == commandSave:
			err := programState.World.Save()
			fmt.Println("Save:", err)
		case command == commandLoad:
			world, err := world.Load()
			fmt.Println("Load:", err)
			if err == nil {
				programState.World = *world
			}
		}
	}
	return programState
}
