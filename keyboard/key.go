package keyboard

import (
	"gocv.io/x/gocv"
)

type Key struct {
	start  int
	finish int
	Pixels []bool
	Pitch  string
	done chan bool
	LeftWhiteKeyIndex int

	bgrThresholds []int
}

func newKey(start, finish int, pitch string, bgrThresholds []int, leftWhiteKeyIndex int) *Key {
	return &Key{
		start:  start,
		finish: finish,
		Pitch:  pitch,
		Pixels: make([]bool, 0),
		done: make(chan bool),
		bgrThresholds: bgrThresholds,
		LeftWhiteKeyIndex: leftWhiteKeyIndex,
	}
}

func (k *Key) readFrame(img *gocv.Mat, row, testRowEnd int) {
	//Read from row to top of window threshold and look for keys being played
	for j := row; j > testRowEnd; j-- {
		if k.checkRow(img, j) {
			k.Pixels = append(k.Pixels, true)
		} else {
			k.Pixels = append(k.Pixels, false)
		}
	}
	k.done <- true
}

func (k *Key) checkRow(img *gocv.Mat, row int) bool {
	passed := 0
	for i := k.start; i < k.finish; i++ {
		pixel := img.GetVecbAt(row, i)

		p := true
		for j, threshold := range k.bgrThresholds {
			if threshold < 0 {
				modified := threshold * -1
				if pixel[j] >= uint8(modified) {
					p = false
				}
			} else if pixel[j] <= uint8(threshold) {
				p = false
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