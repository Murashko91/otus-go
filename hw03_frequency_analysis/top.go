package hw03frequencyanalysis

import (
	"regexp"
	"strings"
)

var (
	pMarksOnlyRegexp  = regexp.MustCompile(`^[!?",.:;_-][!?",.:_-]*[!?",.:;_-]$`)
	pMarksCleanRegexp = regexp.MustCompile(`[0-9a-zа-я]+([!?",.:;_-]*[0-9a-zа-я])*`)
)

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
	// Length of result slice
	maxLen := 10

	if len(wordsCountMap) < 10 {
		maxLen = len(wordsCountMap)
	}

	result := make([]string, 0, maxLen)

	for i := 0; i < maxLen; i++ {
		topKey := ""
		topValue := 0

		// Get top value
		for k, v := range wordsCountMap {
			if v > topValue {
				topKey = k
				topValue = v
			} else if v == topValue && k < topKey {
				topKey = k
				topValue = v
			}
		}

		// Exclude top values from next iterations
		delete(wordsCountMap, topKey)
		result = append(result, topKey)
	}
	return result
}

func processWord(word string) string {
	// Check if word contains punctuation marks only (2+ symbols)
	if pMarksOnlyRegexp.MatchString(word) {
		return word
	}

	// Clean punctuation marks before and after the word
	return pMarksCleanRegexp.FindString(word)
}
