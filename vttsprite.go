package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"jasper-zsh/vttsprite/wrapper"
	"os"
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
	var img image.Image
	img, err = videoReader.Read()
	if err != nil {
		fmt.Printf("Failed to read frame. %s", err.Error())
		os.Exit(1)
	}
	f, _ := os.Create("out.jpg")
	defer f.Close()

	jpeg.Encode(f, img, &jpeg.Options{Quality: 80})

	videoReader.Release()
}
