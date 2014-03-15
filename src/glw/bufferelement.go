package glw

import (
	"github.com/go-gl/gl"
)

type ElementsUbyte struct {
	baseBuffer
	data []gl.GLubyte
}

func NewElementsUbyte(usage gl.GLenum) *ElementsUbyte {
	buffer := new(ElementsUbyte)
	buffer.target = gl.ELEMENT_ARRAY_BUFFER
	buffer.usage = usage
	return buffer
}

func (buffer *ElementsUbyte) SetUp(program gl.Program) {
	// The ebo binding is part of the VAO state.  The only setup we need to
	// do is to bind the ebo and leave it bound.
	buffer.gen()
	buffer.bind()
}

func (buffer *ElementsUbyte) SetData(ed []gl.GLubyte) {
	buffer.data = make([]gl.GLubyte, len(ed), len(ed))
	copy(buffer.data, ed)
	buffer.bufferdataClean = false
}

func (buffer *ElementsUbyte) Update() {
	buffer.update(buffer.data)
}
