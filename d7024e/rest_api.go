package main

import (
	"encoding/json"
	"github.com/gorilla/mux" //glöm inte att köra "github.com/gorilla/mux" om krash
	"log"
	"net/http"
)

// our main function
func Mainrest() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/file/{id}", GetFile).Methods("GET") //handle gefile
	router.HandleFunc("/file", SaveFile).Methods("POST")    //handle savefile
	//router.HandleFunc("/file/{id}", DeleteFile).Methods("DELETE") //handle deletefile
	router.HandleFunc("/file/{id}", PinFile).Methods("PATCH") //handle pinfile
	log.Fatal(http.ListenAndServe(":8080", router))

}

type StoreResponse struct {
	FileID     string `json:"fileID"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type GetFileResponse struct {
	File       string `json:"file"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type PinResponse struct {
	FileID     string `json:"fileID"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type FileContent struct {
	Data string `json:"data"`
}

type PinStatus struct {
	PinType string `json:"pintype"`
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	response := GetFileResponse{}
	if len(id) != 40 {
		response.File = ""
		response.Successful = false
		response.Message = "Invalid id"
	} else {
		fileID := NewKademliaID(id)
		file, fileExist := FS.getFile(fileID)

		if fileExist {
			response.File = string(file.content)
			response.Successful = true
			response.Message = "File found"
		} else {
			foundFile := KademliaObj.LookupData(fileID)
			if foundFile != nil {
				response.File = string(foundFile)
				response.Successful = true
				response.Message = "File found"
			} else {
				response.File = ""
				response.Successful = false
				response.Message = "File not found"
			}
		}
	}

	json.NewEncoder(w).Encode(response)
}

func SaveFile(w http.ResponseWriter, r *http.Request) {

	var file FileContent
	_ = json.NewDecoder(r.Body).Decode(&file)

	fileBytes := []byte(file.Data)

	//FS.StoreFile(fileBytes, &RT.me)
	KademliaObj.Store(fileBytes, &RT.me)
	fileID := NewRandomHash(file.Data)
	_, success := FS.getFile(fileID)
	response := StoreResponse{}
	if success {
		response.FileID = fileID.String()
		response.Successful = true
		response.Message = "File was stored successfully"
	} else {
		response.Successful = false
		response.Message = "File could not be stored"
	}
	json.NewEncoder(w).Encode(response)
}

func PinFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	response := PinResponse{}
	if len(id) != 40 {
		response.FileID = id
		response.Successful = false
		response.Message = "Invalid id"
	} else {
		fileID := NewKademliaID(id)
		_, fileExist := FS.getFile(fileID)
		if fileExist {
			var pin PinStatus
			_ = json.NewDecoder(r.Body).Decode(&pin)
			if pin.PinType == "pin" {
				KademliaObj.Pin(fileID)
			} else {
				KademliaObj.Unpin(fileID)
			}
			response.FileID = id
			response.Successful = true
			response.Message = "File was " + pin.PinType + "ned successfully"
		} else {
			response.FileID = id
			response.Successful = false
			response.Message = "File not found"
		}
	}
	json.NewEncoder(w).Encode(response)
}
