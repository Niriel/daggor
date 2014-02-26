package main

import (
	"batch"
	"fmt"
	"github.com/go-gl/gl"
	"glw"
)

// The batch package provides generic batches.  However, some are too
// specific to a game to be included in that package.  Notably, anything
// that has to do with uniforms.  Therefore, we define uniform batches here
// that correspond to our very specific needs.

type UniformBatch struct {
	batch.BaseBatch
	modelMatrixUniform gl.UniformLocation
	modelMatrix        [16]float32
}

func MakeUniformBatch(
	modelMatrixUniform gl.UniformLocation,
	modelMatrix [16]float32,
) UniformBatch {
	return UniformBatch{modelMatrix: modelMatrix}
}

func (ubatch UniformBatch) Enter() {
	ubatch.modelMatrixUniform.UniformMatrix4f(false, &ubatch.modelMatrix)

	batch.GlobalGlState.Program.Validate()
	if err := glw.CheckGlError(); err != nil {
		err.Description = "Program.Validate failed."
		panic(err)
	}
	status := batch.GlobalGlState.Program.Get(gl.VALIDATE_STATUS)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "Program.Get(VALIDATE_STATUS) failed."
		panic(err)
	}
	if status == gl.FALSE {
		infolog := batch.GlobalGlState.Program.GetInfoLog()
		gl.GetError() // Clear error flag if infolog derped.
		panic(fmt.Errorf("Program validation failed. Log: %v", infolog))
	}
}
func (batch UniformBatch) Exit() {}
