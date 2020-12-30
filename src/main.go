package main

import (
	// "encoding/json"

	"fmt"
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/vsoltan/moodle-sync/src/config"
	"github.com/vsoltan/moodle-sync/src/gdrive"
	"google.golang.org/api/drive/v3"
)

const smallFileLimit = 5 << (10 * 2)

// uploads a file to google drive
func uploadToDrive(srv *drive.Service, content []byte) {
	fmt.Printf("file content: %s, len: %v\n", content, len(content))
	if len(content) <= smallFileLimit {
		fmt.Println("Uploading to gdrive")
	} else {
		fmt.Println("Large file, not implemented yet")
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

	log.Printf("Listening to changes in directory: %v\n", appConfig.Local.SyncFolder)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Fatal("Something went wrong reading event channel")
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("File added to directory!")
					content, err := ioutil.ReadFile(event.Name)
					if err != nil {
						log.Fatalf("Could not read file path supplied by callback")
					}
					uploadToDrive(srv, content)
				}
			}
		}
	}()
	err = watcher.Add(appConfig.Local.SyncFolder)
	<-done
}
