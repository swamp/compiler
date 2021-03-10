/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

/*
type CustomTypeVariantParameterExpand struct {
	identifier *ast.VariableIdentifier
	variant    *dectype.CustomTypeVariant
	index      int
	references []*CustomTypeVariantParameterExpandReference
	foundType  dtype.Type
}

func NewCustomTypeVariantParameterExpand(identifier *ast.VariableIdentifier, foundType dtype.Type, variant *dectype.CustomTypeVariant, index int) *CustomTypeVariantParameterExpand {
	return &CustomTypeVariantParameterExpand{identifier: identifier, variant: variant, index: index, foundType: foundType}
}

func (a *CustomTypeVariantParameterExpand) Identifier() *ast.VariableIdentifier {
	return a.identifier
}

func (a *CustomTypeVariantParameterExpand) Type() dtype.Type {
	return a.foundType
}

func (a *CustomTypeVariantParameterExpand) String() string {
	return fmt.Sprintf("[customtypevariantparam %v = %v]", a.identifier, a.variant, a.index)
}

func (a *CustomTypeVariantParameterExpand) HumanReadable() string {
	return "Function Parameter"
}

func (a *CustomTypeVariantParameterExpand) FetchPositionLength() token.SourceFileReference {
	return a.identifier.Symbol().SourceFileReference
}

func (a *CustomTypeVariantParameterExpand) AddReferee(ref *CustomTypeVariantParameterExpandReference) {
	a.references = append(a.references, ref)
}

func (a *CustomTypeVariantParameterExpand) References() []*CustomTypeVariantParameterExpandReference {
	return a.references
}
*/
