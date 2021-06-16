package utils

import (
	"log"
	"os"
	"path"
	//"strconv"
	"sync"
)

type (
    HighlightMedia struct {
        Title string
        Media []Media
    }
)

// Highlights
// Get Highlights
func (insta *Instagram) getHighlights() error {
    highlights, err := insta.SUser.Highlights()
    if err != nil {
        return err
    }

    // Highlights
    for i := range highlights {
        items := highlights[i].Items

        stories := &HighlightMedia{
            Title: highlights[i].Title,
        }

        // Media in Highlights
        for is := range items {
            image := Media{
                URL: items[is].Images.GetBest(),
                MediaType: "image",
            }
            for iv := range items[is].Videos {
                link := items[is].Videos[iv].URL
                if link != "" {
                    video := Media{
                        URL: link,
                        MediaType: "video",
                    }
                    stories.Media = append(stories.Media, video)
                }
            }
            stories.Media = append(stories.Media, image)
        }
        insta.Highlights = append(insta.Highlights, *stories)
    }
    return nil
}


// Download Highlights
func (insta *Instagram) DownloadHighlights() error {
    wgM := &sync.WaitGroup{}

    err := insta.getHighlights()
    if err != nil {
        return err
    }

    cwd, _ := os.Getwd()
    highlightsRoot := path.Join(cwd, *&insta.SUser.Username,"Highlights")
    err = CreateFolder(highlightsRoot)
    if err != nil {
        return err
    }

    noHgls := len(insta.Highlights)
    log.Println("[Highlights]: Found", noHgls, "Highlights.")

    for i := range insta.Highlights {
        section := insta.Highlights[i].Title
        sectionRoot := path.Join(highlightsRoot,  section)

        noMda := len(insta.Highlights[i].Media)
        log.Println("[Highlights]: Found", noMda, "Medias in", section)

        err = CreateFolder(sectionRoot)
        if err != nil {
            return err
        }

        for is := range insta.Highlights[i].Media {
            wgM.Add(1)
            go func(is, i int, wgM *sync.WaitGroup) {
                //file := section + strconv.Itoa(is + 1) + strconv.Itoa(i + 1)
                file := ""
                err, file = insta.Highlights[i].Media[is].Download(sectionRoot, file)
                if err != nil {
                    log.Println("[Highlights]", err, file)
                }
                wgM.Done()
                log.Printf("[Highlights]:(%d/%d) Progress<%d/%d>", i + 1, noHgls, is + 1, noMda)
            }(is, i, wgM)
        }
    }
    wgM.Wait()
    return nil
}

