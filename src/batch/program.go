package batch

import (
	"fmt"
	"github.com/go-gl/gl"
	"glw"
)

type ProgramBatch struct {
	BaseBatch
	context    *GlContext
	shaderRefs glw.ShaderRefs
}

func MakeProgramBatch(context *GlContext, shaderRefs glw.ShaderRefs) ProgramBatch {
	return ProgramBatch{context: context, shaderRefs: shaderRefs}
}

func (batch ProgramBatch) Enter() {
	program, err := batch.context.Programs.Serve(batch.shaderRefs)
	if err != nil {
		panic(err)
	}

	program.Use()
	if err := glw.CheckGlError(); err != nil {
		err.Description = fmt.Sprintf("program.Use() for ProgramBatch %v", program)
	}
	batch.context.Program = program
}

func (batch ProgramBatch) Exit() {
	gl.ProgramUnuse()
	batch.context.Program = 0
}
