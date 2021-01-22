package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"gocv.io/x/gocv"
)

// var subsections = map[string]image.Rectangle{
// 	"Name":    image.Rect(0, 90, 500, 180),
// 	"Kills":   image.Rect(0, 500, 250, 610),
// 	"Berries": image.Rect(250, 500, 500, 610),
// 	"Deaths":  image.Rect(0, 610, 250, 730),
// 	"Snail":   image.Rect(250, 610, 500, 730),
// 	"Queen":   image.Rect(0, 10, 500, 120),
// }

var subsections = map[string]image.Rectangle{
	"Name":    image.Rect(0, 45, 250, 90),
	"Kills":   image.Rect(0, 250, 125, 305),
	"Berries": image.Rect(125, 250, 250, 305),
	"Deaths":  image.Rect(0, 305, 125, 365),
	"Snail":   image.Rect(125, 305, 250, 365),
	"Queen":   image.Rect(0, 5, 250, 60),
}

var queenSubsection = image.Rect(0, 20, 500, 120)

var blueRects, goldRects []image.Rectangle

// detectText gets text from the Vision API for an image at the given file path.
func detectText(w io.Writer, f io.Reader, outSet *Set) error {
	var teamRects = map[string][]image.Rectangle{
		"Blue": blueRects,
		"Gold": goldRects,
	}
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	// f, err := os.Open(file)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()

	img, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectTexts(ctx, img, nil, 70)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No text found.")
	} else {
		fmt.Fprintln(w, "Text:")
		for color, team := range teamRects {
			teamStats := Team{Color: color}
			for idx, playerRect := range team {
				currentPlayer := Player{}
				if idx == 0 {
					currentPlayer.Queen = 1
				}
				for _, annotation := range annotations {
					// fmt.Println(annotation.Description)
					subStatRect := image.Rect(472, 100, 1720, 844)
					p1 := image.Point{int(annotation.BoundingPoly.Vertices[0].X), int(annotation.BoundingPoly.Vertices[0].Y)}
					p2 := image.Point{int(annotation.BoundingPoly.Vertices[2].X), int(annotation.BoundingPoly.Vertices[2].Y)}
					boundingRect := image.Rectangle{p1, p2}
					if IsInBox(subStatRect, p1) {
						// fmt.Printf("Annotation: %s, Rect: %v\n", annotation.Description, boundingRect)
						for stat, statRect := range subsections {
							offset := image.Point{playerRect.Min.X, playerRect.Min.Y}
							rectOffset := statRect.Add(offset)
							// fmt.Printf("For %s, comparing %v to %v", annotation.Description, p, rect)
							if rectOffset.Overlaps(boundingRect) {
								// log.Println("Found Stat", stat, annotation.Description)
								switch stat {
								// case "Name":
								// 	currentPlayer.Name = annotation.Description
								case "Kills":
									currentPlayer.Kills = ProcessOCRText(annotation.Description)
								case "Berries":
									currentPlayer.Berries = ProcessOCRText(annotation.Description)
								case "Deaths":
									currentPlayer.Deaths = ProcessOCRText(annotation.Description)
								case "Snail":
									currentPlayer.Snail = ProcessOCRText(annotation.Description)
								case "Queen":
									if currentPlayer.Name != "" {
										currentPlayer.Name += " " + annotation.Description
									} else {
										currentPlayer.Name = annotation.Description
									}
								default:
									break
								}
							}
						}
					}
				}
				teamStats.Players = append(teamStats.Players, currentPlayer)
			}
			outSet.Teams = append(outSet.Teams, teamStats)
		}
	}

	return nil
}

func IsInBox(rect image.Rectangle, p image.Point) bool {
	if rect.Min.X < p.X && rect.Max.X > p.X {
		if rect.Min.Y < p.Y && rect.Max.Y > p.Y {
			return true
		}
	}
	return false
}

func ProcessImage(mat *gocv.Mat) {
	loaded := *mat
	width := float64(loaded.Cols())
	height := float64(loaded.Rows())
	// fmt.Printf("Rows %d, Cols %d", height, width)
	rescaleWidth := 1920.0 / width
	rescaleHeight := 1080.0 / height
	gocv.Resize(loaded, mat, image.Point{}, rescaleWidth, rescaleHeight, gocv.InterpolationLanczos4)
	// gocv.EqualizeHist(loaded, mat)
	eroderMat := gocv.GetStructuringElement(gocv.MorphRect, image.Point{2, 2})

	gocv.Dilate(loaded, mat, eroderMat)
	// gocv.Erode(loaded, mat, eroderMat)
	_ = gocv.Threshold(loaded, mat, 160, 255, gocv.ThresholdBinaryInv)
	// gocv.AdaptiveThreshold(loaded, mat, 255, gocv.AdaptiveThresholdMean, gocv.ThresholdBinary, 5, 5)
}

func SetupPlayerRects() {
	xDiff := image.Point{325, 0}
	yDiff := image.Point{0, 375}
	minStart := image.Point{472, 100}
	maxStart := image.Point{722, 465}
	goldMinStart := minStart.Add(yDiff)
	goldMaxStart := maxStart.Add(yDiff)
	// blueRects := make([]image.Rectangle, 0)
	nums := []int{0, 1, 2, 3}
	for _, num := range nums {
		blueRects = append(blueRects, image.Rectangle{minStart.Add(xDiff.Mul(num)), maxStart.Add(xDiff.Mul(num))})
	}

	// goldRects := make([]image.Rectangle, 0)
	for _, num := range nums {
		goldRects = append(goldRects, image.Rectangle{goldMinStart.Add(xDiff.Mul(num)), goldMaxStart.Add(xDiff.Mul(num))})
	}
}

func ProcessOCRText(text string) int {
	text = strings.Replace(text, "O", "0", -1)
	text = strings.Replace(text, "o", "0", -1)
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(text, "")
	if processedString == "" {
		return 0
	}
	processedInt, err := strconv.Atoi(processedString)
	if err != nil {
		log.Fatal(err)
	}
	return processedInt
}

func RecieveHTTPImage(imageData []byte) ([]byte, error) {
	loaded, err := gocv.IMDecode(imageData, gocv.IMReadGrayScale)
	if err != nil {
		fmt.Println("Could not decode image", err)
		return nil, err
	}
	ProcessImage(&loaded)
	// written := gocv.IMWrite("sample.png", loaded)
	// if written == true {
	// 	fmt.Println("Successful Write")
	// }

	written, err := gocv.IMEncode(gocv.PNGFileExt, loaded)
	if err != nil {
		fmt.Println("Could not encode image", err)
		return nil, err
	}

	imageBuf := bytes.NewBuffer(written)
	set := Set{}
	err = detectText(os.Stdout, imageBuf, &set)
	if err != nil {
		fmt.Println("Error calling vision API", err)
		return nil, err
	}
	output, _ := json.MarshalIndent(set, "", "    ")
	return output, nil
}

// func GetOCRText(image string) string {
// 	client := gosseract.NewClient()
// 	client.SetWhitelist("01234567890")
// 	err := client.SetPageSegMode(gosseract.PSM_AUTO_OSD)
// 	if err != nil {
// 		fmt.Println("Received Error", err)
// 	}
// 	defer client.Close()
// 	err = client.SetImage(image)
// 	if err != nil {
// 		fmt.Println("Received Error", err)
// 	}
// 	text, _ := client.Text()
// 	if err != nil {
// 		fmt.Println("Received error", err)
// 	}
// 	fmt.Println(text)
// 	// Hello, World!
// 	return text
// }
