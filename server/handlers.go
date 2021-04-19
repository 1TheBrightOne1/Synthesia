package server

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/1TheBrightOne1/Synthesia/keyboard"
	"github.com/1TheBrightOne1/Synthesia/video"
	"gocv.io/x/gocv"
)

func index(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("server/static/index.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(body))
}

func download(w http.ResponseWriter, r *http.Request) {
	youtubeLink, ok := r.URL.Query()["URL"]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	dir := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	os.Mkdir("/outputs/sessions/"+dir, 0666)

	err = video.DownloadYoutube(youtubeLink[0], "/outputs/sessions/"+string(dir)+"/video.mp4")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = video.GenerateThumbnail("/outputs/sessions/"+string(dir)+"/video.mp4", 2400, "/outputs/sessions/"+string(dir)+"/frame.png")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/setup/?video=%s", dir), http.StatusTemporaryRedirect)
}

func setup(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./server/templates/setup.html")
	video := r.URL.Query()["video"]

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(b), fmt.Sprintf("/sessions/%s/frame.png", video[0]), 640, 360, video[0]) //TODO: get width and height
}

func generate(w http.ResponseWriter, r *http.Request) {
	video, ok := r.URL.Query()["video"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	startingKey, ok := r.URL.Query()["startingKey"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	keyboardRowStr, ok := r.URL.Query()["keyboardRow"]
	testRowStr, ok := r.URL.Query()["testRow"]
	testRowEndStr, ok := r.URL.Query()["testRowEnd"]

	bThresholdStr, ok :=r.URL.Query()["bthreshold"]
	bThreshold, _ := strconv.ParseInt(bThresholdStr[0], 10, 32)
	gThresholdStr, ok :=r.URL.Query()["gthreshold"]
	gThreshold, _ := strconv.ParseInt(gThresholdStr[0], 10, 32)
	rThresholdStr, ok :=r.URL.Query()["rthreshold"]
	rThreshold, _ := strconv.ParseInt(rThresholdStr[0], 10, 32)

	keyBordersStr, ok := r.URL.Query()["keyBorders"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v, err := gocv.VideoCaptureFile("/outputs/sessions/" + video[0] + "/video.mp4")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	img := gocv.NewMat()

	v.Grab(2400)

	v.Read(&img)

	keyList := strings.Split(keyBordersStr[0], ",")
	keyListInt := make([]int, len(keyList))

	for i, val := range keyList {
		key, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			fmt.Printf("Unable to parse %s as int\n", val)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		keyListInt[i] = int(key)
	}

	keyboardRow, _ := strconv.ParseInt(keyboardRowStr[0], 10, 32)
	testRowStart, _ := strconv.ParseInt(testRowStr[0], 10, 32) //TODO: change these to area boundaries to be more clear
	testRowEnd, _ := strconv.ParseInt(testRowEndStr[0], 10, 32)
	testAreaLength := int(testRowStart - testRowEnd)
	
	k := keyboard.NewKeyboard(&img, startingKey[0], keyListInt, []uint8{uint8(bThreshold), uint8(gThreshold), uint8(rThreshold)}, int(keyboardRow), int(testRowStart), int(testRowEnd))

	fmt.Println("Calibrating")
	offset := calibrate(v, k, int(testRowStart))
	testRow := int(testRowEnd) + int(testAreaLength / offset) * offset
	//length = 150 - 100 (50)
	//150 - (50 / 3)floor

	v.Close()

	v, err = gocv.VideoCaptureFile("/outputs/sessions/" + video[0] + "/video.mp4")
	defer v.Close()

	count := 0
	for v.Read(&img) {
		count++
		fmt.Println(count)
		k.CheckFrame(&img, testRow)

		v.Grab(int(testAreaLength / offset))
	}

	fmt.Println("Writing frames")
	k.WriteFrames(w, false)
}
