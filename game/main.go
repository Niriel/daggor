// game project main.go
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/niriel/daggor/glm"
	"runtime"
	"unsafe"
)

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

func (self *Drawable) Draw(mvp_matrix *[16]float32) {
	// Bindind the VAO each time is not efficient but
	// it is correct.
	self.vao.Bind()
	self.program.Use()
	self.mvp.UniformMatrix4f(false, mvp_matrix)
	gl.DrawElements(gl.TRIANGLE_STRIP, self.n_elements, gl.UNSIGNED_BYTE, nil)
}

func CubeMesh() Drawable {
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

func PyramidMesh() Drawable {
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
)

func Commands(events []GlfwKeyEvent) []Command {
	if len(events) == 0 {
		return nil
	}
	result := make([]Command, 0, 1)
	for _, event := range events {
		switch event.key {
		case glfw.KeyW:
			if event.action == glfw.Press {
				result = append(result, COMMAND_FORWARD)
			}
		case glfw.KeyS:
			if event.action == glfw.Press {
				result = append(result, COMMAND_BACKWARD)
			}
		}
	}
	return result
}

type ProgramState struct {
	PlayerPos int
}

func NewPlayerPos(program_state ProgramState, command Command) int {
	delta := 0
	switch command {
	case COMMAND_BACKWARD:
		delta = -1
	case COMMAND_FORWARD:
		delta = 1
	}
	return program_state.PlayerPos + delta
}

func NewProgramState(program_state ProgramState, commands []Command) ProgramState {
	if len(commands) == 0 {
		return program_state
	}
	new_state := program_state // Copy.
	for _, command := range commands {
		new_state.PlayerPos = NewPlayerPos(new_state, command)
	}
	return new_state
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

	cube := CubeMesh()
	pyramid := PyramidMesh()

	shapes := [...]Drawable{
		cube, cube, pyramid, cube, pyramid, pyramid, cube,
	}
	positions := [len(shapes)]glm.Vector3{
		glm.Vector3{0, 4, 0},
		glm.Vector3{1, 3, 0},
		glm.Vector3{-1, 3, 0},
		glm.Vector3{0, 5, 2},
		glm.Vector3{-3, 1, 0},
		glm.Vector3{2, 2, -1},
		glm.Vector3{7, 3, 0},
	}

	p := glm.PerspectiveProj(80, 640./480., .1, 100).Mult(glm.ZUP)

	gl.Enable(gl.CULL_FACE)
	//gl.CullFace(gl.BACK)
	//gl.FrontFace(gl.CCW)

	var program_state ProgramState
	for !window.ShouldClose() {
		keys := glfwKeyEventList.Freeze()
		commands := Commands(keys)
		program_state = NewProgramState(program_state, commands)
		v := glm.Vector3{float64(program_state.PlayerPos), 0, 0}.Translation()
		fmt.Println(program_state)
		gl.ClearColor(0.0, 0.0, 0.4, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		for i := 0; i < len(shapes); i++ {
			shape := shapes[i]
			position := positions[i]
			m := position.Translation()
			mvp := p.Mult(v).Mult(m).Gl()
			shape.Draw(&mvp)
		}
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
