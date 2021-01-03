# moodle-sync

Problem: my university uses a number of platforms to distribute content to students (Blackboard, Moodle, Gradescope, etc). I like to keep a backup of all course content in Google Drive to keep everything organized in one place, but it's tedious to download, upload, and organize everything any time a new document or resource is posted. This is the solution: 

## Setup 

1. Create a Google Cloud Project: `https://developers.google.com/drive/api/v3/quickstart/go`
2. Download the client configuration `credentials.json` to your working directory (root of the cloned repo)
3. Create a `config.yml` file in the `config` folder, see config section for necessary fields 
4. Run: `go run main.go`

> Note: first time setup includes authenticating using the credentials of the target Google Drive account. 

## Config 

```
drive: 
  folders: [
    "Underwater Basket Weaving",
    "Gopher Pictures",
    ...
  ]

local:
  syncfolder: "/Users/vsoltan/Downloads/listeningFolder"
  credpath: "/Users/vsoltan/Documents/src/moodle-sync/credentials.json"
```

> Note: all paths are absolute 


