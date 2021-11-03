package lspservice

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/token"
)

/*
   "tokenTypes": [
       "namespace",
       "type",
       "class",
       "enum",
       "interface",
       "struct",
       "typeParameter",
       "parameter",
       "variable",
       "property",
       "enumMember",
       "event",
       "function",
       "method",
       "macro",
       "keyword",
       "modifier",
       "comment",
       "string",
       "number",
       "regexp",
       "operator"
   ],
   "tokenModifiers": [

*/
type SemanticBuilder struct {
	tokenTypes      []string
	tokenModifiers  []string
	lastRange       token.Range
	lastDebug       string
	encodedIntegers []uint
}

func NewSemanticBuilder() *SemanticBuilder {
	self := &SemanticBuilder{
		tokenTypes: []string{
			"namespace",
			"type",
			"class",
			"enum",
			"interface",
			"struct",
			"typeParameter",
			"parameter",
			"variable",
			"property",
			"enumMember",
			"event",
			"function",
			"method",
			"macro",
			"keyword",
			"modifier",
			"comment",
			"string",
			"number",
			"regexp",
			"operator",
		},
		tokenModifiers: []string{
			"declaration",
			"definition",
			"readonly",
			"static",
			"deprecated",
			"abstract",
			"async",
			"modification",
			"documentation",
			"defaultLibrary",
		},
		lastRange: token.NewPositionLength(token.MakePosition(0, 0, 0), 0),
	}
	return self
}

func FindInStrings(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}

	return -1
}

func (s *SemanticBuilder) EncodedValues() []uint {
	return s.encodedIntegers
}

func (s *SemanticBuilder) EncodeSymbol(tokenRange token.Range, tokenType string, modifiers []string) error {
	log.Printf("encoding symbol %v %v %v", tokenRange, tokenType, modifiers)

	if !tokenRange.IsAfter(s.lastRange) {
		return fmt.Errorf("they must be in order! %v to %v and \n%v", s.lastRange, tokenRange, s.lastDebug)
	}
	// log.Printf("adding symbol %v '%v'\n", tokenRange, debugString)

	tokenTypeId := FindInStrings(s.tokenTypes, tokenType)
	if tokenTypeId < 0 {
		return fmt.Errorf("unknown token type %v", tokenType)
	}

	var modifierBitMask uint
	for _, modifier := range modifiers {
		modifierId := FindInStrings(s.tokenModifiers, modifier)
		if modifierId < 0 {
			return fmt.Errorf("unknown token type %v", tokenType)
		}
		modifierBitMask |= 1 << modifierId
	}

	lastLine := s.lastRange.Position().Line()
	lastStartColumn := s.lastRange.Position().Column()

	deltaLine := uint(tokenRange.Position().Line() - lastLine)
	if deltaLine != 0 {
		lastStartColumn = 0
	}
	deltaColumnFromLastStartColumn := uint(tokenRange.Position().Column() - lastStartColumn)

	tokenLength := tokenRange.OctetCount()
	if tokenLength <= 0 {
		log.Printf(fmt.Errorf("problem with token length %v", tokenType).Error())
		return nil
	}

	encodedIntegers := [5]uint{deltaLine, deltaColumnFromLastStartColumn, uint(tokenLength), uint(tokenTypeId), modifierBitMask}

	s.encodedIntegers = append(s.encodedIntegers, encodedIntegers[:]...)
	s.lastRange = tokenRange

	return nil
}
