/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/instruction_sp"
)

type CaseConsequencePatternMatchingZeroPointer struct {
	zeroPointer SourceDynamicMemoryPos
	label       *Label
}

func NewCaseConsequencePatternMatchingZeroPointer(zeroPointer SourceDynamicMemoryPos, label *Label) *CaseConsequencePatternMatchingZeroPointer {
	return &CaseConsequencePatternMatchingZeroPointer{zeroPointer: zeroPointer, label: label}
}

func (c *CaseConsequencePatternMatchingZeroPointer) Label() *Label {
	return c.label
}

func (c *CaseConsequencePatternMatchingZeroPointer) SourceDynamicMemoryPos() SourceDynamicMemoryPos {
	return c.zeroPointer
}

func (c *CaseConsequencePatternMatchingZeroPointer) String() string {
	return fmt.Sprintf("[caseconpmzp %v (%d) %v]", c.zeroPointer, c.label)
}

type CasePatternMatchingZeroPointer struct {
	test               SourceStackPos
	consequences       []*CaseConsequencePatternMatchingZeroPointer
	defaultConsequence *CaseConsequencePatternMatchingZeroPointer
	matchingType       instruction_sp.PatternMatchingType
}

func (o *CasePatternMatchingZeroPointer) String() string {
	return fmt.Sprintf("[casepmzp %v (%d) and then jump %v (%v)]", o.test, o.matchingType, o.consequences, o.defaultConsequence)
}
