package glw

import (
	"fmt"
	"github.com/go-gl/gl"
)

type Drawable struct {
	vao                  gl.VertexArray
	model_matrix_uniform gl.UniformLocation
	shaders_refs         ShaderRefs
	drawer               Drawer
}

func (self *Drawable) Draw(programs Programs, model_matrix *[16]float32) {
	// Bindind the VAO each time is not efficient but
	// it is correct.
	program, err := programs.Serve(self.shaders_refs)
	if err != nil {
		panic(err)
	}
	self.vao.Bind()
	program.Use()
	self.model_matrix_uniform.UniformMatrix4f(false, model_matrix)

	program.Validate()
	if err := CheckGlError(); err != nil {
		err.Description = "Program.Validate failed."
		panic(err)
	}
	status := program.Get(gl.VALIDATE_STATUS)
	if err := CheckGlError(); err != nil {
		err.Description = "Program.Get(VALIDATE_STATUS) failed."
		panic(err)
	}
	if status == gl.FALSE {
		infolog := program.GetInfoLog()
		gl.GetError() // Clear error flag if infolog derped.
		panic(fmt.Errorf("Program validation failed. Log: %v", infolog))
	}
	self.drawer.Draw()
	//gl.DrawElements(self.primitive, self.n_elements, gl.UNSIGNED_BYTE, nil)
}

type StreamDrawable struct {
	Drawable
	vbo     gl.Buffer
	vbosize int
}

func (self *StreamDrawable) UpdateMesh(mesh []float64) {
	data := make([]gl.GLfloat, len(mesh), len(mesh))
	for i, value := range mesh {
		data[i] = gl.GLfloat(value)
	}
	self.vao.Bind()
	self.vbo.Bind(gl.ARRAY_BUFFER)
	// I don't know what I'm doing.  Setting the data to nil does some magic.
	// It is knows as "orphaning".
	gl.BufferData(gl.ARRAY_BUFFER, self.vbosize, nil, gl.DYNAMIC_DRAW)
	gl.BufferData(gl.ARRAY_BUFFER, self.vbosize, &data[0], gl.DYNAMIC_DRAW)
	self.vbo.Unbind(gl.ARRAY_BUFFER)
}
