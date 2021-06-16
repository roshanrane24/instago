package utils

import (
	"log"
	"os"
	"path"
	"sync"
)

type (
    FeedMedia struct {
        Media []Media
    }
)


func (insta *Instagram) getFeed() {
    feed := insta.SUser.Feed()
    feed.Next()

    for i := range feed.Items {
        image := &Media{
            URL: feed.Items[i].Images.GetBest(),
            MediaType: "image",
        }
        insta.Feed.Media = append(insta.Feed.Media, *image)

        for is := range feed.Items[i].Videos {
        video := &Media{
            URL: feed.Items[i].Videos[is].URL,
            MediaType: "video",
        }
        insta.Feed.Media = append(insta.Feed.Media, *video)
        }
    }
}

func (insta *Instagram) DownloadFeed() error {
    wgM := &sync.WaitGroup{}

    insta.getFeed()


    cwd, _ := os.Getwd()
    feedRoot := path.Join(cwd, *&insta.SUser.Username, "Stories")
    err = CreateFolder(feedRoot)
    if err != nil {
        return err
    }

    noFd := len(insta.Stories.Media)
    log.Println("[Feed]: Found", noFd, "Feed Medias.")

    for i := range insta.Feed.Media {
        wgM.Add(1)
        go func(i int, wgM *sync.WaitGroup) {
            //file := insta.Stories.Title + strconv.Itoa(i + 1)
            file := ""

            err, file = insta.Feed.Media[i].Download(feedRoot, file)
            if err != nil {
                log.Println("[Feed]", err, file)
            }
            wgM.Done()
            log.Printf("[Feed]: Progress<%d/%d>", i + 1, noFd)
        }(i, wgM)
    }
    wgM.Wait()
    return nil

}
