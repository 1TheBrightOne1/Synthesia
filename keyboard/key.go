package keyboard

import (
	"gocv.io/x/gocv"
)

type key struct {
	start  int
	finish int
	pixels []bool
	pitch  string
	windowThreshold int
	done chan bool
}

func newKey(start, finish, windowThreshold int, pitch string) *key {
	return &key{
		start:  start,
		finish: finish,
		pitch:  pitch,
		pixels: make([]bool, 0),
		done: make(chan bool),
		windowThreshold: windowThreshold,
	}
}

func (k *key) readFrame(img *gocv.Mat, row int) {
	//Read from row to top of window threshold and look for keys being played
	for j := row; j <= k.windowThreshold; j-- {
		if k.checkRow(img, j) {
			k.pixels = append(k.pixels, true)
		} else {
			k.pixels = append(k.pixels, false)
		}
	}
	k.done <- true
}

func (k *key) checkRow(img *gocv.Mat, row int) bool {
	passed := 0
	for i := k.start; i < k.finish; i++ {
		pixel := img.GetVecbAt(row, i)

		if pixel[0] > 100 && pixel[2] > 100 && pixel[1] < 100 {
			passed++
		}
	}
	if (k.finish-k.start)-passed < 3 {
		return true
	} else {
		return false
	}
}
