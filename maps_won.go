package main

import (
	"log"

	"gocv.io/x/gocv"
)

var blueMapTemplates map[int]string = map[int]string{
	1: "./templates/blue-1.png",
	2: "./templates/blue-2.png",
	3: "./templates/blue-3.png",
}

var goldMapTemplates map[int]string = map[int]string{
	1: "./templates/gold-1.png",
	2: "./templates/gold-2.png",
	3: "./templates/gold-3.png",
}

func FindMapsWon(mat gocv.Mat) (int, int) {
	blueMapsWon := 0
	goldMapsWon := 0
	for num, templateFile := range blueMapTemplates {
		matchedRects := MatchImage(mat, templateFile, 0.85)
		if len(matchedRects) > 0 {
			log.Printf("Found Blue %d map match", num)
			blueMapsWon = num
		}
	}

	for num, templateFile := range goldMapTemplates {
		matchedRects := MatchImage(mat, templateFile, 0.85)
		if len(matchedRects) > 0 {
			log.Printf("Found Gold %d map match", num)
			goldMapsWon = num
		}
	}

	return blueMapsWon, goldMapsWon
}
