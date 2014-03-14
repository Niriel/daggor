package glw

import (
	"fmt"
	"github.com/go-gl/gl"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

var globaltexture gl.Texture

func LoadImage(filename string) (image.Image, error) {
	f, err := os.Open(filepath.Join(filepath.Dir(os.Args[0]), "..", "..", "textures", filename))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, format, err := image.Decode(f)
	fmt.Println(filename, "decoded as", format)
	return img, err
}

// forceRGBA prevents alpha-premultiplication by making Go believe it's already
// applied.
//
// Go thinks that everybody wants to premultiply their alpha values.  This is
// absolutely not true for us because:
//
// * Alpha premultiplying makes sense in linear, not in sRGB.
// * If my alpha is a refractive index I don't want it to be seen as some kind
//   of opacity.
// * Using alpha as an opacity parameter makes no physical sense anyway.
// * If I put something in a file, I want to get it back the way I put it.
//
// So, I do not want Go to do the conversion for me.
// The RGBA color model assumes that the alpha premultiplication is already
// done, and then Go doesn't try to do it itself.  Our trick is to take a
// NRGBA image content and rewrap it into an RGBA one.  By lying to Go about the
// fact that we applied the premultiplication ourselves, we ensure that Go does
// not do it.
//
// The specs of PNG insist on the fact that the data in the PNG is NOT
// premultiplied.  This is a good thing.  It lets us do whatever we like with
// out alpha and does not force its shit on us.  So, with this forceRGBA
// function, any png is file.  jpg are fine as well since jpg does not store
// alpha without hacks.
func forceRGBA(img image.Image) *image.RGBA {
	switch imgtype := img.(type) {
	case *image.RGBA:
		return imgtype
	case *image.NRGBA:
		result := image.NewRGBA(imgtype.Rect)
		result.Pix = imgtype.Pix
		result.Stride = imgtype.Stride
		return result
	default:
		panic("unexpected image type")
	}
}

// matrix2 Rotation matrix, can be improper to enable symmetry.
// It is column-major, like in opengl.
type matrix23 [6]int

func (m matrix23) mul(x, y int) (X, Y int) {
	X = m[0]*x + m[2]*y + m[4]
	Y = m[1]*x + m[3]*y + m[5]
	return
}

func GlRgba(img image.Image, rot matrix23) []gl.GLubyte {
	bounds := img.Bounds()
	sizeX := bounds.Size().X
	sizeY := bounds.Size().Y
	data := make([]gl.GLubyte, sizeX*sizeY*4)
	// x and y are in the source frame.
	// X and Y are in the destination frame.
	// x and y are not garanteed to start at 0 because that is how images work
	// in Go.  However, X and Y are because my output buffer is built the way
	// I want it.
	Y := -1
	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
		Y++
		X := -1
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			X++
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			rX, rY := rot.mul(X, Y)
			i := (rY*sizeX + rX) * 4 // Assumes square texture.
			//fmt.Println(X, Y, rX, rY, i)
			data[i+0] = gl.GLubyte(r >> 8)
			data[i+1] = gl.GLubyte(g >> 8)
			data[i+2] = gl.GLubyte(b >> 8)
			data[i+3] = gl.GLubyte(a >> 8)
		}
	}
	return data
}

func LoadTexture() {
	gl.ActiveTexture(gl.TEXTURE0)
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	t := gl.GenTexture()
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	const target = gl.TEXTURE_2D

	t.Bind(target)
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	img, err := LoadImage("test.jpg")
	if err != nil {
		panic(err)
	}

	data := GlRgba(img, matrix23{1, 0, 0, 1, 0, 0})
	size := img.Bounds().Size()

	gl.TexImage2D(
		target,
		0,                // Mipmap level.
		gl.SRGB8_ALPHA8,  // Format inside OpenGL.
		size.X,           // Width.
		size.Y,           // Height.
		0,                // Border.  Doc says it must be 0.
		gl.RGBA,          // Format of the data that...
		gl.UNSIGNED_BYTE, // ... I give to OpenGL.
		data,             // And the data itself.
	)
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	gl.GenerateMipmap(target)
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	if err := CheckGlError(); err != nil {
		panic(err)
	}
	gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	if err := CheckGlError(); err != nil {
		panic(err)
	}
	gl.TexParameteri(target, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	if err := CheckGlError(); err != nil {
		panic(err)
	}
	gl.TexParameteri(target, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	t.Unbind(target)

	globaltexture = t
}

func LoadSkybox() {
	gl.ActiveTexture(gl.TEXTURE0)
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	t := gl.GenTexture()
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	t.Bind(gl.TEXTURE_CUBE_MAP)
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	const filename = "Skybox_tut13_384x256.png"
	img, err := LoadImage(filename)
	if err != nil {
		panic(err)
	}

	size := img.Bounds().Size()
	w := size.X
	h := size.Y
	if w/3 != h/2 {
		panic("incorrect aspect ratio")
	}
	s := h / 2 // Size of one face of the cube.

	targets := [...]gl.GLenum{
		gl.TEXTURE_CUBE_MAP_NEGATIVE_X,
		gl.TEXTURE_CUBE_MAP_POSITIVE_X,
		gl.TEXTURE_CUBE_MAP_NEGATIVE_Y,
		gl.TEXTURE_CUBE_MAP_POSITIVE_Y,
		gl.TEXTURE_CUBE_MAP_NEGATIVE_Z,
		gl.TEXTURE_CUBE_MAP_POSITIVE_Z,
	}
	rs := []image.Rectangle{
		image.Rect(0*s, 0*s, 1*s, 1*s), // 0 side.
		image.Rect(1*s, 0*s, 2*s, 1*s), // 1 side.
		image.Rect(2*s, 0*s, 3*s, 1*s), // 2 side.
		image.Rect(0*s, 1*s, 1*s, 2*s), // 3 bottom.
		image.Rect(1*s, 1*s, 2*s, 2*s), // 4 top.
		image.Rect(2*s, 1*s, 3*s, 2*s), // 5 side.
	}
	// redirect must end with 3 and 4.  This places the floor at -z and the
	// sky at +z.
	// All the sides look properly oriented (grass down) on the +y side.
	// Considering that the sky is properly oriented the way it is, then
	// rectangle 5 contains the picture that matches it on the +y side.
	// Therefore, redirect must end with 5, 3, 4.
	// Problem now: all the other sides, and maybe the floor too, need to be
	// rotated in order to line up, that is having their sky on top.
	redirect := [...]int{0, 2, 1, 5, 3, 4}

	ID := matrix23{1, 0, 0, 1, 0, 0}
	ms := [...]matrix23{
		matrix23{0, -1, 1, 0, 0, s - 1},
		matrix23{0, 1, -1, 0, s - 1, 0},
		matrix23{-1, 0, 0, -1, s - 1, s - 1},
		ID,
		matrix23{-1, 0, 0, -1, s - 1, s - 1},
		ID,
	}

	rgba := forceRGBA(img)
	for i, target := range targets {
		subimage := rgba.SubImage(rs[redirect[i]])
		data := GlRgba(subimage, ms[i])
		gl.TexImage2D(
			target,
			0,                // Mipmap level.
			gl.SRGB8_ALPHA8,  // Format inside OpenGL.
			s,                // Width.
			s,                // Height.
			0,                // Border.  Doc says it must be 0.
			gl.RGBA,          // Format of the data that...
			gl.UNSIGNED_BYTE, // ... I give to OpenGL.
			data,             // And the data itself.
		)
		if err := CheckGlError(); err != nil {
			panic(err)
		}
	}

	const target = gl.TEXTURE_CUBE_MAP
	gl.GenerateMipmap(target)
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	if err := CheckGlError(); err != nil {
		panic(err)
	}
	gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	if err := CheckGlError(); err != nil {
		panic(err)
	}
	gl.TexParameteri(target, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	if err := CheckGlError(); err != nil {
		panic(err)
	}
	gl.TexParameteri(target, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	if err := CheckGlError(); err != nil {
		panic(err)
	}
	gl.TexParameteri(target, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	if err := CheckGlError(); err != nil {
		panic(err)
	}

	t.Unbind(target)

	globaltexture = t
}
