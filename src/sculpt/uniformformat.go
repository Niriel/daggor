package sculpt

import (
	"fmt"
	"github.com/go-gl/gl"
	"glm"
	"glw"
)

type Uniforms interface {
	SetGl()
	SetUpVao(program gl.Program)
}

type Modeler interface {
	Model() glm.Matrix4
	SetModel(glm.Matrix4)
}

// For now we have one type of uniform setup only.  In the future we will need
// more and we'll probably need an interface.
type UniformsLoc struct {
	// Filled at the creation of the object.
	globalMatricesUbb uint // Uniform Block Binding.
	// Filled when calling the SetUpVao method with a program.
	globalMatricesUbi gl.UniformBlockIndex
	modelLoc          gl.UniformLocation
	// Filled when updating the uniform values.
	model glm.Matrix4
}

func MakeUniformsLoc(globalMatricesUbb uint) UniformsLoc {
	return UniformsLoc{globalMatricesUbb: globalMatricesUbb}
}

func (unif *UniformsLoc) SetUpVao(program gl.Program) {
	const modelLocName = "model_matrix"
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
		unif.globalMatricesUbb)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "UniformBlockBinding"
		panic(err)
	}
}

func (unif UniformsLoc) Model() glm.Matrix4 {
	return unif.model
}

func (unif *UniformsLoc) SetModel(model glm.Matrix4) {
	unif.model = model
}

func (unif UniformsLoc) SetGl() {
	glMat := unif.model.Gl()
	unif.modelLoc.UniformMatrix4f(false, &glMat)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "unif.modelLoc.UniformMatrix4f(false, &glMat)"
		panic(err)
	}
}
