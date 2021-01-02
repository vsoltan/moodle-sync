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

const progressBarLen = 20

// RenderProgressBar implementation of ProgressUpdater callback for ResumableMedia
func RenderProgressBar(current, total int64) {
	ratio := float64(current) / float64(total)
	width := int(ratio * progressBarLen)
	percentComplete := int(ratio * 100)
	progressBar := ("(" + strings.Repeat("#", width) + strings.Repeat(" ", progressBarLen-width) + ")")
	if width == progressBarLen {
		fmt.Printf("%v, %02d%%, %d / %d bytes\n", progressBar, percentComplete, current, total)
	} else {
		fmt.Printf("%v, %02d%%, %d / %d bytes\r", progressBar, percentComplete, current, total)
	}
}

// ChooseUploadFolder opens a dialog to select the upload destination folder
func ChooseUploadFolder(folderList []string) (folderName string) {
	var input string
	for {
		for idx, folder := range folderList {
			fmt.Printf("[%v] %v\n", idx, folder)
		}
		fmt.Print("Upload file to folder: ")
		fmt.Scan(&input)
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
	fmt.Printf("Delete file with path %v? (y/n):\n", filePath)
	for {
		fmt.Scan(&input)
		input = strings.ToLower(input)
		if input == "y" {
			err := os.RemoveAll(filePath)
			if err != nil {
				log.Println("Could not delete local file copy ", err)
			} else {
				log.Println("Successfully deleted file at ", filePath)
				fmt.Println("Listening...")
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
