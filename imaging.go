package main

import (
	"bytes"

	"gocv.io/x/gocv"

	"github.com/disintegration/imaging"
)

func ImagingProcess(imageBytes []byte) {
	reader := bytes.NewReader(imageBytes)
	img, _ := imaging.Decode(reader)
	dstImage := imaging.AdjustSaturation(img, 15)
	dstImageFill := imaging.Fill(dstImage, 1920, 1080, imaging.Center, imaging.Lanczos)
	mat, _ := gocv.ImageToMatRGB(dstImageFill)
	window := gocv.NewWindow("Output")
	for {
		window.IMShow(mat)
		window.WaitKey(1)
	}
}
