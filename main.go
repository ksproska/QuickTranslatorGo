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
			leftLines = append(leftLines, line+" ")
		}
		if line == "" {
			leftLines = append(leftLines, "\n")
		}
		if strings.Contains(line, "-->") {
			leftLines = append(leftLines, strings.Split(line, " -->")[0]+"\t")
		}
	}
	return strings.Join(leftLines, "")
}

func getWebsitePattern() string {
	fmt.Println(os.Getwd())
	readFile, err := os.ReadFile("youtube_subs.html")
	if err != nil {
		panic(err)
	}
	return string(readFile)
}

func addP(content string) string {
	lines := strings.Split(content, "\n")
	newLines := []string{}
	for i, line := range lines {
		timestamp := strings.Split(line, "\t")[0]
		timeSplited := strings.Split(timestamp, ":")
		if len(timeSplited) >= 3 {
			timestamp = fmt.Sprintf("%sh%sm%ss", timeSplited[0], timeSplited[1], timeSplited[2])
			duration, _ := time.ParseDuration(timestamp)
			timestamp = fmt.Sprintf("%v", duration.Seconds())
		}

		if strings.Contains(line, "\t") {
			line = strings.Split(line, "\t")[1]
		}
		newLines = append(newLines, fmt.Sprintf("<p hidden id=\"%d_timestamp\">%s</p><p id=\"%d\">%s</p>", i, timestamp, i, line))
	}
	return strings.Join(newLines, "\n")
}

func main() {
	r := mux.NewRouter()
	websitePattern := getWebsitePattern()

	r.Handle("/youtube/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		content, err := getSubtitlesForVideoId(id)
		if err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf(websitePattern, err.Error(), id)))
		}
		content = prepareContent(content)
		_, _ = w.Write([]byte(fmt.Sprintf(websitePattern, content, id)))
	}))

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}

func prepareContent(content string) string {
	content = removeLinesContainingTimestamps(content)
	content = strings.Join(strings.Split(content, "\n")[1:], "\n")
	content = addP(content)
	return content
}
