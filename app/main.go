package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/grep-starter-go/app/match"
	"github.com/codecrafters-io/grep-starter-go/app/parse"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	patternStr := os.Args[2]
	input, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}
	parser := parse.NewParser(patternStr)
	pattern := parser.Parse()
	if parse.HadParseError {
		os.Exit(2)
	}
	matcher := match.NewMatcherFromPattern([]rune(string(input)), pattern)

	if ok := matcher.Match(); !ok {
		os.Exit(1)
	}
	// default exit code is 0 which means success
}
