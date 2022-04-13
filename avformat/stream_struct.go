package avformat

import (
	"jasper-zsh/vttsprite/avcodec"
	"jasper-zsh/vttsprite/avutil"
	"unsafe"
)

func (stream *AVStream) CodecParameters() *avcodec.AVCodecParameters {
	return (*avcodec.AVCodecParameters)(unsafe.Pointer(stream.codecpar))
}

func (stream *AVStream) NbFrames() int64 {
	return int64(stream.nb_frames)
}

func (stream *AVStream) RFrameRate() *avutil.AVRational {
	return (*avutil.AVRational)(unsafe.Pointer(&stream.r_frame_rate))
}

func (stream *AVStream) TimeBase() *avutil.AVRational {
	return (*avutil.AVRational)(unsafe.Pointer(&stream.time_base))
}

func (stream *AVStream) Index() int {
	return int(stream.index)
}
