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
)

type Command int

const (
	COMMAND_FORWARD = Command(iota)
	COMMAND_BACKWARD
	COMMAND_STRAFE_LEFT
	COMMAND_STRAFE_RIGHT
	COMMAND_TURN_LEFT
	COMMAND_TURN_RIGHT
	COMMAND_PLACE_CUBE
	COMMAND_PLACE_PYRAMID
	COMMAND_PLACE_FLOOR
	COMMAND_PLACE_WALL
	COMMAND_ROTATE_SHAPE_DIRECT
	COMMAND_ROTATE_SHAPE_RETROGRADE
	COMMAND_REMOVE_SHAPE
	COMMAND_REMOVE_WALL
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
			case glfw.KeyC:
				result = append(result, COMMAND_PLACE_CUBE)
			case glfw.KeyP:
				result = append(result, COMMAND_PLACE_PYRAMID)
			case glfw.KeyF:
				result = append(result, COMMAND_PLACE_FLOOR)
			case glfw.KeyR:
				result = append(result, COMMAND_PLACE_WALL)
			case glfw.KeyDelete:
				result = append(result, COMMAND_REMOVE_SHAPE)
			case glfw.KeyBackspace:
				result = append(result, COMMAND_REMOVE_WALL)
			case glfw.KeyLeftBracket:
				result = append(result, COMMAND_ROTATE_SHAPE_DIRECT)
			case glfw.KeyRightBracket:
				result = append(result, COMMAND_ROTATE_SHAPE_RETROGRADE)
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
	R := glm.RotZ(float64(-90 * pos.F))
	T := glm.Vector3{float64(-pos.X), float64(-pos.Y), -.6}.Translation()
	return R.Mult(T)
}

type GlState struct {
	Window           *glfw.Window
	GlfwKeyEventList *GlfwKeyEventList
	Programs         glw.Programs
	P                glm.Matrix4
	Shapes           [4]glw.Drawable
}

type ProgramState struct {
	Gl    GlState     // Highly mutable, impure.
	World world.World // Immutable, pure.
}

func PlayerCommand(player world.Player, command Command) world.Player {
	switch command {
	case COMMAND_TURN_LEFT:
		player.Pos = player.Pos.TurnLeft()
	case COMMAND_TURN_RIGHT:
		player.Pos = player.Pos.TurnRight()
	case COMMAND_BACKWARD:
		player.Pos = player.Pos.Backward()
	case COMMAND_FORWARD:
		player.Pos = player.Pos.Forward()
	case COMMAND_STRAFE_LEFT:
		player.Pos = player.Pos.StrafeLeft()
	case COMMAND_STRAFE_RIGHT:
		player.Pos = player.Pos.StrafeRight()
	}
	return player
}

func LevelCommand(level world.Level, player world.Player, command Command) world.Level {
	here := player.Pos
	here_x, here_y := world.Coord(here.X), world.Coord(here.Y)
	there := here.Forward()
	x, y := world.Coord(there.X), world.Coord(there.Y)
	switch command {
	case COMMAND_PLACE_CUBE:
		level.Floors = level.Floors.Set(x, y, world.MakeBaseBuilding(CUBE_ID))
	case COMMAND_PLACE_PYRAMID:
		level.Floors = level.Floors.Set(x, y, world.MakeBaseBuilding(PYRAMID_ID))
	case COMMAND_PLACE_FLOOR:
		level.Floors = level.Floors.Set(x, y, world.MakeOrientedBuilding(FLOOR_ID, 0))
	case COMMAND_PLACE_WALL:
		{
			// If the player faces North, then the wall must face South in order
			// to face the player.
			facing := (player.Pos.F + 2) % 4
			level.Walls[facing] = level.Walls[facing].Set(here_x, here_y, world.MakeBaseBuilding(WALL_ID))
		}
	case COMMAND_ROTATE_SHAPE_DIRECT, COMMAND_ROTATE_SHAPE_RETROGRADE:
		{
			var offset int
			if command == COMMAND_ROTATE_SHAPE_DIRECT {
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
	case COMMAND_REMOVE_SHAPE:
		level.Floors = level.Floors.Delete(x, y)
	case COMMAND_REMOVE_WALL:
		{
			// If the player faces North, then the wall must face South in order
			// to face the player.
			fmt.Println("wwewew")
			facing := (player.Pos.F + 2) % 4
			level.Walls[facing] = level.Walls[facing].Delete(here_x, here_y)
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
			program_state.World.Player = PlayerCommand(program_state.World.Player, command)
		case command <= COMMAND_REMOVE_WALL:
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
	glfw.WindowHint(glfw.SrgbCapable, 1)
	glfw.WindowHint(glfw.Resizable, 0)
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
	err = glw.CheckGlError()
	if err != nil {
		err.(*glw.GlError).Description = "OpenGL has this error right after init for some reason."
		fmt.Println(err)
	}

	program_state.Gl.Programs = glw.MakePrograms()

	program_state.Gl.Shapes[CUBE_ID] = glw.Cube(program_state.Gl.Programs)
	program_state.Gl.Shapes[PYRAMID_ID] = glw.Pyramid(program_state.Gl.Programs)
	program_state.Gl.Shapes[FLOOR_ID] = glw.Floor(program_state.Gl.Programs)
	program_state.Gl.Shapes[WALL_ID] = glw.Wall(program_state.Gl.Programs)
	tiles := map[[2]int]world.ModelId{
		[2]int{0, 4}: CUBE_ID,
		[2]int{1, 3}: CUBE_ID,
		[2]int{0, 3}: PYRAMID_ID,
		[2]int{0, 5}: CUBE_ID,
		[2]int{1, 6}: PYRAMID_ID,
		[2]int{2, 2}: PYRAMID_ID,
		[2]int{7, 3}: CUBE_ID,
	}
	for coords, floor := range tiles {
		x, y := world.Coord(coords[0]), world.Coord(coords[1])
		program_state.World.Level.Floors = program_state.World.Level.Floors.Set(x, y, world.MakeBaseBuilding(floor))
	}

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
	MainLoop(program_state)
}

func MainLoop(program_state ProgramState) ProgramState {
	ticker := time.NewTicker(15 * time.Millisecond)
	keep_ticking := true
	for keep_ticking {
		select {
		case _, ok := <-ticker.C:
			{
				if ok {
					program_state, keep_ticking = OnTick(program_state)
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

func OnTick(program_state ProgramState) (ProgramState, bool) {
	glfw.PollEvents()
	keep_ticking := !program_state.Gl.Window.ShouldClose()
	if keep_ticking {
		// Read raw inputs.
		keys := program_state.Gl.GlfwKeyEventList.Freeze()
		// Analyze the inputs, see what they mean.
		commands := Commands(keys)
		// Evolve the program one step.
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
		program_state.Gl.Shapes[floor.Model()].Draw(&mvp)
	}
	for facing := 0; facing < 4; facing++ {
		R := glm.RotZ(float64(90 * facing))
		for coords, wall := range program_state.World.Level.Walls[facing] {
			m := glm.Vector3{float64(coords.X), float64(coords.Y), 0}.Translation()
			m = m.Mult(R)
			mvp := (program_state.Gl.P).Mult(v).Mult(m).Gl()
			program_state.Gl.Shapes[wall.Model()].Draw(&mvp)
		}
	}
}
