package sculpt

import (
	"github.com/go-gl/gl"
	"glw"
)

// Drawer interface abstracts all the possible OpenGL Draw calls.
// The Draw method takes no argument, it is assume that all the arguments are
// already curried.
type InstancedDrawer interface {
	Draw(primcount int)
}

type UninstancedDrawer interface {
	Draw()
}

// DrawElement contains the parameters required by gl.DrawElements.
type DrawElement struct {
	mode    gl.GLenum
	count   int
	typ     gl.GLenum
	indices interface{}
}

// This Draw method is a wrapper to gl.DrawElements.
func (drawer DrawElement) Draw() {
	gl.DrawElements(drawer.mode, drawer.count, drawer.typ, drawer.indices)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "gl.DrawElements"
		panic(err)
	}
}

// DrawElementInstance curries gl.DrawElementsInstanced.
type DrawElementInstanced struct {
	mode    gl.GLenum
	count   int
	typ     gl.GLenum
	indices interface{}
}

func (drawer DrawElementInstanced) Draw(primcount int) {
	gl.DrawElementsInstanced(
		drawer.mode,
		drawer.count,
		drawer.typ,
		drawer.indices,
		primcount,
	)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "gl.DrawElementsInstanced"
		panic(err)
	}
}
