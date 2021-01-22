package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error Retrieving the File", err)
		http.Error(w, "Couldn't receive your file", 400)
		return
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	mimeType := handler.Header.Get("Content-Type")
	if mimeType != "image/png" && mimeType != "image/jpeg" {
		http.Error(w, "Unsupported file type", 400)
		return
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", 400)
		return
	}

	output, err := RecieveHTTPImage(fileBytes)
	if err != nil {
		log.Println("Could not process image, returning http error", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	jsonOutput := string(output)

	parsedTemplate, err := template.ParseFiles("./static/output.html")
	if err != nil {
		log.Println("Cannot parse template file", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = parsedTemplate.Execute(w, jsonOutput)
	if err != nil {
		log.Println("Error executing template :", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	data := TodoPageData{
// 		PageTitle: "My TODO list",
// 		Todos: []Todo{
// 			{Title: "Task 1", Done: false},
// 			{Title: "Task 2", Done: true},
// 			{Title: "Task 3", Done: true},
// 		},
// 	}
// 	tmpl.Execute(w, data)
// })

func SetupWebServer() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/upload", HandleUpload)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	http.ListenAndServe(":"+port, nil)
}
