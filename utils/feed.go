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

    feedMedia := &FeedMedia{}
    for feed.Next() {
        for i := range feed.Items {
            if len(feed.Items[i].Images.Versions) > 0 {
                image := Media{
                    URL: feed.Items[i].Images.GetBest(),
                    MediaType: "image",
                }
                feedMedia.Media = append(feedMedia.Media, image)
            } else {
                for ic := range feed.Items[i].CarouselMedia {
                    image := Media{
                        URL: feed.Items[i].CarouselMedia[ic].Images.GetBest(),
                        MediaType: "image",
                    }
                    feedMedia.Media = append(feedMedia.Media, image)
                }
            }
            

            if len(feed.Items[i].Videos) > 0 {
                for iv := range feed.Items[i].Videos {
                video := Media{
                    URL: feed.Items[i].Videos[iv].URL,
                    MediaType: "video",
                }
                feedMedia.Media = append(feedMedia.Media, video)
                }
            }
        }
    }
    insta.Feed = feedMedia
}

func (insta *Instagram) DownloadFeed() error {
    wgM := &sync.WaitGroup{}

    insta.getFeed()

    cwd, _ := os.Getwd()
    feedRoot := path.Join(cwd, *&insta.SUser.Username, "Feed")
    err = CreateFolder(feedRoot)
    if err != nil {
        return err
    }

    noFd := len(insta.Feed.Media)
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
