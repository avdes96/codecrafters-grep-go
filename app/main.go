package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/grep-starter-go/app/pattern"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	patternStr := os.Args[2]

	input, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	line := []rune(string(input))

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	patternList := []rune(patternStr)
	pattern, err := pattern.Parse(patternList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: parsing pattern: %v\n", err)
		os.Exit(2)
	}

	if ok := matchLine(line, pattern); !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []rune, pattern *pattern.Pattern) bool {
	if pattern.IsEmpty() {
		return true
	}
	if pattern.MustMatchStart {
		return match(line, 0, pattern, pattern.Head)
	}
	for lineIdx := range len(line) {
		if ok := match(line, lineIdx, pattern, pattern.Head); ok {
			return true
		}
	}
	return false
}

func match(line []rune, lineIdx int, pattern *pattern.Pattern, current *pattern.Node) bool {
	if current.IsEmpty() {
		if !pattern.MustMatchEnd {
			return true
		}
		// Matching end of string anchor can probably be made more efficient, fine for now
		return lineIdx == len(line)

	}

	if lineIdx == len(line) {
		return false
	}

	if !current.Match(line[lineIdx]) {
		if !current.IsOptional() {
			return false
		}
		return match(line, lineIdx, pattern, current.Next)
	}

	if current.IsOneOrMore() {
		return matchOneOrMore(line, lineIdx, pattern, current)
	}
	return match(line, lineIdx+1, pattern, current.Next)
}

func matchOneOrMore(line []rune, lineIdx int, pattern *pattern.Pattern, current *pattern.Node) bool {
	if lineIdx == len(line) {
		return current == nil
	}
	if current.Match(line[lineIdx]) && matchOneOrMore(line, lineIdx+1, pattern, current) {
		return true
	}
	return match(line, lineIdx, pattern, current.Next)
}
