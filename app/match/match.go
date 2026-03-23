package match

import (
	"unicode"

	"github.com/codecrafters-io/grep-starter-go/app/parse"
	"github.com/codecrafters-io/grep-starter-go/app/token"
)

type Matcher struct {
	str            []rune
	tokens         []token.Token
	mustMatchStart bool
	mustMatchEnd   bool
}

func NewMatcherFromPattern(str []rune, p *parse.Pattern) *Matcher {
	return &Matcher{
		str:            str,
		tokens:         p.Tokens,
		mustMatchStart: p.MustMatchStart,
		mustMatchEnd:   p.MustMatchEnd,
	}
}

func (m *Matcher) newSubMatcher(newTokens []token.Token) *Matcher {
	return &Matcher{
		str:    m.str,
		tokens: newTokens,
	}
}

func (m *Matcher) Match() bool {
	if m.mustMatchStart {
		return m.match(0, 0) != -1
	}
	for idx := range m.str {
		if m.match(idx, 0) != -1 {
			return true
		}
	}
	return false
}

func (m *Matcher) match(strIdx, patternIdx int) int {
	if strIdx >= len(m.str) && patternIdx < len(m.tokens) {
		return -1
	}
	if patternIdx >= len(m.tokens) {
		if m.matchPatternEnd(strIdx) {
			return strIdx
		}
		return -1
	}
	currentChar := rune(m.str[strIdx])
	switch t := m.tokens[patternIdx].(type) {
	case *token.OneOrMore:
		return m.matchOneOrMore(strIdx, patternIdx+1, t)
	case *token.Optional:
		return m.matchOptional(strIdx, patternIdx+1, t)
	default:
		if !m.matchSingular(t, currentChar) {
			return -1
		}
		return m.match(strIdx+1, patternIdx+1)
	}

}

func (m *Matcher) matchOptional(strIdx, patternIdx int, t *token.Optional) int {
	subMatcher := m.newSubMatcher(t.Tokens)
	if idx := subMatcher.match(strIdx, 0); idx != -1 {
		return m.match(idx, patternIdx)
	}
	return m.match(strIdx, patternIdx)
}

func (m *Matcher) matchOneOrMore(strIdx, patternIdx int, t *token.OneOrMore) int {
	subMatcher := m.newSubMatcher(t.Tokens)
	idx := subMatcher.match(strIdx, 0)
	if idx == -1 {
		return -1
	}
	return subMatcher.matchGreedy(idx, m, patternIdx)
}

func (m *Matcher) matchGreedy(strIdx int, parentMatcher *Matcher, parentPatternIdx int) int {
	if idx := m.match(strIdx, 0); idx != -1 {
		if idx := m.matchGreedy(idx, parentMatcher, parentPatternIdx); idx != -1 {
			return idx
		}
	}
	return parentMatcher.match(strIdx, parentPatternIdx)
}

func (m *Matcher) matchSingular(currentToken token.Token, char rune) bool {
	switch t := currentToken.(type) {
	case *token.Literal:
		return m.matchLiteral(t, char)
	case *token.Digit:
		return m.matchDigit(char)
	case *token.WordCharacter:
		return m.matchWordCharacter(char)
	case *token.PosCharacterGroup:
		return m.matchPosCharacterGroup(t, char)
	case *token.NegCharacterGroup:
		return m.matchNegCharacterGroup(t, char)
	case *token.WildCard:
		return m.matchWildcard(t, char)
	}
	return false
}

func (m *Matcher) matchLiteral(t *token.Literal, char rune) bool {
	return t.Value == char
}

func (m *Matcher) matchDigit(char rune) bool {
	return unicode.IsDigit(char)
}

func (m *Matcher) matchWordCharacter(char rune) bool {
	return char == '_' || unicode.IsDigit(char) || unicode.IsUpper(char) || unicode.IsLower(char)
}

func (m *Matcher) matchPosCharacterGroup(t *token.PosCharacterGroup, char rune) bool {
	return t.Chars.Contains(char)
}

func (m *Matcher) matchNegCharacterGroup(t *token.NegCharacterGroup, char rune) bool {
	return !t.Chars.Contains(char)
}

func (m *Matcher) matchWildcard(_ *token.WildCard, _ rune) bool {
	return true
}

func (m *Matcher) matchPatternEnd(strIdx int) bool {
	if m.mustMatchEnd {
		return strIdx >= len(m.str)
	}
	return true
}
