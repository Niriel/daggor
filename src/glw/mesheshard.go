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
	return Drawable{vao, mvp, srefs, len(indices)}
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
	return Drawable{vao, mvp, srefs, len(indices)}
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
	return Drawable{vao, mvp, srefs, len(indices)}
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
	return Drawable{vao, mvp, srefs, len(indices)}
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
	return Drawable{vao, mvp, srefs, len(indices)}
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
	return Drawable{vao, mvp, srefs, len(indices)}
}
