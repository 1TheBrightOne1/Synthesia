package main

import (
	"os"

	"github.com/1TheBrightOne1/Synthesia/musicxml"
)

/*func main() {
	server.Start()
}*/

func main() {
	b := musicxml.NewBuilder(48, 8, 4, 4)

	b.BuildXML(os.Stdout)
}

/*func main() {
	keys := []int{0, 12, 23, 34, 46, 58, 69, 81, 93, 105, 116, 129, 141, 153, 165, 178, 190, 201, 215, 227, 240, 252, 265, 278, 290, 303, 315, 328, 341, 353, 366, 379, 392, 404, 417, 430, 442, 455, 467, 480, 493, 505, 517, 530, 542, 555, 567, 579, 591, 603, 616, 628, 640}
	frame := gocv.IMRead("./server/static/frames/wat.png", gocv.IMReadUnchanged)
	keyboard.NewKeyboard(&frame, "A", keys, 197, 50)
}*/

/*import (
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

func writeKeysToXML(keys []key) {
	beatsPerMeasure := 4
	framesPerMeasure := 35

	heldFor := 0
	for frame, val := range keys[25].frames {
		if val && !keys[25].frames[frame-1] {
			//Key has been pressed
			heldFor++
		} else if !val && keys[25].frames[frame-1] {
			//Key has been released
			note := fmt.Sprintf("<note><pitch><step>E</step><octave>4</octave></pitch><duration>%d</duration><type>%s</type></note>",
				heldFor/framesPerMeasure*beatsPerMeasure,
				"Eighth")

			heldFor = 0
		} else if val {
			heldFor++
		}
	}
}
*/
