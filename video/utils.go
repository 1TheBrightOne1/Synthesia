package video

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	youtube "github.com/kkdai/youtube/v2"
	"gocv.io/x/gocv"
)

type Metadata struct {
	Height int
	Width  int
}

func DownloadYoutube(url, saveTo string) error {
	client := youtube.Client{}

	log.Info("getting video for ", url)
	video, err := client.GetVideo(url)
	if err != nil {
		log.Error(err)
		return err
	}

	resp, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Close()

	log.Info("saving to ", saveTo)
	file, err := os.Create(saveTo)
	if err != nil {
		log.Error(err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("done saving file for ", url)
	return nil
}

func GenerateThumbnail(videoFilePath string, skipFrames int, frameFilePath string) error {
	log.Info("generating thumbnail")
	video, err := gocv.VideoCaptureFile(videoFilePath)
	if err != nil {
		log.Error(err)
		return err
	}

	defer video.Close()

	img := gocv.NewMat()
	defer img.Close()

	video.Grab(skipFrames)

	ok := video.Read(&img)
	if !ok {
		err := errors.New("not ok, calling video.Read")
		log.Error(err)
		return err
	}

	ok = gocv.IMWrite(frameFilePath, img)
	if !ok {
		return errors.New("unable to save frame")
	}

	meta := Metadata{
		Height: img.Rows(),
		Width:  img.Cols(),
	}

	b, err := json.Marshal(meta)
	if err != nil {
		log.Error(err)
		return err
	}

	dir := filepath.Dir(videoFilePath)
	log.Info(fmt.Sprintf("saving metadata to %s: %v", dir, meta))
	f, err := os.Create(dir + "/metadata.txt")
	if err != nil {
		log.Error(err)
		return err
	}
	defer f.Close()
	f.Write(b)

	return nil
}
