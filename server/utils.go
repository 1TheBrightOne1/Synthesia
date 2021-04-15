package server

import (
	"github.com/1TheBrightOne1/Synthesia/keyboard"
	"gocv.io/x/gocv"
)

func calibrate(v *gocv.VideoCapture, k *keyboard.Keyboard, scanRow int) int {
	requiredCalibrations := 10
	matched := 0
	sum := 0
	//Using a keyboard go until a key is detected for 2 frames
	img := gocv.NewMat()
	for v.Read(&img) {
		key, offset1 := k.CalibrateFrame(&img, -1, scanRow)
		if offset1 >= 0 {
			v.Read(&img)
			_, offset2 := k.CalibrateFrame(&img, key, scanRow)
			if offset2 >= 0 {
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
