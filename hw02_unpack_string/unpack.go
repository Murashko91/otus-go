package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	result := strings.Builder{}

	s := strings.Split(input, "")

	// Ignore digit and escape validation error
	needEscape := false

	// Used to skip digit symbol once it has been already utilized for Repeat action
	needSkip := false

	for i, symbol := range s {
		if needSkip {
			needSkip = false
			continue
		}
		_, digitError := strconv.Atoi(symbol)
		isLast := i == len(s)-1

		if digitError == nil && !needEscape {
			return "", ErrInvalidString
		}

		if isLast {
			result.WriteString(symbol)
			return result.String(), nil
		}

		nextSymbol := s[i+1]
		nextDigit, nextDigitError := strconv.Atoi(nextSymbol)

		fmt.Println(strconv.Atoi(nextSymbol))

		if symbol == "\\" && nextSymbol != "\\" && nextDigitError != nil && !needEscape {
			return "", ErrInvalidString
		}
		needEscape = !needEscape && symbol == "\\" && (nextSymbol == "\\" || nextDigitError == nil)
		if needEscape {
			continue
		}
		if nextDigitError == nil {
			result.WriteString(strings.Repeat(symbol, nextDigit))
			needSkip = true
		} else {
			result.WriteString(symbol)
		}
	}
	return result.String(), nil
}
