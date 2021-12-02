/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
)

type ImportStatementReference struct {
	importStatement *ImportStatement
}

// source *decorated.Module, sourceMountedModuleName dectype.PackageRelativeModuleName, exposeAll bool

func NewImportStatementReference(importStatement *ImportStatement) *ImportStatementReference {
	ref := &ImportStatementReference{importStatement: importStatement}
	importStatement.AddReference(ref)
	return ref
}

func (l *ImportStatementReference) String() string {
	return fmt.Sprintf("[importref %v]", l.importStatement)
}

func (l *ImportStatementReference) StatementString() string {
	return fmt.Sprintf("[import %v]", l.importStatement)
}

func (l *ImportStatementReference) ImportStatement() *ImportStatement {
	return l.importStatement
}
