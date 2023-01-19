/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package semantic

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

type SemanticNode struct {
	encodedIntegers [5]uint
	tokenType       string
	tokenRange      token.Range
	debug           interface{}
}

type SemanticBuilder struct {
	tokenTypes      []string
	tokenModifiers  []string
	lastRange       token.Range
	lastDebug       string
	encodedIntegers []uint
	nodes           []SemanticNode
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

func (s *SemanticBuilder) EncodeSymbol(tokenRange token.Range, tokenType string, modifiers []string, node fmt.Stringer) error {
	if !tokenRange.IsAfter(s.lastRange) {
		/*
			for _, existingNode := range s.nodes {
				log.Printf("  node: %v : %v (%T\n%v)", existingNode.tokenRange, existingNode.tokenType, existingNode.debug, existingNode.debug)
			}
			log.Printf("--> added node: %v : %v (%T\n%v)", tokenRange, tokenType, node, node)
		*/
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

	s.nodes = append(s.nodes, SemanticNode{
		encodedIntegers: encodedIntegers,
		tokenType:       tokenType,
		tokenRange:      tokenRange,
		debug:           node,
	})
	s.encodedIntegers = append(s.encodedIntegers, encodedIntegers[:]...)
	s.lastRange = tokenRange

	return nil
}
