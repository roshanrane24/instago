package utils

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Downloader interface {
    Download(folder, file string) error
}


// Media Downloader
func (m *Media) Download(folder, file string) error {
    switch m.MediaType {
    case "image":
        err := downloadMedia(m.URL, path.Join(folder, file + ".jpg"))
        if err != nil {
            return err
        }
    case "video":
        err := downloadMedia(m.URL, path.Join(folder, file + ".mp4"))
        if err != nil {
            return err
        }
    default:
        err := errors.New("InvalidMediaType")
        log.Println("[WARNING]:", err)
        return err
    }
    return errors.New("DownloadError")
}

// Image Downloader
func downloadMedia(URL, path string) error{

    response, err := http.Get(URL)
    if err != nil {
        return err
    }

    defer response.Body.Close()

    if response.StatusCode != 200{
        return errors.New(strconv.Itoa(response.StatusCode) + "HTTPError")
    }

    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.Copy(file, response.Body)
    if err != nil {
        return err
    }

    return nil
}

// Video Downloader
//func downloadVideo() {}
