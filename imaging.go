package main

import (
	"bytes"

	"gocv.io/x/gocv"

	"github.com/disintegration/imaging"
)

func ImagingProcess(imageBytes []byte) (mat gocv.Mat) {
	reader := bytes.NewReader(imageBytes)
	img, _ := imaging.Decode(reader)
	dstImage := imaging.AdjustGamma(img, 1.3)
	dstImage = imaging.AdjustSaturation(dstImage, 1.5)
	dstImage = imaging.Sharpen(dstImage, 20)
	dstImageFill := imaging.Resize(dstImage, 1920, 1080, imaging.Lanczos)
	mat, _ = gocv.ImageToMatRGB(dstImageFill)
	// window := gocv.NewWindow("Output")
	// for {
	// 	window.IMShow(mat)
	// 	window.WaitKey(1)
	// }
	return
}
