package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/grep-starter-go/app/matcher"
	mapset "github.com/deckarep/golang-set/v2"
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
	pattern := []rune(patternStr)

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

func matchLine(line []rune, pattern []rune) (bool, error) {
	for lineIdx := range len(line) {
		m, err := match(line, pattern, lineIdx, 0)
		if err != nil {
			return false, err
		}
		if m {
			return true, nil
		}
	}
	return false, nil
}

func match(line []rune, pattern []rune, lineIdx int, patternIdx int) (bool, error) {
	if patternIdx == len(pattern) {
		return true, nil
	}
	if lineIdx == len(line) {
		return false, nil
	}
	lineChar := line[lineIdx]
	patternMatcher, newPatternIdx, err := getMatcher(pattern, patternIdx)
	if err != nil {
		return false, err
	}
	if !patternMatcher.Match(lineChar) {
		return false, nil
	}
	match, err := match(line, pattern, lineIdx+1, newPatternIdx)
	if err != nil {
		return false, err
	}
	return match, nil
}

func getMatcher(pattern []rune, patternIdx int) (matcher.Matcher, int, error) {
	switch pattern[patternIdx] {
	case '\\':
		m, newPatternIdx, err := getCharacterClassMatcher(pattern, patternIdx+1)
		if err != nil {
			return nil, -1, err
		}
		return m, newPatternIdx, nil
	case '[':
		m, newPatternIdx, err := getCharacterGroupMatcher(pattern, patternIdx+1)
		if err != nil {
			return nil, -1, err
		}
		return m, newPatternIdx, nil
	}
	return matcher.NewLiteralCharMatcher(rune(pattern[patternIdx])), patternIdx + 1, nil
}

func getCharacterClassMatcher(pattern []rune, patternIdx int) (matcher.Matcher, int, error) {
	if patternIdx == len(pattern) {
		return nil, -1, fmt.Errorf("invalid pattern: incomplete character class")
	}
	switch pattern[patternIdx] {
	case 'd':
		return matcher.NewDigitMatcher(), patternIdx + 1, nil
	case 'w':
		return matcher.NewWordCharacterMatcher(), patternIdx + 1, nil
	case '\\':
		return matcher.NewLiteralCharMatcher(rune('\\')), patternIdx + 1, nil
	}
	return nil, -1, fmt.Errorf("invalid pattern: invalid character class %q", pattern[patternIdx])

}

func getCharacterGroupMatcher(pattern []rune, patternIdx int) (matcher.Matcher, int, error) {
	if patternIdx == len(pattern) {
		return nil, -1, fmt.Errorf("invalid pattern: unmatched [")
	}
	if pattern[patternIdx] == '^' {
		return getNegCharacterGroupMatcher(pattern, patternIdx+1)
	}
	return getPosCharacterGroupMatcher(pattern, patternIdx)
}

func getPosCharacterGroupMatcher(pattern []rune, patternIdx int) (matcher.Matcher, int, error) {
	chars, newPatternIdx, err := extractChars(pattern, patternIdx)
	if err != nil {
		return nil, -1, err
	}
	return matcher.NewPosCharacterGroupMatcher(chars), newPatternIdx, nil
}

func getNegCharacterGroupMatcher(pattern []rune, patternIdx int) (matcher.Matcher, int, error) {
	chars, newPatternIdx, err := extractChars(pattern, patternIdx)
	if err != nil {
		return nil, -1, err
	}
	return matcher.NewNegCharacterGroupMatcher(chars), newPatternIdx, nil
}

func extractChars(pattern []rune, patternIdx int) (mapset.Set[rune], int, error) {
	charSet := mapset.NewSet[rune]()
	newPatternIdx := patternIdx
	for newPatternIdx < len(pattern) {
		if pattern[newPatternIdx] == ']' {
			if charSet.IsEmpty() {
				return nil, -1, fmt.Errorf("invalid pattern: character group contains no characters")
			}
			return charSet, newPatternIdx + 1, nil
		}
		charSet.Add(rune(pattern[newPatternIdx]))
		newPatternIdx++
	}
	return nil, -1, fmt.Errorf("invalid pattern: unmatched [")
}
