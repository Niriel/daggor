package glw

import (
	//"fmt"
	"github.com/go-gl/gl"
	"glm"
)

const nbLightsMax = 100

type Light struct {
	Color  glm.Vector4
	Origin glm.Vector4 // w=0 for directional light, 1 for point light.
}

type GlLight struct {
	Color  [4]gl.GLfloat
	Origin [4]gl.GLfloat
}

func (light Light) ToGl() GlLight {
	return GlLight{
		Color:  light.Color.GlFloats(),
		Origin: light.Origin.GlFloats(),
	}
}

type LightBuffer struct {
	baseBuffer
	data struct {
		lights   [nbLightsMax]GlLight
		nbLights gl.GLuint
	}
	bindingPoint uint
	bound        bool
}

func NewLightBuffer(usage gl.GLenum, bindingPoint uint) *LightBuffer {
	buffer := new(LightBuffer)
	buffer.target = gl.UNIFORM_BUFFER
	buffer.usage = usage
	buffer.bindingPoint = bindingPoint
	return buffer
}

func (buffer *LightBuffer) Update() {
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

func (buffer *LightBuffer) SetLight(i uint32, light Light) {
	buffer.data.lights[i] = light.ToGl()
	buffer.bufferdataClean = false
}

func (buffer *LightBuffer) SetNbLights(nbLights uint32) {
	if nbLights > nbLightsMax {
		panic("too many lights")
	}
	buffer.data.nbLights = gl.GLuint(nbLights)
	buffer.bufferdataClean = false
}

func (buffer *LightBuffer) SetLights(lights []Light) {
	buffer.SetNbLights(uint32(len(lights)))
	for i, light := range lights {
		buffer.SetLight(uint32(i), light)
	}
	// No need to set bufferdataClean to false, SetNbLight and SetLight did it.
}
