package avutil

//#cgo pkg-config: libavutil
//
//#include <libavutil/avutil.h>
//#include <libavutil/frame.h>
//#include <libavutil/imgutils.h>
//#include <stdlib.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type (
	AVMediaType   C.enum_AVMediaType
	AVRational    C.struct_AVRational
	AVDictionary  C.struct_AVDictionary
	AVFrame       C.struct_AVFrame
	AVPixelFormat C.enum_AVPixelFormat
)

const (
	AVMEDIA_TYPE_VIDEO = C.AVMEDIA_TYPE_VIDEO
	AV_TIME_BASE       = C.AV_TIME_BASE
	AV_PIX_FMT_RGB24   = C.AV_PIX_FMT_RGB24
)

var (
	AV_TIME_BASE_Q = (AVRational)(C.AV_TIME_BASE_Q)
)

func (ra *AVRational) Den() int {
	return int(ra.den)
}

func (ra *AVRational) Num() int {
	return int(ra.num)
}

func (ra *AVRational) ToFloat64() float64 {
	return float64(ra.num) / float64(ra.den)
}

func AvFrameAlloc() *AVFrame {
	return (*AVFrame)(C.av_frame_alloc())
}

func AvFrameFree(frame *AVFrame) {
	var t *C.struct_AVFrame
	ptr := C.malloc(C.size_t(unsafe.Sizeof(t)))
	defer C.free(ptr)
	C.av_frame_free((**C.struct_AVFrame)(ptr))
}

func AvMalloc(size uint32) unsafe.Pointer {
	return C.av_malloc(C.ulong(size))
}

func AvRescaleQ(a int64, bq AVRational, cq AVRational) int64 {
	return int64(C.av_rescale_q(
		C.longlong(a),
		(C.struct_AVRational)(bq),
		(C.struct_AVRational)(cq),
	))
}

func (f *AVFrame) Linesize() [8]int {
	var ret [8]int
	for i := range f.linesize {
		ret[i] = int(f.linesize[i])
	}
	return ret
}

func (f *AVFrame) Width() int {
	return int(f.width)
}

func (f *AVFrame) Height() int {
	return int(f.height)
}

func (f *AVFrame) Data() [8]*C.uchar {
	return f.data
}

func FromCPtr(buf unsafe.Pointer, size int) (ret []uint8) {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	hdr.Cap = size
	hdr.Len = size
	hdr.Data = uintptr(buf)
	return
}

func AvImageGetBufferSize(pixFmt AVPixelFormat, width int, height int, align int) int {
	return int(C.av_image_get_buffer_size(
		int32(pixFmt),
		C.int(width),
		C.int(height),
		C.int(align),
	))
}

func AvImageAlloc(pointers unsafe.Pointer, linesizes []int, w int, h int, pixFmt AVPixelFormat, align int) int {
	return int(C.av_image_alloc(
		(**C.uchar)(pointers),
		(*C.int)(unsafe.Pointer(&linesizes)),
		C.int(w),
		C.int(h),
		int32(pixFmt),
		C.int(align),
	))
}

func AvImageFillArrays(dstData unsafe.Pointer, dstLinesize []int, src unsafe.Pointer, pixFmt AVPixelFormat, width int, height int, align int) int {
	return int(C.av_image_fill_arrays(
		(**C.uchar)(dstData),
		(*C.int)(unsafe.Pointer(&dstLinesize)),
		(*C.uchar)(src),
		int32(pixFmt),
		C.int(width),
		C.int(height),
		C.int(align),
	))
}
