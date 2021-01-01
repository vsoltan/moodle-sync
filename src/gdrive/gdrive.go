package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/vsoltan/moodle-sync/src/cli"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var mimeTypes = map[string]string{
	"txt": "text/plain",
	"pdf": "application/pdf",
	"png": "image/png",
	"csv": "text/csv",
	"doc": "application/msword",
}

// ValidateCredentials reads the client secret file and parses it to generate a config struct
func ValidateCredentials(filepath string) (config *oauth2.Config) {
	log.Println("Authenticating Drive API...")
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err = google.ConfigFromJSON(b, drive.DriveScope, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	log.Println("Success!")
	return
}

// GetService generates an access client to Google Drive API and creates a new service instance
func GetService(config *oauth2.Config) *drive.Service {
	log.Println("Generating client...")
	client := getClient(config)
	log.Println("Success!")

	log.Println("Creating Drive service...")
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	log.Println("Success!")
	return srv
}

// GetOrCreateFolder retrieves a folder's ID if it exists or creates a new one in Google Drive
func GetOrCreateFolder(srv *drive.Service, foldername string) (folderID string, err error) {
	folderID, found := findFolder(srv, foldername)
	if !found {
		log.Printf("Folder with name %v not found!", foldername)
		folderID, err = createFolder()
		if err != nil {
			log.Println("Could not create a new folder: ", err)
		}
	}
	return
}

// Upload uploads a file to Google Drive
func Upload(srv *drive.Service, file *os.File, localPath, folderID string) {
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		// do something to handle directories
	} else {
		fmt.Println("uploading file...")
		// show progress bar

		mimeType := getContentType(localPath)
		ok := createFile(srv, file, fileInfo.Name(), mimeType, folderID)
		if ok {
			cli.DeleteLocalFileDialog(localPath)
		}
	}
}

// dynamically determines the content type by looking at the file's extension
func getContentType(path string) (contentType string) {
	splitPath := strings.Split(path, ".")
	if len(splitPath) == 0 {
		return "application/octet-stream"
	}
	ext := splitPath[len(splitPath)-1]
	fmt.Println("ext", ext)
	contentType, ok := mimeTypes[ext]
	if !ok {
		return "application/octet-stream"
	}
	fmt.Println("contentType", contentType)
	return
}

func createFile(service *drive.Service, file *os.File,
	filename string, mimeType string, parentID string) (ok bool) {
	fileMetadata := &drive.File{
		Name:     filename,
		MimeType: mimeType,
	}
	if parentID != "" {
		fileMetadata.Parents = []string{parentID}
	}
	_, err := service.Files.Create(fileMetadata).Media(file).Context(context.Background()).Do()

	if err != nil {
		log.Println("Could not create file in Google Drive: ", err)
	} else {
		log.Println("Success!")
		ok = true
	}
	return
}

func findFolder(srv *drive.Service, foldername string) (folderID string, found bool) {
	q := fmt.Sprintf("name=\"%s\" and mimeType=\"application/vnd.google-apps.folder\"", foldername)
	fileList, err := srv.Files.List().Q(q).Do()
	if err != nil {
		log.Println("Unable to find folder with name: ", foldername)
	} else {
		found = true
		folderID = fileList.Files[0].Id
	}
	return
}

func createFolder() (folderID string, err error) {
	// TODO
	return "", err
}

/*
 * API authentication: https://developers.google.com/drive/api/v3/quickstart/go
 */

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Opening authentication dialogue in browser")

	err := cli.OpenBrowser(authURL)
	if err != nil {
		fmt.Printf("Go to the following link in your browser then type the "+
			"authorization code: \n%v\n", authURL)
	}

	fmt.Println("Paste generated code to continue: ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
