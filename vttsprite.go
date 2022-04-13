package main

import (
	"fmt"
	"jasper-zsh/vttsprite/avcodec"
	"jasper-zsh/vttsprite/avformat"
	"jasper-zsh/vttsprite/avutil"
	"os"
)

func main() {
	inputFile := "sample.mp4"
	ctx := avformat.AvformatAllocContext()
	if avformat.AvformatOpenInput(&ctx, inputFile, nil, nil) != 0 {
		fmt.Printf("Failed to open input video file %s\n", inputFile)
		os.Exit(1)
	}
	if avformat.AvformatFindStreamInfo(ctx, nil) != 0 {
		fmt.Printf("Failed to read stream info.\n")
		os.Exit(1)
	}
	var stream *avformat.AVStream
	foundVideoStream := false
	for i, _ := range ctx.Streams() {
		stream = ctx.Streams()[i]
		if stream.CodecParameters().CodecType() == avutil.AVMEDIA_TYPE_VIDEO {
			foundVideoStream = true
			break
		}
	}
	if !foundVideoStream {
		fmt.Printf("Failed to find video stream.")
		os.Exit(1)
	}
	totalFrames := stream.NbFrames()
	fps := float64(stream.RFrameRate().Num()) / float64(stream.RFrameRate().Den())
	fmt.Printf("Width: %d\nHeight: %d\nTotal frames: %d\nFrame rate: %.2f\nDuration: %.2f Timebase: %d/%d\n",
		stream.CodecParameters().Width(),
		stream.CodecParameters().Height(),
		totalFrames,
		fps,
		float64(totalFrames)/fps,
		stream.TimeBase().Num(),
		stream.TimeBase().Den(),
	)

	codec := avcodec.AvcodecFindDecoder(stream.CodecParameters().CodecId())
	if codec == nil {
		fmt.Printf("Failed to find codec.\n")
	}

	codecCtx := avcodec.AvcodecAllocContext3(codec)
	if avcodec.AvcodecParametersToContext(codecCtx, stream.CodecParameters()) != 0 {
		fmt.Printf("Failed to create codec context.\n")
		os.Exit(1)
	}

	if avcodec.AvcodecOpen2(codecCtx, codec, nil) < 0 {
		fmt.Printf("Failed to open codec.\n")
		os.Exit(1)
	}

	ts := avutil.AvRescaleQ(1*avutil.AV_TIME_BASE, avutil.AV_TIME_BASE_Q, *stream.TimeBase())
	if avformat.AvSeekFrame(ctx, stream.Index(), ts, 1) != 0 {
		fmt.Printf("Failed to seek.\n")
		os.Exit(1)
	}

	frame := avutil.AvFrameAlloc()
	packet := avcodec.AvPacketAlloc()
	for avformat.AvReadFrame(ctx, packet) >= 0 {
		if packet.StreamIndex() == stream.Index() {
			res := avcodec.AvcodecSendPacket(codecCtx, packet)
			if res < 0 {
				fmt.Printf("Failed to send packet to decoder. %s\n", avutil.ErrorFromCode(res))
				continue
			}
			for res >= 0 {
				res = avcodec.AvcodecReceiveFrame(codecCtx, frame)
				if res < 0 {
					fmt.Printf("Failed to receive frame from decoder. %s\n", avutil.ErrorFromCode(res))
					break
				}
			}
		}

		avcodec.AvPacketUnref(packet)
	}
}
