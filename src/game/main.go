// game project main.go
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"glm"
	"glw"
	"runtime"
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

type GlfwKeyEvent struct {
	key      glfw.Key
	scancode int
	action   glfw.Action
	mods     glfw.ModifierKey
}

const EVENT_LIST_CAP = 4

type GlfwKeyEventList struct {
	list []GlfwKeyEvent
}

func MakeGlfwKeyEventList() *GlfwKeyEventList {
	return &GlfwKeyEventList{
		make([]GlfwKeyEvent, 0, EVENT_LIST_CAP),
	}
}

func (self *GlfwKeyEventList) Freeze() []GlfwKeyEvent {
	// The list of key events is double buffered.  This allows the application
	// to process events during a frame without having to worry about new
	// events arriving and growing the list.
	frozen := self.list
	self.list = make([]GlfwKeyEvent, 0, EVENT_LIST_CAP)
	return frozen
}

func (self *GlfwKeyEventList) Callback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	event := GlfwKeyEvent{key, scancode, action, mods}
	self.list = append(self.list, event)
}

const (
	CUBE_ID = iota
	PYRAMID_ID
	FLOOR_ID
	WALL_ID
	COLUMN_ID
	CEILING_ID
)

type Command int

const (
	COMMAND_FORWARD = Command(iota)
	COMMAND_BACKWARD
	COMMAND_STRAFE_LEFT
	COMMAND_STRAFE_RIGHT
	COMMAND_TURN_LEFT
	COMMAND_TURN_RIGHT
	COMMAND_PLACE_FLOOR
	COMMAND_PLACE_CEILING
	COMMAND_PLACE_WALL
	COMMAND_PLACE_COLUMN
	COMMAND_REMOVE_FLOOR
	COMMAND_REMOVE_CEILING
	COMMAND_REMOVE_WALL
	COMMAND_REMOVE_COLUMN
	COMMAND_ROTATE_FLOOR_DIRECT
	COMMAND_ROTATE_FLOOR_RETROGRADE
	COMMAND_ROTATE_CEILING_DIRECT
	COMMAND_ROTATE_CEILING_RETROGRADE
	COMMAND_ROTATE_COLUMN_DIRECT
	COMMAND_ROTATE_COLUMN_RETROGRADE
	COMMAND_SAVE
	COMMAND_LOAD
)

func Commands(events []GlfwKeyEvent) []Command {
	if len(events) == 0 {
		return nil
	}
	result := make([]Command, 0, 1)
	for _, event := range events {
		if event.action == glfw.Press {
			switch event.key {
			case glfw.KeyW:
				result = append(result, COMMAND_FORWARD)
			case glfw.KeyS:
				result = append(result, COMMAND_BACKWARD)
			case glfw.KeyA:
				result = append(result, COMMAND_STRAFE_LEFT)
			case glfw.KeyD:
				result = append(result, COMMAND_STRAFE_RIGHT)
			case glfw.KeyQ:
				result = append(result, COMMAND_TURN_LEFT)
			case glfw.KeyE:
				result = append(result, COMMAND_TURN_RIGHT)
			case glfw.KeyF:
				if event.mods&glfw.ModShift == 0 {
					result = append(result, COMMAND_PLACE_FLOOR)
				} else {
					result = append(result, COMMAND_REMOVE_FLOOR)
				}
			case glfw.KeyC:
				if event.mods&glfw.ModShift == 0 {
					result = append(result, COMMAND_PLACE_CEILING)
				} else {
					result = append(result, COMMAND_REMOVE_CEILING)
				}
			case glfw.KeyR:
				if event.mods&glfw.ModShift == 0 {
					result = append(result, COMMAND_PLACE_WALL)
				} else {
					result = append(result, COMMAND_REMOVE_WALL)
				}
			case glfw.KeyK:
				if event.mods&glfw.ModShift == 0 {
					result = append(result, COMMAND_PLACE_COLUMN)
				} else {
					result = append(result, COMMAND_REMOVE_COLUMN)
				}
			case glfw.KeyLeftBracket:
				result = append(result, COMMAND_ROTATE_CEILING_DIRECT)
			case glfw.KeyRightBracket:
				result = append(result, COMMAND_ROTATE_CEILING_RETROGRADE)
			case glfw.KeySemicolon:
				result = append(result, COMMAND_ROTATE_COLUMN_DIRECT)
			case glfw.KeyApostrophe:
				result = append(result, COMMAND_ROTATE_COLUMN_RETROGRADE)
			case glfw.KeyPeriod:
				result = append(result, COMMAND_ROTATE_FLOOR_DIRECT)
			case glfw.KeySlash:
				result = append(result, COMMAND_ROTATE_FLOOR_RETROGRADE)
			case glfw.KeyF4:
				result = append(result, COMMAND_SAVE)
			case glfw.KeyF5:
				result = append(result, COMMAND_LOAD)
			}
		}
	}
	return result
}

func ViewMatrix(pos world.Position) glm.Matrix4 {
	R := glm.RotZ(float64(-90 * pos.F.Value()))
	T := glm.Vector3{float64(-pos.X), float64(-pos.Y), -.5}.Translation()
	return R.Mult(T)
}

type GlState struct {
	Window           *glfw.Window
	GlfwKeyEventList *GlfwKeyEventList
	Programs         glw.Programs
	P                glm.Matrix4
	Shapes           [6]glw.Drawable
	DynaPyramid      glw.StreamDrawable
	Monster          glw.Drawable
}

type ProgramState struct {
	Gl    GlState     // Highly mutable, impure.
	World world.World // Immutable, pure.
}

func CommandToAction(command Command, subject_id world.ActorId) world.Action {
	var action world.Action
	switch command {
	// If an action is what an actor does when it's its turn to play, then
	// maybe we don't want turning to be one.  Moving yes, turning no.  We'll
	// see.
	case COMMAND_TURN_LEFT:
		action = world.ActionTurn{
			Subject_id: subject_id,
			Direction:  world.LEFT(),
			Steps:      1,
		}
	case COMMAND_TURN_RIGHT:
		action = world.ActionTurn{
			Subject_id: subject_id,
			Direction:  world.RIGHT(),
			Steps:      1,
		}
	case COMMAND_FORWARD:
		action = world.ActionMoveRelative{
			Subject_id: subject_id,
			Direction:  world.FRONT(),
			Steps:      1,
		}
	case COMMAND_STRAFE_LEFT:
		action = world.ActionMoveRelative{
			Subject_id: subject_id,
			Direction:  world.LEFT(),
			Steps:      1,
		}
	case COMMAND_BACKWARD:
		action = world.ActionMoveRelative{
			Subject_id: subject_id,
			Direction:  world.BACK(),
			Steps:      1,
		}
	case COMMAND_STRAFE_RIGHT:
		action = world.ActionMoveRelative{
			Subject_id: subject_id,
			Direction:  world.RIGHT(),
			Steps:      1,
		}
	}
	return action
}

func CommandsToAction(commands []Command, subject_id world.ActorId) (world.Action, []Command) {
	var action_result world.Action
	commands_result := make([]Command, 0, cap(commands))
	for _, command := range commands {
		action := CommandToAction(command, subject_id)
		if action == nil {
			commands_result = append(commands_result, command)
		} else {
			if action_result == nil {
				// Keep the first action only, the other are discarded.  It should
				// not be a big loss anyway as this function is called every frame.
				// How many keys can you hope to press in 15 milliseconds?
				action_result = action
			} else {
				fmt.Println("Discarded action ", action)
			}
		}
	}
	return action_result, commands_result
}

func LevelCommand(level world.Level, position world.Position, command Command) world.Level {
	here_x, here_y := position.X, position.Y
	there := position.MoveForward(1)
	there_x, there_y := there.X, there.Y
	switch command {
	case COMMAND_PLACE_FLOOR:
		floor := world.MakeFloor(FLOOR_ID, world.EAST(), true)
		level.Floors = level.Floors.Set(there_x, there_y, floor)
	case COMMAND_PLACE_CEILING:
		ceiling := world.MakeOrientedBuilding(CEILING_ID, world.EAST())
		level.Ceilings = level.Ceilings.Set(there_x, there_y, ceiling)
	case COMMAND_PLACE_WALL:
		{
			// If the player faces North, then the wall must face South in order
			// to face the player.
			facing := position.F.Add(world.BACK())
			index := facing.Value()
			wall := world.MakeWall(WALL_ID, false)
			level.Walls[index] = level.Walls[index].Set(here_x, here_y, wall)
		}
	case COMMAND_PLACE_COLUMN:
		level.Columns = level.Columns.Set(there_x, there_y, world.MakeOrientedBuilding(COLUMN_ID, world.EAST()))
	case COMMAND_ROTATE_FLOOR_DIRECT, COMMAND_ROTATE_FLOOR_RETROGRADE:
		{
			building := level.Floors[world.Location{X: there_x, Y: there_y}]
			floor, ok := building.(world.Floor)
			if ok {
				var rel_dir world.RelativeDirection
				if command == COMMAND_ROTATE_FLOOR_DIRECT {
					rel_dir = world.LEFT()
				} else {
					rel_dir = world.RIGHT()
				}
				floor.F = floor.F.Add(rel_dir)
				level.Floors = level.Floors.Set(there_x, there_y, floor)
			} else {
				fmt.Println("You cannot rotate that.")
			}
		}
	case COMMAND_ROTATE_COLUMN_DIRECT, COMMAND_ROTATE_COLUMN_RETROGRADE:
		{
			var rel_dir world.RelativeDirection
			if command == COMMAND_ROTATE_COLUMN_DIRECT {
				rel_dir = world.LEFT()
			} else {
				rel_dir = world.RIGHT()
			}
			column := level.Columns[world.Location{X: there_x, Y: there_y}]
			orientable, ok := column.(world.OrientedBuilding)
			if ok {
				orientable.F = orientable.F.Add(rel_dir)
				level.Columns = level.Columns.Set(there_x, there_y, orientable)
			} else {
				fmt.Println("You cannot rotate that.")
			}
		}
	case COMMAND_ROTATE_CEILING_DIRECT, COMMAND_ROTATE_CEILING_RETROGRADE:
		{
			var rel_dir world.RelativeDirection
			if command == COMMAND_ROTATE_CEILING_DIRECT {
				rel_dir = world.LEFT()
			} else {
				rel_dir = world.RIGHT()
			}
			ceiling := level.Ceilings[world.Location{X: there_x, Y: there_y}]
			orientable, ok := ceiling.(world.OrientedBuilding)
			if ok {
				orientable.F = orientable.F.Add(rel_dir)
				level.Ceilings = level.Ceilings.Set(there_x, there_y, orientable)
			} else {
				fmt.Println("You cannot rotate that.")
			}
		}
	case COMMAND_REMOVE_FLOOR:
		level.Floors = level.Floors.Delete(there_x, there_y)
	case COMMAND_REMOVE_COLUMN:
		level.Columns = level.Columns.Delete(there_x, there_y)
	case COMMAND_REMOVE_CEILING:
		level.Ceilings = level.Ceilings.Delete(there_x, there_y)
	case COMMAND_REMOVE_WALL:
		{
			// If the player faces North, then the wall must face South in order
			// to face the player.
			facing := position.F.Add(world.BACK())
			index := facing.Value()
			level.Walls[index] = level.Walls[index].Delete(here_x, here_y)
		}
	}
	return level
}

func ExecuteCommands(program_state ProgramState, commands []Command) ProgramState {
	for _, command := range commands {
		switch {
		case command < COMMAND_SAVE:
			player_id := program_state.World.Player_id
			player := program_state.World.Level.Actors[player_id]
			position := player.Pos
			program_state.World.Level = LevelCommand(program_state.World.Level, position, command)
		case command == COMMAND_SAVE:
			err := program_state.World.Save()
			fmt.Println("Save:", err)
		case command == COMMAND_LOAD:
			world, err := world.Load()
			fmt.Println("Load:", err)
			if err == nil {
				program_state.World = *world
			}
		}
	}
	return program_state
}

func main() {
	var program_state ProgramState
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
	program_state.Gl.Window, err = glfw.CreateWindow(640, 480, "Daggor", nil, nil)
	if err != nil {
		panic(err)
	}
	defer program_state.Gl.Window.Destroy()

	program_state.Gl.GlfwKeyEventList = MakeGlfwKeyEventList()
	program_state.Gl.Window.SetKeyCallback(program_state.Gl.GlfwKeyEventList.Callback)

	program_state.Gl.Window.MakeContextCurrent()
	ec := gl.Init()
	if ec != 0 {
		panic(fmt.Sprintf("OpenGL initialization failed with code %v.", ec))
	}
	// For some reason, here, the OpenGL error flag for me contains "Invalid enum".
	// This is weird since I have not done anything yet.  I imagine that something
	// goes wrong in gl.Init.  Reading the error flag clears it, so I do it.
	// Here's the reason:
	//     https://github.com/go-gl/glfw3/issues/50
	// Maybe I should not even ask for a core profile anyway.
	// What are the advantages are asking for a core profile?
	err = glw.CheckGlError()
	if err != nil {
		err.(*glw.GlError).Description = "OpenGL has this error right after init for some reason."
		fmt.Println(err)
	}
	major := program_state.Gl.Window.GetAttribute(glfw.ContextVersionMajor)
	minor := program_state.Gl.Window.GetAttribute(glfw.ContextVersionMinor)
	fmt.Printf("OpenGL version %v.%v.\n", major, minor)
	if (major < 3) || (major == 3 && minor < 3) {
		panic("OpenGL version 3.3 required, your video card/driver does not seem to support it.")
	}

	program_state.Gl.Programs = glw.MakePrograms()

	program_state.Gl.Shapes[CUBE_ID] = glw.Cube(program_state.Gl.Programs)
	program_state.Gl.Shapes[PYRAMID_ID] = glw.Pyramid(program_state.Gl.Programs)
	program_state.Gl.Shapes[FLOOR_ID] = glw.Floor(program_state.Gl.Programs)
	program_state.Gl.Shapes[WALL_ID] = glw.Wall(program_state.Gl.Programs)
	program_state.Gl.Shapes[COLUMN_ID] = glw.Column(program_state.Gl.Programs)
	program_state.Gl.Shapes[CEILING_ID] = glw.Ceiling(program_state.Gl.Programs)
	program_state.Gl.DynaPyramid = glw.DynaPyramid(program_state.Gl.Programs)
	program_state.Gl.Monster = glw.Monster(program_state.Gl.Programs)

	// I do not like the default reference frame of OpenGl.
	// By default, we look in the direction -z, and y points up.
	// I want z to point up, and I want to look in the direction +x
	// by default.  That way, I move on an xy plane where z is the
	// altitude, instead of having the altitude stuffed between
	// the two things I use the most.  And my reason for pointing
	// toward +x is that I use the convention for trigonometry:
	// an angle of 0 points to the right (east) of the trigonometric
	// circle.  Bonus point: this matches Blender's reference frame.
	my_frame := glm.ZUP.Mult(glm.RotZ(90))
	program_state.Gl.P = glm.PerspectiveProj(110, 640./480., .1, 100).Mult(my_frame)

	gl.Enable(gl.FRAMEBUFFER_SRGB)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)

	program_state.World = world.MakeWorld()

	MainLoop(program_state)
}

func MainLoop(program_state ProgramState) ProgramState {
	const tick_period = 1000000000 / 60
	ticker := time.NewTicker(tick_period * time.Nanosecond)
	keep_ticking := true
	for keep_ticking {
		select {
		case _, ok := <-ticker.C:
			{
				if ok {
					program_state, keep_ticking = OnTick(program_state, tick_period)
					if !keep_ticking {
						fmt.Println("No more ticks.")
						ticker.Stop()
					}
				} else {
					fmt.Println("Ticker closed, weird.")
					keep_ticking = false
				}
			}
		}
	}
	return program_state
}

func OnTick(program_state ProgramState, dt uint64) (ProgramState, bool) {
	glfw.PollEvents()
	keep_ticking := !program_state.Gl.Window.ShouldClose()
	if keep_ticking {
		// Read raw inputs.
		keys := program_state.Gl.GlfwKeyEventList.Freeze()
		// Analyze the inputs, see what they mean.
		commands := Commands(keys)
		// One of these commands may correspond to an action of the player's actor.
		// We take it out so that we can process it in the IA phase.
		// The remaining commands are kept for further processing.
		player_action, commands := CommandsToAction(commands, program_state.World.Player_id)
		// Evolve the program one step.
		program_state.World.Time += dt // No side effect, we own a copy.
		program_state = ExecuteCommands(program_state, commands)
		if player_action != nil {
			var err error
			program_state.World, err = player_action.Execute(program_state.World)
			if err != nil {
				fmt.Println("PlayerAction error:", err)
			}
		}
		// Render on screen.
		Render(program_state)
		program_state.Gl.Window.SwapBuffers()
	}
	return program_state, keep_ticking
}

func Render(program_state ProgramState) {
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	v := ViewMatrix(program_state.World.Level.Actors[program_state.World.Player_id].Pos)
	vp := (program_state.Gl.P).Mult(v)
	RenderBuildings(
		&program_state.World.Level.Floors,
		0, 0,
		&vp, nil,
		&program_state.Gl,
	)
	RenderBuildings(
		&program_state.World.Level.Ceilings,
		0, 0,
		&vp, nil,
		&program_state.Gl,
	)
	RenderBuildings(
		&program_state.World.Level.Columns,
		-.5, -.5,
		&vp, nil,
		&program_state.Gl,
	)
	for facing := 0; facing < 4; facing++ {
		default_r := glm.RotZ(float64(90 * facing))
		RenderBuildings(
			&program_state.World.Level.Walls[facing],
			0, 0,
			&vp, &default_r,
			&program_state.Gl,
		)
	}

	//// Stupid render of the one dynamic object.
	//dyn := program_state.World.Level.Dynamic
	//m := dyn.ModelMat(program_state.World.Time)
	//mvp := (program_state.Gl.P).Mult(v).Mult(m).Gl()
	//pyr := program_state.Gl.DynaPyramid
	//pyr.UpdateMesh(dyn.Mesh(program_state.World.Time))
	//pyr.Draw(program_state.Gl.Programs, &mvp)
	//// Draw a monster.
	//mvp = (program_state.Gl.P).Mult(v).Gl()
	//program_state.Gl.Monster.Draw(program_state.Gl.Programs, &mvp)
}

func RenderBuildings(
	buildings *world.Buildings,
	offset_x, offset_y float64,
	vp, default_r *glm.Matrix4,
	gl_state *GlState,
) {
	for coords, building := range *buildings {
		m := glm.Vector3{
			float64(coords.X) + offset_x,
			float64(coords.Y) + offset_y,
			0,
		}.Translation()
		facer, ok := building.(world.Facer)
		if ok {
			// We obey the facing of the buildings that have one.
			r := glm.RotZ(float64(90 * facer.Facing().Value()))
			m = m.Mult(r)
		} else {
			// Buildings without facing receive the provided default facing.
			// It is given as a precalculated rotation matrix `default_r`.
			m = m.Mult(*default_r)
		}
		mvp := vp.Mult(m).Gl()
		gl_state.Shapes[building.Model()].Draw(gl_state.Programs, &mvp)
	}
}
