package glw

import (
	"github.com/go-gl/gl"
)

var ERROR_NAMES = map[gl.GLenum]string{
	gl.NO_ERROR:                      "no error",
	gl.INVALID_ENUM:                  "invalid enum",
	gl.INVALID_VALUE:                 "invalid value",
	gl.INVALID_OPERATION:             "invalid operation",
	gl.STACK_OVERFLOW:                "stack overflow",
	gl.STACK_UNDERFLOW:               "stack underflow",
	gl.OUT_OF_MEMORY:                 "out of memory",
	gl.INVALID_FRAMEBUFFER_OPERATION: "invalid frame buffer operation",
	gl.TABLE_TOO_LARGE:               "table too large",
}
