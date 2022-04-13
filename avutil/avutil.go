package avutil

//#cgo pkg-config: libavutil
//
//#include <libavutil/avutil.h>
//#include <libavutil/frame.h>
import "C"
import "unsafe"

type (
	AVMediaType  C.enum_AVMediaType
	AVRational   C.struct_AVRational
	AVDictionary C.struct_AVDictionary
	AVFrame      C.struct_AVFrame
)

const (
	AVMEDIA_TYPE_VIDEO = C.AVMEDIA_TYPE_VIDEO
	AV_TIME_BASE       = C.AV_TIME_BASE
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
