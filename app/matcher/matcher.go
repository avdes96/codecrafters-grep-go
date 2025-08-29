package matcher

import (
	"unicode"

	mapset "github.com/deckarep/golang-set/v2"
)

type Matcher interface {
	Match(char rune) bool
}

type LiteralCharMatcher struct {
	toMatch rune
}

func NewLiteralCharMatcher(toMatch rune) *LiteralCharMatcher {
	return &LiteralCharMatcher{toMatch: toMatch}
}

func (m *LiteralCharMatcher) Match(char rune) bool {
	return char == m.toMatch
}

type DigitMatcher struct{}

func NewDigitMatcher() *DigitMatcher {
	return &DigitMatcher{}
}

func (m *DigitMatcher) Match(char rune) bool {
	return unicode.IsDigit(char)
}

type WordCharacterMatcher struct{}

func NewWordCharacterMatcher() *WordCharacterMatcher {
	return &WordCharacterMatcher{}
}

func (m *WordCharacterMatcher) Match(char rune) bool {
	if char == '_' {
		return true
	}
	if unicode.IsDigit(char) {
		return true
	}
	if unicode.IsUpper(char) {
		return true
	}
	if unicode.IsLower(char) {
		return true
	}
	return false
}

type PosCharacterGroupMatcher struct {
	chars mapset.Set[rune]
}

func NewPosCharacterGroupMatcher(chars mapset.Set[rune]) *PosCharacterGroupMatcher {
	return &PosCharacterGroupMatcher{chars: chars}
}

func (m *PosCharacterGroupMatcher) Match(char rune) bool {
	return m.chars.Contains(char)
}

type NegCharacterGroupMatcher struct {
	chars mapset.Set[rune]
}

func NewNegCharacterGroupMatcher(chars mapset.Set[rune]) *NegCharacterGroupMatcher {
	return &NegCharacterGroupMatcher{chars: chars}
}

func (m *NegCharacterGroupMatcher) Match(char rune) bool {
	return !m.chars.Contains(char)
}

type WildCardMatcher struct{}

func NewWildcardMatcher() *WildCardMatcher {
	return &WildCardMatcher{}
}

func (m *WildCardMatcher) Match(char rune) bool {
	return true
}
