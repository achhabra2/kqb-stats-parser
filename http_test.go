package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"
)

func CreateImageFormFile(w *multipart.Writer, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "images", filename))
	h.Set("Content-Type", "image/png")
	return w.CreatePart(h)
}

func TestUploadHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	cases := []struct {
		description string
		filepath    string
		statusCode  int
	}{
		{
			description: "Valid PNG File",
			filepath:    "./example/input.png",
			statusCode:  200,
		},
	}

	for _, c := range cases {
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		file, errFile1 := os.Open(c.filepath)
		if errFile1 != nil {
			t.Fatal(errFile1)
		}
		defer file.Close()
		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
		// part1, errFile1 := writer.CreateFormFile("images", filepath.Base(c.filepath))
		part, err := CreateImageFormFile(writer, filepath.Base(c.filepath))
		if err != nil {
			t.Fatal(err)
		}
		part.Write(fileContents)
		err = writer.Close()
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/api", payload)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleUpload)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)
		// resp := rr.Result()
		// body, _ := ioutil.ReadAll(resp.Body)

		// fmt.Println(resp.StatusCode)
		// fmt.Println(string(body))
		// Check the status code is what we expect.
		if status := rr.Code; status != c.statusCode {
			t.Errorf("handler returned wrong status code: got %v want %v, message %s",
				status, c.statusCode, rr.Body.String())
		}
		// var sets []Set
		// err = json.Unmarshal(rr.Body.Bytes(), &sets)
		// if err != nil {
		// 	t.Errorf("Invalid response: %v", err.Error())
		// }
		// if len(sets) != 1 {
		// 	t.Errorf("Did not receive Set data back")
		// }
	}

	// // Check the response body is what we expect.
	// expected := `{"alive": true}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
