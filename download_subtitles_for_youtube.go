package main

import (
	"fmt"
	"os/exec"
)

func downloadSubtitles(id string) error {
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
	cmd := exec.Command("youtube-dl", "--all-subs", "--skip-download", url)
	err := cmd.Run()
	return err
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
