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

func MakeFakeAstTypeReference(name string) *ast.TypeReference {
	return ast.NewTypeReference(MakeFakeTypeIdentifier(name), nil)
}

func MakeFakeAstTypeReferenceWithArguments(name string, arguments []string) *ast.TypeReference {
	var astTypes []ast.Type
	for _, argumentTypeName := range arguments {
		astTypes = append(astTypes, MakeFakeAstTypeReference(argumentTypeName))
	}
	return ast.NewTypeReference(MakeFakeTypeIdentifier(name), astTypes)
}

func MakeFakeAstTypeReferenceWithLocalTypeNames(name string, arguments []string) *ast.TypeReference {
	var astTypes []ast.Type
	for _, argumentTypeName := range arguments {
		astTypes = append(astTypes, MakeFakeLocalTypeName(argumentTypeName))
	}
	return ast.NewTypeReference(MakeFakeTypeIdentifier(name), astTypes)
}

func MakeFakeVariable(name string) *ast.VariableIdentifier {
	return ast.NewVariableIdentifier(token.NewVariableSymbolToken(name, MakeFakeSourceFileReference(), -1))
}

func MakeFakeLocalTypeName(name string) *ast.LocalTypeName {
	return ast.NewLocalTypeName(MakeFakeVariable(name))
}

func MakeFakeList(localTypeName string) *dectype.PrimitiveAtom {
	listIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("List", token.NewInternalSourceFileReference(), 0))

	localTypeVariable := MakeFakeLocalTypeName(localTypeName)
	astNameDef := ast.NewLocalTypeNameDefinition(localTypeVariable)
	localType := ast.NewLocalTypeNameReference(localTypeVariable, astNameDef)

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

func MakeFakeNameContext(names []string) *ast.LocalTypeNameDefinitionContext {
	var generics []*ast.LocalTypeName
	for _, name := range names {
		generic := MakeFakeLocalTypeName(name)
		generics = append(generics, generic)
	}

	return ast.NewLocalTypeNameContext(generics, nil)
}

func MakeIntegerType() *dectype.PrimitiveAtom {
	return dectype.NewPrimitiveType(MakeFakeTypeIdentifier("Int"), nil)
}

func MakeStringType() *dectype.PrimitiveAtom {
	return dectype.NewPrimitiveType(MakeFakeTypeIdentifier("String"), nil)
}

func TestCustomTypeVariant(t *testing.T) {
	nameOnlyContext := MakeFakeNameContext([]string{"a"})

	firstNameDefRef, err := nameOnlyContext.ParseReferenceFromName(nameOnlyContext.LocalTypeNames()[0])
	if err != nil {
		t.Error(err)
	}
	generics := []ast.Type{firstNameDefRef}
	typeIdentifierForVariant := MakeFakeTypeIdentifier("SomeVariant")
	astVariant := ast.NewCustomTypeVariant(42, typeIdentifierForVariant, generics, nil)

	decoratedTypeNames := dectype.NewLocalTypeNameContext()
	aName := MakeFakeLocalTypeName("a")
	argName := dtype.NewLocalTypeName(aName)
	aNameDef := decoratedTypeNames.AddDef(argName)

	aNameRef := dectype.NewLocalTypeNameReference(firstNameDefRef, aNameDef)

	listA := MakeFakeList("a")
	decoratedGenericsForVariant := []dtype.Type{listA, aNameRef}
	reference := dectype.NewCustomTypeVariant(42, nil, astVariant, decoratedGenericsForVariant)

	integerType := MakeIntegerType()
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

func TestFunction(t *testing.T) {
	nameOnlyContext := MakeFakeNameContext([]string{"a", "b"})

	firstNameDefRef, err := nameOnlyContext.ParseReferenceFromName(nameOnlyContext.LocalTypeNames()[0])
	if err != nil {
		t.Error(err)
	}

	secondNameDefRef, err2 := nameOnlyContext.ParseReferenceFromName(nameOnlyContext.LocalTypeNames()[1])
	if err2 != nil {
		t.Error(err2)
	}

	astListB := MakeFakeAstTypeReferenceWithLocalTypeNames("List", []string{"b"})
	astIntReference := MakeFakeAstTypeReference("Int")
	astFunctionParameters := []ast.Type{astIntReference, astListB, firstNameDefRef, secondNameDefRef}
	astFunction := ast.NewFunctionType(astFunctionParameters)

	decoratedTypeNames := dectype.NewLocalTypeNameContext()
	aName := MakeFakeLocalTypeName("a")
	aArgName := dtype.NewLocalTypeName(aName)
	aNameDef := decoratedTypeNames.AddDef(aArgName)
	aNameRef := dectype.NewLocalTypeNameReference(firstNameDefRef, aNameDef)

	bName := MakeFakeLocalTypeName("b")
	bArgName := dtype.NewLocalTypeName(bName)
	bNameDef := decoratedTypeNames.AddDef(bArgName)
	bNameRef := dectype.NewLocalTypeNameReference(secondNameDefRef, bNameDef)

	integerType := MakeIntegerType()
	listOfB := MakeFakeListOf(bNameRef)
	decoratedFunctionParameters := []dtype.Type{integerType, listOfB, aNameRef, bNameRef}
	functionWithLocalTypeNames := dectype.NewFunctionAtom(astFunction, decoratedFunctionParameters)

	stringType := MakeStringType()
	listOfString := MakeFakeListOf(stringType)
	concreteArguments := []dtype.Type{integerType, listOfString, stringType, dectype.NewAnyType()}

	resolver := dectype.NewTypeParameterContext()
	resolver.AddExpectedDefs([]*dtype.LocalTypeName{aArgName, bArgName})

	decoratedTypeNames.SetType(functionWithLocalTypeNames)

	concretizedVariant, concreteErr := concretize.ConcreteArguments(decoratedTypeNames, concreteArguments)
	if concreteErr != nil {
		t.Error(concreteErr)
	}

	log.Printf("functionWithLocalTypeNames variant: %v", functionWithLocalTypeNames)
	log.Printf("concretized function: %v", concretizedVariant)
}
