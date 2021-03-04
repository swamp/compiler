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

type Alias struct {
	name             *ast.Alias
	referencedType   dtype.Type
	artifactTypeName ArtifactFullyQualifiedTypeName
}

func (u *Alias) String() string {
	return fmt.Sprintf("[alias %v %v]", u.name.Name(), u.referencedType)
}

func (u *Alias) AstAlias() *ast.Alias {
	return u.name
}

func (u *Alias) HumanReadable() string {
	return fmt.Sprintf("%v", u.name.Name())
}

func (u *Alias) StatementString() string {
	return fmt.Sprintf("%v", u.name.Name())
}

func (u *Alias) TypeIdentifier() *ast.TypeIdentifier {
	return u.name.Identifier()
}

func (u *Alias) FetchPositionLength() token.SourceFileReference {
	return u.name.FetchPositionLength()
}

func (u *Alias) ArtifactTypeName() ArtifactFullyQualifiedTypeName {
	return u.artifactTypeName
}

func (u *Alias) ShortString() string {
	return fmt.Sprintf("[alias %v %v]", u.name.Name(), u.referencedType.ShortString())
}

func (u *Alias) DecoratedName() string {
	return u.name.Name()
}

func (u *Alias) ShortName() string {
	return u.DecoratedName()
}

func (u *Alias) ParameterCount() int {
	return u.referencedType.ParameterCount()
}

func (u *Alias) Resolve() (dtype.Atom, error) {
	return u.referencedType.Resolve()
}

func (u *Alias) Next() dtype.Type {
	return u.referencedType
}

func (u *Alias) Generate(params []dtype.Type) (dtype.Type, error) {
	return u.referencedType.Generate(params)
}

func NewAliasType(name *ast.Alias, artifactTypeName ArtifactFullyQualifiedTypeName,
	referencedType dtype.Type) *Alias {
	return &Alias{name: name, artifactTypeName: artifactTypeName, referencedType: referencedType}
}
