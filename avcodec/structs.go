package avcodec

import "jasper-zsh/vttsprite/avutil"

func (par *AVCodecParameters) CodecType() avutil.AVMediaType {
	return avutil.AVMediaType(par.codec_type)
}

func (par *AVCodecParameters) Height() int {
	return int(par.height)
}

func (par *AVCodecParameters) Width() int {
	return int(par.width)
}

func (par *AVCodecParameters) CodecId() AVCodecID {
	return AVCodecID(par.codec_id)
}

func (pkt *AVPacket) StreamIndex() int {
	return int(pkt.stream_index)
}

func (ctx *AVCodecContext) PixFmt() avutil.AVPixelFormat {
	return avutil.AVPixelFormat(ctx.pix_fmt)
}
