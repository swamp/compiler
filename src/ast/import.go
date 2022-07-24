/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type Import struct {
	moduleReference     *ModuleReference
	optionalAlias       *TypeIdentifier
	typesToExpose       []*TypeIdentifier
	definitionsToExpose []*VariableIdentifier
	exposeAll           bool
	keywordImport       token.Keyword
	keywordExposing     *token.Keyword
	keywordAs           *token.Keyword
	precedingComments   *MultilineComment
	inclusive           token.SourceFileReference
}

func NewImport(keywordImport token.Keyword, keywordAs *token.Keyword, keywordExposing *token.Keyword,
	moduleReference *ModuleReference,
	optionalAlias *TypeIdentifier, typesToExpose []*TypeIdentifier,
	definitionsToExpose []*VariableIdentifier,
	exposeAll bool, precedingComments *MultilineComment) *Import {
	lastSourceRef := moduleReference.Parts()[len(moduleReference.Parts())-1].FetchPositionLength()
	if optionalAlias != nil {
		lastSourceRef = optionalAlias.FetchPositionLength()
	}
	if definitionsToExpose != nil {
		lastSourceRef = definitionsToExpose[len(definitionsToExpose)-1].FetchPositionLength()
	}
	inclusive := token.MakeInclusiveSourceFileReference(keywordImport.FetchPositionLength(), lastSourceRef)
	return &Import{
		keywordImport:   keywordImport,
		keywordExposing: keywordExposing,
		keywordAs:       keywordAs,
		moduleReference: moduleReference, optionalAlias: optionalAlias,
		exposeAll:           exposeAll,
		typesToExpose:       typesToExpose,
		definitionsToExpose: definitionsToExpose,
		precedingComments:   precedingComments,
		inclusive:           inclusive,
	}
}

func (i *Import) KeywordImport() token.Keyword {
	return i.keywordImport
}

func (i *Import) KeywordExposing() *token.Keyword {
	return i.keywordExposing
}

func (i *Import) KeywordAs() *token.Keyword {
	return i.keywordAs
}

func (i *Import) ExposeAll() bool {
	return i.exposeAll
}

func (i *Import) Alias() *TypeIdentifier {
	return i.optionalAlias
}

func (i *Import) ModuleName() *ModuleReference {
	return i.moduleReference
}

func (i *Import) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *Import) String() string {
	s := fmt.Sprintf("[Import %v", i.moduleReference)
	if i.optionalAlias != nil {
		s += fmt.Sprintf(" as %v", i.optionalAlias)
	}
	if len(i.typesToExpose) > 0 || len(i.definitionsToExpose) > 0 {
		s += fmt.Sprintf(" exposing (%v %v)", i.typesToExpose, i.definitionsToExpose)
	} else if i.exposeAll {
		s += " exposing (..)"
	}
	s += "]"

	return s
}

func (i *Import) DebugString() string {
	return "[Import]"
}
