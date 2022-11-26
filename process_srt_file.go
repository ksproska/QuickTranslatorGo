package main

import (
	"fmt"
	"strings"
	"time"
)

type subtitle struct {
	index     int
	startTime string
	endTime   string
	text      string
}
type subtitles []subtitle

func removeDescriptors(content string) string {
	return strings.Join(strings.Split(content, "\n\n")[1:], "\n\n")
}

func getSecFromTime(timeStr string) (string, error) {
	timeSplited := strings.Split(timeStr, ":")
	if len(timeSplited) < 3 {
		return "", fmt.Errorf("wrong input: %s", timeStr)
	}
	timeStr = fmt.Sprintf("%sh%sm%ss", timeSplited[0], timeSplited[1], timeSplited[2])
	duration, err := time.ParseDuration(timeStr)
	if err != nil {
		return "", err
	}
	timeStr = fmt.Sprintf("%v", duration.Seconds())
	return timeStr, nil
}

func readSubtitles(content string) subtitles {
	var subs []subtitle
	content = removeDescriptors(content)
	for i, line := range strings.Split(content, "\n\n") {
		split := strings.Split(line, "-->")
		startTime := strings.Trim(split[0], " ")
		startTime, _ = getSecFromTime(startTime)

		if len(split) > 1 {
			split2 := strings.SplitN(split[1], "\n", 2)
			if len(split2) > 1 {
				endTime := strings.Trim(split2[0], " ")
				endTime, _ = getSecFromTime(endTime)
				subs = append(subs, subtitle{i, startTime, endTime, strings.ReplaceAll(split2[1], "\n", " ")})
			}
		}
	}
	return subs
}

func addPSubs(subs subtitles) string {
	var newLines []string
	for _, s := range subs {
		newLines = append(newLines,
			fmt.Sprintf(
				`
<p hidden id="%d_timestamp">
	%s
</p>
<p id="%d">
	%s
</p>`,
				s.index, s.startTime, s.index, s.text),
		)
	}
	return strings.Join(newLines, "\n")
}

func prepareContent(content string) string {
	subs := readSubtitles(content)
	content = addPSubs(subs)
	return content
}
