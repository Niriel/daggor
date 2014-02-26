package batch

import (
	"fmt"
	"github.com/go-gl/gl"
	"glw"
)

type ProgramBatch struct {
	BaseBatch
	shaderRefs glw.ShaderRefs
}

func MakeProgramBatch(shaderRefs glw.ShaderRefs) ProgramBatch {
	return ProgramBatch{shaderRefs: shaderRefs}
}

func (batch ProgramBatch) Enter() {
	program, err := GlobalGlState.Programs.Serve(batch.shaderRefs)
	if err != nil {
		panic(err)
	}

	program.Use()
	if err := glw.CheckGlError(); err != nil {
		err.Description = fmt.Sprintf("program.Use() for ProgramBatch %v", program)
	}
	GlobalGlState.Program = program
}

func (batch ProgramBatch) Exit() {
	gl.ProgramUnuse()
	GlobalGlState.Program = 0
}
