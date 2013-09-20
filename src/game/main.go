// game project main.go
package main

import (
	"encoding/gob"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"glm" // Local import, make sure Daggor is in your gopath.
	"glw"
	"os"
	"runtime"
	"time"
	"unsafe"
	"world"
)

var Programs glw.Programs

func init() {
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

func MakeGlfwKeyEventList() GlfwKeyEventList {
	return GlfwKeyEventList{
		make([]GlfwKeyEvent, 0, EVENT_LIST_CAP),
	}
}

func (self *GlfwKeyEventList) Freeze() []GlfwKeyEvent {
	// The list of key events is double buffered.  This allows the application
	// to process events during a frame without having to worry about new
	// events arriving and growing the list.
	result := self.list
	self.list = make([]GlfwKeyEvent, 0, EVENT_LIST_CAP)
	return result
}

func (self *GlfwKeyEventList) Callback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	self.list = append(self.list, GlfwKeyEvent{key, scancode, action, mods})
}

type Drawable struct {
	vao        gl.VertexArray
	mvp        gl.UniformLocation
	program    gl.Program
	n_elements int
}

const (
	CUBE_ID = iota
	PYRAMID_ID
	FLOOR_ID
)

func (self *Drawable) Draw(mvp_matrix *[16]float32) {
	// Bindind the VAO each time is not efficient but
	// it is correct.
	self.vao.Bind()
	self.program.Use()
	self.mvp.UniformMatrix4f(false, mvp_matrix)

	self.program.Validate()
	if err := glw.CheckGlError(); err != nil {
		err.Description = "Program.Validate failed."
		panic(err)
	}
	status := self.program.Get(gl.VALIDATE_STATUS)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "Program.Get(VALIDATE_STATUS) failed."
		panic(err)
	}
	if status == gl.FALSE {
		infolog := self.program.GetInfoLog()
		gl.GetError() // Clear error flag if infolog derped.
		panic(fmt.Errorf("Program validation failed. Log: %v", infolog))
	}

	gl.DrawElements(gl.TRIANGLE_STRIP, self.n_elements, gl.UNSIGNED_BYTE, nil)
}

func Cube() Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := [...]gl.GLfloat{
		m, m, 0,
		m, m, 1,
		m, p, 0,
		m, p, 1,
		p, m, 0,
		p, m, 1,
		p, p, 0,
		p, p, 1,
	}
	// Indices for triangle strip adapted from
	// http://www.cs.umd.edu/gvil/papers/av_ts.pdf .
	// I mirrored their cube to have CCW, and I used a natural order to
	// number the vertices (see above, it's binary code).
	indices := [...]gl.GLubyte{
		6, 2, 7, 3, 1, 2, 0, 6, 4, 7, 5, 1, 4, 0,
	}
	vao := gl.GenVertexArray()
	vao.Bind()
	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)
	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices)), &indices, gl.STATIC_DRAW)

	srefs := glw.ShaderRefs{glw.VSH_POS3, glw.FSH_ZRED}
	program, err := Programs.Serve(srefs)
	if err != nil {
		panic(err)
	}

	program.Use()

	att := program.GetAttribLocation("vpos")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		nil)
	mvp := program.GetUniformLocation("mvp")
	vbuf.Unbind(gl.ARRAY_BUFFER)
	return Drawable{vao, mvp, program, len(indices)}
}

func Pyramid() Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := [...]gl.GLfloat{
		m, m, 0,
		m, p, 0,
		p, m, 0,
		p, p, 0,
		0, 0, 1,
	}
	indices := [...]gl.GLubyte{
		1, 4, 3, 2, 1, 0, 4, 2,
	}
	vao := gl.GenVertexArray()
	vao.Bind()
	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)

	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices)), &indices, gl.STATIC_DRAW)

	srefs := glw.ShaderRefs{glw.VSH_POS3, glw.FSH_ZGREEN}
	program, err := Programs.Serve(srefs)
	if err != nil {
		panic(err)
	}

	att := program.GetAttribLocation("vpos")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		nil)
	mvp := program.GetUniformLocation("mvp")
	vbuf.Unbind(gl.ARRAY_BUFFER)
	return Drawable{vao, mvp, program, len(indices)}
}

func Floor() Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := [...]gl.GLfloat{
		// x y z r v b
		m, m, 0, .1, .1, .5,
		m, p, 0, .1, .1, .5,
		p, m, 0, 0, 1, 0,
		p, p, 0, 1, 0, 0,
	}
	indices := [...]gl.GLubyte{
		0, 2, 1, 3,
	}
	vao := gl.GenVertexArray()
	vao.Bind()

	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)

	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices)), &indices, gl.STATIC_DRAW)

	srefs := glw.ShaderRefs{glw.VSH_COL3, glw.FSH_VCOL}
	program, err := Programs.Serve(srefs)
	if err != nil {
		panic(err)
	}

	att := program.GetAttribLocation("vpos")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		6*4,
		uintptr(0))
	att = program.GetAttribLocation("vcol")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		6*4,
		uintptr(3*4))
	mvp := program.GetUniformLocation("mvp")

	vbuf.Unbind(gl.ARRAY_BUFFER)
	return Drawable{vao, mvp, program, len(indices)}
}

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
	COMMAND_REMOVE_SHAPE
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
			case glfw.KeyDelete:
				result = append(result, COMMAND_REMOVE_SHAPE)
			case glfw.KeyF4:
				result = append(result, COMMAND_SAVE)
			case glfw.KeyF5:
				result = append(result, COMMAND_LOAD)
			}
		}
	}
	return result
}

var SIN = [...]int{0, 1, 0, -1}
var COS = [...]int{1, 0, -1, 0}

type Player struct {
	X int // Position.
	Y int // Position.
	F int // Facing.
}

func (self Player) TurnLeft() Player {
	self.F = (self.F + 1) % 4
	return self
}
func (self Player) TurnRight() Player {
	// Here I add 3, because if I subtract 1 I get the stupid
	// go result: -1 % 4 = -1 (go) instead of -1 % 4 = 3 (python).
	// See this discussion:
	// https://code.google.com/p/go/issues/detail?id=448
	self.F = (self.F + 3) % 4
	return self
}
func (self Player) Forward() Player {
	self.X += COS[self.F]
	self.Y += SIN[self.F]
	return self
}
func (self Player) Backward() Player {
	self.X -= COS[self.F]
	self.Y -= SIN[self.F]
	return self
}
func (self Player) StrafeLeft() Player {
	self.X -= SIN[self.F]
	self.Y += COS[self.F]
	return self
}
func (self Player) StrafeRight() Player {
	self.X += SIN[self.F]
	self.Y -= COS[self.F]
	return self
}
func (self Player) ViewMatrix() glm.Matrix4 {
	R := glm.RotZ(float64(-90 * self.F))
	T := glm.Vector3{float64(-self.X), float64(-self.Y), -.6}.Translation()
	return R.Mult(T)
}

type World struct {
	Player Player
	Level  world.Level
}

type ProgramState struct {
	Shapes [3]Drawable
	World  World
}

func PlayerCommand(player Player, command Command) Player {
	switch command {
	case COMMAND_TURN_LEFT:
		return player.TurnLeft()
	case COMMAND_TURN_RIGHT:
		return player.TurnRight()
	case COMMAND_BACKWARD:
		return player.Backward()
	case COMMAND_FORWARD:
		return player.Forward()
	case COMMAND_STRAFE_LEFT:
		return player.StrafeLeft()
	case COMMAND_STRAFE_RIGHT:
		return player.StrafeRight()
	}
	return player
}

func LevelCommand(level world.Level, player Player, command Command) world.Level {
	where := player.Forward()
	x, y := world.Coord(where.X), world.Coord(where.Y)
	switch command {
	case COMMAND_PLACE_CUBE:
		level.Floors = level.Floors.Set(x, y, world.MakeBaseBuilding(CUBE_ID))
	case COMMAND_PLACE_PYRAMID:
		level.Floors = level.Floors.Set(x, y, world.MakeBaseBuilding(PYRAMID_ID))
	case COMMAND_PLACE_FLOOR:
		level.Floors = level.Floors.Set(x, y, world.MakeOrientedBuilding(FLOOR_ID, 1))
	case COMMAND_REMOVE_SHAPE:
		level.Floors = level.Floors.Delete(x, y)
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
		case command <= COMMAND_REMOVE_SHAPE:
			program_state.World.Level = LevelCommand(program_state.World.Level, program_state.World.Player, command)
		case command == COMMAND_SAVE:
			err := Save(program_state.World)
			fmt.Println("Save:", err)
		case command == COMMAND_LOAD:
			world, err := Load()
			fmt.Println("Load:", err)
			if err == nil {
				program_state.World = *world
			}
		}
	}
	return program_state
}

func Save(world World) error {
	f, err := os.Create("quicksave.sav")
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(f)
	return encoder.Encode(world)
}

func Load() (*World, error) {
	f, err := os.Open("quicksave.sav")
	if err != nil {
		return nil, err
	}
	var world World
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&world)
	return &world, err
}

func main() {
	gob.Register(world.MakeBaseBuilding(0))
	gob.Register(world.MakeOrientedBuilding(0, 0))
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("GLFW initialization failed.")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.SrgbCapable, 1)
	glfw.WindowHint(glfw.Resizable, 0)
	window, err := glfw.CreateWindow(640, 480, "Daggor", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	glfwKeyEventList := MakeGlfwKeyEventList()
	window.SetKeyCallback(glfwKeyEventList.Callback)

	window.MakeContextCurrent()
	ec := gl.Init()
	if ec != 0 {
		panic(fmt.Sprintf("OpenGL initialization failed with code %v.", ec))
	}
	// For some reason, here, the OpenGL error flag for me contains "Invalid enum".
	// This is weird since I have not done anything yet.  I imagine that something
	// goes wrong in gl.Init.  Reading the error flag clears it, so I do it.
	gl.GetError()

	var program_state ProgramState
	Programs = glw.MakePrograms()

	program_state.Shapes[CUBE_ID] = Cube()
	program_state.Shapes[PYRAMID_ID] = Pyramid()
	program_state.Shapes[FLOOR_ID] = Floor()
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
	p := glm.PerspectiveProj(110, 640./480., .1, 100).Mult(my_frame)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
	for !window.ShouldClose() {
		keys := glfwKeyEventList.Freeze()
		commands := Commands(keys)
		program_state = NewProgramState(program_state, commands)
		v := program_state.World.Player.ViewMatrix()
		gl.ClearColor(0.0, 0.0, 0.4, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
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
			mvp := p.Mult(v).Mult(m).Gl()
			program_state.Shapes[floor.Model()].Draw(&mvp)
		}
		window.SwapBuffers()
		time.Sleep(15 * time.Millisecond)
		glfw.PollEvents()
	}
}
