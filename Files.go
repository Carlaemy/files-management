package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Document struct {
	Id   int
	Name string
	Size int
}

func main() {
	router := mux.NewRouter()
	router.Handle("/", http.FileServer(http.Dir("./public")))
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	router.HandleFunc("/documents/{Id}", getDocumentById).Methods("GET")
	router.HandleFunc("/delete/{Id}", deleteDocument).Methods("DELETE")
	router.HandleFunc("/add", addDocument).Methods("POST")
	log.Printf("Listening in http://localhost:9000")
	http.ListenAndServe(":9000", router)
}

func addDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, handle, err := r.FormFile("UploadFile")
		if err != nil {
			log.Printf("Error loading the file %v", err)
			fmt.Fprintf(w, "Error loading the file %v", err)
			return
		}

		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("Error reading the file %v", err)
			fmt.Fprintf(w, "Error reading the file %v", err)
			return
		}

		err = ioutil.WriteFile("./Files/"+handle.Filename, data, 0666)
		if err != nil {
			log.Printf("Error writing the file %v", err)
			fmt.Fprintf(w, "Error writing the file %v", err)
			return
		}

		fmt.Fprintf(w, "Loaded successful!!")
	}
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Generate_List())
}

func getDocumentById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["Id"])

	for _, doc := range Generate_List() {
		if doc.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(doc)
			return
		}
	}
	http.Error(w, "NotFound", 404)
}

func deleteDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["Id"])

	for _, fil := range Generate_List() {
		if fil.Id == id {
			os.Remove("./Files/" + fil.Name)
			fmt.Fprintf(w, "Removed successful!!")
			return
		}
	}
	http.Error(w, "NotFound", 404)
}

func Generate_List() []Document {
	var docs []Document
	id := 1

	f, err := os.Open("./Files")
	if err != nil {
		panic(err)
	}

	files, err := f.Readdir(0)
	if err != nil {
		panic(err)
	}

	for _, v := range files {
		docs = append(docs, Document{Id: id, Name: v.Name(), Size: int(v.Size())})
		id++
	}
	return docs
}
