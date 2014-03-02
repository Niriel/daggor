package batch

import (
	"sculpt"
)

type DrawBatch struct {
	drawer sculpt.Drawer
}

func MakeDrawBatch(drawer sculpt.Drawer) DrawBatch {
	return DrawBatch{drawer: drawer}
}

func (batch DrawBatch) Enter() {}
func (batch DrawBatch) Exit()  {}
func (batch DrawBatch) Run() {
	batch.drawer.Draw()
}
