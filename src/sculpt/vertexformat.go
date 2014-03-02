package sculpt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-gl/gl"
	"glw"
	"unsafe"
)

//=============================================================================

// This section defines the various vertex formats used in the game.
// It is at the level of a single vertex, not a collection of vertices.

// VertexXyz defines vertices that have a location and nothing else.
// No color, UV or any other parameter.
type VertexXyz struct {
	X, Y, Z gl.GLfloat
}

// VertexXyzRgb defines vertices that have a location and a color.
// No UV information.  Note that there is no Alpha component to the color.
type VertexXyzRgb struct {
	X, Y, Z gl.GLfloat
	R, G, B gl.GLfloat
}

//=============================================================================

// This section defines collections of vertices.

// baseVertices is common to every collection of vertices.
// This base class is responsible for converting Go-friendly slices of
// vertices into an OpenGL buffer.
type baseVertices struct {
	// vbo is the Vertex Buffer Object in which we send the vertices.
	vbo gl.Buffer
	// bufferdata is the binary data to send to the OpenGL buffer.
	bufferdata []byte
	// bufferdataClean signals whether bufferdata must be reconstructed from
	// the Go-friendly vertex data or if it's still good.
	bufferdataClean bool
	// OpenGL hint about how often the buffer is expected to be updated.
	usage gl.GLenum
}

func (vertices baseVertices) bind() {
	vertices.vbo.Bind(gl.ARRAY_BUFFER)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "vbo.Bind(gl.ARRAY_BUFFER)"
		panic(err)
	}
}

func (vertices baseVertices) unbind() {
	vertices.vbo.Unbind(gl.ARRAY_BUFFER)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "vbo.Unbind(gl.ARRAY_BUFFER)"
		panic(err)
	}
}

func (vertices baseVertices) BufferName() gl.Buffer {
	return vertices.vbo
}

// updateBuffer creates and fill the OpenGL buffer IF NEEDED.
// It creates the buffer if none is created yet, and it updates it if the
// bufferdataClean flag is false.  It is safe to call this method every frame,
// via the concrete classes method UpdateBuffer: most of the time it will just
// return immediately.
func (vertices *baseVertices) updateBuffer(vertexdata interface{}) {
	if vertices.bufferdataClean {
		return
	}
	const target = gl.ARRAY_BUFFER

	// Convert the Go-friendly data into OpenGL-friendly data.
	bufferdata := new(bytes.Buffer)
	err := binary.Write(bufferdata, endianness, vertexdata)
	if err != nil {
		panic(err)
	}
	vertices.bufferdata = bufferdata.Bytes()
	vertices.bufferdataClean = true

	// Has a VBO been created yet?
	isVboNew := vertices.vbo == 0

	if isVboNew {
		vertices.vbo = gl.GenBuffer()
		if vertices.vbo == 0 {
			panic("gl.GenBuffer returned 0")
		}
	}

	vertices.bind()

	if isVboNew {
		gl.BufferData(
			target,
			len(vertices.bufferdata),
			&vertices.bufferdata[0],
			vertices.usage,
		)
		if err := glw.CheckGlError(); err != nil {
			err.Description = "gl.BufferData"
			panic(err)
		}
	} else {
		gl.BufferSubData(
			target,
			0,
			len(vertices.bufferdata),
			&vertices.bufferdata[0],
		)
		if err := glw.CheckGlError(); err != nil {
			err.Description = "gl.BufferSubData"
			panic(err)
		}
	}

	vertices.unbind()
}

// Delete the OpenGL vertex buffer object.
func (vertices *baseVertices) DeleteBuffer() {
	vertices.vbo.Delete()
	vertices.vbo = 0
}

// Concrete classes of Vertices satisfy this public interface.
type Vertices interface {
	bind()
	unbind()
	SetUpVao(gl.Program)
	BufferName() gl.Buffer
	// UpdateBuffer can be called every frame.  It does nothing if it has
	// nothing to do.  If needed, it creates a VBO on the fly, and/or fills
	// it with the OpenGL-friendly version of the Go-friendly vertexdata
	// passed with SetVertexData.
	UpdateBuffer()
	// Delete the OpenGL vertex buffer object.
	DeleteBuffer()
}

// Concrete classes of Vertices derive from baseVertices and correspond to
// a specific vertex format.

type VerticesXyz struct {
	baseVertices
	vertexdata []VertexXyz
}

type VerticesXyzRgb struct {
	baseVertices
	vertexdata []VertexXyzRgb
}

func (vertices *VerticesXyz) SetVertexData(vd []VertexXyz) {
	vertices.vertexdata = make([]VertexXyz, len(vd), len(vd))
	copy(vertices.vertexdata, vd)
	vertices.bufferdataClean = false
}

func (vertices *VerticesXyzRgb) SetVertexData(vd []VertexXyzRgb) {
	vertices.vertexdata = make([]VertexXyzRgb, len(vd), len(vd))
	copy(vertices.vertexdata, vd)
	vertices.bufferdataClean = false
}

func (vertices *VerticesXyz) UpdateBuffer() {
	vertices.updateBuffer(vertices.vertexdata)
}

func (vertices *VerticesXyzRgb) UpdateBuffer() {
	vertices.updateBuffer(vertices.vertexdata)
}

//=============================================================================

// This section is about configuring the Vertex Array Object of a mesh.
// The mesh calls SetUpVao which is a method of Vertices.
// SetUpVao needs a gl Program object in order to query variable parameter
// names and all.

// Since there is a lot of code common to all the Vertices object about how
// to set up a VAO, we just ask each Vertices object to satisfy an interface
// containing what is specific to that Vertices object.  Then we pass it to
// a generic function.

func (vertices VerticesXyz) SetUpVao(program gl.Program) {
	vertices.bind()
	verticesSetUpVao(vertices, program)
	vertices.unbind()
}
func (vertices VerticesXyzRgb) SetUpVao(program gl.Program) {
	vertices.bind()
	verticesSetUpVao(vertices, program)
	vertices.unbind()
}

// The verticesSetUpVaoInt contains everything that is needed by the function
// verticesSetUpVao.
type verticesSetUpVaoInt interface {
	names() []string
	attribPointers([]gl.AttribLocation)
}

func verticesSetUpVao(vertices verticesSetUpVaoInt, program gl.Program) {
	// Collect the attrib locations for each attrib name.
	atts_names := vertices.names() // Expected GLSL variable names.
	atts := make([]gl.AttribLocation, len(atts_names))
	for i, att_name := range atts_names {
		atts[i] = program.GetAttribLocation(att_name)
		if err := glw.CheckGlError(); err != nil {
			err.Description = fmt.Sprintf("program.GetAttribLocation(%#v)", att_name)
			panic(err)
		}
		if atts[i] == -1 {
			panic(fmt.Sprintf("attrib location %#v not found", att_name))
		}
		atts[i].EnableArray()
		if err := glw.CheckGlError(); err != nil {
			err.Description = "atts[i].EnableArray()"
			panic(err)
		}
	}
	// Now that the locations are known, we can relate them to vertex data.
	vertices.attribPointers(atts)
}

func (vertices VerticesXyz) names() []string {
	return []string{"vpos"}
}
func (vertices VerticesXyzRgb) names() []string {
	return []string{"vpos", "vcol"}
}

func (vertices VerticesXyz) attribPointers(atts []gl.AttribLocation) {
	const FLOATSIZE = unsafe.Sizeof(gl.GLfloat(0))
	const NB_COORDS = 3 // x y and z.
	const COORDS_SIZE = NB_COORDS * FLOATSIZE
	const COORDS_OFS = uintptr(0)
	const TOTAL_SIZE = int(COORDS_SIZE)
	atts[0].AttribPointer(NB_COORDS, gl.FLOAT, false, TOTAL_SIZE, COORDS_OFS)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "VerticesXyz atts[0].AttribPointer"
		panic(err)
	}
}
func (vertices VerticesXyzRgb) attribPointers(atts []gl.AttribLocation) {
	const FLOATSIZE = unsafe.Sizeof(gl.GLfloat(0))
	const NB_COORDS = 3 // x y and z.
	const NB_COLORS = 3 // r g and b.
	const COORDS_SIZE = NB_COORDS * FLOATSIZE
	const COLORS_SIZE = NB_COLORS * FLOATSIZE
	const COORDS_OFS = uintptr(0)
	const COLORS_OFS = uintptr(COORDS_SIZE)
	const TOTAL_SIZE = int(COORDS_SIZE + COLORS_SIZE)
	atts[0].AttribPointer(NB_COORDS, gl.FLOAT, false, TOTAL_SIZE, COORDS_OFS)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "VerticesXyzRgb atts[0].AttribPointer"
		panic(err)
	}
	atts[1].AttribPointer(NB_COLORS, gl.FLOAT, false, TOTAL_SIZE, COLORS_OFS)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "VerticesXyzRgb atts[1].AttribPointer"
		panic(err)
	}
}
