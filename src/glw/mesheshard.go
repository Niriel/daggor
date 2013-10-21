// This module contains the code for hard coded meshes used during
// the development and debug phases of the game.

package glw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-gl/gl"
	"unsafe"
)

var ENDIANNES binary.ByteOrder

func isLittleEndian() bool {
	// Credit matt kane, taken from his gosndfile project.
	// https://groups.google.com/forum/#!msg/golang-nuts/3GEzwKfRRQw/D1bMbFP-ClAJ
	// https://github.com/mkb218/gosndfile
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return (b == 0x04)
}

func init() {
	// Here, somehow detect the endianness of the system.
	// We need to encode vertex data before sending it to OpenGL, and OpenGL
	// uses the endianness of the system on which we are running.
	// Why encode? Because we do not want to care about the way Go allign its
	// data in memory.  Encoding in binary ensures that the data has no gap or
	// padding, which is easier to manage.
	if isLittleEndian() {
		ENDIANNES = binary.LittleEndian
	} else {
		ENDIANNES = binary.BigEndian
	}
}

type GlEncoder interface {
	GlEncode() ([]byte, error)
}

type ElementIndexFormat interface {
	GlEncoder
	Len() int
	GlType() gl.GLenum
}

type VertexFormat interface {
	GlEncoder
	AttribPointers([]gl.AttribLocation)
	Names() []string
}

type VertexFormatXyz struct {
	x, y, z gl.GLfloat
}
type VertexFormatXyzRgb struct {
	x, y, z gl.GLfloat
	r, g, b gl.GLfloat
}
type VerticesFormatXyz []VertexFormatXyz
type VerticesFormatXyzRgb []VertexFormatXyzRgb

func (self VerticesFormatXyz) Names() []string {
	return []string{"vpos"}
}
func (self VerticesFormatXyzRgb) Names() []string {
	return []string{"vpos", "vcol"}
}
func (self VerticesFormatXyz) AttribPointers(atts []gl.AttribLocation) {
	const FLOATSIZE = unsafe.Sizeof(gl.GLfloat(0))
	const NB_COORDS = 3
	const COORDS_SIZE = NB_COORDS * FLOATSIZE
	const COORDS_OFS = uintptr(0)
	const TOTAL_SIZE = int(COORDS_SIZE)
	atts[0].AttribPointer(NB_COORDS, gl.FLOAT, false, TOTAL_SIZE, COORDS_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesFormatXyz atts[0].AttribPointer"
		panic(err)
	}
}
func (self VerticesFormatXyzRgb) AttribPointers(atts []gl.AttribLocation) {
	const FLOATSIZE = unsafe.Sizeof(gl.GLfloat(0))
	const NB_COORDS = 3
	const NB_COLORS = 3
	const COORDS_SIZE = NB_COORDS * FLOATSIZE
	const COLORS_SIZE = NB_COLORS * FLOATSIZE
	const COORDS_OFS = uintptr(0)
	const COLORS_OFS = uintptr(COORDS_SIZE)
	const TOTAL_SIZE = int(COORDS_SIZE + COLORS_SIZE)
	atts[0].AttribPointer(NB_COORDS, gl.FLOAT, false, TOTAL_SIZE, COORDS_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesFormatXyzRgb atts[0].AttribPointer"
		panic(err)
	}
	atts[1].AttribPointer(NB_COLORS, gl.FLOAT, false, TOTAL_SIZE, COLORS_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesFormatXyzRgb atts[1].AttribPointer"
		panic(err)
	}
}

func (self VerticesFormatXyz) GlEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, ENDIANNES, self)
	return buf.Bytes(), err
}
func (self VerticesFormatXyzRgb) GlEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, ENDIANNES, self)
	return buf.Bytes(), err
}

type ElementIndicesUbyte []gl.GLubyte

func (self ElementIndicesUbyte) GlEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, ENDIANNES, self)
	return buf.Bytes(), err
}

func (self ElementIndicesUbyte) Len() int {
	return len(self)
}
func (self ElementIndicesUbyte) GlType() gl.GLenum {
	return gl.UNSIGNED_BYTE
}

type Drawer interface {
	Draw()
}

type DrawElements struct {
	Primitive gl.GLenum
	Elements  ElementIndexFormat
}

func (self DrawElements) Draw() {
	gl.DrawElements(
		self.Primitive,
		self.Elements.Len(),
		self.Elements.GlType(),
		nil, // Indices are in a buffer, never give them directly.
	)
}

func MakeDrawable(programs Programs, srefs ShaderRefs, vertices VertexFormat, indices ElementIndexFormat, binding_point uint) Drawable {
	vao := gl.GenVertexArray()
	vao.Bind()

	vbo := gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)
	if data, err := vertices.GlEncode(); err != nil {
		panic(err)
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, len(data), &data[0], gl.STATIC_DRAW)
	}

	ebo := gl.GenBuffer()
	ebo.Bind(gl.ELEMENT_ARRAY_BUFFER)
	if data, err := indices.GlEncode(); err != nil {
		panic(err)
	} else {
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(data), &data[0], gl.STATIC_DRAW)
	}

	program, err := programs.Serve(srefs)
	if err != nil {
		panic(err)
	}

	atts_names := vertices.Names()
	atts := make([]gl.AttribLocation, len(atts_names))
	for i, att_name := range atts_names {
		atts[i] = program.GetAttribLocation(att_name)
		atts[i].EnableArray()
	}
	vertices.AttribPointers(atts)

	model_matrix_uniform := program.GetUniformLocation("model_matrix")
	ubi := program.GetUniformBlockIndex("GlobalMatrices")
	if err := CheckGlError(); err != nil {
		err.Description = "GetUniformBlockIndex"
		panic(err)
	}
	if ubi == gl.INVALID_INDEX {
		fmt.Println("INVALID_INDEX")
	}
	program.UniformBlockBinding(ubi, binding_point)
	if err := CheckGlError(); err != nil {
		err.Description = "UniformBlockBinding"
		panic(err)
	}
	vbo.Unbind(gl.ARRAY_BUFFER)
	return Drawable{gl.TRIANGLE_STRIP, vao, model_matrix_uniform, srefs, indices.Len()}
}

func Cube(programs Programs, binding_point uint) Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := VerticesFormatXyz{
		VertexFormatXyz{m, m, 0},
		VertexFormatXyz{m, m, 1},
		VertexFormatXyz{m, p, 0},
		VertexFormatXyz{m, p, 1},
		VertexFormatXyz{p, m, 0},
		VertexFormatXyz{p, m, 1},
		VertexFormatXyz{p, p, 0},
		VertexFormatXyz{p, p, 1},
	}
	// Indices for triangle strip adapted from
	// http://www.cs.umd.edu/gvil/papers/av_ts.pdf .
	// I mirrored their cube to have CCW, and I used a natural order to
	// number the vertices (see above, it's binary code).
	indices := ElementIndicesUbyte{6, 2, 7, 3, 1, 2, 0, 6, 4, 7, 5, 1, 4, 0}
	srefs := ShaderRefs{VSH_POS3, FSH_ZRED}
	return MakeDrawable(programs, srefs, vertices, indices, binding_point)
}

func Pyramid(programs Programs, binding_point uint) Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := VerticesFormatXyz{
		VertexFormatXyz{m, m, 0},
		VertexFormatXyz{m, p, 0},
		VertexFormatXyz{p, m, 0},
		VertexFormatXyz{p, p, 0},
		VertexFormatXyz{0, 0, 1},
	}
	indices := ElementIndicesUbyte{1, 4, 3, 2, 1, 0, 4, 2}
	srefs := ShaderRefs{VSH_POS3, FSH_ZGREEN}
	return MakeDrawable(programs, srefs, vertices, indices, binding_point)
}

func Floor(programs Programs, binding_point uint) Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := VerticesFormatXyzRgb{
		// x y z r v b
		VertexFormatXyzRgb{m, m, 0, .1, .1, .5},
		VertexFormatXyzRgb{m, p, 0, .1, .1, .5},
		VertexFormatXyzRgb{p, m, 0, 0, 1, 0},
		VertexFormatXyzRgb{p, p, 0, 1, 0, 0},
	}
	indices := ElementIndicesUbyte{
		0, 2, 1, 3,
	}
	srefs := ShaderRefs{VSH_COL3, FSH_VCOL}
	return MakeDrawable(programs, srefs, vertices, indices, binding_point)
}

func Ceiling(programs Programs, binding_point uint) Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := VerticesFormatXyzRgb{
		VertexFormatXyzRgb{m, m, 1, .1, .1, .5},
		VertexFormatXyzRgb{m, p, 1, .1, .1, .5},
		VertexFormatXyzRgb{p, m, 1, 0, 1, 0},
		VertexFormatXyzRgb{p, p, 1, 1, 0, 0},
	}
	indices := ElementIndicesUbyte{0, 1, 2, 3}
	srefs := ShaderRefs{VSH_COL3, FSH_VCOL}
	return MakeDrawable(programs, srefs, vertices, indices, binding_point)
}

func Wall(programs Programs, binding_point uint) Drawable {
	// The wall meshes are relative to the center of the tile to which they belong.
	// They are given for a facing of 0 (east), therefore this mesh depicts a
	// western wall.
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := VerticesFormatXyzRgb{
		VertexFormatXyzRgb{-.4, m, 0, .1, .1, .5},
		VertexFormatXyzRgb{-.4, m, 1, 0, 1, 0},
		VertexFormatXyzRgb{-.4, p, 0, .1, .1, .5},
		VertexFormatXyzRgb{-.4, p, 1, 1, 0, 0},
	}
	indices := ElementIndicesUbyte{0, 2, 1, 3}
	srefs := ShaderRefs{VSH_COL3, FSH_VCOL}
	return MakeDrawable(programs, srefs, vertices, indices, binding_point)
}

func Column(programs Programs, binding_point uint) Drawable {
	const p = .15 // Plus sign.
	const m = -p  // Minus sign.
	vertices := VerticesFormatXyz{
		VertexFormatXyz{m, m, 0},
		VertexFormatXyz{m, m, 1},
		VertexFormatXyz{m, p, 0},
		VertexFormatXyz{m, p, 1},
		VertexFormatXyz{p, m, 0},
		VertexFormatXyz{p, m, 1},
		VertexFormatXyz{p, p, 0},
		VertexFormatXyz{p, p, 1},
	}
	// Indices for triangle strip adapted from
	// http://www.cs.umd.edu/gvil/papers/av_ts.pdf .
	// I mirrored their cube to have CCW, and I used a natural order to
	// number the vertices (see above, it's binary code).
	indices := ElementIndicesUbyte{1, 0, 5, 4, 7, 6, 3, 2, 1, 0}
	srefs := ShaderRefs{VSH_POS3, FSH_ZRED}
	return MakeDrawable(programs, srefs, vertices, indices, binding_point)
}

func DynaPyramid(programs Programs) StreamDrawable {
	var vertices [5 * 3]gl.GLfloat
	indices := [...]gl.GLubyte{
		1, 4, 3, 2, 1, 0, 4, 2,
	}
	vao := gl.GenVertexArray()
	vao.Bind()
	vbo := gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), nil, gl.DYNAMIC_DRAW)

	ebo := gl.GenBuffer()
	ebo.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices)), &indices, gl.STATIC_DRAW)

	srefs := ShaderRefs{VSH_POS3, FSH_ZGREEN}
	program, err := programs.Serve(srefs)
	if err != nil {
		panic(err)
	}

	att := program.GetAttribLocation("vpos")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		nil)
	model_matrix_uniform := program.GetUniformLocation("model_matrix")
	vbo.Unbind(gl.ARRAY_BUFFER)
	var result StreamDrawable
	result.Drawable.primitive = gl.TRIANGLE_STRIP
	result.Drawable.vao = vao
	result.Drawable.model_matrix_uniform = model_matrix_uniform
	result.Drawable.shaders_refs = srefs
	result.Drawable.n_elements = len(indices)
	result.vbo = vbo
	result.vbosize = int(unsafe.Sizeof(vertices))
	return result
}

func Monster(programs Programs) Drawable {
	const h = .5
	const w = .2
	const d = .15
	const l = w / 2       // left
	const r = -w / 2      // right
	const b = -d / 2      // back
	const f = d / 2       // front
	const c = h * 4 / 5   // chin height
	const n = f * 3       // Nose front
	const N = (h + c) / 2 // Nose height
	const n_verts = 7     // Number of vertices.
	const cpv = 6         // Components per vertex.
	vertices := [n_verts * cpv]gl.GLfloat{
		b, l, 0, 0.40, 0.20, 0.00,
		b, r, 0, 0.40, 0.20, 0.00,
		b, r, h, 0.80, 0.40, 0.00,
		b, l, h, 0.80, 0.40, 0.00,
		f, 0, 0, 0.40, 0.20, 0.00,
		f, 0, c, 0.40, 0.20, 0.00,
		n, 0, N, 1.00, 0.80, 0.00,
	}
	indices := [...]gl.GLubyte{
		// Bottom.
		0, 4, 1,
		// Back.
		0, 1, 2,
		0, 2, 3,
		// Left body.
		5, 2, 1,
		5, 1, 4,
		// Right body.
		5, 0, 3,
		5, 4, 0,
		// Left face.
		5, 6, 2,
		// Right face.
		5, 3, 6,
		// Top head.
		2, 6, 3,
	}
	vao := gl.GenVertexArray()
	vao.Bind()
	vbo := gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)

	ebo := gl.GenBuffer()
	ebo.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices)), &indices, gl.STATIC_DRAW)

	srefs := ShaderRefs{VSH_COL3, FSH_VCOL}
	program, err := programs.Serve(srefs)
	if err != nil {
		panic(err)
	}

	att := program.GetAttribLocation("vpos")
	if err := CheckGlError(); err != nil {
		err.Description = "get attrib location 1"
		panic(err)
	}
	att.EnableArray()
	if err := CheckGlError(); err != nil {
		err.Description = "enable array 1"
		panic(err)
	}
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		6*4,
		uintptr(0))
	if err := CheckGlError(); err != nil {
		err.Description = "attrib pointer 1"
		panic(err)
	}
	att = program.GetAttribLocation("vcol")
	if err := CheckGlError(); err != nil {
		err.Description = "get attrib location 2"
		panic(err)
	}
	att.EnableArray()
	if err := CheckGlError(); err != nil {
		err.Description = "enable array 2"
		panic(err)
	}
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		6*4,
		uintptr(3*4))
	if err := CheckGlError(); err != nil {
		err.Description = "attrib pointer 2"
		panic(err)
	}
	model_matrix_uniform := program.GetUniformLocation("model_matrix")
	if err := CheckGlError(); err != nil {
		err.Description = "Get uniform mvp"
		panic(err)
	}

	vbo.Unbind(gl.ARRAY_BUFFER)

	if err := CheckGlError(); err != nil {
		err.Description = "vbo unbind"
		panic(err)
	}
	return Drawable{gl.TRIANGLES, vao, model_matrix_uniform, srefs, len(indices)}
}
