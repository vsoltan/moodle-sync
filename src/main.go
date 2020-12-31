package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/vsoltan/moodle-sync/src/config"
	"github.com/vsoltan/moodle-sync/src/gdrive"
)

const smallFileLimit = 5 << (10 * 2)

// TODO: make folderIdx an int
func chooseUploadFolder(folderList []string) (folderIdx int64) {
	var input string
	for {
		fmt.Println("Upload material to folder: ")
		for idx, course := range folderList {
			fmt.Printf("[%v] %v\n", idx, course)
		}
		fmt.Scanln(&input)
		folderIdx, err := strconv.ParseInt(input, 10, 2)
		if err != nil {
			log.Printf("Invalid input: %v\n", err)
		} else if 0 > folderIdx || int(folderIdx) >= len(folderList) {
			log.Println("Invalid input: index out of bounds")
		} else {
			return folderIdx
		}
	}
}

func main() {

	log.Println("Starting moodle-sync service...")

	appConfig := config.Parse()

	driveConfig := gdrive.ValidateCredentials(appConfig.Local.CredPath)

	srv := gdrive.GetService(driveConfig)

	done := make(chan bool)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer watcher.Close()

	syncDir := appConfig.Local.SyncFolder
	courseList := appConfig.Drive.Courses

	log.Printf("Listening to changes in directory: %v\n", syncDir)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Fatal("Something went wrong reading event channel")
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("File %v has been added to directory!\n", event.Name)
					file, err := os.Open(event.Name)
					if err != nil {
						log.Fatalf("Could not read file path supplied by callback")
					}
					folderID := chooseUploadFolder(courseList)
					log.Println("Selected folder: ", courseList[folderID])
					gdrive.Upload(srv, file)
				}
			}
		}
	}()
	err = watcher.Add(syncDir)
	<-done
}
