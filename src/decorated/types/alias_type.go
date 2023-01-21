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
	references       []*AliasReference
	inclusive        token.SourceFileReference
	wasReferenced    bool
}

func (u *Alias) String() string {
	return fmt.Sprintf("[Alias %v %v]", u.name.Name(), u.referencedType)
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
	return u.inclusive
}

func (u *Alias) ArtifactTypeName() ArtifactFullyQualifiedTypeName {
	return u.artifactTypeName
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

func (c *Alias) AddReferee(ref *AliasReference) {
	c.references = append(c.references, ref)
}

func (c *Alias) References() []*AliasReference {
	return c.references
}

func (c *Alias) WasReferenced() bool {
	return c.wasReferenced || len(c.references) > 0
}

func (c *Alias) MarkAsReferenced() {
	c.wasReferenced = true
}

func NewAliasType(name *ast.Alias, artifactTypeName ArtifactFullyQualifiedTypeName,
	referencedType dtype.Type) *Alias {
	_, wasAlias := referencedType.(*Alias)
	if wasAlias {
		panic("we cant have alias inside alias")
	}
	if artifactTypeName.String() == "" {
		panic("must have complete alias name")
	}

	inclusive := token.MakeInclusiveSourceFileReference(name.FetchPositionLength(), referencedType.FetchPositionLength())

	return &Alias{name: name, artifactTypeName: artifactTypeName, referencedType: referencedType, inclusive: inclusive}
}
