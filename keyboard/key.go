package keyboard

import (
	"fmt"

	"gocv.io/x/gocv"
)

type Key struct {
	start             int
	finish            int
	Pixels            []bool
	Pitch             string
	done              chan bool
	LeftWhiteKeyIndex int

	bgrThresholds []int
}

func newKey(start, finish int, pitch string, bgrThresholds []int, leftWhiteKeyIndex int) *Key {
	return &Key{
		start:             start,
		finish:            finish,
		Pitch:             pitch,
		Pixels:            make([]bool, 0),
		done:              make(chan bool),
		bgrThresholds:     bgrThresholds,
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
	if (k.finish-k.start)-passed < 1 {
		if k.LeftWhiteKeyIndex > 0 {
			if k.validateBlackKey(img, row) {
				return true
			}
			return false
		}
		return true
	} else {
		return false
	}
}

func (k *Key) validateBlackKey(img *gocv.Mat, row int) bool {
	//Check if a white key is playing to the left
	fakeKey := newKey(k.start-(k.finish-k.start), k.start, "", k.bgrThresholds, -1)
	if fakeKey.checkRow(img, row) {
		//If it is get the pixel average for the black key vs the partial white key. They should be different.
		summer := func(img *gocv.Mat, row, start, finish int) []int {
			sums := make([]int, 3)
			for i := start; i < finish; i++ {
				pixel := img.GetVecbAt(row, i)
				fmt.Println(pixel)
				sums[0] += int(pixel[0])
				sums[1] += int(pixel[1])
				sums[2] += int(pixel[2])
			}
			sums[0] /= finish - start
			sums[1] /= finish - start
			sums[2] /= finish - start
			return sums
		}
		fmt.Println("fake")
		fakeSums := summer(img, row, fakeKey.start, fakeKey.finish)
		fmt.Printf("AVG %v\n", fakeSums)
		fmt.Println("blackkey")
		blackKeySums := summer(img, row, k.start, k.finish)
		fmt.Printf("AVG %v\n", blackKeySums)

		for i, _ := range fakeSums {
			diff := fakeSums[i] - blackKeySums[i]
			fmt.Println(diff)
			if diff < 0 {
				diff *= -1
			}
			if diff > 4 {
				return true
			}
		}

		return false
	}
	return true
}
