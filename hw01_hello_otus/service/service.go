package service

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func OutputResult(input string) {

	fmt.Println(reverse.String(input))
}
