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

type ExternalFunctionValue struct {
	forcedFunctionType dectype.FunctionTypeLike
	parameters         []*FunctionParameterDefinition
	commentBlock       *ast.MultilineComment

	sourceFileReference token.SourceFileReference
	references          []*FunctionReference
	annotation          *AnnotationStatement
}

func NewExternalFunctionValue(annotation *AnnotationStatement, forcedFunctionType dectype.FunctionTypeLike, parameters []*FunctionParameterDefinition, commentBlock *ast.MultilineComment) *ExternalFunctionValue {
	if len(parameters) != (forcedFunctionType.ParameterCount() - 1) {
		panic("not great. different parameters")
	}

	sourceFileReference := annotation.FetchPositionLength()

	return &ExternalFunctionValue{annotation: annotation, forcedFunctionType: forcedFunctionType, parameters: parameters, commentBlock: commentBlock, sourceFileReference: sourceFileReference}
}

func (f *ExternalFunctionValue) Annotation() *AnnotationStatement {
	return f.annotation
}

func (f *ExternalFunctionValue) Parameters() []*FunctionParameterDefinition {
	return f.parameters
}

func (f *ExternalFunctionValue) ForcedFunctionType() dectype.FunctionTypeLike {
	return f.forcedFunctionType
}

func (f *ExternalFunctionValue) String() string {
	return fmt.Sprintf("[externalfunction (%v) -> %v]", f.parameters)
}

func (f *ExternalFunctionValue) HumanReadable() string {
	return fmt.Sprintf("[externalfunction (%v) -> %v]", f.parameters)
}

func (f *ExternalFunctionValue) Type() dtype.Type {
	return f.annotation.Type()
}

func (f *ExternalFunctionValue) Next() dtype.Type {
	return f.forcedFunctionType
}

func (f *ExternalFunctionValue) Resolve() (dtype.Atom, error) {
	return f.forcedFunctionType.Resolve()
}

func (f *ExternalFunctionValue) FetchPositionLength() token.SourceFileReference {
	return f.annotation.FetchPositionLength()
}

func (f *ExternalFunctionValue) CommentBlock() *ast.MultilineComment {
	return f.commentBlock
}

func (f *ExternalFunctionValue) AddReferee(ref *FunctionReference) {
	f.references = append(f.references, ref)
}

func (f *ExternalFunctionValue) References() []*FunctionReference {
	return f.references
}
