package glw

import (
	"encoding/binary"
	"unsafe"
)

var endianness binary.ByteOrder

func init() {
	// Here, somehow detect the endianness of the system.
	// We need to encode vertex data before sending it to OpenGL, and OpenGL
	// uses the endianness of the system on which we are running.
	// Why encode? Because we do not want to care about the way Go allign its
	// data in memory.  Encoding in binary ensures that the data has no gap or
	// padding, which is easier to manage.

	// Endianness detection.
	// Credit matt kane, taken from his gosndfile project.
	// https://groups.google.com/forum/#!msg/golang-nuts/3GEzwKfRRQw/D1bMbFP-ClAJ
	// https://github.com/mkb218/gosndfile
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	isLittleEndian := b == 0x04

	if isLittleEndian {
		endianness = binary.LittleEndian
	} else {
		endianness = binary.BigEndian
	}
}
