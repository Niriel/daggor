// game project main.go
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
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

	vao := gl.GenVertexArray()
	vao.Bind()
	vertices := [...]gl.GLfloat{
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		0.0, 1.0, 0.0,
	}
	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)

	vsh_code := `
	#version 330 core
	in vec3 vpos;
	void main(){
		gl_Position.xyz = vpos;
        gl_Position.w = 1.0;
	}
	`
	fsh_code := `
	#version 330 core
    out vec3 color;

    void main(){
        color = vec3(1,0,0);
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

	var program_state ProgramState
	for !window.ShouldClose() {
		keys := glfwKeyEventList.Freeze()
		commands := Commands(keys)
		program_state = NewProgramState(program_state, commands)
		fmt.Println(program_state)
		gl.ClearColor(0.0, 0.0, 0.4, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
