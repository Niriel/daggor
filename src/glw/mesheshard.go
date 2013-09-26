// This module contains the code for hard coded meshes used during
// the development and debug phases of the game.

package glw

import (
	"github.com/go-gl/gl"
	"unsafe"
)

func Cube(programs Programs) Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := [...]gl.GLfloat{
		m, m, 0,
		m, m, 1,
		m, p, 0,
		m, p, 1,
		p, m, 0,
		p, m, 1,
		p, p, 0,
		p, p, 1,
	}
	// Indices for triangle strip adapted from
	// http://www.cs.umd.edu/gvil/papers/av_ts.pdf .
	// I mirrored their cube to have CCW, and I used a natural order to
	// number the vertices (see above, it's binary code).
	indices := [...]gl.GLubyte{
		6, 2, 7, 3, 1, 2, 0, 6, 4, 7, 5, 1, 4, 0,
	}
	vao := gl.GenVertexArray()
	vao.Bind()
	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)
	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices)), &indices, gl.STATIC_DRAW)

	srefs := ShaderRefs{VSH_POS3, FSH_ZRED}
	program, err := programs.Serve(srefs)
	if err != nil {
		panic(err)
	}

	program.Use()

	att := program.GetAttribLocation("vpos")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		nil)
	mvp := program.GetUniformLocation("mvp")
	vbuf.Unbind(gl.ARRAY_BUFFER)
	return Drawable{gl.TRIANGLE_STRIP, vao, mvp, srefs, len(indices)}
}

func Pyramid(programs Programs) Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := [...]gl.GLfloat{
		m, m, 0,
		m, p, 0,
		p, m, 0,
		p, p, 0,
		0, 0, 1,
	}
	indices := [...]gl.GLubyte{
		1, 4, 3, 2, 1, 0, 4, 2,
	}
	vao := gl.GenVertexArray()
	vao.Bind()
	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)

	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
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
	mvp := program.GetUniformLocation("mvp")
	vbuf.Unbind(gl.ARRAY_BUFFER)
	return Drawable{gl.TRIANGLE_STRIP, vao, mvp, srefs, len(indices)}
}

func Floor(programs Programs) Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := [...]gl.GLfloat{
		// x y z r v b
		m, m, 0, .1, .1, .5,
		m, p, 0, .1, .1, .5,
		p, m, 0, 0, 1, 0,
		p, p, 0, 1, 0, 0,
	}
	indices := [...]gl.GLubyte{
		0, 2, 1, 3,
	}
	vao := gl.GenVertexArray()
	vao.Bind()

	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)

	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices)), &indices, gl.STATIC_DRAW)

	srefs := ShaderRefs{VSH_COL3, FSH_VCOL}
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
		6*4,
		uintptr(0))
	att = program.GetAttribLocation("vcol")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		6*4,
		uintptr(3*4))
	mvp := program.GetUniformLocation("mvp")

	vbuf.Unbind(gl.ARRAY_BUFFER)
	return Drawable{gl.TRIANGLE_STRIP, vao, mvp, srefs, len(indices)}
}

func Ceiling(programs Programs) Drawable {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := [...]gl.GLfloat{
		// x y z r v b
		m, m, 1, .1, .1, .5,
		m, p, 1, .1, .1, .5,
		p, m, 1, 0, 1, 0,
		p, p, 1, 1, 0, 0,
	}
	indices := [...]gl.GLubyte{
		0, 1, 2, 3,
	}
	vao := gl.GenVertexArray()
	vao.Bind()

	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)

	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices)), &indices, gl.STATIC_DRAW)

	srefs := ShaderRefs{VSH_COL3, FSH_VCOL}
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
		6*4,
		uintptr(0))
	att = program.GetAttribLocation("vcol")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		6*4,
		uintptr(3*4))
	mvp := program.GetUniformLocation("mvp")

	vbuf.Unbind(gl.ARRAY_BUFFER)
	return Drawable{gl.TRIANGLE_STRIP, vao, mvp, srefs, len(indices)}
}

func Wall(programs Programs) Drawable {
	// The wall meshes are relative to the center of the tile to which they belong.
	// They are given for a facing of 0 (east), therefore this mesh depicts a
	// western wall.
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertices := [...]gl.GLfloat{
		// x y z r v b
		-.4, m, 0, .1, .1, .5,
		-.4, m, 1, 0, 1, 0,
		-.4, p, 0, .1, .1, .5,
		-.4, p, 1, 1, 0, 0,
	}
	indices := [...]gl.GLubyte{
		0, 2, 1, 3,
	}
	vao := gl.GenVertexArray()
	vao.Bind()

	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)

	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices)), &indices, gl.STATIC_DRAW)

	srefs := ShaderRefs{VSH_COL3, FSH_VCOL}
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
		6*4,
		uintptr(0))
	att = program.GetAttribLocation("vcol")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		6*4,
		uintptr(3*4))
	mvp := program.GetUniformLocation("mvp")

	vbuf.Unbind(gl.ARRAY_BUFFER)
	return Drawable{gl.TRIANGLE_STRIP, vao, mvp, srefs, len(indices)}
}

func Column(programs Programs) Drawable {
	const p = .15 // Plus sign.
	const m = -p  // Minus sign.
	vertices := [...]gl.GLfloat{
		m, m, 0,
		m, m, 1,
		m, p, 0,
		m, p, 1,
		p, m, 0,
		p, m, 1,
		p, p, 0,
		p, p, 1,
	}
	// Indices for triangle strip adapted from
	// http://www.cs.umd.edu/gvil/papers/av_ts.pdf .
	// I mirrored their cube to have CCW, and I used a natural order to
	// number the vertices (see above, it's binary code).
	indices := [...]gl.GLubyte{
		1, 0, 5, 4, 7, 6, 3, 2, 1, 0,
	}
	vao := gl.GenVertexArray()
	vao.Bind()
	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)
	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices)), &indices, gl.STATIC_DRAW)

	srefs := ShaderRefs{VSH_POS3, FSH_ZRED}
	program, err := programs.Serve(srefs)
	if err != nil {
		panic(err)
	}

	program.Use()

	att := program.GetAttribLocation("vpos")
	att.EnableArray()
	att.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		nil)
	mvp := program.GetUniformLocation("mvp")
	vbuf.Unbind(gl.ARRAY_BUFFER)
	return Drawable{gl.TRIANGLE_STRIP, vao, mvp, srefs, len(indices)}
}

func DynaPyramid(programs Programs) StreamDrawable {
	var vertices [5 * 3]gl.GLfloat
	indices := [...]gl.GLubyte{
		1, 4, 3, 2, 1, 0, 4, 2,
	}
	vao := gl.GenVertexArray()
	vao.Bind()
	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), nil, gl.DYNAMIC_DRAW)

	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
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
	mvp := program.GetUniformLocation("mvp")
	vbuf.Unbind(gl.ARRAY_BUFFER)
	var result StreamDrawable
	result.Drawable.primitive = gl.TRIANGLE_STRIP
	result.Drawable.vao = vao
	result.Drawable.mvp = mvp
	result.Drawable.shaders_refs = srefs
	result.Drawable.n_elements = len(indices)
	result.vbo = vbuf
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
	vbuf := gl.GenBuffer()
	vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), &vertices, gl.STATIC_DRAW)

	ebuf := gl.GenBuffer()
	ebuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
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
	mvp := program.GetUniformLocation("mvp")
	if err := CheckGlError(); err != nil {
		err.Description = "Get uniform mvp"
		panic(err)
	}

	vbuf.Unbind(gl.ARRAY_BUFFER)

	if err := CheckGlError(); err != nil {
		err.Description = "vbo unbind"
		panic(err)
	}
	return Drawable{gl.TRIANGLES, vao, mvp, srefs, len(indices)}
}
