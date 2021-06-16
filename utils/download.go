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
func (m *Media) Download(folder, file string) (error, string) {
    f := ""
    switch m.MediaType {
    case "image":
        err, f := downloadMedia(m.URL, folder, file, ".jpg")
        if err != nil {
            return err, f
        }
    case "video":
        err, f := downloadMedia(m.URL, folder, file, ".mp4")
        if err != nil {
            return err, f
        }
    default:
        err := errors.New("InvalidMediaType")
        log.Println("[WARNING]:", err)
        return err, ""
    }
    return nil, f
}

// Image Downloader
func downloadMedia(URL, folder, file, ext string) (error, string) {
    response, err := http.Get(URL)
    if err != nil {
        return err, URL
    }

    if file  == "" {
        file = path.Base(response.Request.URL.Path)
    } else {
        file = file + ext
    }

    path := path.Join(folder, file)

    defer response.Body.Close()

    if response.StatusCode != 200{
        return errors.New(strconv.Itoa(response.StatusCode) + "HTTPError"), file
    }

    openFile, err := os.Create(path)
    if err != nil {
        return err, file
    }
    defer openFile.Close()

    _, err = io.Copy(openFile, response.Body)
    if err != nil {
        return err, file
    }

    return nil, file
}
