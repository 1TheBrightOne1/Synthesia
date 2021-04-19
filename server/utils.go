package server

import (
	"github.com/1TheBrightOne1/Synthesia/keyboard"
	"gocv.io/x/gocv"
	"fmt"
)

func calibrate(v *gocv.VideoCapture, k *keyboard.Keyboard, scanRow int) int {
	requiredCalibrations := 100
	matched := 0
	sum := 0
	//Using a keyboard go until a key is detected for 2 frames
	img := gocv.NewMat()
	for v.Read(&img) {
		key, offset1 := k.CalibrateFrame(&img, scanRow, -1)
		if offset1 >= 0 {
			gocv.IMWrite("/out1.png", img)
			v.Read(&img)
			_, offset2 := k.CalibrateFrame(&img, scanRow, key)
			if offset2 >= 0 && offset2 > offset1 {
				gocv.IMWrite("/out2.png", img)
				fmt.Printf("%d %d\n", offset1, offset2)
				matched++
				sum += offset2 - offset1

				if matched == requiredCalibrations {
					break
				}
				continue
			}
		}
	}
	return int(sum / requiredCalibrations)
}
