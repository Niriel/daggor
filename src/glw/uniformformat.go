package glw

import (
	"fmt"
	"github.com/go-gl/gl"
	"glm"
)

type Uniforms interface {
	SetLocation(glm.Matrix4)
	SetGl()
	SetUp(program gl.Program)
}

//----------------------------------------------------------------------------

// To use where there is no uniform to deal with.
type UniformsNone struct{}

func (unif *UniformsNone) SetUp(program gl.Program) {}
func (unif *UniformsNone) SetGl()                   {}

//----------------------------------------------------------------------------
type UniformsLoc struct {
	// Filled when calling the SetUp method with a program.
	globalMatricesUbi gl.UniformBlockIndex
	modelLoc          gl.UniformLocation
	// Filled when updating the uniform values.
	location glm.Matrix4
}

func (unif *UniformsLoc) SetUp(program gl.Program) {
	const modelLocName = "model_to_eye"
	const matricesUbiName = "GlobalMatrices"

	unif.modelLoc = program.GetUniformLocation(modelLocName)
	if err := CheckGlError(); err != nil {
		err.Description = fmt.Sprintf("program.GetUniformLocation(%#v)", modelLocName)
		panic(err)
	}
	if unif.modelLoc == -1 {
		panic(fmt.Sprintf("uniform %#v not found", modelLocName))
	}

	unif.globalMatricesUbi = program.GetUniformBlockIndex(matricesUbiName)
	if err := CheckGlError(); err != nil {
		err.Description = fmt.Sprintf("GetUniformBlockIndex(%#v)", matricesUbiName)
		panic(err)
	}
	if unif.globalMatricesUbi == gl.INVALID_INDEX {
		panic(fmt.Sprintf("GetUniformBlockIndex(%#v) returned INVALID_INDEX", matricesUbiName))
	}
	program.UniformBlockBinding(unif.globalMatricesUbi, CameraUboBindingPoint)
	if err := CheckGlError(); err != nil {
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
	if err := CheckGlError(); err != nil {
		err.Description = "unif.modelLoc.UniformMatrix4f(false, &glMat)"
		panic(err)
	}
}

//----------------------------------------------------------------------------

type UniformsLocInstanced struct {
	globalMatricesUbi gl.UniformBlockIndex
	textureLoc        gl.UniformLocation
}

func (unif *UniformsLocInstanced) SetUp(program gl.Program) {
	const matricesUbiName = "GlobalMatrices"

	unif.globalMatricesUbi = program.GetUniformBlockIndex(matricesUbiName)
	if err := CheckGlError(); err != nil {
		err.Description = fmt.Sprintf("GetUniformBlockIndex(%#v)", matricesUbiName)
		panic(err)
	}
	if unif.globalMatricesUbi == gl.INVALID_INDEX {
		panic(fmt.Sprintf("GetUniformBlockIndex(%#v) returned INVALID_INDEX", matricesUbiName))
	}

	program.UniformBlockBinding(unif.globalMatricesUbi, CameraUboBindingPoint)
	if err := CheckGlError(); err != nil {
		err.Description = "UniformBlockBinding"
		panic(err)
	}

	unif.textureLoc = program.GetUniformLocation("environment_map")
	if err := CheckGlError(); err != nil {
		err.Description = "GetUniformLocation(environment)"
		panic(err)
	}
	if unif.textureLoc == -1 {
		panic("environment uniform not found")
	}
}

func (unif *UniformsLocInstanced) SetLocation(location glm.Matrix4) {}
func (unif *UniformsLocInstanced) SetGl() {
	gl.ActiveTexture(gl.TEXTURE0)
	if err := CheckGlError(); err != nil {
		err.Description = "gl.ActiveTexture(gl.TEXTURE0)"
		panic(err)
	}

	globaltexture.Bind(gl.TEXTURE_CUBE_MAP)
	if err := CheckGlError(); err != nil {
		err.Description = "gl.Texture(1).Bind(gl.TEXTURE_CUBE_MAP)"
		panic(err)
	}

	unif.textureLoc.Uniform1i(0)
	if err := CheckGlError(); err != nil {
		err.Description = "unif.textureloc.Uniform1i(0)"
		panic(err)
	}
}
