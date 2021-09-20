/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type CaseConsequence struct {
	caseValue uint8
	label     *Label
}

func NewCaseConsequence(caseValue uint8, label *Label) *CaseConsequence {
	return &CaseConsequence{caseValue: caseValue, label: label}
}

func (c *CaseConsequence) Label() *Label {
	return c.label
}

func (c *CaseConsequence) InternalEnumIndex() uint8 {
	return c.caseValue
}

func (c *CaseConsequence) String() string {
	return fmt.Sprintf("[casecon %v %v]", c.caseValue, c.label)
}

type Case struct {
	test               SourceStackPos
	consequences       []*CaseConsequence
	defaultConsequence *CaseConsequence
}

func (o *Case) String() string {
	return fmt.Sprintf("[case %v and then jump %v (%v)]", o.test, o.consequences, o.defaultConsequence)
}
