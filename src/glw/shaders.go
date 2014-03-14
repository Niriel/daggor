package glw

import (
	"fmt"
	"github.com/go-gl/gl"
)

type ShaderRef uint16 // A unique ID, could also be a file name.
type ShaderType gl.GLenum
type ShaderSeed struct {
	Type   ShaderType
	Source string
}

const (
	VERTEX_SHADER   = ShaderType(gl.VERTEX_SHADER)
	GEOMETRY_SHADER = ShaderType(gl.GEOMETRY_SHADER)
	FRAGMENT_SHADER = ShaderType(gl.FRAGMENT_SHADER)
)

const (
	VSH_NOR_UV_INSTANCED = ShaderRef(iota)
	FSH_NOR_UV
)

// Programs are uniquely identified by their shaders.  I need to be able to sort
// shader references in order to create program references.  Here I implement the
// interface required by the `sort` package.
type ShaderRefs []ShaderRef

func (self ShaderRefs) Len() int {
	return len(self)
}
func (self ShaderRefs) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
func (self ShaderRefs) Less(i, j int) bool {
	return self[i] < self[j]
}

// This keeps track of the shaders that are in OpenGL.
// Each shader is known by the application by a unique identifier called
// "Shader reference", of type ShaderRef.
type Shaders map[ShaderRef]gl.Shader

func MakeShaders() Shaders {
	return make(Shaders)
}

func (shaders Shaders) Serve(ref ShaderRef) (gl.Shader, error) {
	shader, ok := shaders[ref]
	if ok {
		return shader, nil
	}

	shader_seed, ok := SHADER_SOURCES[ref]
	if !ok {
		return 0, fmt.Errorf("Shader Reference '%#v' not found.", ref)
	}

	shader = gl.CreateShader(gl.GLenum(shader_seed.Type))
	if err := CheckGlError(); err != nil {
		if shader == 0 {
			err.Description = "CreateShader failed."
		} else {
			err.Description = "CreateShader succeeded but OpenGL reports an error."
		}
		return shader, err
	} else {
		if shader == 0 {
			return 0, fmt.Errorf("CreateShader failed but OpenGL reports no error.")
		}
	}

	shader.Source(shader_seed.Source)
	if err := CheckGlError(); err != nil {
		shader.Delete()
		gl.GetError() // Delete may also raise an error, ignore.
		err.Description = "Shader.Source failed."
		return 0, err
	}

	shader.Compile()
	if err := CheckGlError(); err != nil {
		infolog := shader.GetInfoLog()
		shader.Delete()
		gl.GetError() // GetInfoLog and delete may raise an error.
		err.Description = fmt.Sprintf("Shader.Compile failed. Log: %v", infolog)
		return 0, err
	}
	compileStatus := shader.Get(gl.COMPILE_STATUS)
	if err := CheckGlError(); err != nil {
		shader.Delete()
		gl.GetError()
		err.Description = "shader.Get(gl.COMPILE_STATUS) failed."
		return 0, err
	}
	if compileStatus == gl.FALSE {
		infolog := shader.GetInfoLog()
		shader.Delete()
		gl.GetError()
		err := fmt.Errorf("Shader compiling failed.  Log: %v", infolog)
		return 0, err
	}
	shaders[ref] = shader
	return shader, nil
}
