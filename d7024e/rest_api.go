package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux" //glöm inte att köra "github.com/gorilla/mux" om krash
	"log"
	"net/http"
	//"github.com/gorilla/mux"
)

// our main function
func Mainrest() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/file/{id}", GetFile).Methods("GET")       //handle gefile
	router.HandleFunc("/file/{id}", SaveFile).Methods("POST")     //handle savefile
	router.HandleFunc("/file/{id}", DeleteFile).Methods("DELETE") //handle deletefile
	router.HandleFunc("/file/{id}", PinFile).Methods("POST")      //handle pinfile
	router.HandleFunc("/file/{id}", UnpinFile).Methods("DELETE")  //handle unpinfile
	log.Fatal(http.ListenAndServe(":8080", router))

}

type File struct {
	Files  string `json:"File"`
	FileID string `json:"FileID"`
}

var files []File //byt från lista till stukt när bygger ihop med david

func GetFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "im in GetFile")
	params := mux.Vars(r)
	var file File
	for _, item := range files {
		if item.FileID == params["FileID"] {
			file = item //bygg om när bygger ihop med david
		}
	}
	json.NewEncoder(w).Encode(file)
}

func SaveFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "im in SaveFile")
	params := mux.Vars(r)
	file := File{
		Files:  params["File"],
		FileID: params["FileID"],
	}
	files = append(files, file)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "im in DeleteFile")
	params := mux.Vars(r)
	var templist []File
	for _, item := range files {
		if item.FileID != params["FileID"] {
			templist = append(templist, item) //bygg om när bygger ihop med david
		}
	}
	files = templist
}

func PinFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "im in PinFile")
}

func UnpinFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "im in UnpinFile")
}
