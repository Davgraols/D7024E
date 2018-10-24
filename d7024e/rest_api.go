package main

import (
	"encoding/json"
	"fmt"
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
	router.HandleFunc("/file/{id}", PinFile).Methods("PATCH") //handle pinfile and unpinfile
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

type UnpinResponse struct {
	FileID     string `json:"fileID"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type FileContent struct {
	Data string `json:"data"`
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	response := GetFileResponse{
		File:       params["id"],
		Successful: true,
		Message:    "this is a message",
	}
	json.NewEncoder(w).Encode(response)
}

func SaveFile(w http.ResponseWriter, r *http.Request) {

	var file FileContent
	_ = json.NewDecoder(r.Body).Decode(&file)

	response := StoreResponse{
		FileID:     "A new id",
		Successful: true,
		Message:    "file: " + file.Data,
	}
	json.NewEncoder(w).Encode(response)
}

/*
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "im in DeleteFile")
	params := mux.Vars(r)
	response :=
	files = templist
}
*/

func PinFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "im in PinFile")
}

func UnpinFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "im in UnpinFile")
}
