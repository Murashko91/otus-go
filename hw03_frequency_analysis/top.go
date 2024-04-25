package hw03frequencyanalysis

import (
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`([-]*[a-zа-я]+[-]*[a-zа-я]*[-]*)|([--]{2,})`)

func Top10(input string) []string {

	input = strings.ToLower(input)

	words := strings.Fields(input)

	wordsCountMap := make(map[string]int)

	for _, word := range words {

		fixedWord := processWord(word)

		if fixedWord == "" {
			continue
		}

		value, isPresent := wordsCountMap[fixedWord]

		if isPresent {
			wordsCountMap[fixedWord] = value + 1
		} else {
			wordsCountMap[fixedWord] = 1
		}
	}

	return sortTop10(wordsCountMap)
}

func sortTop10(wordsCountMap map[string]int) []string {

	result := make([]string, 0, 10)

	maxLen := 10

	if maxLen > len(wordsCountMap) {
		maxLen = len(wordsCountMap)
	}

	for i := 0; i < maxLen; i++ {

		topKey := ""
		topValue := 0

		for k, v := range wordsCountMap {
			if v > topValue {
				topKey = k
				topValue = v
			} else if v == topValue && k < topKey {
				topKey = k
				topValue = v
			}
		}

		delete(wordsCountMap, topKey)
		result = append(result, topKey)
	}
	return result
}

func processWord(word string) string {

	return reg.FindString(word)
}
