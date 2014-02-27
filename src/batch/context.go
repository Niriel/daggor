package batch

import (
	"github.com/go-gl/gl"
	"glm"
	"glw"
	"unsafe"
)

const CameraUboBindingPoint = 0

type GlContext struct {
	cameraProj glm.Matrix4
	cameraView glm.Matrix4
	cameraUbo  gl.Buffer
	Programs   glw.Programs
	// The Program in use.  Set by ProgramBatch.  Usable by uniform batches
	// when they need to validate their inputs before a draw call.
	Program gl.Program
}

func MakeGlContext() GlContext {
	cameraUbo := gl.GenBuffer()
	if err := glw.CheckGlError(); err != nil {
		err.Description = "cameraUbo := gl.GenBuffer()"
		panic(err)
	}

	cameraUbo.Bind(gl.UNIFORM_BUFFER)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "cameraUbo.Bind(gl.UNIFORM_BUFFER)"
		panic(err)
	}

	gl.BufferData(
		gl.UNIFORM_BUFFER,
		int(unsafe.Sizeof(gl.GLfloat(0))*16*2), // Two matrices of 16 floats.
		nil,
		gl.STREAM_DRAW,
	)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "gl.BufferData(...) for camera UBO"
		panic(err)
	}

	cameraUbo.BindBufferBase(gl.UNIFORM_BUFFER, CameraUboBindingPoint)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "cameraUbo.BindBufferBase"
		panic(err)
	}

	cameraUbo.Unbind(gl.UNIFORM_BUFFER)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "cameraUbo.Unbind(gl.UNIFORM_BUFFER)"
		panic(err)
	}

	programs := glw.MakePrograms()

	return GlContext{
		cameraUbo: cameraUbo,
		Programs:  programs,
	}
}

func (context *GlContext) CameraProj() glm.Matrix4 {
	return context.cameraProj
}

func (context *GlContext) CameraView() glm.Matrix4 {
	return context.cameraView
}

func (context *GlContext) SetCameraProj(projMatrix glm.Matrix4) {
	const projMatrixStartOffset = 0
	context.cameraProj = projMatrix
	context.updateMatrix(projMatrix, projMatrixStartOffset)
}

func (context *GlContext) SetCameraView(viewMatrix glm.Matrix4) {
	const viewMatrixStartOffset = unsafe.Sizeof(gl.GLfloat(0)) * 16
	context.cameraView = viewMatrix
	context.updateMatrix(viewMatrix, viewMatrixStartOffset)
}

func (context *GlContext) updateMatrix(matrix glm.Matrix4, offset uintptr) {
	glmatrix := matrix.Gl()
	context.cameraUbo.Bind(gl.UNIFORM_BUFFER)
	gl.BufferSubData(
		gl.UNIFORM_BUFFER,
		int(offset),
		int(unsafe.Sizeof(glmatrix)),
		&glmatrix,
	)
	context.cameraUbo.Unbind(gl.UNIFORM_BUFFER)
}
