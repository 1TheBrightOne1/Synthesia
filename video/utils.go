package video

import (
	"errors"
	"io"
	"os"

	youtube "github.com/kkdai/youtube/v2"
	"gocv.io/x/gocv"
)

func DownloadYoutube(url, saveTo string) error {
	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		return err
	}

	resp, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(saveTo)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func GenerateThumbnail(videoFilePath string, skipFrames int, frameFilePath string) error {
	video, err := gocv.VideoCaptureFile(videoFilePath)
	if err != nil {
		return err
	}

	defer video.Close()

	img := gocv.NewMat()
	defer img.Close()

	video.Grab(skipFrames)

	video.Read(&img)

	ok := gocv.IMWrite(frameFilePath, img)
	if !ok {
		return errors.New("Unable to save frame")
	}

	return nil
}
