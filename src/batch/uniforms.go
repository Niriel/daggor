package batch

import (
	"fmt"
	"github.com/go-gl/gl"
	"glw"
)

// The batch package provides generic batches.  However, some are too
// specific to a game to be included in that package.  Notably, anything
// that has to do with uniforms.  Therefore, we define uniform batches here
// that correspond to our very specific needs.

type UniformBatch struct {
	BaseBatch
	context            *GlContext
	modelMatrixUniform gl.UniformLocation
	modelMatrix        [16]float32
}

func MakeUniformBatch(
	context *GlContext,
	modelMatrixUniform gl.UniformLocation,
	modelMatrix [16]float32,
) UniformBatch {
	return UniformBatch{
		context:            context,
		modelMatrixUniform: modelMatrixUniform,
		modelMatrix:        modelMatrix,
	}
}

func (batch UniformBatch) Enter() {
	batch.modelMatrixUniform.UniformMatrix4f(false, &batch.modelMatrix)

	batch.context.Program.Validate()
	if err := glw.CheckGlError(); err != nil {
		err.Description = "Program.Validate failed."
		panic(err)
	}
	status := batch.context.Program.Get(gl.VALIDATE_STATUS)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "Program.Get(VALIDATE_STATUS) failed."
		panic(err)
	}
	if status == gl.FALSE {
		infolog := batch.context.Program.GetInfoLog()
		gl.GetError() // Clear error flag if infolog derped.
		panic(fmt.Errorf("Program validation failed. Log: %v", infolog))
	}
}
func (batch UniformBatch) Exit() {}
