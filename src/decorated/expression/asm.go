/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type AsmConstant struct {
	asm            *ast.Asm
	doNotCheckType dtype.Type
}

func NewAsmConstant(asm *ast.Asm) *AsmConstant {
	return &AsmConstant{asm: asm, doNotCheckType: dectype.NewAnyType()}
}

func (i *AsmConstant) Type() dtype.Type {
	return i.doNotCheckType
}

func (i *AsmConstant) Asm() *ast.Asm {
	return i.asm
}

func (i *AsmConstant) String() string {
	return fmt.Sprintf("[asm %v]", i.asm.Asm())
}

func (i *AsmConstant) FetchPositionAndLength() token.PositionLength {
	return i.asm.PositionLength()
}
