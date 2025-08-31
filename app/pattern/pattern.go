package pattern

import (
	"fmt"

	"github.com/codecrafters-io/grep-starter-go/app/matcher"
	mapset "github.com/deckarep/golang-set/v2"
)

type Node struct {
	Matcher   matcher.Matcher
	Optional  bool
	OneOrMore bool
	Next      *Node
}

func (n *Node) IsEmpty() bool {
	return n.Matcher == nil
}

func (n *Node) Match(c rune) bool {
	return n.Matcher.Match(c)
}

func (n *Node) IsOptional() bool {
	return n.Optional
}

func (n *Node) IsOneOrMore() bool {
	return n.OneOrMore
}

type Pattern struct {
	Head           *Node
	MustMatchStart bool
	MustMatchEnd   bool
}

func (p *Pattern) IsEmpty() bool {
	return p.Head.IsEmpty()
}

func Parse(pattern []rune) (*Pattern, error) {
	head := Node{}
	current := &head
	p := Pattern{
		Head: &head,
	}
	if len(pattern) == 0 {
		return &p, nil
	}
	idx := 0
	if pattern[idx] == '^' {
		p.MustMatchStart = true
		idx = 1
	}
	if pattern[len(pattern)-1] == '$' {
		p.MustMatchEnd = true
		pattern = pattern[:len(pattern)-1]
	}
	var err error
	for idx < len(pattern) {
		idx, err = addMatcher(pattern, idx, current)
		if err != nil {
			return nil, err
		}
		current.Next = &Node{}
		current = current.Next
	}
	return &p, nil
}

func addMatcher(pattern []rune, patternIdx int, Node *Node) (int, error) {
	m, newIdx, err := getBaseMatcher(pattern, patternIdx)
	if err != nil {
		return -1, err
	}
	Node.Matcher = m
	if newIdx >= len(pattern) {
		return newIdx, nil
	}
	switch pattern[newIdx] {
	case '?':
		Node.Optional = true
		newIdx++
	case '+':
		Node.OneOrMore = true
		newIdx++
	}
	return newIdx, nil
}

func getBaseMatcher(pattern []rune, patternIdx int) (matcher.Matcher, int, error) {
	switch pattern[patternIdx] {
	case '+':
		return nil, -1, fmt.Errorf("no previous char with +")
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
	case '.':
		return matcher.NewWildcardMatcher(), patternIdx + 1, nil
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
