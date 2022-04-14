package main

import (
	"fmt"
	"image/jpeg"
	"jasper-zsh/vttsprite/wrapper"
	"math"
	"os"
	"time"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

const (
	ROWS  = 5
	COLS  = 3
	WIDTH = 300
)

func main() {
	inputFile := "sample.mp4"
	videoReader := wrapper.VideoReader{
		FileName: inputFile,
	}
	err := videoReader.Open()
	if err != nil {
		fmt.Printf("Failed to open video. %s", err.Error())
		os.Exit(1)
	}

	targetHeight := int(math.Round(float64(WIDTH) / float64(videoReader.VideoInfo().Width) * float64(videoReader.VideoInfo().Height)))
	everyNSeconds := videoReader.VideoInfo().Duration / ROWS / COLS
	everyNFrames := float64(videoReader.VideoInfo().FrameCount) / ROWS / COLS

	spriteCtx := gg.NewContext(WIDTH*COLS, targetHeight*ROWS)

	curTs := 0.0
	curFrameIdx := 0.0
	idx := 0
	execTime := time.Now().Unix()
	for idx < ROWS*COLS {
		T1 := PerfTimer()
		videoReader.SeekSeconds(curTs)
		T2 := PerfTimer()
		img, err := videoReader.Read()
		T3 := PerfTimer()
		if err != nil {
			fmt.Printf("Failed to extract frame. %s", err.Error())
			break
		}
		row := idx / COLS
		col := idx % COLS
		x := col * WIDTH
		y := row * targetHeight
		scaled := resize.Resize(WIDTH, uint(targetHeight), img, resize.Bilinear)
		T4 := PerfTimer()
		spriteCtx.DrawImage(scaled, x, y)
		T5 := PerfTimer()
		now := time.Now().Unix()
		if now-execTime >= 1 {
			fmt.Printf("Timestamp: %.3fs Perf(ms) Seek: %.3f Read: %.3f Resize: %.3f Draw: %.3f\n", curTs, T2-T1, T3-T2, T4-T3, T5-T4)
			execTime = now
		}

		curTs += everyNSeconds
		curFrameIdx += everyNFrames
		idx += 1
	}

	f, _ := os.Create("out.jpg")
	defer f.Close()

	jpeg.Encode(f, spriteCtx.Image(), &jpeg.Options{Quality: 80})

	// videoReader.Release()
}

func PerfTimer() float64 {
	return float64(time.Now().UnixMicro()) / float64(1e3)
}
