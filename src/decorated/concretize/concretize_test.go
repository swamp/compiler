package concretize_test

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/concretize"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
	"log"
	"testing"
)

func MakeFakeSourceFileReference() token.SourceFileReference {
	return token.NewInternalSourceFileReference()
}

func MakeFakeTypeIdentifier(name string) *ast.TypeIdentifier {
	return ast.NewTypeIdentifier(token.NewTypeSymbolToken(name, MakeFakeSourceFileReference(), -1))
}

func MakeFakeVariable(name string) *ast.VariableIdentifier {
	return ast.NewVariableIdentifier(token.NewVariableSymbolToken(name, MakeFakeSourceFileReference(), -1))
}

func MakeFakeLocalTypeName(name string) *ast.LocalTypeName {
	return ast.NewLocalTypeName(MakeFakeVariable(name))
}

func MakeFakeListA() *dectype.PrimitiveAtom {
	listIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("List", token.NewInternalSourceFileReference(), 0))

	localTypeVariable := MakeFakeLocalTypeName("a")
	astNameDef := ast.NewLocalTypeNameDefinition(localTypeVariable)
	localType := ast.NewLocalTypeNameReference(astNameDef)

	decLocalTypeName := dtype.NewLocalTypeName(localTypeVariable)
	decDef := dectype.NewLocalTypeNameDefinition(decLocalTypeName)

	decNameRef := dectype.NewLocalTypeNameReference(localType, decDef)
	listType := dectype.NewPrimitiveType(listIdentifier, []dtype.Type{decNameRef})

	return listType
}

func MakeFakeListOf(typeParam dtype.Type) *dectype.PrimitiveAtom {
	listIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("List", token.NewInternalSourceFileReference(), 0))
	listType := dectype.NewPrimitiveType(listIdentifier, []dtype.Type{typeParam})
	return listType
}

func TestCustomTypeVariant(t *testing.T) {
	typeIdentifierForVariant := MakeFakeTypeIdentifier("SomeVariant")

	generic := MakeFakeLocalTypeName("a")
	localTypeNames := []*ast.LocalTypeName{generic}
	nameOnlyContext := ast.NewTypeParameterIdentifierContext(localTypeNames, nil)
	firstNameDefRef, err := nameOnlyContext.ParseReferenceFromName(localTypeNames[0])
	if err != nil {
		t.Error(err)
	}

	generics := []ast.Type{firstNameDefRef}
	astVariant := ast.NewCustomTypeVariant(42, typeIdentifierForVariant, generics, nil)

	decoratedTypeNames := dectype.NewLocalTypeNameContext()
	aName := MakeFakeLocalTypeName("a")
	argName := dtype.NewLocalTypeName(aName)
	aNameDef := decoratedTypeNames.AddDef(argName)

	aNameRef := dectype.NewLocalTypeNameReference(firstNameDefRef, aNameDef)

	listA := MakeFakeListA()
	decoratedGenericsForVariant := []dtype.Type{listA, aNameRef}
	reference := dectype.NewCustomTypeVariant(42, nil, astVariant, decoratedGenericsForVariant)

	integerType := dectype.NewPrimitiveType(MakeFakeTypeIdentifier("Int"), nil)
	listOfInt := MakeFakeListOf(integerType)

	concreteArguments := []dtype.Type{listOfInt, integerType}

	resolver := dectype.NewTypeParameterContext()
	resolver.AddExpectedDef(argName)

	decoratedTypeNames.SetType(reference)

	concretizedVariant, concreteErr := concretize.ConcreteArguments(decoratedTypeNames, concreteArguments)
	if concreteErr != nil {
		t.Error(concreteErr)
	}

	log.Printf("reference variant: %v", reference)
	log.Printf("concretized variant: %v", concretizedVariant)
}
