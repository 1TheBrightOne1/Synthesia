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
	video.Grab(30 * 80)

	img := gocv.NewMat()
	if !video.Read(&img) {
		fmt.Println("WHOOPS")
		return
	}

	defer video.Close()

	for x := 0; x < img.Cols(); x++ {
		for y := 0; y < img.Rows(); y++ {
			val0 := img.GetUCharAt3(x, y, 0)
			val1 := img.GetUCharAt3(x, y, 1)
			val2 := img.GetUCharAt3(x, y, 2)

			if val0 != val1 || val0 != val2 {
				fmt.Printf("x: %d, y: %d, b: %d, g: %d, r: %d\n", x, y, val0, val1, val2)
			}
		}
	}

	return

	/*for key := 0; key < img.Rows(); key++ {
		fmt.Printf("%d\t%d\n", key, t[key])
	}
	for {
		x, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		y, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		xI, err := strconv.ParseInt(strings.TrimSpace(x), 10, 32)
		if err != nil {
			fmt.Println(err.Error())
		}
		yI, err := strconv.ParseInt(strings.TrimSpace(y), 10, 32)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(xI + yI)

		fmt.Println(img.GetUCharAt3(int(xI), int(yI), 0))
		fmt.Println(img.GetUCharAt3(int(xI), int(yI), 1))
		fmt.Println(img.GetUCharAt3(int(xI), int(yI), 2))
	}
	return
	for i := 0; i < int(whiteKeyCount); i++ {
		video, _ := gocv.VideoCaptureFile("/out.mp4")
		result := monitorKey(video, int64(i), skip, row)
		video.Close()

		fmt.Printf("Key: %d V: %d\n", i, result)
	}*/
	monitorKey(video, key, skip, row)
}

func monitorKey(video *gocv.VideoCapture, keyIndex, skipFrames, row int64) {
	video.Grab(int(skipFrames))

	img := gocv.NewMat()

	var counter int64
	for video.Read(&img) {
		whiteKeyWidth := int64(img.Cols() / int(whiteKeyCount))
		ch := img.Channels()

		var sum, count int
		avg := make([]int, ch)
		for i := int64(0); i < whiteKeyWidth; i++ {
			for z := 0; z < ch; z++ {
				val := img.GetUCharAt3(int(whiteKeyWidth*keyIndex+i), int(row), z)
				sum += int(val)
				avg[z] += int(val)
			}
			count++ //This is stupid, use keywidth
		}

		fmt.Printf("frame: %d v: %v 0: %d 1: %d 2: %d\n", counter+skipFrames, int(float64(sum)/float64(count*ch)), avg[0]/(count*ch), avg[1]/(count*ch), avg[2]/(count*ch))
		counter++
	}
}
