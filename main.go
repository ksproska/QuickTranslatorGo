package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
	"runtime"
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

	//go func() {
	//	time.Sleep(3 * time.Second)
	//	_ = open("http://localhost:3000/youtube/u1yUxpC0xgs")
	//}()
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Println(err)
	}
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
