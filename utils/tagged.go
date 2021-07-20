package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
)

type (
    TaggedMedia struct {
        Media []Media
    }
)


func (insta *Instagram) getTagged() {
    tagged, err := insta.SUser.Tags([]byte(insta.SUser.Username))
    if err != nil {
        fmt.Println("Error Getting Tagged Images for user " + insta.SUser.Username)
        return
    }

    taggedMedia := &TaggedMedia{}
    for tagged.Next() {
        for i := range tagged.Items {
            if len(tagged.Items[i].Images.Versions) > 0 {
                image := Media{
                    URL: tagged.Items[i].Images.GetBest(),
                    MediaType: "image",
                }
                taggedMedia.Media = append(taggedMedia.Media, image)
            } else {
                for ic := range tagged.Items[i].CarouselMedia {
                    image := Media{
                        URL: tagged.Items[i].CarouselMedia[ic].Images.GetBest(),
                        MediaType: "image",
                    }
                    taggedMedia.Media = append(taggedMedia.Media, image)
                }
            }
            if len(tagged.Items[i].Videos) > 0 {
                for iv := range tagged.Items[i].Videos {
                video := Media{
                    URL: tagged.Items[i].Videos[iv].URL,
                    MediaType: "video",
                }
                taggedMedia.Media = append(taggedMedia.Media, video)
                }
            }
        }
    }
    insta.Tagged = taggedMedia
}

func (insta *Instagram) DownloadTagged() error {
    wgM := &sync.WaitGroup{}

    insta.getTagged()

    cwd, _ := os.Getwd()
    taggedRoot := path.Join(cwd, insta.SUser.Username, "Tagged")
    err = CreateFolder(taggedRoot)
    if err != nil {
        return err
    }

    noTg := len(insta.Tagged.Media)
    log.Println("[Tagged]: Found", noTg, "Tagged Medias.")

    for i := range insta.Tagged.Media {
        wgM.Add(1)
        go func(i int, wgM *sync.WaitGroup) {
            //file := insta.Stories.Title + strconv.Itoa(i + 1)
            file := ""

            err, file = insta.Tagged.Media[i].Download(taggedRoot, file)
            if err != nil {
                log.Println("[Tagged]", err, file)
            }
            wgM.Done()
            log.Printf("[Tagged]: Progress<%d/%d>", i + 1, noTg)
        }(i, wgM)
    }
    wgM.Wait()
    return nil
}
