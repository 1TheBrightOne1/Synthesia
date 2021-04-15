package keyboard

import (
	"gocv.io/x/gocv"
)

type key struct {
	start  int
	finish int
	pixels []bool
	pitch  string
}

func newKey(start, finish int, pitch string) *key {
	return &key{
		start:  start,
		finish: finish,
		pitch:  pitch,
		pixels: make([]bool, 0),
	}
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
		k.pixels = append(k.pixels, true)
	} else {
		k.pixels = append(k.pixels, false)
	}
}
