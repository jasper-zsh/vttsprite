package swscale

//#cgo pkg-config: libswscale
//
//#include <libswscale/swscale.h>
import "C"
import (
	"jasper-zsh/vttsprite/avutil"
	"unsafe"
)

type (
	SwsContext C.struct_SwsContext
	SwsFilter  C.struct_SwsFilter
)

const (
	SWS_BILINEAR = 2
)

func SwsGetContext(
	srcW int, srcH int, srcFormat avutil.AVPixelFormat,
	dstW int, dstH int, dstFormat avutil.AVPixelFormat,
	flags int, srcFilter *SwsFilter,
	dstFilter *SwsFilter, param *float64,
) *SwsContext {
	return (*SwsContext)(C.sws_getContext(
		C.int(srcW), C.int(srcH), int32(srcFormat),
		C.int(dstW), C.int(dstH), int32(dstFormat),
		C.int(flags), (*C.struct_SwsFilter)(unsafe.Pointer(srcFilter)),
		(*C.struct_SwsFilter)(unsafe.Pointer(dstFilter)), (*C.double)(param),
	))
}

func SwsScale(
	c *SwsContext, srcSlice unsafe.Pointer,
	srcStride []int, srcSliceY int, srcSliceH int,
	dst unsafe.Pointer, dstStride []int,
) int {
	return int(C.sws_scale(
		(*C.struct_SwsContext)(unsafe.Pointer(c)),
		(**C.uchar)(srcSlice),
		(*C.int)(unsafe.Pointer(&srcStride)),
		C.int(srcSliceY),
		C.int(srcSliceH),
		(**C.uchar)(dst),
		(*C.int)(unsafe.Pointer(&dstStride)),
	))
}
