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
	WhiteKeys []*Key
	BlackKeys []*Key

	testRow        int //TODO: remove, no longer used
	testRowEnd     int
	keyboardHeight int //The pixel row height used to identify the location of the black keys

	startingKey string
}

func NewKeyboard(frame *gocv.Mat, startingKey string, whiteKeyBorders, bgrThresholds []int, keyboardHeight, testRow, testRowEnd int) *Keyboard {
	k := &Keyboard{
		testRow:        testRow,
		testRowEnd:     testRowEnd,
		keyboardHeight: keyboardHeight,
		startingKey:    startingKey,
	}
	k.initWhiteKeys(whiteKeyBorders, bgrThresholds)
	k.initBlackKeys(frame, bgrThresholds)
	return k
}

func (k *Keyboard) WriteFrames(w io.Writer, mode int) {
	if mode >= 0 {
		for i := 0; i < len(k.WhiteKeys[0].Pixels); i++ {
			for j := 0; j < len(k.WhiteKeys); j++ {
				if k.WhiteKeys[j].Pixels[i] {
					fmt.Fprint(w, "1 ")
				} else {
					fmt.Fprint(w, "0 ")
				}
			}
			fmt.Fprintf(w, "\n")
		}
	}
	if mode == 0 {
		fmt.Fprintf(w, "--------------------------------------")
	}
	if mode <= 0 {
		for i := 0; i < len(k.BlackKeys[0].Pixels); i++ {
			for j := 0; j < len(k.BlackKeys); j++ {
				if k.BlackKeys[j].Pixels[i] {
					fmt.Fprint(w, "1 ")
				} else {
					fmt.Fprint(w, "0 ")
				}
			}
			fmt.Fprintf(w, "\n")
		}
	}
}

func (k *Keyboard) CheckFrame(frame *gocv.Mat, testRow int) {
	for _, key := range k.WhiteKeys {
		go key.readFrame(frame, testRow, k.testRowEnd)
	}
	for _, key := range k.BlackKeys {
		go key.readFrame(frame, testRow, k.testRowEnd)
	}

	for _, key := range k.WhiteKeys {
		<-key.done
	}
	for _, key := range k.BlackKeys {
		<-key.done
	}
}

func (k *Keyboard) CalibrateFrame(frame *gocv.Mat, row, keyIndex int) (key, offset int) {
	if keyIndex >= 0 {
		if k.WhiteKeys[keyIndex].checkRow(frame, row) {
			for j := row - 1; j >= row-40; /*TODO: 40 could be too far out of bounds. Also DRY*/ j-- {
				if !k.WhiteKeys[keyIndex].checkRow(frame, j) {
					return keyIndex, j
				}
			}
		}
		return -1, -1
	}
	for i, key := range k.WhiteKeys {
		if key.checkRow(frame, row) {
			for j := row - 1; j >= row-40; /*TODO: 40 could be too far out of bounds*/ j-- {
				if !key.checkRow(frame, j) {
					return i, j
				}
			}
		}
	}
	return -1, -1
}

func (k *Keyboard) initWhiteKeys(borders, bgrThresholds []int) {
	k.WhiteKeys = make([]*Key, len(borders)-1)

	pitch := k.startingKey
	pitchIndex := 0
	for key, val := range octave {
		if val == pitch {
			pitchIndex = key
			break
		}
	}
	for i := range k.WhiteKeys {
		k.WhiteKeys[i] = newKey(borders[i], borders[i+1], octave[pitchIndex%len(octave)], bgrThresholds, -1)
		pitchIndex++
	}
}

func (k *Keyboard) initBlackKeys(frame *gocv.Mat, bgrThresholds []int) {
	start := -1
	finish := -1
	for x := 0; x < frame.Cols()-1; x++ {
		if getPixelAverage(frame, x, k.keyboardHeight) < 60 && getPixelAverage(frame, x-1, k.keyboardHeight) > 60 { //TODO: Make configurable
			start = x
		}

		if getPixelAverage(frame, x, k.keyboardHeight) < 60 && getPixelAverage(frame, x+1, k.keyboardHeight) > 60 && start > 0 {
			finish = x

			if finish-start > 4 { //TODO: Make configurable
				for i, whiteKey := range k.WhiteKeys {
					if start > whiteKey.start && start < whiteKey.finish {
						k.BlackKeys = append(k.BlackKeys, newKey(start, finish+1, whiteKey.Pitch+"#", bgrThresholds, i))
						break
					}
				}
			}

			start = -1
			finish = -1
		}
	}

	for _, key := range k.BlackKeys {
		fmt.Printf("%d - %d\n", key.start, key.finish)
	}
}

func getPixelAverage(frame *gocv.Mat, x, y int) uint32 {
	pixel := frame.GetVecbAt(y, x)

	avg := (uint32(pixel[0]) + uint32(pixel[1]) + uint32(pixel[2])) / uint32((3)) //TODO: Make this based on channel count

	return avg
}
