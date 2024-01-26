package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type FileApiRequest struct {
	filepath string
}

type FileserverApiResponse struct {
	Status  int
	Message string
	Body    string
}

var FILESERVER_ROOT_PATH = "/usr/local/bin/milkyteadrop-fs/"

func main() {
	fmt.Println("» Starting Milkyteadrop Fileserver...")
	fmt.Println("» Fileserver should be located at: ", FILESERVER_ROOT_PATH+"milkyteadrop-fileserver")
	fmt.Println("» Now checking for location of binary...")
	checkExecutionLocation()

	println("» Current files:")
	listFiles()

	setupRestApi()
}

func listFiles() {
	dirList, err := os.ReadDir(FILESERVER_ROOT_PATH)
	if err != nil {
		fmt.Println("Error occurred while trying to list files in fileserver root", err.Error())
	}

	for i, entry := range dirList {
		println("» Dir = Directory, F = File")
		fmt.Print(string(rune(i)) + ". \n")
		if entry.IsDir() {
			fmt.Println("Dir: ", entry.Name())
		} else {
			fmt.Println("F: ", entry.Name())
		}
	}
	println("—————————————————————————————————————————————————————")
}

// --- REST API --- //

func setupRestApi() {
	println("» Starting up Rest-Api on port :7777...")
	http.HandleFunc("/api/v1/file/", handleFileRetrieval)

	http.ListenAndServe(":7777", nil)
}

// Handlers

func handleFileRetrieval(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		statusCode := 405
		res, _ := json.Marshal(createApiErrorResponse(statusCode, "Method for /file not allowed"))
		w.WriteHeader(statusCode)
		w.Write(res)
	}

	fmt.Println("» GET /file")

	http.Handle(req.RequestURI, http.FileServer(http.Dir(FILESERVER_ROOT_PATH)))
	w.Write([]byte{})
}

func createApiErrorResponse(status int, errorMsg string) FileserverApiResponse {
	return FileserverApiResponse{
		Status:  status,
		Message: errorMsg,
		Body:    "",
	}
}

func checkExecutionLocation() {
	path := os.Args[0]
	expectedPath := FILESERVER_ROOT_PATH + "milkyteadrop-fileserver"
	if path != expectedPath {
		fmt.Println("» ERROR: Fileserver is not located in :", expectedPath)
		fmt.Println("» ERROR: Location is :", path)
		os.Exit(1)
	}

	fmt.Println("» Location is OK")
}
