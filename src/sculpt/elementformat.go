package sculpt

import (
	"bytes"
	"encoding/binary"
	"github.com/go-gl/gl"
	"glw"
)

type baseElements struct {
	// ebo is the Element Buffer Object in which we send the elements.
	ebo gl.Buffer
	// bufferdata is the binary data to send to the OpenGL buffer.
	bufferdata []byte
	// bufferdataClean signals whether bufferdata must be reconstructed from
	// the Go-friendly vertex data or if it's still good.
	bufferdataClean bool
}

func (elements *baseElements) updateBuffer(elementdata interface{}) {
	const target = gl.ELEMENT_ARRAY_BUFFER
	const usage = gl.STATIC_DRAW
	if elements.bufferdataClean {
		return
	}
	bufferdata := new(bytes.Buffer)
	err := binary.Write(bufferdata, endianness, elementdata)
	if err != nil {
		panic(err)
	}
	elements.bufferdata = bufferdata.Bytes()
	elements.bufferdataClean = true

	// Has an EBO been created yet?
	isEboNew := elements.ebo == 0

	if isEboNew {
		elements.ebo = gl.GenBuffer()
	}

	elements.ebo.Bind(target)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "ebo.Bind(gl.ELEMENT_ARRAY_BUFFER)"
		panic(err)
	}

	if isEboNew {
		gl.BufferData(
			target,
			len(elements.bufferdata),
			&elements.bufferdata[0],
			usage,
		)
		if err := glw.CheckGlError(); err != nil {
			err.Description = "gl.BufferData"
			panic(err)
		}
	} else {
		gl.BufferSubData(
			target,
			0,
			len(elements.bufferdata),
			&elements.bufferdata[0],
		)
		if err := glw.CheckGlError(); err != nil {
			err.Description = "gl.BufferSubData"
			panic(err)
		}
	}

	elements.ebo.Unbind(target)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "ebo.Unbind(gl.Element_ARRAY_BUFFER)"
		panic(err)
	}
}

func (elements *baseElements) BufferName() gl.Buffer {
	return elements.ebo
}

func (elements *baseElements) DeleteBuffer() {
	elements.ebo.Delete()
	elements.ebo = 0
}

func (elements *baseElements) SetUpVao() {
	// The ebo binding is part of the VAO state.  The only setup we need to
	// do is to bind the ebo.
	const target = gl.ELEMENT_ARRAY_BUFFER
	elements.ebo.Bind(target)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "ebo.Bind(gl.ELEMENT_ARRAY_BUFFER)"
		panic(err)
	}
}

type Elements interface {
	SetUpVao()
	BufferName() gl.Buffer
	UpdateBuffer()
	DeleteBuffer()
}

type ElementsUbyte struct {
	baseElements
	elementdata []gl.GLubyte
}

func (elements *ElementsUbyte) SetElementData(ed []gl.GLubyte) {
	elements.elementdata = make([]gl.GLubyte, len(ed), len(ed))
	copy(elements.elementdata, ed)
	elements.bufferdataClean = false
}

func (elements *ElementsUbyte) UpdateBuffer() {
	elements.updateBuffer(elements.elementdata)
}
