package main

import (
	"fmt"
	"os"
)

func main() {
	var protocol, host, port string
	for i, arg := range os.Args {
		switch i {
		case 1:
			protocol = arg
		case 2:
			host = arg
		case 3:
			port = arg
		}

		fmt.Println(protocol, port, host)

	}

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
