package main

import (
	"fmt"
	gt "github.com/bas24/googletranslatefree"
	"strings"
	"sync"
	"time"
)

type subtitle struct {
	index       int
	startTime   string
	endTime     string
	text        string
	translation string
}
type subtitles []*subtitle

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

func readSubtitles(content string, sourceLang, targetLang string) subtitles {
	var subs subtitles
	content = removeDescriptors(content)
	lines := strings.Split(content, "\n\n")
	wg := sync.WaitGroup{}
	for i, line := range lines {
		split := strings.Split(line, "-->")
		startTime := strings.Trim(split[0], " ")
		startTime, _ = getSecFromTime(startTime)

		if len(split) > 1 {
			split2 := strings.SplitN(split[1], "\n", 2)
			if len(split2) > 1 {
				endTime := strings.Trim(split2[0], " ")
				endTime, _ = getSecFromTime(endTime)
				text := strings.ReplaceAll(split2[1], "\n", " ")
				sub := subtitle{
					i,
					startTime,
					endTime,
					text,
					"",
				}
				wg.Add(1)
				go func(s *subtitle) {
					translation, _ := gt.Translate(text, sourceLang, targetLang)
					s.translation = translation
					wg.Done()
				}(&sub)
				subs = append(subs, &sub)
			}
		}
	}
	wg.Wait()
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
</p>
<p class="translation" id="%d_translation">
	&emsp;&emsp;%s
</p>`,
				s.index, s.startTime, s.index, s.text, s.index, s.translation),
		)
	}
	return strings.Join(newLines, "\n")
}

func getLanguage(content string) (string, error) {
	splitted := strings.Split(content, "Language: ")
	if len(splitted) < 2 {
		return "", fmt.Errorf("language not found")
	}
	return strings.Split(splitted[1], "\n")[0], nil
}

func prepareContent(content string) string {
	sourcelang, _ := getLanguage(content)
	subs := readSubtitles(content, sourcelang, "pl")
	content = addPSubs(subs)
	return content
}
