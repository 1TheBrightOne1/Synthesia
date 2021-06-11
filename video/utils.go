package video

import (
	"errors"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	youtube "github.com/kkdai/youtube/v2"
	"gocv.io/x/gocv"
)

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

	return nil
}
