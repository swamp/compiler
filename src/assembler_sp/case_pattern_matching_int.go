/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"
)

type CaseConsequencePatternMatchingInt struct {
	constantInteger int32
	label           *Label
}

func NewCaseConsequencePatternMatchingInt(constantInteger int32, label *Label) *CaseConsequencePatternMatchingInt {
	return &CaseConsequencePatternMatchingInt{constantInteger: constantInteger, label: label}
}

func (c *CaseConsequencePatternMatchingInt) Label() *Label {
	return c.label
}

func (c *CaseConsequencePatternMatchingInt) ConstantInteger() int32 {
	return c.constantInteger
}

func (c *CaseConsequencePatternMatchingInt) String() string {
	return fmt.Sprintf("[caseconpmi %v (%d)]", c.constantInteger, c.label)
}

type CasePatternMatchingInt struct {
	test               SourceStackPos
	consequences       []*CaseConsequencePatternMatchingInt
	defaultConsequence *Label
}

func (o *CasePatternMatchingInt) String() string {
	return fmt.Sprintf("[casepmi %v and then jump %v (%v)]", o.test, o.consequences, o.defaultConsequence)
}
