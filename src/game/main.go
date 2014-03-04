// game project main.go
package main

import (
	"batch"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"glm"
	"glw"
	"ia"
	"runtime"
	"sculpt"
	"time"
	"world"
)

func init() {
	// OpenGL and GLFW want to run on the main thread.
	// Or at least, want to run always from the same thread.
	runtime.LockOSThread()
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

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

const (
	cubeID = iota
	pyramidID
	floorID
	wallID
	columnID
	ceilingID
	monsterID
)

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

func viewMatrix(pos world.Position) glm.Matrix4 {
	R := glm.RotZ(float64(-90 * pos.F.Value()))
	T := glm.Vector3{float64(-pos.X), float64(-pos.Y), -.5}.Translation()
	return R.Mult(T)
}

type glState struct {
	Window           *glfw.Window
	glfwKeyEventList *glfwKeyEventList
	Shapes           [7]*sculpt.Mesh
	context          *glw.GlContext
}

type programState struct {
	Gl    glState     // Highly mutable, impure.
	World world.World // Immutable, pure.
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

func main() {
	var programState programState
	var err error
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("GLFW initialization failed.")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.SrgbCapable, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	programState.Gl.Window, err = glfw.CreateWindow(640, 480, "Daggor", nil, nil)
	if err != nil {
		panic(err)
	}
	defer programState.Gl.Window.Destroy()

	programState.Gl.glfwKeyEventList = makeGlfwKeyEventList()
	programState.Gl.Window.SetKeyCallback(programState.Gl.glfwKeyEventList.Callback)

	programState.Gl.Window.MakeContextCurrent()
	if ec := gl.Init(); ec != 0 {
		panic(fmt.Sprintf("OpenGL initialization failed with code %v.", ec))
	}
	// For some reason, here, the OpenGL error flag for me contains "Invalid enum".
	// This is weird since I have not done anything yet.  I imagine that something
	// goes wrong in gl.Init.  Reading the error flag clears it, so I do it.
	// Here's the reason:
	//     https://github.com/go-gl/glfw3/issues/50
	// Maybe I should not even ask for a core profile anyway.
	// What are the advantages are asking for a core profile?
	if err := glw.CheckGlError(); err != nil {
		err.Description = "OpenGL has this error right after init for some reason."
		//fmt.Println(err)
	}
	major := programState.Gl.Window.GetAttribute(glfw.ContextVersionMajor)
	minor := programState.Gl.Window.GetAttribute(glfw.ContextVersionMinor)
	fmt.Printf("OpenGL version %v.%v.\n", major, minor)
	if (major < 3) || (major == 3 && minor < 3) {
		panic("OpenGL version 3.3 required, your video card/driver does not seem to support it.")
	}

	programState.Gl.context = glw.NewGlContext()

	//programState.Gl.Shapes[cubeID] = glw.Cube(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[pyramidID] = glw.Pyramid(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[floorID] = glw.Floor(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[wallID] = glw.Wall(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[columnID] = glw.Column(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[ceilingID] = glw.Ceiling(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[monsterID] = glw.Monster(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.DynaPyramid = glw.DynaPyramid(programState.Gl.context.Programs)
	programState.Gl.Shapes[floorID] = sculpt.FloorInst(programState.Gl.context.Programs)
	programState.Gl.Shapes[floorID].SetUpVao()

	// I do not like the default reference frame of OpenGl.
	// By default, we look in the direction -z, and y points up.
	// I want z to point up, and I want to look in the direction +x
	// by default.  That way, I move on an xy plane where z is the
	// altitude, instead of having the altitude stuffed between
	// the two things I use the most.  And my reason for pointing
	// toward +x is that I use the convention for trigonometry:
	// an angle of 0 points to the right (east) of the trigonometric
	// circle.  Bonus point: this matches Blender's reference frame.
	myFrame := glm.ZUP.Mult(glm.RotZ(90))
	projectionMatrix := glm.PerspectiveProj(110, 640./480., .1, 100).Mult(myFrame)
	programState.Gl.context.SetCameraProj(projectionMatrix)

	gl.Enable(gl.FRAMEBUFFER_SRGB)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)

	//blah := make([]int32, 1)
	//for _, thing := range []gl.GLenum{gl.UNIFORM_BUFFER_OFFSET_ALIGNMENT,
	//	gl.UNIFORM_BLOCK_DATA_SIZE,
	//} {
	//	gl.GetIntegerv(thing, blah)
	//	fmt.Println(blah)
	//}

	programState.World = world.MakeWorld()
	mainLoop(programState)
}

func mainLoop(programState programState) programState {
	const tickPeriod = 1000000000 / 60
	ticker := time.NewTicker(tickPeriod * time.Nanosecond)
	keepTicking := true
	for keepTicking {
		select {
		case _, ok := <-ticker.C:
			{
				if ok {
					programState, keepTicking = onTick(programState, tickPeriod)
					if !keepTicking {
						fmt.Println("No more ticks.")
						ticker.Stop()
					}
				} else {
					fmt.Println("Ticker closed, weird.")
					keepTicking = false
				}
			}
		}
	}
	return programState
}

func onTick(programState programState, dt uint64) (programState, bool) {
	glfw.PollEvents()
	keepTicking := !programState.Gl.Window.ShouldClose()
	if keepTicking {
		// Read raw inputs.
		keys := programState.Gl.glfwKeyEventList.Freeze()
		// Analyze the inputs, see what they mean.
		commands := commands(keys)
		// One of these commands may correspond to an action of the player's actor.
		// We take it out so that we can process it in the IA phase.
		// The remaining commands are kept for further processing.
		playerAction, commands := commandsToAction(commands, programState.World.Player_id)
		// Evolve the program one step.
		programState.World.Time += dt // No side effect, we own a copy.
		// $$$ THERE COULD BE SIDE EFFECTS HERE ACTUALLY:  IF I GAVE A POINTER
		// TO THE WORLD OR PROGRAM STATE TO SOMETHING.  NEED TO CORRECT THAT.
		programState = executeCommands(programState, commands)
		//
		programState.World = runAI(programState.World, playerAction)
		// render on screen.
		render(programState)
		programState.Gl.Window.SwapBuffers()
	}
	return programState, keepTicking
}

func runAI(w world.World, playerAction ia.Action) world.World {
	var action ia.Action
	// It's like on a board game.  Every one plays when it is their turn.
	// This function is called every frame.

	// Temporary: Any creature that is not scheduled yet is added to the
	// scheduler.
	schedule := w.Level.ActorSchedule
	for actorID := range w.Level.Actors.Content() {
		index := schedule.PosActorID(actorID)
		if index == -1 {
			fmt.Println("Force scheduling", actorID)
			schedule = schedule.Add(actorID, w.Time)
		}
	}
	w = w.SetActorSchedule(schedule)

	for {
		actorTime, ok := w.Level.ActorSchedule.Next(w.Time)
		if !ok {
			// Actions can modify the list of actors, so I cannot loop over
			// all the actors.  This is why I break the loop this way.
			break // No more actors to process.
		}
		newSchedule, ok := w.Level.ActorSchedule.Remove(actorTime)
		if !ok {
			panic("Could not find actor to remove from scheduler")
		}
		w = w.SetActorSchedule(newSchedule)
		if actorTime.Actor_id == w.Player_id {
			action = playerAction
		} else {
			action = ia.DecideAction(actorTime.Actor_id)
		}
		if action != nil {
			var err error
			newSchedule = newSchedule.Add(actorTime.Actor_id, actorTime.Time+100000000)
			w = w.SetActorSchedule(newSchedule)
			w, err = action.Execute(w)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// Nil actions should only happen for the player.  The player is the
			// only actor who can decide not to act.  All other actors decide
			// an action, even if it is just a waiting action.
			if actorTime.Actor_id == w.Player_id {
				// Reschedule the player for next turn.
				newSchedule = newSchedule.Add(actorTime.Actor_id, w.Time+1)
				w = w.SetActorSchedule(newSchedule)
			} else {
				panic("Only the player is allowed to idle.")
			}
		}
	}
	return w
}

func render(programState programState) {
	clearBatch := batch.MakeClearBatch(
		gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT,
		[4]gl.GLclampf{0.0, 0.0, 0.4, 0.0},
		1,
	)

	actorID := programState.World.Player_id
	position, ok := programState.World.Level.ActorPosition(actorID)
	if !ok {
		panic("Could not find player's character position.")
	}

	view := viewMatrix(position)

	camBatch := batch.MakeCameraBatch(
		programState.Gl.context,
		view,
		programState.Gl.context.CameraProj(),
	)

	clearBatch.Batches = append(clearBatch.Batches, camBatch)

	clearBatch.Enter()
	clearBatch.Run()
	clearBatch.Exit()
	renderBuildings(
		programState.World.Level.Floors,
		0, 0,
		nil,
		view,
		programState.Gl,
	)
}

func renderBuildings(
	buildings world.Buildings,
	offsetX, offsetY float64,
	defaultR *glm.Matrix4, // Can be nil.
	view glm.Matrix4,
	glState glState,
) {
	locations := make(map[world.ModelId][]glm.Matrix4)
	for coords, building := range buildings {
		m := glm.Vector3{
			float64(coords.X) + offsetX,
			float64(coords.Y) + offsetY,
			0,
		}.Translation()
		facer, ok := building.(world.Facer)
		if ok {
			// We obey the facing of the buildings that have one.
			r := glm.RotZ(float64(90 * facer.Facing().Value()))
			m = m.Mult(r)
		} else {
			// Buildings without facing receive the provided default facing.
			// It is given as a precalculated rotation matrix `defaultR`.
			m = m.Mult(*defaultR)
		}
		m = view.Mult(m) // Shaders work in view space.
		modelID := building.Model()
		meshlocs := locations[modelID]
		meshlocs = append(meshlocs, m)
		locations[modelID] = meshlocs
	}
	//
	for modelID, locs := range locations {
		mesh := glState.Shapes[modelID]
		if len(locs) == 0 {
			panic("empty list of locations for model")
		}
		unif, ok := mesh.Uniforms.(*sculpt.UniformsLoc)
		if ok {
			for _, loc := range locs {
				unif.SetModel(loc)
				mesh.Draw()
			}
		}
		inst, ok := mesh.Instances.(*sculpt.ModelMatInstances)
		if ok {
			gllocs := make([]sculpt.ModelMatInstance, len(locs), len(locs))
			for i, mat := range locs {
				gllocs[i] = sculpt.ModelMatInstance{mat.GlFloats()}
			}
			inst.SetData(gllocs)
		}
		draw, ok := mesh.Drawer.(*sculpt.DrawElementInstanced)
		if ok {
			draw.Primcount = len(locs)
			mesh.Draw()
		}
	}
}
