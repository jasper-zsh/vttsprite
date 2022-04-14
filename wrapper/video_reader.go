package wrapper

//#include <string.h>
import "C"
import (
	"errors"
	"image"
	"jasper-zsh/vttsprite/avcodec"
	"jasper-zsh/vttsprite/avformat"
	"jasper-zsh/vttsprite/avutil"
	"unsafe"
)

type VideoReader struct {
	FileName    string
	fmtCtx      *avformat.AVFormatContext
	codec       *avcodec.AVCodec
	codecCtx    *avcodec.AVCodecContext
	videoStream *avformat.AVStream
	rawFrame    *avutil.AVFrame
	packet      *avcodec.AVPacket
	videoInfo   *VideoInfo
}

func (reader *VideoReader) Open() error {
	reader.fmtCtx = avformat.AvformatAllocContext()
	ret := avformat.AvformatOpenInput(&reader.fmtCtx, reader.FileName, nil, nil)
	if ret != 0 {
		return avutil.ErrorFromCode(ret)
	}
	ret = avformat.AvformatFindStreamInfo(reader.fmtCtx, nil)
	if ret != 0 {
		return avutil.ErrorFromCode(ret)
	}
	videoStreamIdx := -1
	for i := range reader.fmtCtx.Streams() {
		if reader.fmtCtx.Streams()[i].CodecParameters().CodecType() == avutil.AVMEDIA_TYPE_VIDEO {
			videoStreamIdx = i
			break
		}
	}
	if videoStreamIdx < 0 {
		return errors.New("no video stream found")
	}
	reader.videoStream = reader.fmtCtx.Streams()[videoStreamIdx]
	reader.codec = avcodec.AvcodecFindDecoder(reader.videoStream.CodecParameters().CodecId())
	if reader.codec == nil {
		return errors.New("no codec found")
	}
	reader.codecCtx = avcodec.AvcodecAllocContext3(reader.codec)
	ret = avcodec.AvcodecParametersToContext(reader.codecCtx, reader.videoStream.CodecParameters())
	if ret != 0 {
		return avutil.ErrorFromCode(ret)
	}
	ret = avcodec.AvcodecOpen2(reader.codecCtx, reader.codec, nil)
	if ret < 0 {
		return avutil.ErrorFromCode(ret)
	}
	reader.rawFrame = avutil.AvFrameAlloc()
	reader.packet = avcodec.AvPacketAlloc()
	return nil
}

func (reader *VideoReader) SeekSeconds(seconds float64) error {
	ts := avutil.AvRescaleQ(int64(seconds*avutil.AV_TIME_BASE), avutil.AV_TIME_BASE_Q, *reader.videoStream.TimeBase())
	ret := avformat.AvSeekFrame(reader.fmtCtx, reader.videoStream.Index(), ts, 0)
	if ret != 0 {
		return avutil.ErrorFromCode(ret)
	}
	avcodec.AvcodecFlushBuffers(reader.codecCtx)
	return nil
}

func (reader *VideoReader) Read() (image.Image, error) {
	for avformat.AvReadFrame(reader.fmtCtx, reader.packet) >= 0 {
		if reader.packet.StreamIndex() == reader.videoStream.Index() {
			ret := avcodec.AvcodecSendPacket(reader.codecCtx, reader.packet)
			if ret < 0 {
				return nil, avutil.ErrorFromCode(ret)
			}
			for ret >= 0 {
				ret = avcodec.AvcodecReceiveFrame(reader.codecCtx, reader.rawFrame)
				if ret >= 0 {
					w := reader.rawFrame.Width()
					h := reader.rawFrame.Height()
					ys := reader.rawFrame.Linesize()[0]
					cs := reader.rawFrame.Linesize()[1]
					img := image.YCbCr{
						Y:              avutil.FromCPtr(unsafe.Pointer(reader.rawFrame.Data()[0]), ys*h),
						Cb:             avutil.FromCPtr(unsafe.Pointer(reader.rawFrame.Data()[1]), cs*h/2),
						Cr:             avutil.FromCPtr(unsafe.Pointer(reader.rawFrame.Data()[2]), cs*h/2),
						YStride:        ys,
						CStride:        cs,
						SubsampleRatio: image.YCbCrSubsampleRatio420,
						Rect:           image.Rect(0, 0, w, h),
					}

					return &img, nil
				} else if ret != avutil.EAGAIN1 {
					return nil, avutil.ErrorFromCode(ret)
				}
			}
		}
		avcodec.AvPacketUnref(reader.packet)
	}
	return nil, nil
}

func (r *VideoReader) Release() error {
	if r.packet != nil {
		avcodec.AvPacketFree(&r.packet)
	}
	if r.rawFrame != nil {
		avutil.AvFrameFree(&r.rawFrame)
	}
	if r.fmtCtx != nil {
		avformat.AvformatCloseInput(&r.fmtCtx)
	}
	if r.codecCtx != nil {
		avcodec.AvcodecFreeContext(&r.codecCtx)
	}
	return nil
}

func (r *VideoReader) VideoInfo() *VideoInfo {
	if r.videoInfo == nil {
		fps := float64(r.videoStream.RFrameRate().Num()) / float64(r.videoStream.RFrameRate().Den())
		r.videoInfo = &VideoInfo{
			Width:      r.videoStream.CodecParameters().Width(),
			Height:     r.videoStream.CodecParameters().Height(),
			FrameCount: r.videoStream.NbFrames(),
			Fps:        fps,
			Duration:   float64(r.videoStream.NbFrames()) / fps,
		}
	}
	return r.videoInfo
}
