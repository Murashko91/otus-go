package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {

	result := strings.Builder{}

	// digit regexp
	re, _ := regexp.Compile(`[0-9]{1}`)

	// Find all substring indexes devided by digits
	indexes := re.FindAllStringSubmatchIndex(input, -1)

	lastIndex := 0

	for _, indexPair := range indexes {
		substring := input[lastIndex:indexPair[1]]

		processedSubString, err := processSubstring(substring)

		if err != nil {

			if processedSubString != "" {
				// Handle case if substing ends on escaped digit, (join with next substring).
				continue
			} else {
				return "", err
			}
		}
		lastIndex = indexPair[1]
		result.WriteString(processedSubString)
	}

	// Parse last part if exists
	if lastIndex < len(input) {

		processedSubString, err := processSubstring(input[lastIndex:])
		if err != nil {

			// Handle case if substing ends on escaped digit, (no need join with next substring because of it's last one)
			if processedSubString != "" {
				result.WriteString(processedSubString)
			} else {
				return "", err
			}

		} else {
			result.WriteString(processedSubString)
		}
	}

	return result.String(), nil
}

// handle substing by symbols
func processSubstring(substring string) (string, error) {

	result := strings.Builder{}

	s := strings.Split(substring, "")

	// Used if there is escaped char, so we skip adding / char to stringBuilder
	needEscape := false

	for i, symbol := range s {

		_, digitError := strconv.Atoi(symbol)
		isNextLast := i == len(s)-2
		isLast := i == len(s)-1

		// Handle case if non-escaped digit is located not in the end of substrings
		if digitError == nil && !needEscape {
			return "", ErrInvalidString
		}

		if isLast {
			if needEscape && digitError == nil {
				result.WriteString(symbol)
				// Handle case if there are escaped digit in the end of substring
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
			// Skip writing the symbol to StringBuilder
			continue
		} else {
			// Handle non-escaped digit in the end of substrings
			if isNextLast && nextDigitError == nil {
				result.WriteString(strings.Repeat(symbol, nextDigit))
				return result.String(), nil
			}
			result.WriteString(symbol)
		}

	}

	return result.String(), nil
}
