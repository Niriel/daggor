package batch

import (
	"glw"
)

type DrawBatch struct {
	drawer glw.MeshDrawer
}

func MakeDrawBatch(drawer glw.MeshDrawer) DrawBatch {
	return DrawBatch{drawer: drawer}
}

func (batch DrawBatch) Enter() {}
func (batch DrawBatch) Exit()  {}
func (batch DrawBatch) Run() {
	batch.drawer.DrawMesh(nil)
}
