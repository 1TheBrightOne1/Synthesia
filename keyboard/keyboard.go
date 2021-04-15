package keyboard

import (
	"fmt"
	"io"

	"gocv.io/x/gocv"
)

var octave = map[int]string{
	0: "A",
	1: "B",
	2: "C",
	3: "D",
	4: "E",
	5: "F",
	6: "G",
}

type Keyboard struct {
	whiteKeys []*key
	blackKeys []*key

	testRow        int //The pixel row height where we check for whether a key is pressed
	keyboardHeight int //The pixel row height used to identify the location of the black keys

	startingKey string
}

func NewKeyboard(frame *gocv.Mat, startingKey string, whiteKeyBorders []int, keyboardHeight, testRow int) *Keyboard {
	k := &Keyboard{
		testRow:        testRow,
		keyboardHeight: keyboardHeight,
		startingKey:    startingKey,
	}
	k.initWhiteKeys(whiteKeyBorders)
	k.initBlackKeys(frame)
	return k
}

func (k *Keyboard) WriteFrames(w io.Writer, includeBlackKeys bool) {
	for i := 0; i < len(k.whiteKeys[0].pixels); i++ {
		for j := 0; j < len(k.whiteKeys); j++ {
			if k.whiteKeys[j].pixels[i] {
				fmt.Fprint(w, "1 ")
			} else {
				fmt.Fprint(w, "0 ")
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

func (k *Keyboard) CheckFrame(frame *gocv.Mat) {
	for _, key := range k.whiteKeys {
		go key.readFrame(frame, k.testRow)
	}
	for _, key := range k.blackKeys {
		go key.readFrame(frame, k.testRow)
	}

	for _, key := range k.whiteKeys {
		<-key.done
	}
	for _, key := range k.blackKeys {
		<-key.done
	}
}

func (k *Keyboard) CalibrateFrame(frame *gocv.Mat, row, keyIndex int) (key, offset int) {
	if keyIndex >= 0 {
		if k.whiteKeys[keyIndex].checkRow(frame, row) {
			for j := row + 1; j >= row-40; /*TODO: 40 could be too far out of bounds. Also DRY*/ j-- {
				if !k.whiteKeys[keyIndex].checkRow(frame, j) {
					return keyIndex, j
				}
			}
		}
		return -1, -1
	}
	for i, key := range k.whiteKeys {
		if key.checkRow(frame, row) {
			for j := row + 1; j >= row-40; /*TODO: 40 could be too far out of bounds*/ j-- {
				if !key.checkRow(frame, i) {
					return i, j
				}
			}
		}
	}
	return -1, -1
}

func (k *Keyboard) initWhiteKeys(borders []int) {
	k.whiteKeys = make([]*key, len(borders)-1)

	pitch := k.startingKey
	pitchIndex := 0
	for key, val := range octave {
		if val == pitch {
			pitchIndex = key
			break
		}
	}
	for i := range k.whiteKeys {
		k.whiteKeys[i] = newKey(borders[i], borders[i+1], 0, octave[pitchIndex%len(octave)])
		pitchIndex++
	}
}

func (k *Keyboard) initBlackKeys(frame *gocv.Mat) {
	start := -1
	finish := -1
	for x := 0; x < frame.Cols()-1; x++ {
		if getPixelAverage(frame, x, k.keyboardHeight) < 60 && getPixelAverage(frame, x-1, k.keyboardHeight) > 60 { //TODO: Make configurable
			start = x
		}

		if getPixelAverage(frame, x, k.keyboardHeight) < 60 && getPixelAverage(frame, x+1, k.keyboardHeight) > 60 && start > 0 {
			finish = x

			if finish-start > 4 { //TODO: Make configurable
				for _, whiteKey := range k.whiteKeys {
					if start > whiteKey.start && start < whiteKey.finish {
						k.blackKeys = append(k.blackKeys, newKey(start, finish, 0, whiteKey.pitch+"#"))
						break
					}
				}
			}

			start = -1
			finish = -1
		}
	}

	for _, key := range k.blackKeys {
		fmt.Printf("%d - %d\n", key.start, key.finish)
	}
}

func getPixelAverage(frame *gocv.Mat, x, y int) uint32 {
	pixel := frame.GetVecbAt(y, x)

	avg := (uint32(pixel[0]) + uint32(pixel[1]) + uint32(pixel[2])) / uint32((3)) //TODO: Make this based on channel count

	return avg
}
