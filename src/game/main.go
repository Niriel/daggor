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

func MaybeMove(level world.Level, position world.Position, rel_dir world.RelativeDirection) (world.Position, bool) {
	wall_passable := true   // By default, no wall is good.
	floor_passable := false // By default, no floor is bad.
	// Direction of the movement relative to facing.
	direction := position.F.Add(rel_dir)
	// Western walls face East.
	wall_facing := rel_dir.Add(world.BACK())
	wall_index := wall_facing.Value()

	building, ok := level.Walls[wall_index].Get(position.X, position.Y)
	if ok {
		wall_passable = building.(world.Wall).IsPassable()
	}

	new_pos := position.SetF(direction).MoveForward(1)
	building, ok = level.Floors.Get(new_pos.X, new_pos.Y)
	if ok {
		floor_passable = building.(world.Floor).IsPassable()
	}

	if wall_passable && floor_passable {
		return new_pos.SetF(position.F), true
	}
	return position, false
}

func PlayerCommand(level world.Level, player world.Player, command Command) world.Player {
	switch command {
	case COMMAND_TURN_LEFT:
		player.Pos = player.Pos.TurnLeft()
	case COMMAND_TURN_RIGHT:
		player.Pos = player.Pos.TurnRight()
	default:
		{
			var rel_dir world.RelativeDirection
			switch command {
			case COMMAND_FORWARD:
				rel_dir = world.FRONT()
			case COMMAND_STRAFE_LEFT:
				rel_dir = world.LEFT()
			case COMMAND_BACKWARD:
				rel_dir = world.BACK()
			case COMMAND_STRAFE_RIGHT:
				rel_dir = world.RIGHT()
			default:
				return player
			}
			new_pos, ok := MaybeMove(level, player.Pos, rel_dir)
			if ok {
				player.Pos = new_pos
			} else {
				fmt.Println("You cannot go that way.")
			}
		}
	}
	return player
}

func LevelCommand(level world.Level, player world.Player, command Command) world.Level {
	here := player.Pos
	here_x, here_y := world.Coord(here.X), world.Coord(here.Y)
	there := here.MoveForward(1)
	x, y := world.Coord(there.X), world.Coord(there.Y)
	switch command {
	case COMMAND_PLACE_FLOOR:
		floor := world.MakeFloor(FLOOR_ID, 0, true)
		level.Floors = level.Floors.Set(x, y, floor)
	case COMMAND_PLACE_CEILING:
		level.Ceilings = level.Ceilings.Set(x, y, world.MakeOrientedBuilding(CEILING_ID, 0))
	case COMMAND_PLACE_WALL:
		{
			// If the player faces North, then the wall must face South in order
			// to face the player.
			facing := player.Pos.F.Add(world.BACK())
			index := facing.Value()
			wall := world.MakeWall(WALL_ID, false)
			level.Walls[index] = level.Walls[index].Set(here_x, here_y, wall)
		}
	case COMMAND_PLACE_COLUMN:
		level.Columns = level.Columns.Set(x, y, world.MakeOrientedBuilding(COLUMN_ID, 0))
	case COMMAND_ROTATE_FLOOR_DIRECT, COMMAND_ROTATE_FLOOR_RETROGRADE:
		{
			var offset int
			if command == COMMAND_ROTATE_FLOOR_DIRECT {
				offset = 1
			} else {
				offset = 3
				// Equivalent to -1 with modulo 4.
				// Because Go's modulo is stupid.
			}
			floor := level.Floors[world.Location{X: x, Y: y}]
			orientable, ok := floor.(world.OrientedBuilding)
			if ok {
				orientable.Facing = (orientable.Facing + offset) % 4
				level.Floors = level.Floors.Set(x, y, orientable)
			} else {
				fmt.Println("You cannot rotate that.")
			}
		}
	case COMMAND_ROTATE_COLUMN_DIRECT, COMMAND_ROTATE_COLUMN_RETROGRADE:
		{
			var offset int
			if command == COMMAND_ROTATE_COLUMN_DIRECT {
				offset = 1
			} else {
				offset = 3
				// Equivalent to -1 with modulo 4.
				// Because Go's modulo is stupid.
			}
			column := level.Columns[world.Location{X: x, Y: y}]
			orientable, ok := column.(world.OrientedBuilding)
			if ok {
				orientable.Facing = (orientable.Facing + offset) % 4
				level.Columns = level.Columns.Set(x, y, orientable)
			} else {
				fmt.Println("You cannot rotate that.")
			}
		}
	case COMMAND_ROTATE_CEILING_DIRECT, COMMAND_ROTATE_CEILING_RETROGRADE:
		{
			var offset int
			if command == COMMAND_ROTATE_CEILING_DIRECT {
				offset = 1
			} else {
				offset = 3
				// Equivalent to -1 with modulo 4.
				// Because Go's modulo is stupid.
			}
			ceiling := level.Ceilings[world.Location{X: x, Y: y}]
			orientable, ok := ceiling.(world.OrientedBuilding)
			if ok {
				orientable.Facing = (orientable.Facing + offset) % 4
				level.Ceilings = level.Ceilings.Set(x, y, orientable)
			} else {
				fmt.Println("You cannot rotate that.")
			}
		}
	case COMMAND_REMOVE_FLOOR:
		level.Floors = level.Floors.Delete(x, y)
	case COMMAND_REMOVE_COLUMN:
		level.Columns = level.Columns.Delete(x, y)
	case COMMAND_REMOVE_CEILING:
		level.Ceilings = level.Ceilings.Delete(x, y)
	case COMMAND_REMOVE_WALL:
		{
			// If the player faces North, then the wall must face South in order
			// to face the player.
			facing := player.Pos.F.Add(world.BACK())
			index := facing.Value()
			level.Walls[index] = level.Walls[index].Delete(here_x, here_y)
		}
	}
	return level
}

func NewProgramState(program_state ProgramState, commands []Command) ProgramState {
	if len(commands) == 0 {
		return program_state
	}
	for _, command := range commands {
		switch {
		case command <= COMMAND_TURN_RIGHT:
			program_state.World.Player = PlayerCommand(program_state.World.Level, program_state.World.Player, command)
		case command < COMMAND_SAVE:
			program_state.World.Level = LevelCommand(program_state.World.Level, program_state.World.Player, command)
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

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)

	program_state.World.Player.Pos = program_state.World.Player.Pos.SetF(world.EAST())

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
		// Evolve the program one step.
		program_state.World.Time += dt
		program_state = NewProgramState(program_state, commands)
		// Render on screen.
		Render(program_state)
		program_state.Gl.Window.SwapBuffers()
	}
	return program_state, keep_ticking
}

func Render(program_state ProgramState) {
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	v := ViewMatrix(program_state.World.Player.Pos)
	for coords, floor := range program_state.World.Level.Floors {
		// The default model matrix only needs a translation.
		m := glm.Vector3{float64(coords.X), float64(coords.Y), 0}.Translation()
		switch floor.(type) {
		case world.OrientedBuilding:
			{
				// But Oriented buildig have a rotation as well.
				facing := floor.(world.OrientedBuilding).Facing
				R := glm.RotZ(float64(90 * facing))
				m = m.Mult(R)
			}
		}
		mvp := (program_state.Gl.P).Mult(v).Mult(m).Gl()
		program_state.Gl.Shapes[floor.Model()].Draw(program_state.Gl.Programs, &mvp)
	}
	for coords, ceiling := range program_state.World.Level.Ceilings {
		// The default model matrix only needs a translation.
		m := glm.Vector3{float64(coords.X), float64(coords.Y), 0}.Translation()
		switch ceiling.(type) {
		case world.OrientedBuilding:
			{
				// But Oriented buildig have a rotation as well.
				facing := ceiling.(world.OrientedBuilding).Facing
				R := glm.RotZ(float64(90 * facing))
				m = m.Mult(R)
			}
		}
		mvp := (program_state.Gl.P).Mult(v).Mult(m).Gl()
		program_state.Gl.Shapes[ceiling.Model()].Draw(program_state.Gl.Programs, &mvp)
	}
	for facing := 0; facing < 4; facing++ {
		R := glm.RotZ(float64(90 * facing))
		for coords, wall := range program_state.World.Level.Walls[facing] {
			m := glm.Vector3{float64(coords.X), float64(coords.Y), 0}.Translation()
			m = m.Mult(R)
			mvp := (program_state.Gl.P).Mult(v).Mult(m).Gl()
			program_state.Gl.Shapes[wall.Model()].Draw(program_state.Gl.Programs, &mvp)
		}
	}
	for coords, column := range program_state.World.Level.Columns {
		// The default model matrix only needs a translation.
		m := glm.Vector3{float64(coords.X) - .5, float64(coords.Y) - .5, 0}.Translation()
		switch column.(type) {
		case world.OrientedBuilding:
			{
				// But Oriented buildig have a rotation as well.
				facing := column.(world.OrientedBuilding).Facing
				R := glm.RotZ(float64(90 * facing))
				m = m.Mult(R)
			}
		}
		mvp := (program_state.Gl.P).Mult(v).Mult(m).Gl()
		program_state.Gl.Shapes[column.Model()].Draw(program_state.Gl.Programs, &mvp)
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
