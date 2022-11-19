package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var websitePattern = getWebsitePattern()

const expectedIdLen = 11

func handle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println(fmt.Sprintf("Handling request for %s", id))
	if len(id) != expectedIdLen {
		_, err := w.Write([]byte(fmt.Sprintf("<h1>Incorrect video id (length %d, expected %d)</h1>", len(id), expectedIdLen)))
		logErr(err)
		return
	}
	srtContent, err := getSubtitlesForVideoId(id)
	logErr(err)
	if err != nil {
		_, err = w.Write([]byte(fmt.Sprintf(websitePattern, err.Error(), id)))
		logErr(err)
		return
	}
	htmlContent := prepareContent(srtContent)
	_, err = w.Write([]byte(fmt.Sprintf(websitePattern, htmlContent, id)))
	logErr(err)
}

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.Handle("/youtube/{id}", http.HandlerFunc(handle))

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Println(err)
	}
}
