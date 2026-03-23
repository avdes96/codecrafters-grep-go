package parse

import "fmt"

var HadParseError = false

func PrintErrorMessage(msg string) {
	HadParseError = true
	fmt.Printf("grep: %s\n", msg)
}
