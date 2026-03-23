package token

import mapset "github.com/deckarep/golang-set/v2"

type Token interface{}

type Literal struct {
	Value rune
}

func NewLiteral(value rune) *Literal {
	return &Literal{Value: value}
}

type Digit struct{}

func NewDigit() *Digit {
	return &Digit{}
}

type WordCharacter struct{}

func NewWordCharacter() *WordCharacter {
	return &WordCharacter{}
}

type PosCharacterGroup struct {
	Chars mapset.Set[rune]
}

func NewPosCharacterGroup(chars mapset.Set[rune]) *PosCharacterGroup {
	return &PosCharacterGroup{Chars: chars}
}

type NegCharacterGroup struct {
	Chars mapset.Set[rune]
}

func NewNegCharacterGroup(chars mapset.Set[rune]) *NegCharacterGroup {
	return &NegCharacterGroup{Chars: chars}
}

type WildCard struct{}

func NewWildcard() *WildCard {
	return &WildCard{}
}

type Optional struct {
	Tokens []Token
}

func NewOptional(tokens []Token) *Optional {
	return &Optional{Tokens: tokens}
}

type OneOrMore struct {
	Tokens []Token
}

func NewOneOrMore(tokens []Token) *OneOrMore {
	return &OneOrMore{Tokens: tokens}
}
