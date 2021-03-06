package utils

import (
	"log"
	"os"
)


func CreateFolder(folderName string) error {
        if _, err := os.Stat(folderName); os.IsNotExist(err) {
            err := os.Mkdir(folderName, 0755)
            if err != nil {
                log.Println("Failed to create directory", folderName)
                return err
            } else {
                log.Println("Created folder", folderName)
            }
        }
        return nil
    }
