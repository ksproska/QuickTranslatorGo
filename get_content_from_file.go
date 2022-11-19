package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var fileNotFoundError = fmt.Errorf("file with id not found")

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

func getWebsitePattern() string {
	readFile, err := os.ReadFile("youtube_subs.html")
	if err != nil {
		panic(err)
	}
	return string(readFile)
}
