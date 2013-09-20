package glw

import (
	"fmt"
	"github.com/go-gl/gl"
	"sort"
)

type ProgramRef string

func MakeProgramRef(shaders ShaderRefs) ProgramRef {
	nb := len(shaders)
	sorted := make(ShaderRefs, nb, nb)
	copy(sorted, shaders)
	sort.Sort(sorted)
	return ProgramRef(fmt.Sprintf("%v", sorted))
}

type Programs struct {
	shaders  Shaders
	programs map[ProgramRef]gl.Program
}

func MakePrograms() Programs {
	var programs Programs
	programs.shaders = MakeShaders()
	programs.programs = make(map[ProgramRef]gl.Program)
	return programs
}

func (self Programs) Serve(srefs ShaderRefs) (gl.Program, error) {
	pref := MakeProgramRef(srefs)
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
	if link_status == 0 {
		infolog := program.GetInfoLog()
		program.Delete()
		gl.GetError()
		err := fmt.Errorf("Program linking failed.  Log: %v", infolog)
		return 0, err
	}
	self.programs[pref] = program
	return program, nil
}
