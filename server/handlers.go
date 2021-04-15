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

	os.Mkdir("./sessions/"+dir, 0666)

	video.DownloadYoutube(youtubeLink[0], "./sessions/"+string(dir)+"/video.mp4")

	video.GenerateThumbnail("./sessions/"+string(dir)+"/video.mp4", 2400, "./sessions/"+string(dir)+"/frame.png")

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

	keyBordersStr, ok := r.URL.Query()["keyBorders"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v, err := gocv.VideoCaptureFile("./sessions/" + video[0] + "/video.mp4")

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
	testRow, _ := strconv.ParseInt(testRowStr[0], 10, 32)
	k := keyboard.NewKeyboard(&img, startingKey[0], keyListInt, int(keyboardRow), int(testRow))

	fmt.Println("Calibrating")
	offset := calibrate(v, k, int(testRow))
	w.Write([]byte(strconv.Itoa(offset)))

	v.Close()
	v, err = gocv.VideoCaptureFile("./sessions/" + video[0] + "/video.mp4")
	defer v.Close()

	count := 0
	for v.Read(&img) {
		count++
		fmt.Println(count)
		k.CheckFrame(&img)
	}

	k.WriteFrames(w, false)
}
