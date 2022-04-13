package avcodec

//#cgo pkg-config: libavcodec
//
//#include <libavcodec/codec_id.h>
import "C"

type (
	AVCodecID C.enum_AVCodecID
)
