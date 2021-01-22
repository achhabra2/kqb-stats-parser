package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

func main() {
	SetupPlayerRects()
	SetupWebServer()
	// ProcessLocalFile()
}

func ProcessLocalFile() {
	loaded := gocv.IMRead("kqb.png", gocv.IMReadGrayScale)
	ProcessImage(&loaded)

	written := gocv.IMWrite("sample.png", loaded)
	if written == true {
		fmt.Println("Successful Write")
	}

	// ioutil.TempFile("", "stats.*.png")

	// blueRectColor := color.RGBA{255, 255, 255, 1}
	// goldRectColor := color.RGBA{255, 255, 255, 1}
	// for _, rect := range blueRects {
	// 	// gocv.Rectangle(&loaded, rect, blueRectColor, 5)
	// 	for name, statRect := range subsections {
	// 		if name != "Name" && name != "Queen" {
	// 			offset := image.Point{rect.Min.X, rect.Min.Y}
	// 			draw := statRect.Add(offset)
	// 			gocv.Rectangle(&loaded, draw, blueRectColor, 3)
	// 			fmt.Printf("Blue Rect: %s, Loc: %v\n", name, draw)
	// 		}
	// 	}
	// }
	// for _, rect := range goldRects {
	// 	// gocv.Rectangle(&loaded, rect, goldRectColor, 5)
	// 	for name, statRect := range subsections {
	// 		if name != "Name" && name != "Queen" {
	// 			offset := image.Point{rect.Min.X, rect.Min.Y}
	// 			draw := statRect.Add(offset)
	// 			gocv.Rectangle(&loaded, draw, goldRectColor, 3)
	// 			fmt.Printf("Gold Rect: %s, Loc: %v\n", name, draw)
	// 		}
	// 	}
	// }

	// set := Set{}
	// err := detectText(os.Stdout, "sample.png", &set)
	// if err != nil {
	// 	fmt.Println("Error calling vision API", err)
	// }
	// output, _ := json.Marshal(set)
	// _ = ioutil.WriteFile("stats.json", output, 0644)
}
