package batch

import (
	"sculpt"
)

type DrawBatch struct {
	drawer sculpt.MeshDrawer
}

func MakeDrawBatch(drawer sculpt.MeshDrawer) DrawBatch {
	return DrawBatch{drawer: drawer}
}

func (batch DrawBatch) Enter() {}
func (batch DrawBatch) Exit()  {}
func (batch DrawBatch) Run() {
	batch.drawer.DrawMesh(nil)
}
