package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/vsoltan/moodle-sync/src/config"
	"github.com/vsoltan/moodle-sync/src/gdrive"
	"google.golang.org/api/drive/v3"
)

func uploadToDrive() {

}

func main() {

	log.Println("Starting moodle-sync service...")
	log.Println("Reading config file...")

	appConfig := config.Init()

	log.Println("Success!")
	log.Println("Authenticating Drive API...")

	driveConfig := gdrive.ReadConfig(appConfig.Local.CredPath)
	log.Println("Success!")

	client := gdrive.GetClient(driveConfig)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	// example interfacing with api
	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}

	done := make(chan bool)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Fatal("Something went wrong reading event channel")
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("File added to directory!")
					// upload to gdrive
				}
			}
		}
	}()
	err = watcher.Add(appConfig.Local.SyncFolder)
	<-done
}
