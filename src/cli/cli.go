package cli

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Progress Bar

// ChooseUploadFolder opens a dialog to select the upload destination folder
func ChooseUploadFolder(folderList []string) (folderName string) {
	var input string
	for {
		fmt.Println("Upload file to folder: ")
		for idx, course := range folderList {
			fmt.Printf("[%v] %v\n", idx, course)
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
	fmt.Println("Delete local file? (y/n)")
	for {
		fmt.Scanln(&input)
		input = strings.ToLower(input)
		if input == "y" {
			err := os.Remove(filePath)
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
