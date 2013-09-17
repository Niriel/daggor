// game project main.go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"glm" // Local import, make sure Daggor is in your gopath.
	"runtime"
	"time"
	"unsafe"
)

var HARDDRIVE bytes.Buffer

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

const EMPTY_ID = 0
const CUBE_ID = 1
const PYRAMID_ID = 2

func (self *Drawable) Draw(mvp_matrix *[16]float32) {
	// Bindind the VAO each time is not efficient but
	// it is correct.
	self.vao.Bind()
	self.program.Use()
	self.mvp.UniformMatrix4f(false, mvp_matrix)
	gl.DrawElements(gl.TRIANGLE_STRIP, self.n_elements, gl.UNSIGNED_BYTE, nil)
}

func Cube() Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := [...]gl.GLfloat{
		m, m, m,
		m, m, p,
		m, p, m,
		m, p, p,
		p, m, m,
		p, m, p,
		p, p, m,
		p, p, p,
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
	vsh_code := `
	#version 330 core
	in vec3 vpos;
	out vec3 fpos;
	uniform mat4 mvp;
	void main(){
		gl_Position = mvp * vec4(vpos, 1.0);
		fpos = gl_Position.xyz;
	}
	`
	fsh_code := `
	#version 330 core
    in vec3 fpos;
    out vec3 color;

    void main(){
        color = vec3(1.0 - fpos.z *.1, 0, 0);
    }
	`
	vsh := gl.CreateShader(gl.VERTEX_SHADER)
	vsh.Source(vsh_code)
	vsh.Compile()
	fsh := gl.CreateShader(gl.FRAGMENT_SHADER)
	fsh.Source(fsh_code)
	fsh.Compile()
	program := gl.CreateProgram()
	program.AttachShader(vsh)
	program.AttachShader(fsh)
	program.Link()
	program.Use()
	fmt.Println(program.GetInfoLog())

	att := program.GetAttribLocation("vpos")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		nil)
	mvp := program.GetUniformLocation("mvp")
	return Drawable{vao, mvp, program, len(indices)}
}

func Pyramid() Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := [...]gl.GLfloat{
		m, m, m,
		m, p, m,
		p, m, m,
		p, p, m,
		0, 0, p,
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
	vsh_code := `
	#version 330 core
	in vec3 vpos;
	out vec3 fpos;
	uniform mat4 mvp;
	void main(){
		gl_Position = mvp * vec4(vpos, 1.0);
		fpos = gl_Position.xyz;
	}
	`
	fsh_code := `
	#version 330 core
	in vec3 fpos;
    out vec3 color;

    void main(){
        color = vec3(0,1.0 - fpos.z *.1, 0);
    }
	`
	vsh := gl.CreateShader(gl.VERTEX_SHADER)
	vsh.Source(vsh_code)
	vsh.Compile()
	fsh := gl.CreateShader(gl.FRAGMENT_SHADER)
	fsh.Source(fsh_code)
	fsh.Compile()
	program := gl.CreateProgram()
	program.AttachShader(vsh)
	program.AttachShader(fsh)
	program.Link()
	program.Use()
	fmt.Println(program.GetInfoLog())

	att := program.GetAttribLocation("vpos")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		nil)
	mvp := program.GetUniformLocation("mvp")
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
	T := glm.Vector3{float64(-self.X), float64(-self.Y), 0}.Translation()
	return R.Mult(T)
}

const LANDSCAPE_SIZE = 16

type Landscape struct {
	Tiles [LANDSCAPE_SIZE * LANDSCAPE_SIZE]int
}

func (self *Landscape) Tile(x, y int) int {
	if (x < 0) || (x >= LANDSCAPE_SIZE) {
		panic("Landscape x index out of range.")
	}
	if (y < 0) || (y >= LANDSCAPE_SIZE) {
		panic("Landscape y index out of range.")
	}
	return self.Tiles[y*LANDSCAPE_SIZE+x]
}

func (self Landscape) SetTile(x, y int, shape_id int) Landscape {
	if (x < 0) || (x >= LANDSCAPE_SIZE) {
		panic("Landscape x index out of range.")
	}
	if (y < 0) || (y >= LANDSCAPE_SIZE) {
		panic("Landscape y index out of range.")
	}
	self.Tiles[y*LANDSCAPE_SIZE+x] = shape_id
	return self
}

func (self Landscape) SetTiles(tiles map[[2]int]int) Landscape {
	// Allows to set a bunch of tiles at once, without having
	// to create a bunch of intermediary landscapes.
	for coords, shape_id := range tiles {
		x, y := coords[0], coords[1]
		if (x < 0) || (x >= LANDSCAPE_SIZE) {
			panic("Landscape x index out of range.")
		}
		if (y < 0) || (y >= LANDSCAPE_SIZE) {
			panic("Landscape y index out of range.")
		}
		self.Tiles[y*LANDSCAPE_SIZE+x] = shape_id
	}
	return self
}

type World struct {
	Player    Player
	Landscape Landscape
}

type ProgramState struct {
	Shapes [2]Drawable
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

func LandscapeCommand(landscape Landscape, player Player, command Command) Landscape {
	where := player.Forward()
	x, y := where.X, where.Y
	if x >= 0 && x < LANDSCAPE_SIZE && y >= 0 && y < LANDSCAPE_SIZE {
		switch command {
		case COMMAND_PLACE_CUBE:
			landscape = landscape.SetTile(x, y, CUBE_ID)
		case COMMAND_PLACE_PYRAMID:
			landscape = landscape.SetTile(x, y, PYRAMID_ID)
		case COMMAND_REMOVE_SHAPE:
			landscape = landscape.SetTile(x, y, EMPTY_ID)
		}
	}
	return landscape
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
			program_state.World.Landscape = LandscapeCommand(program_state.World.Landscape, program_state.World.Player, command)
		case command == COMMAND_SAVE:
			err := Save(program_state.World)
			fmt.Println("Save:", err)
		case command == COMMAND_LOAD:
			world, err := Load()
			fmt.Println("Load:", err)
			if err == nil {
				program_state.World = world
			}
		}
	}
	return program_state
}

func Save(world World) error {
	encoder := gob.NewEncoder(&HARDDRIVE)
	return encoder.Encode(world)
}

func Load() (World, error) {
	var world World
	decoder := gob.NewDecoder(&HARDDRIVE)
	err := decoder.Decode(&world)
	return world, err
}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
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
	fmt.Println("OpenGL error code", ec)

	var program_state ProgramState

	program_state.Shapes[CUBE_ID-1] = Cube()
	program_state.Shapes[PYRAMID_ID-1] = Pyramid()
	tiles := map[[2]int]int{
		[2]int{0, 4}: CUBE_ID,
		[2]int{1, 3}: CUBE_ID,
		[2]int{0, 3}: PYRAMID_ID,
		[2]int{0, 5}: CUBE_ID,
		[2]int{1, 6}: PYRAMID_ID,
		[2]int{2, 2}: PYRAMID_ID,
		[2]int{7, 3}: CUBE_ID,
	}
	program_state.World.Landscape = program_state.World.Landscape.SetTiles(tiles)

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
	p := glm.PerspectiveProj(80, 640./480., .1, 100).Mult(my_frame)

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
		for x := 0; x < LANDSCAPE_SIZE; x++ {
			for y := 0; y < LANDSCAPE_SIZE; y++ {
				shape_id := program_state.World.Landscape.Tile(x, y)
				if shape_id != EMPTY_ID {
					m := glm.Vector3{float64(x), float64(y), 0}.Translation()
					mvp := p.Mult(v).Mult(m).Gl()
					program_state.Shapes[shape_id-1].Draw(&mvp)
				}
			}
		}
		window.SwapBuffers()
		time.Sleep(15 * time.Millisecond)
		glfw.PollEvents()
	}
}
