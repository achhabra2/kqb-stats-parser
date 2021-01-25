package main

import (
	"bytes"

	"gocv.io/x/gocv"

	"github.com/disintegration/imaging"
)

func ImagingProcess(imageBytes []byte) (mat gocv.Mat) {
	reader := bytes.NewReader(imageBytes)
	img, _ := imaging.Decode(reader)
	dstImage := imaging.AdjustGamma(img, 1.2)
	dstImage = imaging.AdjustContrast(dstImage, 30)
	dstImage = imaging.AdjustSaturation(dstImage, 30)
	dstImage = imaging.Sharpen(dstImage, 10)
	dstImageFill := imaging.Resize(dstImage, 1920*scaleFactor, 1080*scaleFactor, imaging.Lanczos)
	mat, _ = gocv.ImageToMatRGB(dstImageFill)
	// window := gocv.NewWindow("Output")
	// for {
	// 	window.IMShow(mat)
	// 	window.WaitKey(1)
	// }
	return
}
