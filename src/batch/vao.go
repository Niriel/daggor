package batch

import (
	"github.com/go-gl/gl"
	"glw"
)

type VaoBatch struct {
	BaseBatch
	vao gl.VertexArray
}

func MakeVaoBatch(vao gl.VertexArray) VaoBatch {
	return VaoBatch{vao: vao}
}

func (batch VaoBatch) Enter() {
	batch.vao.Bind()
	if err := glw.CheckGlError(); err != nil {
		err.Description = "VaoBatch.vao.Bind()"
		panic(err)
	}
}

func (batch VaoBatch) Exit() {
	gl.VertexArray(0).Bind()
	// Cannot fail with argument 0.
}
