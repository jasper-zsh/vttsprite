package avcodec

//#cgo pkg-config: libavcodec
//
//#include <libavcodec/avcodec.h>
import "C"

type (
	AVCodecID C.enum_AVCodecID
)
