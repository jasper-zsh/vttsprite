package main

import (
	"fmt"
	"image/jpeg"
	"jasper-zsh/vttsprite/wrapper"
	"math"
	"os"
	"path"
	"time"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

const (
	ROWS  = 100
	COLS  = 4
	WIDTH = 300
)

func main() {
	inputFile := os.Args[1]
	dirPath := path.Dir(inputFile)
	videoFilename := path.Base(inputFile)
	spriteFilename := videoFilename + ".sprite.jpg"
	vttFilename := videoFilename + ".vtt"
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
	vttContent := "WEBVTT\n\n"

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

		vttContent += fmt.Sprintf(
			"%02d:%02d:%02d.%03d --> %02d:%02d:%02d.%03d\n%s#xywh=%d,%d,%d,%d\n\n",
			int(curTs)/3600,
			int(curTs)/60%60,
			int(curTs)%60,
			int(curTs*1000)%1000,
			int(curTs+everyNSeconds)/3600,
			int(curTs+everyNSeconds)/60%60,
			int(curTs+everyNSeconds)%60,
			int((curTs+everyNSeconds)*1000)%1000,
			spriteFilename,
			WIDTH*(col),
			targetHeight*row,
			WIDTH,
			targetHeight,
		)

		curTs += everyNSeconds
		curFrameIdx += everyNFrames
		idx += 1
	}

	f, err := os.Create(path.Join(dirPath, spriteFilename))
	if err != nil {
		fmt.Printf("Failed to create sprite file.")
		panic(err)
	}
	defer f.Close()

	jpeg.Encode(f, spriteCtx.Image(), &jpeg.Options{Quality: 80})

	vttFile, err := os.Create(path.Join(dirPath, vttFilename))
	if err != nil {
		fmt.Printf("Failed to create vtt file.")
		panic(err)
	}
	defer vttFile.Close()

	vttFile.WriteString(vttContent)

	// videoReader.Release()
}

func PerfTimer() float64 {
	return float64(time.Now().UnixMicro()) / float64(1e3)
}
