package concretize_test

import (
	"log"
	"testing"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/concretize"
	"github.com/swamp/compiler/src/decorated/debug"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func MakeFakeSourceFileReference() token.SourceFileReference {
	return token.NewInternalSourceFileReferenceRow(1)
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

func MakeFakeLocalTypeNameRef(name string) *ast.LocalTypeNameReference {
	aName := MakeFakeLocalTypeName(name)
	aNameDef := ast.NewLocalTypeNameDefinition(aName)
	return ast.NewLocalTypeNameReference(aName, aNameDef)
}

func MakeFakeList(localTypeName string) *dectype.PrimitiveAtom {
	listIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("List", token.NewInternalSourceFileReference(), 0))

	localTypeVariable := MakeFakeLocalTypeName(localTypeName)
	astNameDef := ast.NewLocalTypeNameDefinition(localTypeVariable)
	localType := ast.NewLocalTypeNameReference(localTypeVariable, astNameDef)

	decLocalTypeName := dtype.NewLocalTypeName(localTypeVariable)
	decDef := dectype.NewLocalTypeName(decLocalTypeName)

	decNameRef := dectype.NewLocalTypeNameReference(localType, decDef)
	listType := dectype.NewPrimitiveType(listIdentifier, []dtype.Type{decNameRef})

	return listType
}

func MakeFakeListOf(typeParam dtype.Type) *dectype.PrimitiveAtom {
	listIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("List", token.NewInternalSourceFileReferenceRow(1),
		0))
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

func MakeBoolType() *dectype.PrimitiveAtom {
	return dectype.NewPrimitiveType(MakeFakeTypeIdentifier("Bool"), nil)
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

	aNameRef := MakeFakeLocalTypeNameRef("a")
	localTypeNameOnlyContext := dectype.NewLocalTypeNameOnlyContext([]*ast.LocalTypeName{aNameRef.LocalTypeName()})
	aNameDecRef, _ := localTypeNameOnlyContext.LookupNameReference(aNameRef)

	listA := MakeFakeList("a")
	decoratedGenericsForVariant := []dtype.Type{listA, aNameDecRef}
	reference := dectype.NewCustomTypeVariant(42, nil, astVariant, decoratedGenericsForVariant)

	integerType := MakeIntegerType()
	listOfInt := MakeFakeListOf(integerType)

	concreteArguments := []dtype.Type{listOfInt}

	localTypeNameOnlyContext.SetType(reference)

	localTypeNameOnlyContextRef := dectype.NewLocalTypeNameContextReference(nil, localTypeNameOnlyContext)

	concretizedVariant, concreteErr := concretize.ConcretizeLocalTypeContextUsingArguments(
		localTypeNameOnlyContextRef, concreteArguments,
	)
	if concreteErr != nil {
		t.Error(concreteErr)
	}

	log.Printf("reference variant: %v", reference)
	log.Printf("concretized variant: %v", concretizedVariant)
}

func TestFailingMap(t *testing.T) {
	// -----------------------------------------------------------------------------------------------------------------
	/*
	   drawSprite : (Int) -> Bool =
	   	true


	   drawSprites : (sprites: List Int) -> List Bool =
	   	List.map drawSprite sprites
	*/

	//astIntReference := MakeFakeAstTypeReference("Int")
	astIntType := MakeFakeTypeIdentifier("Int")
	decIntType := dectype.NewPrimitiveType(astIntType, nil)

	//astBoolReference := MakeFakeAstTypeReference("Bool")

	astBoolType := MakeFakeTypeIdentifier("Bool")
	decBoolType := dectype.NewPrimitiveType(astBoolType, nil)

	//	astFunctionParameters := []ast.Type{astIntReference, astBoolReference}
	//	astFunctionType := ast.NewFunctionType(astFunctionParameters)
	//	decParameters := []dtype.Type{decIntType, decBoolType}

	// drawSpriteFunctionType := dectype.NewFunctionAtom(astFunctionType, decParameters)

	// -----------------------------------------------------------------------------------------------------------------
	//  map : ((a -> b), List a) -> List b
	// -----------------------------------------------------------------------------------------------------------------
	//astListA := MakeFakeAstTypeReferenceWithLocalTypeNames("List", []string{"a"})
	//astListB := MakeFakeAstTypeReferenceWithLocalTypeNames("List", []string{"b"})

	astNames := []*ast.LocalTypeName{MakeFakeLocalTypeName("a"), MakeFakeLocalTypeName("b")}
	astNameOnlyContext := ast.NewLocalTypeNameContext(astNames, nil)
	astARef, _ := astNameOnlyContext.GetOrCreateReferenceFromName(MakeFakeLocalTypeName("a"))
	astBRef, _ := astNameOnlyContext.GetOrCreateReferenceFromName(MakeFakeLocalTypeName("b"))
	astFnAToB := ast.NewFunctionType([]ast.Type{astARef, astBRef})
	astListRef := MakeFakeAstTypeReference("List")
	astListOfA := ast.NewTypeReference(astListRef.TypeIdentifier(), []ast.Type{astARef})
	astListOfB := ast.NewTypeReference(astListRef.TypeIdentifier(), []ast.Type{astBRef})

	astGenericMapFn := ast.NewFunctionType([]ast.Type{astFnAToB, astListOfA, astListOfB})

	decNameOnlyContext := dectype.NewLocalTypeNameOnlyContext(astNames)
	decRefA, _ := decNameOnlyContext.LookupNameReference(astARef)
	decRefB, _ := decNameOnlyContext.LookupNameReference(astBRef)
	decFnAToB := dectype.NewFunctionAtom(astFnAToB, []dtype.Type{decRefA, decRefB})

	decListA := dectype.NewPrimitiveType(astListRef.TypeIdentifier(), []dtype.Type{decRefA})
	decListB := dectype.NewPrimitiveType(astListRef.TypeIdentifier(), []dtype.Type{decRefB})

	decGenericMapFn := dectype.NewFunctionAtom(astGenericMapFn, []dtype.Type{decFnAToB, decListA, decListB})
	decNameOnlyContext.SetType(decGenericMapFn)

	decNameOnlyContextForMapRef := dectype.NewLocalTypeNameContextReference(nil, decNameOnlyContext)

	// -----------------------------------------------------------------------------------------------------------------
	// Concrete
	// -----------------------------------------------------------------------------------------------------------------

	decIntToBoolFn := dectype.NewFunctionAtom(nil, []dtype.Type{decIntType, decBoolType})

	astListANames := []*ast.LocalTypeName{MakeFakeLocalTypeName("a")}
	astListANameOnlyContext := ast.NewLocalTypeNameContext(astListANames, nil)
	astListARef, _ := astListANameOnlyContext.GetOrCreateReferenceFromName(MakeFakeLocalTypeName("a"))

	astListAAgain := ast.NewTypeReference(MakeFakeTypeIdentifier("List"), []ast.Type{astListARef})
	decListANameOnlyContext := dectype.NewLocalTypeNameOnlyContext(astListANames)
	decRefAAgain, _ := decNameOnlyContext.LookupNameReference(astListARef)
	decListAAgain := dectype.NewPrimitiveType(astListAAgain.TypeIdentifier(), []dtype.Type{decRefAAgain})
	decListANameOnlyContext.SetType(decListAAgain)
	decListANameOnlyContextRef := dectype.NewLocalTypeNameContextReference(nil, decListANameOnlyContext)

	decListOfIntWrappedInResolvedContext, _ := dectype.NewResolvedLocalTypeContext(decListANameOnlyContextRef,
		[]dtype.Type{decIntType})

	argumentsForConretization := []dtype.Type{decIntToBoolFn, decListOfIntWrappedInResolvedContext,
		dectype.NewAnyType()}

	concretizedContext, err := concretize.ConcretizeLocalTypeContextUsingArguments(decNameOnlyContextForMapRef,
		argumentsForConretization)

	if err != nil {
		t.Fatal(err)
	}

	log.Println(debug.TreeString(concretizedContext))

	log.Printf("RESOLVING TO ATOM")

	atomFromContext := dectype.ResolveToAtom(concretizedContext)
	functionAtom, _ := atomFromContext.(*dectype.FunctionAtom)
	log.Printf("atom is %T :%s", functionAtom, debug.TreeString(functionAtom))

	expectedReturn := MakeFakeListOf(MakeBoolType())

	err2 := dectype.CompatibleTypes(expectedReturn, functionAtom.ReturnType())
	if err2 != nil {
		t.Fatal(err2)
	}
}

/*
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

	decoratedTypeNames := dectype.NewLocalTypeNameOnlyContext()
	aName := MakeFakeLocalTypeNameRef("a")
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

	decoratedTypeNames.SetType(functionWithLocalTypeNames)

	concretizedVariant, concreteErr := concretize.ConcretizeLocalTypeContextUsingArguments(
		decoratedTypeNames, concreteArguments,
	)
	if concreteErr != nil {
		t.Error(concreteErr)
	}

	resolved, resolveErr := concretizedVariant.Resolve()
	if resolveErr != nil {
		t.Fatalf("did not work %v", resolveErr)
	}

	log.Printf("functionWithLocalTypeNames variant: %v", functionWithLocalTypeNames)
	log.Printf("concretized function: %v", concretizedVariant)
	log.Printf("resolved %v", resolved)

}

*/
