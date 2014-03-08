package glw

import (
	"bytes"
	"encoding/binary"
	"github.com/go-gl/gl"
)

// Concrete classes of array buffers must satisfy this public interface.
type Buffer interface {
	bind()
	unbind()
	SetUpVao(gl.Program)
	// UpdateBuffer can be called every frame.  It does nothing if it has
	// nothing to do.  If needed, it creates a VBO on the fly, and/or fills
	// it with the OpenGL-friendly version of the Go-friendly vertexdata
	// passed with SetVertexData.
	Update()
	// Delete the OpenGL vertex buffer object.
	Delete()
}

// baseBuffer is common to every collection of vertices.
// This base class is responsible for converting Go-friendly slices of
// vertices into an OpenGL buffer.
type baseBuffer struct {
	// vbo is the Vertex Buffer Object in which we send the vertices.
	name gl.Buffer
	// bufferdata is the binary data to send to the OpenGL buffer.
	bufferdata []byte
	// bufferdataClean signals whether bufferdata must be reconstructed from
	// the Go-friendly vertex data or if it's still good.
	bufferdataClean bool
	// OpenGL hint about how often the buffer is expected to be updated.
	usage  gl.GLenum
	target gl.GLenum
}

func (buffer *baseBuffer) gen(target, usage gl.GLenum) {
	if buffer.name == 0 {
		buffer.name = gl.GenBuffer()
		buffer.target = target
		buffer.usage = usage
	}
}

func (buffer *baseBuffer) bind() {
	if buffer.name == 0 {
		panic("tried to bind buffer 0")
	}
	buffer.name.Bind(buffer.target)
	if err := CheckGlError(); err != nil {
		err.Description = "buffer.name.Bind(buffer.target)"
		panic(err)
	}
}

func (buffer *baseBuffer) unbind() {
	if buffer.name == 0 {
		panic("tried to unbind buffer 0")
	}
	buffer.name.Unbind(buffer.target)
	if err := CheckGlError(); err != nil {
		err.Description = "buffer.name.Unbind(buffer.target)"
		panic(err)
	}
}

// updateBuffer fills the OpenGL buffer IF NEEDED.
// It updates the buffer if the bufferdataClean flag is false.
// It is safe to call this method every frame via the concrete classes method
// UpdateBuffer: most of the time it will just return immediately.
func (buffer *baseBuffer) update(vertexdata interface{}) {
	if buffer.bufferdataClean {
		return
	}
	if buffer.name == 0 {
		panic("tried to update buffer 0")
	}

	// Convert the Go-friendly data into OpenGL-friendly data.
	oldSize := len(buffer.bufferdata)
	bufferdata := new(bytes.Buffer)
	err := binary.Write(bufferdata, endianness, vertexdata)
	if err != nil {
		panic(err)
	}
	buffer.bufferdata = bufferdata.Bytes()
	buffer.bufferdataClean = true
	newSize := len(buffer.bufferdata)

	// Should we make the buffer bigger?
	needBigger := newSize > oldSize

	buffer.bind()

	if needBigger {
		// (Re)allocate a buffer.
		gl.BufferData(
			buffer.target,
			len(buffer.bufferdata),
			&buffer.bufferdata[0],
			buffer.usage,
		)
		if err := CheckGlError(); err != nil {
			err.Description = "gl.BufferData"
			panic(err)
		}
	} else {
		// Re-use existing buffer.
		gl.BufferSubData(
			buffer.target,
			0,
			len(buffer.bufferdata),
			&buffer.bufferdata[0],
		)
		if err := CheckGlError(); err != nil {
			err.Description = "gl.BufferSubData"
			panic(err)
		}
	}

	buffer.unbind()
}

// Delete the OpenGL vertex buffer object.
func (buffer *baseBuffer) Delete() {
	buffer.name.Delete()
	buffer.name = 0
}
