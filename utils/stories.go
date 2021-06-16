package utils

import (
	"log"
	"os"
	"path"
	"strconv"
	"sync"
)

//"log"

type (
    StoryMedia struct {
        Title string
        Media []Media
    }
)

func (insta *Instagram) getStories() error {
    stories := insta.SUser.Stories()
    stories.Next()
    storiesMedia := &StoryMedia{
        Title: stories.Title,
    }

    for i := range stories.Items {
        image := Media{
            URL: stories.Items[i].Images.GetBest(),
            MediaType: "image",
        }
        storiesMedia.Media = append(storiesMedia.Media, image)

        for is := range stories.Items[i].Videos {
            video := Media{
                URL: stories.Items[i].Videos[is].URL,
                MediaType: "video",
            }
            storiesMedia.Media = append(storiesMedia.Media, video)
        }
    }
    insta.Stories = storiesMedia
    return nil
}


// Download Stories
func (insta *Instagram) DownloadStories() error {
    wgM := &sync.WaitGroup{}

    err := insta.getStories()
    if err != nil {
        return err
    }

    cwd, _ := os.Getwd()
    storiesRoot := path.Join(cwd, *&insta.SUser.Username, "Stories")
    err = CreateFolder(storiesRoot)
    if err != nil {
        return err
    }

    noStr := len(insta.Stories.Media)
    log.Println("[Stories]: Found", noStr, "Stories.")

    for i := range insta.Stories.Media {
        wgM.Add(1)
        go func(i int, wgM *sync.WaitGroup) {
            //file := insta.Stories.Title + strconv.Itoa(i + 1)
            file := ""

            err, file = insta.Stories.Media[i].Download(storiesRoot, file)
            if err != nil {
                log.Println("[Stories]", err, file)
            }
            wgM.Done()
            log.Printf("[Stories]: Progress<%d/%d>", i + 1, noStr)
        }(i, wgM)
    }
    wgM.Wait()
    return nil
}
