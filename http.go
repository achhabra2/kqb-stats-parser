package main

import (
	"encoding/json"
	"fmt"
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
	// partySize := r.MultipartForm.Value["party"]
	// partySizeInt, err := strconv.Atoi(partySize[0])
	// if err != nil {
	// 	http.Error(w, "Bad Request, unsupport party size, please enter 1-4", 400)
	// 	return
	// }

	fhs := r.MultipartForm.File["images"]
	var finalOutput []Set
	for _, fh := range fhs {
		f, err := fh.Open()
		defer f.Close()
		log.Printf("Received File: %+v\n", fh.Filename)
		mimeType := fh.Header.Get("Content-Type")
		if mimeType != "image/png" && mimeType != "image/jpeg" {
			http.Error(w, "Unsupported file type", 400)
			return
		}
		// f is one of the files
		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(f)
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
		finalOutput = append(finalOutput, output)
	}

	jsonBytes, err := json.MarshalIndent(finalOutput, "", "    ")
	if err != nil {
		log.Println("Cannot write json file", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	jsonOutput := string(jsonBytes)
	fmt.Printf("Request URI: %v\n", r.RequestURI)
	if r.RequestURI == "/upload" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
	} else if r.RequestURI == "/api" {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(jsonBytes)
		if err != nil {
			log.Println("Could not return parsed json", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}

func SetupWebServer() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/upload", HandleUpload)
	http.HandleFunc("/api", HandleUpload)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Could not start http listener")
	}
}
