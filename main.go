package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gocv.io/x/gocv"
)

type key struct {
	start  int
	finish int
	frames []bool
}

func (k *key) readFrame(img *gocv.Mat, row int) {
	passed := 0
	for i := k.start; i < k.finish; i++ {
		pixel := img.GetVecbAt(row, i)

		if pixel[0] > 100 && pixel[2] > 100 && pixel[1] < 100 {
			passed++
		}
	}
	if (k.finish-k.start)-passed < 3 {
		k.frames = append(k.frames, true)
	} else {
		k.frames = append(k.frames, false)
	}
}

func main() {
	keys := loadKeys()

	rowStr := os.Getenv("row")
	row, _ := strconv.ParseInt(rowStr, 10, 64)

	skipStr := os.Getenv("skipframes")
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	video, _ := gocv.VideoCaptureFile("/out.mp4")

	defer video.Close()

	img := gocv.NewMat()
	defer img.Close()

	video.Grab(int(skip))

	for video.Read(&img) {
		for _, key := range keys {
			key.readFrame(&img, int(row))
		}
	}

	frames := len(keys[0].frames)
	for i := 0; i < frames; i++ {
		for _, val := range keys {
			if val.frames[i] {
				fmt.Print(" 1")
			} else {
				fmt.Print(" 0")
			}
		}
		fmt.Println()
	}
}

func loadKeys() []*key {
	keys := make([]*key, 0)

	f, _ := os.Open("keys.txt")

	r := bufio.NewReader(f)

	index := 0
	for {
		val, err := r.ReadString('\n')
		if err != nil {
			break
		}

		keyIndex, _ := strconv.ParseInt(strings.TrimSpace(val), 10, 32)

		newKey := &key{
			start: int(keyIndex),
		}

		keys = append(keys, newKey)
		if index > 0 {
			keys[index-1].finish = int(keyIndex)
		}
		index++
	}

	return keys[0 : len(keys)-1]
}
