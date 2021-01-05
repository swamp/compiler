/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"testing"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func createVariableIdentifier(name string) *ast.VariableIdentifier {
	symbol := token.NewVariableSymbolToken(name, token.PositionLength{}, 0)
	ident := ast.NewVariableIdentifier(symbol)

	return ident
}

func createLocalType(name string) dtype.Type {
	ident := createVariableIdentifier(name)
	typeParam := ast.NewTypeParameter(ident)
	return dectype.NewLocalType(typeParam)
}

func createTypeIdentifier(name string) *ast.TypeIdentifier {
	symbol := token.NewTypeSymbolToken(name, token.PositionLength{}, 0)
	ident := ast.NewTypeIdentifier(symbol)
	return ident
}

func TestIt(t *testing.T) {
	intType := dectype.NewPrimitiveType(createTypeIdentifier("Int"), nil)

	listA := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{createLocalType("a")})
	listInt := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{intType})

	newType, smashErr := dectype.SmashType(listA, listInt)
	if smashErr != nil {
		t.Fatal(smashErr)
	}
	differentErr := dectype.CompatibleTypes(newType, listInt)
	if differentErr != nil {
		t.Fatal(differentErr)
	}
}

func TestFunc(t *testing.T) {
	intType := dectype.NewPrimitiveType(createTypeIdentifier("Int"), nil)
	stringType := dectype.NewPrimitiveType(createTypeIdentifier("String"), nil)

	a := createLocalType("a")
	listA := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{createLocalType("a")})
	listInt := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{intType})

	genericFunc := dectype.NewFunctionAtom([]dtype.Type{listA, stringType, a})
	callFunc := dectype.NewFunctionAtom([]dtype.Type{listInt, stringType, intType})

	newType, smashErr := dectype.SmashType(genericFunc, callFunc)
	if smashErr != nil {
		t.Fatal(smashErr)
	}
	differentErr := dectype.CompatibleTypes(newType, callFunc)
	if differentErr != nil {
		t.Fatal(differentErr)
	}
}

func TestFuncSecondLevel(t *testing.T) {
	intType := dectype.NewPrimitiveType(createTypeIdentifier("Int"), nil)
	stringType := dectype.NewPrimitiveType(createTypeIdentifier("String"), nil)

	a := createLocalType("a")
	listA := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{createLocalType("a")})

	listInt := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{intType})
	listListInt := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{listInt})

	genericFunc := dectype.NewFunctionAtom([]dtype.Type{listA, stringType, a})
	wrongCallFunc := dectype.NewFunctionAtom([]dtype.Type{listListInt, stringType, intType})

	_, smashErr := dectype.SmashType(genericFunc, wrongCallFunc)
	if smashErr == nil {
		t.Fatalf("should fail here")
	}

	correctCallFunc := dectype.NewFunctionAtom([]dtype.Type{listListInt, stringType, listInt})
	newType, smashErr2 := dectype.SmashType(genericFunc, correctCallFunc)
	if smashErr2 != nil {
		t.Fatal(smashErr2)
	}

	differentErr := dectype.CompatibleTypes(newType, correctCallFunc)
	if differentErr != nil {
		t.Fatal(differentErr)
	}
}

func TestFuncSecondLevel2(t *testing.T) {
	stringType := dectype.NewPrimitiveType(createTypeIdentifier("String"), nil)

	listB := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{createLocalType("a")})
	anotherName := createVariableIdentifier("another")
	fieldAnotherListB := dectype.NewRecordField(anotherName, listB)
	recordParameterB := dtype.NewTypeArgumentName(createVariableIdentifier("a"))
	recordAnotherB := dectype.NewRecordType([]*dectype.RecordField{fieldAnotherListB}, []*dtype.TypeArgumentName{recordParameterB})

	coolName := createVariableIdentifier("cool")
	fieldCoolListA := dectype.NewRecordField(coolName, recordAnotherB)

	recordParameterA := dtype.NewTypeArgumentName(createVariableIdentifier("a"))
	recordA := dectype.NewRecordType([]*dectype.RecordField{fieldCoolListA}, []*dtype.TypeArgumentName{recordParameterA})

	newType, newTypeErr := dectype.CallType(recordA, []dtype.Type{stringType})
	if newTypeErr != nil {
		t.Fatal(newTypeErr)
	}

	listString := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{stringType})
	fieldAnotherListString := dectype.NewRecordField(anotherName, listString)

	// Result
	subRecordResult := dectype.NewRecordType([]*dectype.RecordField{fieldAnotherListString},nil)
	resultCoolName := createVariableIdentifier("cool")
	resultSubCoolField := dectype.NewRecordField(resultCoolName, subRecordResult)
	recordResult := dectype.NewRecordType([]*dectype.RecordField{resultSubCoolField},nil)

	compatibleErr := dectype.CompatibleTypes(recordResult, newType)
	if compatibleErr != nil {
		t.Fatal(compatibleErr)
	}
}


func createMaybe() dtype.Type {
	a := createLocalType("a")
	aName := dtype.NewTypeArgumentName(createVariableIdentifier("a"))
	firstVariant := dectype.NewCustomTypeVariant(0, createTypeIdentifier("Nothing"), nil)
	secondVariant := dectype.NewCustomTypeVariant(1, createTypeIdentifier("Just"), []dtype.Type{a})
	globalMaybe := dectype.NewCustomType(createTypeIdentifier("Maybe"), []*dtype.TypeArgumentName{aName},
		[]*dectype.CustomTypeVariant{firstVariant, secondVariant})

	return globalMaybe
}

// foldlstop : (a -> b -> Maybe b) -> b -> List a -> b

// 1: findAnimationStep : StepInfo -> StepAccumulator -> Maybe StepAccumulator
// 2: { accumulatedTime = 0, targetTime = timeInAnimation, step = { index = -1, time = -1 } } // StepAccumulator
// 3: List <StepInfo> // steps

func TestFuncFoldlStop(t *testing.T) {
	a := createLocalType("a")
	b := createLocalType("b")

	globalMaybe := createMaybe()

	invokedMaybeB, _ := dectype.NewInvokerType(globalMaybe, []dtype.Type{b})

	subFunction := dectype.NewFunctionAtom([]dtype.Type{a, b, invokedMaybeB})

	listA := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{a})
	wholeFunction := dectype.NewFunctionAtom([]dtype.Type{subFunction, b, listA, b})

	step := dectype.NewPrimitiveType(createTypeIdentifier("String"), nil)
	accumulator := dectype.NewPrimitiveType(createTypeIdentifier("Int"), nil)
	steps := dectype.NewPrimitiveType(createTypeIdentifier("List"), []dtype.Type{step})

	maybeAccumulator, _ := dectype.NewInvokerType(globalMaybe, []dtype.Type{accumulator})
	findAnimationStep := dectype.NewFunctionAtom([]dtype.Type{step, accumulator, maybeAccumulator})


	arguments := []dtype.Type{findAnimationStep, accumulator, steps}
	callFunc := dectype.NewFunctionAtom(arguments)

	_, callErr := dectype.SmashFunctions(wholeFunction, callFunc)
	if callErr != nil {
		t.Fatal(callErr)
	}
}


func TestSimple(t *testing.T) {
	stringType := dectype.NewPrimitiveType(createTypeIdentifier("String"), nil)
	intType := dectype.NewPrimitiveType(createTypeIdentifier("Int"), nil)

	recordType := dectype.NewRecordType(nil, nil)

	declaration := dectype.NewFunctionAtom([]dtype.Type{intType, recordType})

	callFunc := dectype.NewFunctionAtom([]dtype.Type{stringType, recordType})

	called, callErr := dectype.SmashFunctions(declaration, callFunc)
	if callErr != nil {
		t.Fatal(callErr)
	}

	first := called.FunctionParameterTypes()[0]
	if first.DecoratedName() != "Int" {
		t.Errorf("wrong first function: %v", first)
	}

	compatibleErr := dectype.CompatibleTypes(called, callFunc)
	if compatibleErr == nil {
		t.Errorf("should have failed")
	}
	fmt.Printf("intentional error:%v\n", compatibleErr)
}
