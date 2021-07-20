package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/roshanrane24/instago/utils"
)

var (
    insta  *utils.Instagram

    // Flags declaration for CLI
    username = flag.String("u", "", "Login Username")
    password = flag.String("p", "", "Login Password")
    all = flag.Bool("all", false, "Download All Media")
    highlights = flag.Bool("highlights", false, "Download User Highlights")
    stories = flag.Bool("stories", false, "Download User Stories")
    feed = flag.Bool("feed", false, "Download User Feed Post")
    tagged = flag.Bool("tagged", false, "Download User Tagged Post")

    target *string
)

func main() {
    flag.Parse()
    if len(flag.Args()) > 0 {
        target = &flag.Args()[0]
    } else {
        log.Println("Please provide target username") 
        return
    }

    if *all {
        stories, feed, highlights = all, all, all
    }

    // Waitgroup
    wg := &sync.WaitGroup{}

    // Check For username. Usersname Needed
    if *username == "" {
        fmt.Println("Please provide login username")
        return
    }

    insta := &utils.Instagram{
        User: username,
        Pass: password,
    }

    err := insta.InstaLogin()
    if err != nil {
        fmt.Println(err)
        return
    }
    defer insta.InstaLogout()
    err = insta.SearchUser(*target)
    if err != nil {
        log.Fatalln("User", *target, "Not Found.\n", err)
        return
    }

    // Profile Picture
    err = insta.GetProfilePic()
    if err != nil {
        log.Println("Failed to Download Profile Picture")
    }

    // Highlights
    if *highlights {
        wg.Add(1)
        go func(wg *sync.WaitGroup) {
            log.Println(">>> Highlights")
            err := insta.DownloadHighlights()
            if err != nil {
                log.Println(err)
                log.Println("Error While Downloading Highlights.")
            }
            log.Println("Highlights <<<")
            wg.Done()
        }(wg)
    }

    // Stories
    if *stories {
        wg.Add(1)
        go func(wg *sync.WaitGroup) {
            log.Println(">>> Stories")
            err := insta.DownloadStories()
            if err != nil {
                log.Println(err)
                log.Println("Error While Downloading Stories.")
            }
            log.Println("Stories <<<")
            wg.Done()
        }(wg)
    }

    // Feed
    if *feed {
        wg.Add(1)
        go func(wg *sync.WaitGroup) {
            log.Println(">>> Feed")
            err := insta.DownloadFeed()
            if err != nil {
                log.Println(err)
                log.Println("Error While Downloading Feed.")
            }
            log.Println("Feed <<<")
            wg.Done()
        }(wg)
    }

    if *tagged {
        wg.Add(1)
        go func(wg *sync.WaitGroup) {
            log.Println(">>> Tagged")
            err := insta.DownloadTagged()
            if err != nil {
                log.Println(err)
                log.Println("Error While Downloading Tagged Media.")
            }
            log.Println("Tagged <<<")
            wg.Done()
        }(wg)
    }

    wg.Wait()
}
