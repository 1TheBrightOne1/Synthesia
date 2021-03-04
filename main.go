package main

import (
	"fmt"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

const (
	whiteKeyCount = int64(52)
)

func main() {
	keyStr := os.Getenv("keyindex")
	key, _ := strconv.ParseInt(keyStr, 10, 64)
	fmt.Println(key)

	rowStr := os.Getenv("row")
	row, _ := strconv.ParseInt(rowStr, 10, 64)

	skipStr := os.Getenv("skipframes")
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	video, _ := gocv.VideoCaptureFile("/out.mp4")

	defer video.Close()

	img := gocv.NewMat()
	defer img.Close()

	video.Grab(int(skip))

	video.Read(&img)
	keys := scanLine(&img, int(row))

	for _, val := range keys {
		if val {
			fmt.Print(" 1")
		} else {
			fmt.Print(" 0")
		}
	}
	fmt.Println()

	gocv.IMWrite("/outputs/wat.png", img)

	/*for i := 24; i < int(whiteKeyCount); i++ {
		fmt.Printf("Key %d\n", i)
		video, _ := gocv.VideoCaptureFile("/out.mp4")
		frames := monitorKey(video, int64(i), skip, row)
		fmt.Println(frames)
		video.Close()
	}*/

	//monitorKey(video, key, skip, row)
}

func scanLine(img *gocv.Mat, row int) []bool {
	keys := make([][]gocv.Vecb, int(whiteKeyCount+1))
	keyWidth := float64(img.Cols()) / float64(whiteKeyCount)

	switchAtFloat := 0.0
	switchAt := 0

	keyIndex := -1
	for i := 0; i < img.Cols(); i++ {
		if i >= switchAt {
			keyIndex++
			switchAtFloat += keyWidth

			switchAt = int(switchAtFloat)
			keys[keyIndex] = make([]gocv.Vecb, 0)
		}

		pixel := img.GetVecbAt(int(row), i)
		fmt.Printf("%d %d %d %d\n", i, pixel[0], pixel[1], pixel[2])
		keys[keyIndex] = append(keys[keyIndex], pixel)
	}

	result := make([]bool, int(whiteKeyCount))

	for i, _ := range result {
		counter := 0
		passed := 0
		for _, pixel := range keys[i] {
			counter++
			if pixel[0] > 100 && pixel[1] < 100 && pixel[2] > 100 {
				passed++
			}
		}
		if counter-passed < 4 {
			result[i] = true
		} else {
			result[i] = false
		}
	}

	return result
}
