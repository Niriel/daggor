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
	VSH_POS3 = ShaderRef(iota)
	VSH_COL3
	FSH_ZRED
	FSH_ZGREEN
	FSH_VCOL
)

// This keeps track of the shaders that are in OpenGL.
// Each shader is known by the application by a unique identifier called
// "Shader reference", of type ShaderRef.
type Shaders map[ShaderRef]gl.Shader

func (shaders Shaders) Serve(ref ShaderRef) (gl.Shader, error) {
	shader, ok := shaders[ref]
	if ok {
		return shader, nil
	}
	// Transparently cache the shader.
	// Also cache shaders that failed the compilation step.
	shader_seed, ok := SHADER_SOURCES[ref]
	if !ok {
		return 0, fmt.Errorf("Shader Reference '%#v' not found.", ref)
	}
	shader = gl.CreateShader(gl.GLenum(shader_seed.Type))
	if shader == 0 {
		return 0, fmt.Errorf("CreateShader failed with [%v].", ERROR_NAMES[gl.GetError()])
	}
	shader.Source(shader_seed.Source)
	if glec := gl.GetError(); glec != gl.NO_ERROR {
		return 0, fmt.Errorf("Shader.Source failed with [%v].", ERROR_NAMES[glec])
	}
	shader.Compile()
	infolog := shader.GetInfoLog()
	if glec := gl.GetError(); glec != gl.NO_ERROR {
		// We do not want to try to compile the same shader over and
		// over again, so once it failed, we put it in the map.  A value
		// of 0 corresponds to an invalid shader in OpenGl.  So next time,
		// Shaders.Serve will return `0, nil`, and because it is 0 we know
		// that this shader is broken.
		shaders[ref] = 0
		return 0, fmt.Errorf("Shader.Compile failed with [%v].\nInfolog:\n%v", ERROR_NAMES[glec], infolog)
	}
	shaders[ref] = shader
	return shader, nil
}
