package glw

import (
	"fmt"
	"github.com/go-gl/gl"
	"sort"
)

type programRef string

// MakeProgramRef creates a unique ID for a program from the unique IDs of its
// shaders.  This is not exposed to the user.  Users are expected to identify
// programs by using ShaderRefs.
func makeProgramRef(shaders ShaderRefs) programRef {
	nb := len(shaders)
	sorted := make(ShaderRefs, nb, nb)
	copy(sorted, shaders)
	sort.Sort(sorted)
	return programRef(fmt.Sprintf("%v", sorted))
}

type Programs struct {
	// This structure is not supposed to be mutable.  That is,
	// the `shaders` and `programs` field will never change after
	// their creation.  They will always point to the same maps.
	// However, the maps can change.  Anyway, it is all hidden
	// from the caller since these fields are not exported.
	// So what I mean is: pass your instance of Programs by
	// value, not by address, unless you think that you need to.
	shaders  Shaders
	programs map[programRef]gl.Program
}

func NewPrograms() Programs {
	var programs Programs
	programs.shaders = MakeShaders()
	programs.programs = make(map[programRef]gl.Program)
	return programs
}

func (self Programs) Serve(srefs ShaderRefs) (gl.Program, error) {
	pref := makeProgramRef(srefs)
	program, ok := self.programs[pref]
	if ok {
		return program, nil
	}

	// Make sure that all the shaders are available and compile.
	nb_shaders := len(srefs)
	shaders := make([]gl.Shader, nb_shaders, nb_shaders)
	for i, sref := range srefs {
		shader, err := self.shaders.Serve(sref)
		if err != nil {
			return 0, err
		}
		shaders[i] = shader
	}

	program = gl.CreateProgram()
	if err := CheckGlError(); err != nil {
		if program == 0 {
			err.Description = "CreateProgram failed."
		} else {
			err.Description = "CreateProgram succeeded but OpenGL reports an error."
		}
		return program, err
	} else {
		if program == 0 {
			return 0, fmt.Errorf("CreateProgram failed but OpenGL reports no error.")
		}
	}

	for _, shader := range shaders {
		program.AttachShader(shader)
		if err := CheckGlError(); err != nil {
			program.Delete()
			gl.GetError() // Ignore Delete error.
			err.Description = "Program.AttachShader failed."
			return 0, err
		}
	}

	program.Link()
	if err := CheckGlError(); err != nil {
		infolog := program.GetInfoLog()
		program.Delete()
		gl.GetError() // Ignore Delete error.
		err.Description = fmt.Sprintf("Program.Link failed. Log: %v", infolog)
		return 0, err
	}
	link_status := program.Get(gl.LINK_STATUS)
	if err := CheckGlError(); err != nil {
		program.Delete()
		gl.GetError()
		err.Description = "Program.Get(LINK_STATUS) failed."
		return 0, err
	}
	if link_status == gl.FALSE {
		infolog := program.GetInfoLog()
		program.Delete()
		gl.GetError()
		err := fmt.Errorf("Program linking failed.  Log: %v", infolog)
		return 0, err
	}
	self.programs[pref] = program
	return program, nil
}
