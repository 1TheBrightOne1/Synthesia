package keyboard

import (
	"gocv.io/x/gocv"
)

type key struct {
	start  int
	finish int
	pixels []bool
	pitch  string
	done   chan bool

	bgrThresholds []int
}

func newKey(start, finish int, pitch string, bgrThresholds []int) *key {
	return &key{
		start:         start,
		finish:        finish,
		pitch:         pitch,
		pixels:        make([]bool, 0),
		done:          make(chan bool),
		bgrThresholds: bgrThresholds,
	}
}

func (k *key) readFrame(img *gocv.Mat, row, testRowEnd int) {
	//Read from row to top of window threshold and look for keys being played
	for j := row; j > testRowEnd; j-- {
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

		p := true
		for j, threshold := range k.bgrThresholds {
			if threshold < 0 {
				modified := threshold * -1
				if !(pixel[j] < uint8(modified)) {

					p = false
					break
				}
			} else {
				if !(pixel[j] > uint8(threshold)) {
					p = false
					break
				}
			}
		}
		if p {
			passed++
		}
	}
	if (k.finish-k.start)-passed < 3 {
		return true
	} else {
		return false
	}
}
