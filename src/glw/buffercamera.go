package glw

import (
	"github.com/go-gl/gl"
	"glm"
)

const (
	cameraBufferEyeToClp = iota
	cameraBufferEyeToWld
	cameraBufferMatrixNb
)

type CameraBuffer struct {
	baseBuffer
	data         [][16]float32 // Bunch of matrices.
	bindingPoint uint
}

func NewCameraBuffer(usage gl.GLenum, bindingPoint uint) *CameraBuffer {
	buffer := new(CameraBuffer)
	buffer.target = gl.UNIFORM_BUFFER
	buffer.usage = usage
	buffer.data = make([][16]float32, cameraBufferMatrixNb)
	buffer.bindingPoint = bindingPoint
	return buffer
}

func (buffer *CameraBuffer) SetUp() {
	buffer.gen()
	buffer.bind()
	// BindBufferBase complains with Invalid_Value if the buffer has no data
	// store.  So we create its datastore by calling Update().
	buffer.Update()
	buffer.name.BindBufferBase(buffer.target, buffer.bindingPoint)
	if err := CheckGlError(); err != nil {
		err.Description = "BindBufferBase"
		panic(err)
	}
	buffer.unbind()
}

func (buffer *CameraBuffer) Update() {
	buffer.update(buffer.data)
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
