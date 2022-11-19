package main

import (
	"fmt"
	"strings"
	"time"
)

func removeLinesContainingTimestamps(content string) string {
	lines := strings.Split(content, "\n")
	var leftLines []string
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

func addP(content string) string {
	lines := strings.Split(content, "\n")
	var newLines []string
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

func prepareContent(content string) string {
	content = removeLinesContainingTimestamps(content)
	content = strings.Join(strings.Split(content, "\n")[1:], "\n")
	content = addP(content)
	return content
}
