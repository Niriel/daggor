package glw

import (
	"github.com/go-gl/gl"
	"glm"
)

const (
	CameraUboBindingPoint = iota
	LightUboBindingPoint
)

type GlContext struct {
	cameraBuffer *CameraBuffer
	lightBuffer  *LightBuffer
	Programs     Programs
	// The Program in use.  Set by ProgramBatch.  Usable by uniform batches
	// when they need to validate their inputs before a draw call.
	Program gl.Program
}

func NewGlContext() *GlContext {
	programs := NewPrograms()
	cameraBuffer := NewCameraBuffer(gl.STREAM_DRAW, CameraUboBindingPoint)
	lightBuffer := NewLightBuffer(gl.STREAM_DRAW, LightUboBindingPoint)
	return &GlContext{
		cameraBuffer: cameraBuffer,
		lightBuffer:  lightBuffer,
		Programs:     programs,
	}
}

func (context *GlContext) SetEyeToClp(matrix glm.Matrix4) {
	context.cameraBuffer.SetEyeToClp(matrix)
}

func (context *GlContext) SetEyeToWld(matrix glm.Matrix4) {
	context.cameraBuffer.SetEyeToWld(matrix)
}

func (context *GlContext) UpdateCamera() {
	context.cameraBuffer.Update()
}

func (context *GlContext) SetLights(lights []Light) {
	context.lightBuffer.SetLights(lights)
}

func (context *GlContext) UpdateLights() {
	context.lightBuffer.Update()
}
