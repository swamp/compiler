/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

/*
type LocalTypeDefinitionContext struct {
	lookup                   map[string]*dectype.LocalTypeDefinition
	typeParameterIdentifiers []*dectype.LocalTypeDefinition
}

func (t *LocalTypeDefinitionContext) arrayToString() string {
	var typeParams []string
	if len(t.typeParameterIdentifiers) == 0 {
		return ""
	}
	for _, v := range t.typeParameterIdentifiers {
		typeParams = append(typeParams, v.Identifier().Name())
	}
	return "[" + strings.Join(typeParams, " ") + "]"
}

func (t *LocalTypeDefinitionContext) String() string {
	return fmt.Sprintf("[type-param-context %s]", t.arrayToString())
}

func (t *LocalTypeDefinitionContext) HasTypeParameter(parameter *dectype.LocalTypeDefinition) bool {
	return t.lookup[parameter.Identifier().Name()] != nil
}

func (t *LocalTypeDefinitionContext) Lookup(parameter *dtype.TypeArgumentName) *dectype.LocalTypeDefinition {
	return t.lookup[parameter.LocalType().Name()]
}

func (t *LocalTypeDefinitionContext) LookupReference(parameter *dectype.LocalTypeDefinitionReference) *dectype.LocalTypeDefinition {
	return t.lookup[parameter.Identifier().Name()]
}

func (t *LocalTypeDefinitionContext) LocalTypeNames() []*dtype.TypeArgumentName {
	var names []*dtype.TypeArgumentName
	for _, x := range t.typeParameterIdentifiers {
		names = append(names, x.Identifier())
	}
	return names
}

func NewLocalTypeDefinitionContext(typeParameterNames []*ast.LocalTypeName) *LocalTypeDefinitionContext {
	lookup := make(map[string]*dectype.LocalTypeDefinition)
	var localTypeDefs []*dectype.LocalTypeDefinition
	for _, typeParameterIdentifier := range typeParameterNames {
		name := dtype.NewTypeArgumentName(typeParameterIdentifier)
		newLocalDef := dectype.NewLocalTypeDefinition(name, dectype.NewAnyType())
		localTypeDefs = append(localTypeDefs, newLocalDef)
		lookup[typeParameterIdentifier.Name()] = newLocalDef
	}
	return &LocalTypeDefinitionContext{lookup: lookup, typeParameterIdentifiers: localTypeDefs}
}


*/
