package glw

import (
	"github.com/go-gl/gl"
	"glm"
)

type Light struct {
	Color  glm.Vector3
	Origin glm.Vector4 // w=0 for directional light, 1 for point light.
}

type LightBuffer struct {
	baseBuffer
	data []Light
}

func NewLightBuffer(usage gl.GLenum) *LightBuffer {
	buffer := new(LightBuffer)
	buffer.target = gl.UNIFORM_BUFFER
	buffer.usage = usage
	return buffer
}

// Satisfy the Buffer interface.

func (buffer *LightBuffer) SetUp(_ gl.Program) {
	buffer.gen()
}

func (buffer *LightBuffer) Update() {
	buffer.update(buffer.data)
}

//
func (buffer *LightBuffer) SetData(lights []Light) {
	buffer.data = make([]Light, len(lights), len(lights))
	copy(buffer.data, lights)
	buffer.bufferdataClean = false
}
