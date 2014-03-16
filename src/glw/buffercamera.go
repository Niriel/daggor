package glw

import (
	"github.com/go-gl/gl"
	"glm"
)

const (
	cameraBufferEyeToClp = iota
	cameraBufferEyeToWld
	cameraBufferMatrixNb // Number of matrices.
)

type CameraBuffer struct {
	baseBuffer
	data         [cameraBufferMatrixNb][16]float32
	bindingPoint uint
	bound        bool
}

func NewCameraBuffer(usage gl.GLenum, bindingPoint uint) *CameraBuffer {
	buffer := new(CameraBuffer)
	buffer.target = gl.UNIFORM_BUFFER
	buffer.usage = usage
	buffer.bindingPoint = bindingPoint
	return buffer
}

func (buffer *CameraBuffer) Update() {
	// Create the buffer if needed.
	buffer.gen() // Does nothing if buffer is already created.
	// Create/fill/update the buffer store.
	buffer.update(buffer.data)
	// BindBufferBase raises an error when called without a buffer store, or
	// with an empty buffer store.
	if len(buffer.bufferdata) > 0 && !buffer.bound {
		buffer.bind()
		buffer.name.BindBufferBase(buffer.target, buffer.bindingPoint)
		if err := CheckGlError(); err != nil {
			err.Description = "BindBufferBase"
			panic(err)
		}
		buffer.unbind()
		buffer.bound = true
	}
}

// Data access.
func (buffer *CameraBuffer) SetEyeToClp(matrix glm.Matrix4) {
	buffer.data[cameraBufferEyeToClp] = matrix.Gl()
	buffer.bufferdataClean = false
}

func (buffer *CameraBuffer) SetEyeToWld(matrix glm.Matrix4) {
	buffer.data[cameraBufferEyeToWld] = matrix.Gl()
	buffer.bufferdataClean = false
}
