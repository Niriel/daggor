package glw

import (
	"github.com/go-gl/gl"
	"glm"
)

const CameraUboBindingPoint = 0

type GlContext struct {
	cameraBuffer *CameraBuffer
	Programs     Programs
	// The Program in use.  Set by ProgramBatch.  Usable by uniform batches
	// when they need to validate their inputs before a draw call.
	Program gl.Program
}

func NewGlContext() *GlContext {
	programs := NewPrograms()
	cameraBuffer := NewCameraBuffer(gl.STREAM_DRAW, CameraUboBindingPoint)
	cameraBuffer.SetUp()
	return &GlContext{
		cameraBuffer: cameraBuffer,
		Programs:     programs,
	}
}

func (context *GlContext) SetCameraProj(matrix glm.Matrix4) {
	context.cameraBuffer.SetEyeToClp(matrix)
}

func (context *GlContext) SetCameraViewI(matrix glm.Matrix4) {
	context.cameraBuffer.SetEyeToWld(matrix)
}

func (context *GlContext) UpdateCamera() {
	context.cameraBuffer.Update()
}
