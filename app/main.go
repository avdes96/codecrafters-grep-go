package main

import (
	"fmt"
	"io"
	"os"
	"unicode"

	mapset "github.com/deckarep/golang-set/v2"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	input, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	line := string(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line string, pattern string) (bool, error) {
	if pattern == "" {
		return true, nil
	}
	if pattern[0] == '[' {
		if len(pattern) == 2 {
			return false, fmt.Errorf("unmatched [")
		}
		if pattern[1] != '^' {
			return matchPositiveCharacterGroup(line, pattern[1:len(pattern)-1]), nil
		}
		if pattern[2] == ']' {
			return false, fmt.Errorf("unmatched [")
		}
		return matchNegativeCharacterGroup(line, pattern[2:len(pattern)-1]), nil
	}
	switch pattern {
	case "\\d":
		return matchDigit(line), nil
	case "\\w":
		return matchWordCharacterClass(line), nil
	default:
		return matchLiteralChar(line, pattern)
	}
}

func matchPositiveCharacterGroup(line string, chars string) bool {
	charSet := newCharSet(chars)
	for _, c := range line {
		if charSet.Contains(c) {
			return true
		}
	}
	return false
}

func matchNegativeCharacterGroup(line string, chars string) bool {
	charSet := newCharSet(chars)
	for _, c := range line {
		if !charSet.Contains(c) {
			return true
		}
	}
	return false
}

func newCharSet(chars string) mapset.Set[rune] {
	charSet := mapset.NewSet[rune]()
	for _, c := range chars {
		charSet.Add(c)
	}
	return charSet
}

func matchDigit(line string) bool {
	for _, c := range line {
		if unicode.IsDigit(rune(c)) {
			return true
		}
	}
	return false
}

func matchWordCharacterClass(line string) bool {
	for _, b := range line {
		r := rune(b)
		if r == '_' {
			return true
		}
		if unicode.IsDigit(r) {
			return true
		}
		if unicode.IsUpper(r) {
			return true
		}
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func matchLiteralChar(line string, pattern string) (bool, error) {
	if len(pattern) != 1 {
		return false, fmt.Errorf("expected pattern of len 1, got %s of len %d", pattern, len(pattern))
	}

	for _, c := range line {
		if string(c) == pattern {
			return true, nil
		}
	}

	return false, nil
}
