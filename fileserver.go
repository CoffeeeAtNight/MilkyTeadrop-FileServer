package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type FileseverApiPostRequest struct {
	Filename    string `json:"filename"`
	Filetype    string `json:"filetype"`
	FileContent string `json:"fileContent"`
}

type FileserverApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Body    string `json:"body"`
}

var FILESERVER_ROOT_PATH = "/usr/local/bin/milkyteadrop-fs/"

func main() {
	fmt.Println("» Starting Milkyteadrop Fileserver...")
	fmt.Println("» Fileserver should be located at: ", FILESERVER_ROOT_PATH+"milkyteadrop-fileserver")
	listFiles()
	setupApi()
}

func listFiles() {
	println("» Dir = Directory, F = File\n")
	println("» Current files:")
	dirList, err := os.ReadDir(FILESERVER_ROOT_PATH)
	if err != nil {
		fmt.Println("Error occurred while trying to list files in fileserver root", err.Error())
		return
	}

	for i, entry := range dirList {
		fmt.Printf("%d. ", i)
		if entry.IsDir() {
			fmt.Println("Dir: ", entry.Name())
		} else {
			fmt.Println("F: ", entry.Name())
		}
	}
	println("—————————————————————————————————————————————————————")
}

// --- API --- //

func setupApi() {
	println("» Starting up on port :7676...")

	// POST
	http.HandleFunc("/api/v1/create/file", handleFileCreation)

	//GET
	http.Handle("/api/v1/file/", http.StripPrefix("/api/v1/file/", http.FileServer(http.Dir(FILESERVER_ROOT_PATH))))

	if err := http.ListenAndServe(":7676", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Handlers
func handleFileCreation(w http.ResponseWriter, req *http.Request) {

	log.Printf("» Received request: Method=%s, URL=%s", req.Method, req.URL.Path)

	if req.Method != "POST" {
		log.Printf("/create/file endpoint was called with wrong HTTP method")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusMethodNotAllowed)
		return
	}

	requestBodyAsBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var fileCreationRequest FileseverApiPostRequest
	if err := json.Unmarshal(requestBodyAsBytes, &fileCreationRequest); err != nil {
		log.Printf("Error while trying to unmarshal response: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := checkIfFileAlreadyExists(fileCreationRequest.Filename); err != nil {
		http.Error(w, "File already exists", http.StatusInternalServerError)
		return
	}

	handleFileWriting(w, fileCreationRequest)
	replySuccess(w, fileCreationRequest)
}

func replySuccess(w http.ResponseWriter, request FileseverApiPostRequest) {
	response := FileserverApiResponse{
		Status:  200,
		Message: "Successfully created file",
		Body:    request.Filename,
	}

	res_json, err := json.Marshal(response)
	if err != nil {
		println("Error marshalling the response json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res_json)
}

func handleFileWriting(w http.ResponseWriter, request FileseverApiPostRequest) {
	reqFileType := request.Filetype
	if reqFileType == "image" {

		fileContentBase64, err := base64.StdEncoding.DecodeString(request.FileContent)
		if err != nil {
			log.Printf("Error decoding base64 data: %v", err)
			http.Error(w, "Invalid base64 data", http.StatusBadRequest)
			return
		}

		if err := os.WriteFile(FILESERVER_ROOT_PATH+request.Filename, fileContentBase64, 0666); err != nil {
			log.Printf("Error writing base64 to file: %v", err)
			http.Error(w, "Error writing base64 to file", http.StatusInternalServerError)
			return
		}
	} else if reqFileType == "text" {
		if err := os.WriteFile(FILESERVER_ROOT_PATH+request.Filename, []byte(request.FileContent), 0666); err != nil {
			log.Printf("Error writing file: %v", err)
			http.Error(w, "Error writing file", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Filetype is not supported")
		http.Error(w, "Filetype is not supported", http.StatusUnsupportedMediaType)
	}
}

func checkIfFileAlreadyExists(filename string) error {
	_, err := os.Stat(FILESERVER_ROOT_PATH + filename)
	if err == nil {
		log.Printf("File with name %s already exists in directory", filename)
		return fmt.Errorf("file %s already exists", filename)
	} else if os.IsNotExist(err) {
		return nil
	}
	return err
}
