package avformat

//#cgo pkg-config: libavutil libavformat libavcodec
//
//#include <libavformat/avformat.h>
//#include <libavcodec/codec_par.h>
import "C"
import (
	"jasper-zsh/vttsprite/avcodec"
	"jasper-zsh/vttsprite/avutil"
	"unsafe"
)

type (
	AVFormatContext C.struct_AVFormatContext
	AVInputFormat   C.struct_AVInputFormat
	AVStream        C.struct_AVStream
)

func AvformatAllocContext() *AVFormatContext {
	return (*AVFormatContext)(C.avformat_alloc_context())
}

func AvformatOpenInput(ps **AVFormatContext, url string, fmt *AVInputFormat, options **avutil.AVDictionary) int {
	return int(C.avformat_open_input((**C.struct_AVFormatContext)(unsafe.Pointer(ps)), C.CString(url), (*C.struct_AVInputFormat)(unsafe.Pointer(fmt)), (**C.struct_AVDictionary)(unsafe.Pointer(options))))
}

func AvformatFindStreamInfo(ic *AVFormatContext, options **avutil.AVDictionary) int {
	return int(C.avformat_find_stream_info((*C.struct_AVFormatContext)(unsafe.Pointer(ic)), (**C.struct_AVDictionary)(unsafe.Pointer(options))))
}

func AvReadFrame(s *AVFormatContext, pkt *avcodec.AVPacket) int {
	return int(C.av_read_frame(
		(*C.struct_AVFormatContext)(unsafe.Pointer(s)),
		(*C.struct_AVPacket)(unsafe.Pointer(pkt)),
	))
}

func AvSeekFrame(s *AVFormatContext, streamIndex int, timestamp int64, flags int) int {
	return int(C.av_seek_frame(
		(*C.struct_AVFormatContext)(unsafe.Pointer(s)),
		C.int(streamIndex),
		C.longlong(timestamp),
		C.int(flags),
	))
}