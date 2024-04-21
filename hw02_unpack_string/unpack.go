package hw02unpackstring

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {

	result := strings.Builder{}

	re, _ := regexp.Compile(`[0-9]{1}`)
	indexes := re.FindAllStringSubmatchIndex(input, -1)

	lastIndex := 0

	fmt.Println(indexes)
	for _, indexPair := range indexes {
		substring := input[lastIndex:indexPair[1]]

		processedSubString, err := processSubstring(substring)

		if err != nil {
			if processedSubString != "" {
				continue
			} else {
				return "", err
			}
		}
		lastIndex = indexPair[1]
		result.WriteString(processedSubString)
	}
	if lastIndex < len(input) {

		processedSubString, err := processSubstring(input[lastIndex:])
		if err != nil {
			result.WriteString(processedSubString)
		} else {
			result.WriteString(processedSubString)
		}
	}

	return result.String(), nil
}

func processSubstring(substring string) (string, error) {
	result := strings.Builder{}

	s := strings.Split(substring, "")
	needEscape := false

	for i, symbol := range s {

		_, digitError := strconv.Atoi(symbol)

		isNextLast := i == len(s)-2
		isLast := i == len(s)-1

		if digitError == nil && !needEscape {

			return "", ErrInvalidString
		}

		if isLast && i >= 0 {
			fmt.Println(i)
			if needEscape && digitError == nil {
				result.WriteString(symbol)
				return result.String(), errors.New("string should be extended by new chank")
			}

			result.WriteString(symbol)
			return result.String(), nil
		}

		nextSymbol := s[i+1]
		nextDigit, nextDigitError := strconv.Atoi(nextSymbol)

		if symbol == "\\" && nextSymbol != "\\" && nextDigitError != nil && !needEscape {
			return "", ErrInvalidString
		}

		needEscape = !needEscape && symbol == "\\" && (nextSymbol == "\\" || nextDigitError == nil)

		if needEscape {
			continue
		} else {
			if isNextLast && nextDigitError == nil {
				result.WriteString(strings.Repeat(symbol, nextDigit))
				return result.String(), nil
			}
			result.WriteString(symbol)
		}

	}

	return result.String(), nil
}
