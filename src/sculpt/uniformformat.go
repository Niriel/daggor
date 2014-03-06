package sculpt

import (
	"fmt"
	"github.com/go-gl/gl"
	"glm"
	"glw"
)

type Uniforms interface {
	SetLocation(glm.Matrix4)
	SetGl()
	SetUpVao(program gl.Program)
}

//----------------------------------------------------------------------------

// To use where there is no uniform to deal with.
type UniformsNone struct{}

func (unif *UniformsNone) SetUpVao(program gl.Program) {}
func (unif *UniformsNone) SetGl()                      {}

//----------------------------------------------------------------------------
type UniformsLoc struct {
	// Filled when calling the SetUpVao method with a program.
	globalMatricesUbi gl.UniformBlockIndex
	modelLoc          gl.UniformLocation
	// Filled when updating the uniform values.
	location glm.Matrix4
}

func (unif *UniformsLoc) SetUpVao(program gl.Program) {
	const modelLocName = "view_matrix"
	const matricesUbiName = "GlobalMatrices"

	// The Model transformation matrix.
	// Soon to be replaced by a ModelView, as this reduces rounding errors.
	unif.modelLoc = program.GetUniformLocation(modelLocName)
	if err := glw.CheckGlError(); err != nil {
		err.Description = fmt.Sprintf("program.GetUniformLocation(%#v)", modelLocName)
		panic(err)
	}
	if unif.modelLoc == -1 {
		panic(fmt.Sprintf("uniform %#v not found", modelLocName))
	}

	// Uniform Block for the View and Projection matrices.
	unif.globalMatricesUbi = program.GetUniformBlockIndex(matricesUbiName)
	if err := glw.CheckGlError(); err != nil {
		err.Description = fmt.Sprintf("GetUniformBlockIndex(%#v)", matricesUbiName)
		panic(err)
	}
	if unif.globalMatricesUbi == gl.INVALID_INDEX {
		panic(fmt.Sprintf("GetUniformBlockIndex(%#v) returned INVALID_INDEX", matricesUbiName))
	}

	program.UniformBlockBinding(
		unif.globalMatricesUbi,
		glw.CameraUboBindingPoint)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "UniformBlockBinding"
		panic(err)
	}
}

func (unif *UniformsLoc) Location() glm.Matrix4 {
	return unif.location
}

func (unif *UniformsLoc) SetLocation(location glm.Matrix4) {
	unif.location = location
}

func (unif *UniformsLoc) SetGl() {
	glMat := unif.location.Gl()
	unif.modelLoc.UniformMatrix4f(false, &glMat)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "unif.modelLoc.UniformMatrix4f(false, &glMat)"
		panic(err)
	}
}

//----------------------------------------------------------------------------

type UniformsLocInstanced struct {
	globalMatricesUbi gl.UniformBlockIndex
}

func (unif *UniformsLocInstanced) SetUpVao(program gl.Program) {
	const matricesUbiName = "GlobalMatrices"

	// Uniform Block for the View and Projection matrices.
	unif.globalMatricesUbi = program.GetUniformBlockIndex(matricesUbiName)
	if err := glw.CheckGlError(); err != nil {
		err.Description = fmt.Sprintf("GetUniformBlockIndex(%#v)", matricesUbiName)
		panic(err)
	}
	if unif.globalMatricesUbi == gl.INVALID_INDEX {
		panic(fmt.Sprintf("GetUniformBlockIndex(%#v) returned INVALID_INDEX", matricesUbiName))
	}

	program.UniformBlockBinding(
		unif.globalMatricesUbi,
		glw.CameraUboBindingPoint)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "UniformBlockBinding"
		panic(err)
	}
}

func (unif *UniformsLocInstanced) SetLocation(location glm.Matrix4) {}
func (unif *UniformsLocInstanced) SetGl()                           {}
