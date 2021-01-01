package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/vsoltan/moodle-sync/src/cli"
	"github.com/vsoltan/moodle-sync/src/config"
	"github.com/vsoltan/moodle-sync/src/gdrive"
)

func main() {

	log.Println("Starting moodle-sync service...")

	appConfig := config.Parse()

	driveConfig := gdrive.ValidateCredentials(appConfig.Local.CredPath)

	srv := gdrive.GetService(driveConfig)

	done := make(chan bool)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	syncDir := appConfig.Local.SyncFolder
	folderList := appConfig.Drive.Folders

	log.Printf("Listening to changes in directory: %v\n", syncDir)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Fatal("Something went wrong reading event channel")
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("File %v has been added to directory!\n", event.Name)
					file, err := os.Open(event.Name)
					defer file.Close()
					if err != nil {
						log.Fatalf("Could not read file path supplied by callback")
					}
					folderName := cli.ChooseUploadFolder(folderList)
					log.Println("Selected folder: ", folderName)
					folderID, err := gdrive.GetOrCreateFolder(srv, folderName)
					if err == nil {
						log.Println("Folder Id: ", folderID)
						gdrive.Upload(srv, file, event.Name, folderID)
					}
				}
			}
		}
	}()
	err = watcher.Add(syncDir)
	<-done
}
