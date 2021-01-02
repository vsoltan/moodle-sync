package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func RenderProgressBar(current, total int64) {

}

// ChooseUploadFolder opens a dialog to select the upload destination folder
func ChooseUploadFolder(folderList []string) (folderName string) {
	var input string
	for {
		fmt.Println("Upload file to folder: ")
		for idx, folder := range folderList {
			fmt.Printf("[%v] %v\n", idx, folder)
		}
		fmt.Scanln(&input)
		folderIdx, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			log.Printf("Invalid input: %v\n", err)
		} else if 0 > folderIdx || int(folderIdx) >= len(folderList) {
			log.Println("Invalid input: index out of bounds")
		} else {
			return folderList[folderIdx]
		}
	}
}

// DeleteLocalFileDialog opens a dialog to delete local copy of the file after upload
func DeleteLocalFileDialog(filePath string) {
	var input string
	fmt.Printf("Delete file with path %v? (y/n)\n", filePath)
	for {
		fmt.Scanln(&input)
		input = strings.ToLower(input)
		if input == "y" {
			err := os.RemoveAll(filePath)
			if err != nil {
				log.Println("Could not delete local file copy ", err)
			} else {
				log.Println("Successfully deleted file at ", filePath)
			}
			break
		} else if input == "n" {
			break
		} else {
			fmt.Println("Invalid input, select Y for yes and N for no...")
		}
	}
}

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenBrowser(url string) (err error) {
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
