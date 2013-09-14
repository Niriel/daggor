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

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	//glfw.WindowHint(glfw.ContextVersionMajor, 3)
	//glfw.WindowHint(glfw.ContextVersionMinor, 3)
	//glfw.WindowHint(glfw.SrgbCapable, 1)
	window, err := glfw.CreateWindow(640, 480, "Daggor", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

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

	for !window.ShouldClose() {
		gl.ClearColor(0.0, 0.0, 0.4, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
