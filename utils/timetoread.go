package utils

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
)

//TimeToRead will scrap the content from the URL and count the number of words between the <article> tags. If no article
// tag is found it will check the words in the page. Then it will calculate how much time it will take to read the
// page based on a calculation of 200 words per minute

func TimeToRead (url string) (int, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error at ", err)
		return 0, err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	text := string(bodyText)

	preIndex := strings.Index(text, "<article>")
	postIndex := 0
	indexLengh := 0
	if preIndex == -1 {
		preIndex = 0
		postIndex = len(text)
	} else {
		postIndex = strings.Index(text, `</article>`)
		indexLengh = 9
	}

	textToAnalyse := text[preIndex:postIndex + indexLengh]
	index := 0
	var words = ""
	lenOfText := postIndex - preIndex
	for ok := true; ok; ok = index < lenOfText {
		// FindUserById the first closer for a html tag
		closerIndex := strings.Index(textToAnalyse[index:], ">") + 1 + index
		// If the next thing that we see is another tag, loop again
		if closerIndex >= lenOfText{
			break
		}
		if textToAnalyse[closerIndex] == '<' {
			index = closerIndex
		} else {
			//If the next thing that we find is text, we want to save from this index to the next open html tag
			openIndex := strings.Index(textToAnalyse[closerIndex:], "<") +closerIndex
			// This means that the last open bracket that we found is the last one
			if openIndex < closerIndex {
				break
			}
			index = openIndex

			extractedString := textToAnalyse[closerIndex:openIndex]
			extractedString = cleanString(extractedString)

			if extractedString == "" {
				continue
			}
			words += extractedString

		}
	}

	wordsCounted :=strings.Split(words, " ")
	wordsPerMinute := float64(len(wordsCounted)) / float64(200)
	min, seconds := math.Modf(wordsPerMinute)
	totalSeconds := seconds * 0.60

	if totalSeconds < 0.30 {
		return int(min), nil
	} else {
		return int(min + 1), nil
	}
}

func cleanString (text string) string {
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\t", "", -1)
	text = strings.TrimSpace(text)
	return text
}
