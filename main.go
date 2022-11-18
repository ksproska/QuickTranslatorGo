package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var fileNotFoundError = fmt.Errorf("file with id not found")

func downloadSubtitles(id string) error {
	url := "https://www.youtube.com/watch?v=" + id
	cmd := exec.Command("youtube-dl", "--all-subs", "--skip-download", url)
	err := cmd.Run()
	return err
}

func getContentFromFileContainingId(id string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	files, err := os.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), id) {
			readFile, err := os.ReadFile(file.Name())
			if err != nil {
				return "", err
			}
			return string(readFile), nil
		}
	}
	return "", fileNotFoundError
}

func getSubtitlesForVideoId(id string) (string, error) {
	err := downloadSubtitles(id)
	if err != nil {
		return "", err
	}
	content, err := getContentFromFileContainingId(id)
	if err != nil {
		return "", err
	}
	return content, nil
}

func removeLinesContainingTimestamps(content string) string {
	lines := strings.Split(content, "\n")
	leftLines := []string{}
	for _, line := range lines {
		if !strings.Contains(line, "-->") && line != "" {
			leftLines = append(leftLines, line)
		}
		if line == "" {
			leftLines = append(leftLines, "\n")
		}
	}
	return strings.Join(leftLines, "")
}

func getWebsitePattern() string {
	readFile, err := os.ReadFile("youtube_subs.html")
	if err != nil {
		panic(err)
	}
	return string(readFile)
}

func main() {
	r := mux.NewRouter()
	websitePattern := getWebsitePattern()

	r.Handle("/youtube/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		content, err := getSubtitlesForVideoId(id) // Ye8mB6VsUHw
		log.Println(err)
		content = removeLinesContainingTimestamps(content)
		content = strings.ReplaceAll(content, "\n", "<br>")
		_, _ = w.Write([]byte(fmt.Sprintf(websitePattern, content, id)))
		time.Sleep(5 * time.Second)
		fmt.Println(w.Header())
	}))

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
