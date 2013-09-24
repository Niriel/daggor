package glw

import (
	"fmt"
	"github.com/go-gl/gl"
)

type Drawable struct {
	vao          gl.VertexArray
	mvp          gl.UniformLocation
	shaders_refs ShaderRefs
	n_elements   int
}

func (self *Drawable) Draw(programs Programs, mvp_matrix *[16]float32) {
	// Bindind the VAO each time is not efficient but
	// it is correct.
	program, err := programs.Serve(self.shaders_refs)
	if err != nil {
		panic(err)
	}
	self.vao.Bind()
	program.Use()
	self.mvp.UniformMatrix4f(false, mvp_matrix)

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

	gl.DrawElements(gl.TRIANGLE_STRIP, self.n_elements, gl.UNSIGNED_BYTE, nil)
}
