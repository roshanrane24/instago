package utils

import (
	"fmt"
	"os"
	"path"
	"github.com/ahmdrz/goinsta/v2"
	"github.com/tcnksm/go-input"
)


var (
    err error
)

type(
    Instagram struct{
        User *string
        Pass *string
        inst *goinsta.Instagram
        SUser *goinsta.User
        Highlights []HighlightMedia
    }

    Media struct {
        URL string
        MediaType string
    }

)


func (insta *Instagram) InstaLogin() error {
    // Check if there is saved  session for username
    cache_dir, err := os.UserCacheDir()
    if err != nil {
        fmt.Println("[ERROR]: Cache", err, cache_dir)
        return err
    }
    session_path := path.Join(cache_dir, "insta_" + *insta.User)
    if _, err = os.Stat(session_path); err == nil {
        fmt.Println("Session Found for", *insta.User)
        insta.inst, err = goinsta.Import(session_path)
        if err != nil {
            fmt.Println("[ERROR]: Importing Session ", err)
            return err
        }
    } else {
        // Ask for password if password isn't provided
        if *insta.Pass == "" {
            fmt.Printf("Enter password for %s: ", *insta.User)
            _, err := fmt.Scanf("%s", insta.Pass)
            if err != nil {
                fmt.Println("[ERROR]: Password Input", err)
                return err
            }
        }

        insta.inst = goinsta.New(*insta.User, *insta.Pass)
        if err := insta.inst.Login(); err != nil {
            insta.inst, err = loginError(err, insta.inst)
            fmt.Printf("logged in as %s \n", insta.inst.Account.Username)
        }

        // Export New Instagram session
        if err := insta.inst.Export(session_path); err != nil {
            fmt.Println("[ERROR]: Exporting Session", err)
            return err
        }
    }
    return nil
}

func (insta *Instagram) InstaLogout() {
    insta.inst.Logout()
}

func loginError(err error, insta *goinsta.Instagram) (*goinsta.Instagram, error) {
        switch v := err.(type) {
        case goinsta.ChallengeError:
            err := insta.Challenge.Process(v.Challenge.APIPath)
            if err != nil {
                fmt.Println("[ERROR]: ChallengeError", err)
                return nil, err
            }

            ui := &input.UI{
                Writer: os.Stdout,
                Reader: os.Stdin,
            }

            query := "What is SMS code for instagram?"
            code, err := ui.Ask(query, &input.Options{
                Default:  "000000",
                Required: true,
                Loop:     true,
            })
            if err != nil {
                fmt.Println("[ERROR]: 2FA Input", err)
                return nil, err
            }

            err = insta.Challenge.SendSecurityCode(code)
            if err != nil {
                fmt.Println("[ERROR]: SendSecurityCode", err)
                return nil, err
            }

            insta.Account = insta.Challenge.LoggedInUser
        default:
            fmt.Println("[ERROR]: ChallengeErrorFatal", err)
            return nil, err
        }
        return insta, nil
}


func (insta *Instagram) SearchUser(username string) error {
    insta.SUser, err = insta.inst.Profiles.ByName(username)
    wd, err := os.Getwd()
    if err != nil {
        return err
    }
    CreateFolder(path.Join(wd, insta.SUser.Username))
    return nil
}



// Profile Pic
func (insta *Instagram) GetProfilePic() {
    fmt.Println(insta.SUser.HdProfilePicURLInfo.URL)
}
