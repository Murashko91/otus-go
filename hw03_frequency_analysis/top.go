package hw03frequencyanalysis

import "strings"

func Top10(input string) []string {

	words := strings.Fields(input)

	wordsCountMap := make(map[string]int)

	for _, word := range words {

		value, isPresent := wordsCountMap[word]

		if isPresent {
			wordsCountMap[word] = value + 1
		} else {
			wordsCountMap[word] = 1
		}
	}

	return nil
}

func sortTop10(wordsCountMap map[string]int) []string {

	result := make([]string, 0, 10)

	for i := 0; i < 10; i++ {

	}
	return nil

}
