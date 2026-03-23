package parse

import (
	"github.com/codecrafters-io/grep-starter-go/app/token"
	mapset "github.com/deckarep/golang-set/v2"
)

const unmatchedErrorMessage string = "Unmatched [, [^, [:, [., or [="
const nullChar rune = '\x00'

type Parser struct {
	input      string
	currentIdx int
	endIdx     int
}

func NewParser(pattern string) *Parser {
	return &Parser{
		input:      pattern,
		currentIdx: 0,
		endIdx:     len(pattern) - 1,
	}
}

type Pattern struct {
	Tokens         []token.Token
	MustMatchStart bool
	MustMatchEnd   bool
}

func newPattern() *Pattern {
	return &Pattern{
		Tokens: make([]token.Token, 0),
	}
}

func (p *Pattern) addToken(t token.Token) {
	p.Tokens = append(p.Tokens, t)
}

func (psr *Parser) Parse() *Pattern {
	pattern := newPattern()
	if psr.input[0] == '^' {
		pattern.MustMatchStart = true
		psr.currentIdx = 1
	}
	if psr.input[psr.endIdx] == '$' {
		pattern.MustMatchEnd = true
		psr.endIdx--
	}
	for !psr.atEnd() {
		char := psr.consume()
		var t token.Token
		switch char {
		case '.':
			t = token.NewWildcard()
		case '\\':
			t = psr.parseMetaCharacter()
		case '[':
			t = psr.parseCharacterGroup()
		default:
			t = token.NewLiteral(char)
		}
		pattern.addToken(psr.parseRepetition(t))
	}
	return pattern
}

func (psr *Parser) parseMetaCharacter() token.Token {
	char := psr.consume()
	switch char {
	case nullChar:
		PrintErrorMessage("Trailing backslash")
		return nil
	case 'w':
		return token.NewWordCharacter()
	case 'd':
		return token.NewDigit()
	default:
		return token.NewLiteral(char)
	}
}

func (psr *Parser) parseCharacterGroup() token.Token {
	if psr.atEnd() || psr.peek() == ']' {
		PrintErrorMessage(unmatchedErrorMessage)
		return nil
	}
	if psr.peek() == '^' {
		psr.consume()
		return psr.parseNegCharacterGroup()
	}
	return psr.parsePosCharacterGroup()
}

func (psr *Parser) parsePosCharacterGroup() token.Token {
	return token.NewPosCharacterGroup(psr.extractChars())
}

func (psr *Parser) parseNegCharacterGroup() token.Token {
	return token.NewNegCharacterGroup(psr.extractChars())
}

func (psr *Parser) extractChars() mapset.Set[rune] {
	charSet := mapset.NewSet[rune]()
	for !psr.atEnd() {
		char := psr.consume()
		if char == ']' {
			return charSet
		}
		charSet.Add(char)
	}
	PrintErrorMessage(unmatchedErrorMessage)
	return nil
}

func (psr *Parser) parseRepetition(t token.Token) token.Token {
	if psr.peek() == '+' {
		psr.consume()
		return token.NewOneOrMore([]token.Token{t})
	}
	if psr.peek() == '?' {
		psr.consume()
		return token.NewOptional([]token.Token{t})
	}
	return t
}

func (psr *Parser) consume() rune {
	if psr.currentIdx > psr.endIdx {
		return nullChar
	}
	char := psr.input[psr.currentIdx]
	psr.currentIdx++
	return rune(char)
}

func (psr *Parser) peek() rune {
	if psr.atEnd() {
		return nullChar
	}
	return rune(psr.input[psr.currentIdx])
}

func (psr *Parser) atEnd() bool {
	return psr.currentIdx > psr.endIdx
}
