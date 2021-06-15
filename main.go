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
    target = flag.String("t", "", "Username")
    highlights = flag.Bool("highlights", false, "Download User Highlights")
    stories = flag.Bool("stories", false, "Download User Stories")
    caption = flag.Bool("caption", false, "Download User Post Captions")
    audio = flag.Bool("audio", false, "Download Audio for User's Video Post")
    bio = flag.Bool("bio", false, "Download User's Bio")
    followers = flag.Bool("followers", false, "Download User's Follower List")
    following = flag.Bool("following", false, "Download User's Following List")
    tag = flag.Bool("tag", false, "Download Post from User's Tagged List")

)

func main() {
    flag.Parse()

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
    insta.SearchUser(*target)

    if *highlights {
        wg.Add(1)
        go func(wg *sync.WaitGroup) {
            log.Println(">>> Highlights")
            err := insta.DownloadHighlights()
            if err != nil {
                log.Println(err)
            }
            log.Println("Highlights <<<")
            wg.Done()
        }(wg)
    }

    //insta.GetProfilePic()

    wg.Wait()
}
