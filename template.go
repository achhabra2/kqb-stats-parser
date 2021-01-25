package main

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func MatchImage(inputMat gocv.Mat, templateFile string) []image.Rectangle {

	templateMat := gocv.IMRead(templateFile, gocv.IMReadColor)
	defer templateMat.Close()
	templateWidth := templateMat.Cols()
	templateHeight := templateMat.Rows()
	outputMat := gocv.NewMat()
	defer outputMat.Close()
	mask := gocv.NewMat()
	defer mask.Close()
	threshold := gocv.NewMat()
	defer threshold.Close()
	gocv.MatchTemplate(inputMat, templateMat, &outputMat, gocv.TmCcoeffNormed, mask)

	// minVal, maxVal, _, _ := gocv.MinMaxLoc(outputMat)
	// fmt.Printf("MinVal: %f, MaxVal: %f\n", minVal, maxVal)
	gocv.Threshold(outputMat, &threshold, 0.50, 1.0, gocv.ThresholdBinary)

	// rows = height, cols = width
	// window := gocv.NewWindow("Output")
	// for {
	// 	window.IMShow(threshold)
	// 	window.WaitKey(1)
	// }

	gocv.Normalize(threshold, &threshold, 0, 255, gocv.NormMinMax)

	threshold.ConvertTo(&threshold, gocv.MatTypeCV8U)

	nonZero := gocv.NewMat()
	defer nonZero.Close()
	gocv.FindNonZero(threshold, &nonZero)
	// text := fmt.Sprintf("Nonzero: row: %d, cols: %d", nonZero.Rows(), nonZero.Cols())
	// fmt.Println(text)
	rects := make([]image.Rectangle, 0)
	for i := 0; i < nonZero.Rows(); i++ {
		vect := nonZero.GetVeciAt(i, 1)
		x := int(vect[0])
		y := int(vect[1])
		if x > 0 && y > 0 {
			// fmt.Println(vect)
			rects = append(rects, image.Rect(x, y, x+templateWidth, y+templateHeight))
		}
	}
	if len(rects) > 0 {
		rects = gocv.GroupRectangles(rects, 1, 0.6)
	}

	return rects
}

func DrawRects(mat *gocv.Mat, rects []image.Rectangle, text string) {
	rectColor := color.RGBA{255, 0, 0, 1}
	for _, rect := range rects {
		gocv.Rectangle(mat, rect, rectColor, 3)
		gocv.PutText(mat, text, rect.Min, gocv.FontHersheySimplex, 0.75, rectColor, 2)
	}
}

func WriteImage(mat gocv.Mat, filename string) {
	written := gocv.IMWrite(filename, mat)
	if written == true {
		fmt.Println("Successful Write")
	}
}
