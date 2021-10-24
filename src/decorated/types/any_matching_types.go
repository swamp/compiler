/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

func HasAnyMatchingTypes(types []dtype.Type) (bool, int) {
	for index, param := range types {
		_, isAnyMatching := param.(*AnyMatchingTypes)
		if isAnyMatching {
			return isAnyMatching, index
		}
	}

	return false, -1
}

type AnyMatchingTypes struct {
	astAnyMatchingType *ast.AnyMatchingType
}

func (u *AnyMatchingTypes) String() string {
	return fmt.Sprintf("[anymatching types %v]", u.astAnyMatchingType.Name())
}

func (u *AnyMatchingTypes) FetchPositionLength() token.SourceFileReference {
	return u.astAnyMatchingType.FetchPositionLength()
}

func (u *AnyMatchingTypes) HumanReadable() string {
	return fmt.Sprintf("%v", u.astAnyMatchingType.Name())
}

func (u *AnyMatchingTypes) Identifier() *ast.AnyMatchingType {
	return u.astAnyMatchingType
}

func (u *AnyMatchingTypes) AtomName() string {
	return u.astAnyMatchingType.Name()
}

func (u *AnyMatchingTypes) IsEqual(_ dtype.Atom) error {
	return nil
}

func (u *AnyMatchingTypes) ParameterCount() int {
	return 0
}

func (u *AnyMatchingTypes) Resolve() (dtype.Atom, error) {
	return u, nil
}

func (u *AnyMatchingTypes) Next() dtype.Type {
	return nil
}

func NewAnyMatchingTypes(identifier *ast.AnyMatchingType) *AnyMatchingTypes {
	return &AnyMatchingTypes{astAnyMatchingType: identifier}
}
