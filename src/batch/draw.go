package batch

import (
	"glw"
)

type DrawBatch struct {
	drawer glw.Renderer
}

func MakeDrawBatch(drawer glw.Renderer) DrawBatch {
	return DrawBatch{drawer: drawer}
}

func (batch DrawBatch) Enter() {}
func (batch DrawBatch) Exit()  {}
func (batch DrawBatch) Run() {
	batch.drawer.Render(nil)
}
