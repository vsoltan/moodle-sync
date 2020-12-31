package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Upload uploads a file to Google Drive
func Upload(srv *drive.Service, file *os.File, folderName string) {
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		// do something to handle directories
	} else {
		fmt.Println("uploading file...")
		folderID := getOrCreateFolder(srv, folderName)
		fmt.Println(folderID)
		createFile(srv, file, fileInfo.Name(), "text/plain", folderID)
		// if fileInfo.Size() <= smallFileLimit {

		// } else {

		// }
	}
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

func getOrCreateFolder(srv *drive.Service, foldername string) (folderID string) {
	folderID, found := findFolder(srv, foldername)
	if !found {
		folderID = createFolder()
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
		fmt.Println("folderID: ", folderID)
	}
	return
}

func createFolder() string {
	return ""
}

func createFile(service *drive.Service, file *os.File,
	filename string, mimeType string, parentID string) {
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
	}
}

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

	err := openbrowser(authURL)
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

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openbrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return
}
