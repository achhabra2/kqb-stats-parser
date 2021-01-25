package main

import (
	"bytes"
	"context"
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
// 	"Name":    image.Rect(0, 45, 250, 90),
// 	"Kills":   image.Rect(0, 250, 125, 305),
// 	"Berries": image.Rect(125, 250, 250, 305),
// 	"Deaths":  image.Rect(0, 305, 125, 365),
// 	"Snail":   image.Rect(125, 305, 250, 365),
// 	"Queen":   image.Rect(0, 5, 250, 60),
// }

const scaleFactor = 3

var subsections = map[string]image.Rectangle{
	"Name":    image.Rect(0*scaleFactor, 45*scaleFactor, 250*scaleFactor, 90*scaleFactor),
	"Kills":   image.Rect(40*scaleFactor, 250*scaleFactor, 145*scaleFactor, 305*scaleFactor),
	"Berries": image.Rect(165*scaleFactor, 250*scaleFactor, 270*scaleFactor, 305*scaleFactor),
	"Deaths":  image.Rect(40*scaleFactor, 305*scaleFactor, 145*scaleFactor, 365*scaleFactor),
	"Snail":   image.Rect(165*scaleFactor, 305*scaleFactor, 270*scaleFactor, 365*scaleFactor),
	"Queen":   image.Rect(0*scaleFactor, 5*scaleFactor, 250*scaleFactor, 60*scaleFactor),
}

// detectText gets text from the Vision API for an image at the given file path.
func detectText(w io.Writer, f io.Reader, outSet *Set, blueRects []image.Rectangle, goldRects []image.Rectangle) error {
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
	annotations, err := client.DetectTexts(ctx, img, nil, 120)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No text found.")
	} else {
		for color, team := range teamRects {
			teamStats := Team{Color: color}
			for idx, playerRect := range team {
				currentPlayer := Player{}
				if idx == 0 {
					currentPlayer.Queen = 1
				}
				for _, annotation := range annotations {
					// fmt.Println(annotation.Description)
					subStatRect := image.Rect(450*scaleFactor, 90*scaleFactor, 1820*scaleFactor, 1000*scaleFactor)
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
										currentPlayer.Name = ProcessOCRName(annotation.Description)
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

func FindOrigin(mat *gocv.Mat) image.Point {
	loaded := *mat
	matchedRects := MatchImage(loaded, "./static/queen.png")
	if len(matchedRects) > 0 {
		log.Printf("Using Origin %v\n", matchedRects[0])
		return matchedRects[0].Min
	} else {
		log.Println("Could not find origin, using default")
		return image.Point{472, 180}
	}
}

func ResizeImage(mat *gocv.Mat) {
	loaded := *mat
	width := float64(loaded.Cols())
	height := float64(loaded.Rows())
	// fmt.Printf("Rows %d, Cols %d", height, width)
	rescaleWidth := 1920.0 / width
	rescaleHeight := 1080.0 / height
	gocv.Resize(loaded, mat, image.Point{}, rescaleWidth, rescaleHeight, gocv.InterpolationLanczos4)
}

func ProcessImage(mat *gocv.Mat) {
	loaded := *mat
	// fmt.Println(matchedRects)
	// DrawRects(mat, matchedRects, "Queen")
	// WriteImage(loaded, "test.png")
	//gocv.BilateralFilter(loaded, mat, 5, 20, 20)
	// gocv.EqualizeHist(test, &test)
	// test := gocv.NewMat()
	// defer test.Close()
	// fmt.Println(loaded.Type().String())
	// filter := gocv.NewMatWithSizeFromScalar(gocv.NewScalar(255, 100, 90, 0), 1080, 1920, gocv.MatTypeCV8UC3)
	// gocv.AddWeighted(loaded, 0.9, filter, 0.25, 0, mat)
	// dilateMat := gocv.GetStructuringElement(gocv.MorphRect, image.Point{2, 2})
	// defer dilateMat.Close()
	// gocv.Erode(loaded, mat, dilateMat)
	_ = gocv.Threshold(loaded, mat, 180, 255, gocv.ThresholdToZero)
	gocv.CvtColor(loaded, mat, gocv.ColorBGRToGray)
	_ = gocv.Threshold(loaded, mat, 110, 255, gocv.ThresholdBinaryInv)
	//gocv.Dilate(loaded, mat, dilateMat)
	// erodeMat := gocv.GetStructuringElement(gocv.MorphRect, image.Point{2, 1})
	// defer erodeMat.Close()
	// gocv.Erode(loaded, mat, erodeMat)
	// gocv.Dilate(loaded, mat, dilateMat)
	// window := gocv.NewWindow("Output")
	// for {
	// 	window.IMShow(test)
	// 	window.WaitKey(1)
	// }
}

func SetupPlayerRects(origin image.Point) (blueRects []image.Rectangle, goldRects []image.Rectangle) {
	xDiff := image.Point{325, 0}.Mul(scaleFactor)
	yDiff := image.Point{0, 375}.Mul(scaleFactor)
	var minStart, maxStart image.Point
	// if partySize < 3 {
	// 	minStart = image.Point{472, 140}
	// 	maxStart = image.Point{722, 505}
	// } else {
	// 	minStart = image.Point{472, 100}
	// 	maxStart = image.Point{722, 465}
	// }

	minStart = image.Point{origin.X, origin.Y - 44*scaleFactor}
	maxStart = image.Point{minStart.X + 250*scaleFactor, minStart.Y + 365*scaleFactor}
	goldMinStart := minStart.Add(yDiff)
	goldMaxStart := maxStart.Add(yDiff)
	blueRects = make([]image.Rectangle, 0)
	nums := []int{0, 1, 2, 3}
	for _, num := range nums {
		blueRects = append(blueRects, image.Rectangle{minStart.Add(xDiff.Mul(num)), maxStart.Add(xDiff.Mul(num))})
	}

	goldRects = make([]image.Rectangle, 0)
	for _, num := range nums {
		goldRects = append(goldRects, image.Rectangle{goldMinStart.Add(xDiff.Mul(num)), goldMaxStart.Add(xDiff.Mul(num))})
	}
	return
}

func ProcessOCRText(text string) int {
	// fmt.Println("Processing Number", text)
	text = strings.Replace(text, "O", "0", -1)
	text = strings.Replace(text, "o", "0", -1)
	text = strings.Replace(text, "l", "1", -1)
	text = strings.Replace(text, "I", "1", -1)

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

func ProcessOCRName(text string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9(){}]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(text, "")
	return processedString
}

func RecieveHTTPImage(imageData []byte) (Set, error) {
	loaded, err := gocv.IMDecode(imageData, gocv.IMReadColor)
	if err != nil {
		fmt.Println("Could not decode image", err)
		return Set{}, err
	}
	defer loaded.Close()
	ResizeImage(&loaded)
	origin := FindOrigin(&loaded)
	origin = origin.Mul(scaleFactor)
	blueRects, goldRects := SetupPlayerRects(origin)
	processedMat := ImagingProcess(imageData)
	defer processedMat.Close()
	ProcessImage(&processedMat)
	written, err := gocv.IMEncode(gocv.PNGFileExt, processedMat)
	if err != nil {
		fmt.Println("Could not encode image", err)
		return Set{}, err
	}

	// boundryMat := gocv.NewMat()
	// defer boundryMat.Close()
	// gocv.CvtColor(processedMat, &boundryMat, gocv.ColorGrayToBGR)
	// DrawRects(&boundryMat, blueRects, "blue")
	// DrawRects(&boundryMat, goldRects, "gold")
	// write := gocv.IMWrite("./internal/step2.png", boundryMat)
	// if write == true {
	// 	fmt.Println("Successful Write")
	// }

	imageBuf := bytes.NewBuffer(written)
	set := Set{}
	err = detectText(os.Stdout, imageBuf, &set, blueRects, goldRects)
	if err != nil {
		fmt.Println("Error calling vision API", err)
		return Set{}, err
	}
	return set, nil
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
