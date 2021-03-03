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

	/*m := make(map[string]bool)

	for x := 0; x < img.Cols(); x++ {
		for y := 0; y < img.Rows(); y++ {
			val := img.GetVecbAt(x, y)

			if val[0] != val[1] || val[0] != val[2] {
				m[fmt.Sprintf("b: %d, g: %d, r: %d\n", val[0], val[1], val[2])] = true
			}
		}
	}

	for key, _ := range m {
		fmt.Println(key)
	}

	return*/

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
	return*/
	for i := 24; i < int(whiteKeyCount); i++ {
		fmt.Printf("Key %d\n", i)
		video, _ := gocv.VideoCaptureFile("/out.mp4")
		frames := monitorKey(video, int64(i), skip, row)
		fmt.Println(frames)
		video.Close()
	}

	//monitorKey(video, key, skip, row)
}

func monitorKey(video *gocv.VideoCapture, keyIndex, skipFrames, row int64) int {
	var total int

	video.Grab(int(skipFrames))

	img := gocv.NewMat()
	defer img.Close()

	once := false
	for video.Read(&img) {
		data := img.ToBytes()
		cols := img.Cols()
		ch := img.Channels()

		whiteKeyWidth := float64(cols) / float64(whiteKeyCount)
		if !once {
			once = true
			fmt.Printf("%d - %d\n", int(whiteKeyWidth*float64(keyIndex)), int(whiteKeyWidth*float64(keyIndex+1)))
		}

		show := false
		for x := int(whiteKeyWidth * float64(keyIndex)); x < int(whiteKeyWidth*float64(keyIndex+1)); x++ {
			pixel := readPixel(data, cols, ch, x, int(row))

			if pixel[0] > 100 && pixel[2] > 100 {
				show = true
			}
		}

		if show {
			total++
		}
	}

	return total
}

func readPixel(data []byte, cols, ch, x, y int) []uint8 {
	idx := x*cols*ch + y*ch

	result := make([]uint8, ch)
	for i := 0; i < ch; i++ {
		result[i] = data[idx+i]
	}

	return result
}
