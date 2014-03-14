package glw

import (
	"fmt"
	"github.com/go-gl/gl"
	"glm"
	"unsafe"
)

//=============================================================================

// This section defines the various vertex formats used in the game.
// It is at the level of a single vertex, not a collection of vertices.

// VertexXyz defines vertices that have a location and nothing else.
// No color, UV or any other parameter.
type VertexXyz struct {
	X, Y, Z gl.GLfloat
}

// VertexXyzRgb defines vertices that have a location and a color.
// No UV information.  Note that there is no Alpha component to the color.
type VertexXyzRgb struct {
	X, Y, Z gl.GLfloat
	R, G, B gl.GLfloat
}

// VertexXyzNor defines vertices that have a position and a normal.
// It would be nice to compress the normals in a gl_int_2_10_10_10 once I figure
// out how to do it.  The signs confuse me, as well as the fact that if
// normalized then no value seems to be able to represent 0.
type VertexXyzNor struct {
	Px, Py, Pz gl.GLfloat
	Nx, Ny, Nz gl.GLfloat
}

func MakeVertexXyzNor(pos, nor glm.Vector3) VertexXyzNor {
	px64, py64, pz64 := pos.Xyz()
	nx64, ny64, nz64 := nor.Xyz()
	return VertexXyzNor{
		Px: gl.GLfloat(px64),
		Py: gl.GLfloat(py64),
		Pz: gl.GLfloat(pz64),
		Nx: gl.GLfloat(nx64),
		Ny: gl.GLfloat(ny64),
		Nz: gl.GLfloat(nz64),
	}
}

type VertexXyzNorUv struct {
	Px, Py, Pz gl.GLfloat
	Nx, Ny, Nz gl.GLfloat
	U, V       gl.GLfloat
}

// ModelMatInstance defines transformation matrices for instanced rendering.
// This is a per-instance attribute, but it is read by the shader as a
// per-vertex attribute.  It requires a VBO.
type ModelMatInstance struct {
	M [16]gl.GLfloat
}

//=============================================================================

// Concrete classes of Vertices derive from baseBuffer and correspond to
// a specific vertex format.

type VerticesXyz struct {
	baseBuffer
	data []VertexXyz
}

type VerticesXyzRgb struct {
	baseBuffer
	data []VertexXyzRgb
}

type VerticesXyzNor struct {
	baseBuffer
	data []VertexXyzNor
}

type VerticesXyzNorUv struct {
	baseBuffer
	data []VertexXyzNorUv
}

type ModelMatInstances struct {
	baseBuffer
	data []ModelMatInstance
}

func NewVerticesXyz(usage gl.GLenum) *VerticesXyz {
	buffer := new(VerticesXyz)
	buffer.gen(gl.ARRAY_BUFFER, usage)
	return buffer
}

func NewVerticesXyzRgb(usage gl.GLenum) *VerticesXyzRgb {
	buffer := new(VerticesXyzRgb)
	buffer.gen(gl.ARRAY_BUFFER, usage)
	return buffer
}

func NewVerticesXyzNor(usage gl.GLenum) *VerticesXyzNor {
	buffer := new(VerticesXyzNor)
	buffer.gen(gl.ARRAY_BUFFER, usage)
	return buffer
}

func NewVerticesXyzNorUv(usage gl.GLenum) *VerticesXyzNorUv {
	buffer := new(VerticesXyzNorUv)
	buffer.gen(gl.ARRAY_BUFFER, usage)
	return buffer
}

func NewModelMatInstances(usage gl.GLenum) *ModelMatInstances {
	buffer := new(ModelMatInstances)
	buffer.gen(gl.ARRAY_BUFFER, usage)
	return buffer
}

type locationDataSetter interface {
	SetLocationData([]glm.Matrix4)
}

//-----------------------------------------------------------------------------

func (buffer *VerticesXyz) SetData(vd []VertexXyz) {
	buffer.data = make([]VertexXyz, len(vd), len(vd))
	copy(buffer.data, vd)
	buffer.bufferdataClean = false
}

func (buffer *VerticesXyzRgb) SetData(vd []VertexXyzRgb) {
	buffer.data = make([]VertexXyzRgb, len(vd), len(vd))
	copy(buffer.data, vd)
	buffer.bufferdataClean = false
}

func (buffer *VerticesXyzNor) SetData(vd []VertexXyzNor) {
	buffer.data = make([]VertexXyzNor, len(vd), len(vd))
	copy(buffer.data, vd)
	buffer.bufferdataClean = false
}

func (buffer *VerticesXyzNorUv) SetData(vd []VertexXyzNorUv) {
	buffer.data = make([]VertexXyzNorUv, len(vd), len(vd))
	copy(buffer.data, vd)
	buffer.bufferdataClean = false
}

func (buffer *ModelMatInstances) SetLocationData(locations []glm.Matrix4) {
	buffer.data = make([]ModelMatInstance, len(locations), len(locations))
	for i, m := range locations {
		buffer.data[i] = ModelMatInstance{m.GlFloats()}
	}
	buffer.bufferdataClean = false
}

func (buffer *VerticesXyz) Update() {
	buffer.update(buffer.data)
}

func (buffer *VerticesXyzRgb) Update() {
	buffer.update(buffer.data)
}

func (buffer *VerticesXyzNor) Update() {
	buffer.update(buffer.data)
}

func (buffer *VerticesXyzNorUv) Update() {
	buffer.update(buffer.data)
}

func (buffer *ModelMatInstances) Update() {
	buffer.update(buffer.data)
}

//=============================================================================

// This section is about configuring the Vertex Array Object of a mesh.
// The mesh calls SetUpVao which is a method of Vertices.
// SetUpVao needs a gl Program object in order to query variable parameter
// names and all.

// Since there is a lot of code common to all the Vertices object about how
// to set up a VAO, we just ask each Vertices object to satisfy an interface
// containing what is specific to that Vertices object.  Then we pass it to
// a generic function.

func (buffer *VerticesXyz) SetUpVao(program gl.Program) {
	bufferSetUpVao(buffer, program)
}
func (buffer *VerticesXyzRgb) SetUpVao(program gl.Program) {
	bufferSetUpVao(buffer, program)
}
func (buffer *VerticesXyzNor) SetUpVao(program gl.Program) {
	bufferSetUpVao(buffer, program)
}
func (buffer *VerticesXyzNorUv) SetUpVao(program gl.Program) {
	bufferSetUpVao(buffer, program)
}
func (buffer *ModelMatInstances) SetUpVao(program gl.Program) {
	bufferSetUpVao(buffer, program)
}

// The bufferSetUpVaoInt contains everything that is needed by the function
// bufferSetUpVao.
type bufferSetUpVaoInt interface {
	names() []string
	attribPointers([]gl.AttribLocation)
	bind()
	unbind()
}

func bufferSetUpVao(buffer bufferSetUpVaoInt, program gl.Program) {
	buffer.bind()
	// Collect the attrib locations for each attrib name.
	atts_names := buffer.names() // Expected GLSL variable names.
	atts := make([]gl.AttribLocation, len(atts_names))
	for i, att_name := range atts_names {
		atts[i] = program.GetAttribLocation(att_name)
		if err := CheckGlError(); err != nil {
			err.Description = fmt.Sprintf("program.GetAttribLocation(%#v)", att_name)
			panic(err)
		}
		if atts[i] == -1 {
			panic(fmt.Sprintf("attrib location %#v not found", att_name))
		}
	}
	// Now that the locations are known, we can relate them to vertex data.
	buffer.attribPointers(atts)
	buffer.unbind()
}

func (buffer *VerticesXyz) names() []string {
	return []string{"vpos"}
}
func (buffer *VerticesXyzRgb) names() []string {
	return []string{"vpos", "vcol"}
}
func (buffer *VerticesXyzNor) names() []string {
	return []string{"vpos", "vnor"}
}
func (buffer *VerticesXyzNorUv) names() []string {
	return []string{"vpos", "vnor", "vuv"}
}
func (buffer *ModelMatInstances) names() []string {
	return []string{"model_to_eye"}
}

func (buffer *VerticesXyz) attribPointers(atts []gl.AttribLocation) {
	const FLOATSIZE = unsafe.Sizeof(gl.GLfloat(0))
	const NB_COORDS = 3 // x y and z.
	const COORDS_SIZE = NB_COORDS * FLOATSIZE
	const COORDS_OFS = uintptr(0)
	const TOTAL_SIZE = int(COORDS_SIZE)
	atts[0].AttribPointer(NB_COORDS, gl.FLOAT, false, TOTAL_SIZE, COORDS_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesXyz atts[0].AttribPointer"
		panic(err)
	}
	for _, att := range atts {
		att.EnableArray()
		if err := CheckGlError(); err != nil {
			err.Description = fmt.Sprintf("atts[%v].EnableArray()\n", att)
			panic(err)
		}
	}
}
func (buffer *VerticesXyzRgb) attribPointers(atts []gl.AttribLocation) {
	const FLOATSIZE = unsafe.Sizeof(gl.GLfloat(0))
	const NB_COORDS = 3 // x y and z.
	const NB_COLORS = 3 // r g and b.
	const COORDS_SIZE = NB_COORDS * FLOATSIZE
	const COLORS_SIZE = NB_COLORS * FLOATSIZE
	const COORDS_OFS = uintptr(0)
	const COLORS_OFS = uintptr(COORDS_SIZE)
	const TOTAL_SIZE = int(COORDS_SIZE + COLORS_SIZE)
	atts[0].AttribPointer(NB_COORDS, gl.FLOAT, false, TOTAL_SIZE, COORDS_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesXyzRgb atts[0].AttribPointer"
		panic(err)
	}
	atts[1].AttribPointer(NB_COLORS, gl.FLOAT, false, TOTAL_SIZE, COLORS_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesXyzRgb atts[1].AttribPointer"
		panic(err)
	}
	for _, att := range atts {
		att.EnableArray()
		if err := CheckGlError(); err != nil {
			err.Description = fmt.Sprintf("atts[%v].EnableArray()\n", att)
			panic(err)
		}
	}
}
func (buffer *VerticesXyzNor) attribPointers(atts []gl.AttribLocation) {
	const FLOATSIZE = unsafe.Sizeof(gl.GLfloat(0))
	const NB_POS = 3 // px py and pz.
	const NB_NOR = 3 // nx ny and nz.
	const POS_SIZE = NB_POS * FLOATSIZE
	const NOR_SIZE = NB_NOR * FLOATSIZE
	const POS_OFS = uintptr(0)
	const NOR_OFS = uintptr(POS_SIZE)
	const TOTAL_SIZE = int(POS_SIZE + NOR_SIZE)
	atts[0].AttribPointer(NB_POS, gl.FLOAT, false, TOTAL_SIZE, POS_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesXyzNor atts[0].AttribPointer"
		panic(err)
	}
	atts[1].AttribPointer(NB_NOR, gl.FLOAT, false, TOTAL_SIZE, NOR_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesXyzNor atts[1].AttribPointer"
		panic(err)
	}
	for _, att := range atts {
		att.EnableArray()
		if err := CheckGlError(); err != nil {
			err.Description = fmt.Sprintf("atts[%v].EnableArray()\n", att)
			panic(err)
		}
	}
}
func (buffer *VerticesXyzNorUv) attribPointers(atts []gl.AttribLocation) {
	const FLOATSIZE = unsafe.Sizeof(gl.GLfloat(0))
	const NB_POS = 3 // px py and pz.
	const NB_NOR = 3 // nx ny and nz.
	const NB_UV = 2  // u and v.
	const POS_SIZE = NB_POS * FLOATSIZE
	const NOR_SIZE = NB_NOR * FLOATSIZE
	const UV_SIZE = NB_UV * FLOATSIZE
	const POS_OFS = uintptr(0)
	const NOR_OFS = uintptr(POS_SIZE)
	const UV_OFS = uintptr(POS_SIZE + NOR_SIZE)
	const TOTAL_SIZE = int(POS_SIZE + NOR_SIZE + UV_SIZE)
	atts[0].AttribPointer(NB_POS, gl.FLOAT, false, TOTAL_SIZE, POS_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesXyzNorUv atts[0].AttribPointer"
		panic(err)
	}
	atts[1].AttribPointer(NB_NOR, gl.FLOAT, false, TOTAL_SIZE, NOR_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesXyzNorUv atts[1].AttribPointer"
		panic(err)
	}
	atts[2].AttribPointer(NB_UV, gl.FLOAT, false, TOTAL_SIZE, UV_OFS)
	if err := CheckGlError(); err != nil {
		err.Description = "VerticesXyzNorUv atts[2].AttribPointer"
		panic(err)
	}
	for _, att := range atts {
		att.EnableArray()
		if err := CheckGlError(); err != nil {
			err.Description = fmt.Sprintf("atts[%v].EnableArray()\n", att)
			panic(err)
		}
	}
}
func (buffer *ModelMatInstances) attribPointers(atts []gl.AttribLocation) {
	const FLOATSIZE = unsafe.Sizeof(gl.GLfloat(0))
	const NB_COORDS = 4 // 4 floats per matrix row.
	const NB_ATTS = 4   // Matrix is 4 rows or 4 columns.
	const COORDS_SIZE = NB_COORDS * FLOATSIZE
	const COORDS_OFS = uintptr(0)
	const TOTAL_SIZE = int(COORDS_SIZE) * NB_ATTS
	for i := 0; i < NB_ATTS; i++ {
		// We pass each column of the matrix separately.
		// Because that's how OpenGL does matrix vertex attributes.
		att := atts[0] + gl.AttribLocation(i)
		offset := COORDS_OFS + uintptr(i)*COORDS_SIZE
		att.AttribPointer(NB_COORDS, gl.FLOAT, false, TOTAL_SIZE, offset)
		if err := CheckGlError(); err != nil {
			err.Description = "ModelMatInstances att.AttribPointer"
			panic(err)
		}
		// 1 here means that we switch to a new matrix every 1 instance.
		// This AttribDivisor call with a non-zero value is what makes the
		// attribute instanced.
		att.AttribDivisor(1)
		if err := CheckGlError(); err != nil {
			err.Description = "ModelMatInstances att.AttribDivisor"
			panic(err)
		}
		// Each column of the matrix must be enabled.
		att.EnableArray()
		if err := CheckGlError(); err != nil {
			err.Description = fmt.Sprintf("atts[%v].EnableArray()\n", att)
			panic(err)
		}
	}
}
