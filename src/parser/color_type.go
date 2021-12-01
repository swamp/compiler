/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/swamp/compiler/src/coloring"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func colorAlias(alias *dectype.Alias, colorer coloring.Colorer) {
	colorer.AliasNameSymbol(alias.TypeIdentifier().Symbol())
}

func colorFunctionParameters(functionParameterTypes []dtype.Type, indentation int, inside bool, colorer coloring.Colorer) {
	for index, parameterType := range functionParameterTypes {
		if index > 0 {
			colorer.OneSpace()
			colorer.RightArrow()
			colorer.OneSpace()
		}
		ColorType(parameterType, indentation, inside, colorer)
	}
}

func colorFunctionParametersWithAlias(functionParameterTypes []dtype.Type, indentation int, inside bool, colorer coloring.Colorer) {
	colorFunctionParameters(functionParameterTypes, indentation, inside, colorer)
	userInstruction("Explanation", indentation+1, colorer)
	ColorTypesWithAtom(functionParameterTypes, indentation+2, inside, colorer)
}

func colorFunctionType(functionType *dectype.FunctionAtom, indentation int, inside bool, colorer coloring.Colorer) {
	leftToken := token.NewOperatorToken(token.LeftParen, token.SourceFileReference{}, "(", "(")
	colorer.Operator(leftToken)
	colorFunctionParameters(functionType.FunctionParameterTypes(), indentation, inside, colorer)

	rightToken := token.NewOperatorToken(token.RightParen, token.SourceFileReference{}, ")", ")")
	colorer.Operator(rightToken)
}

func newDoubleLine(indentation int, colorer coloring.Colorer) {
	colorer.NewLine(0)
	colorer.NewLine(indentation)
}

func colorRecordType(recordType *dectype.RecordAtom, indentation int, inside bool, colorer coloring.Colorer) {
	colorer.OperatorString("{")
	colorer.OneSpace()
	continuationIndentation := indentation
	if inside {
		continuationIndentation++
	}
	for index, fieldInType := range recordType.ParseOrderedFields() {
		if index > 0 {
			colorer.NewLine(continuationIndentation)
			colorer.OperatorString(",")
			colorer.OneSpace()
		}
		colorer.VariableSymbol(fieldInType.VariableIdentifier().Symbol())
		operatorToken := token.NewOperatorToken(token.Colon, token.SourceFileReference{}, ":", ":")
		colorer.OneSpace()
		colorer.Operator(operatorToken)
		colorer.OneSpace()
		ColorType(fieldInType.Type(), continuationIndentation, true, colorer)
	}
	colorer.OneSpace()
	colorer.OperatorString("}")
}

func colorCustomType(recordType *dectype.CustomTypeAtom, indentation int, inside bool, colorer coloring.Colorer) {
	s := coloring.ColorTypeSymbol(recordType.TypeIdentifier().Symbol())
	indentation++
	colorer.NewLine(indentation)
	for index, fieldInType := range recordType.Variants() {
		if index > 0 {
			colorer.NewLine(indentation)
		}
		colorer.TypeSymbol(fieldInType.Name().Symbol())

		for index, parameterType := range fieldInType.ParameterTypes() {
			if index > 0 {
				s += " "
			}
			ColorType(parameterType, indentation, inside, colorer)
		}
	}
}

func colorTypeEmbed(recordType *dectype.InvokerType, colorer coloring.Colorer) {
	shortName := recordType.TypeGenerator().String()
	typeSymbolToken := token.NewTypeSymbolToken(shortName, token.SourceFileReference{}, 0)
	colorer.TypeGeneratorName(typeSymbolToken)
	colorer.OperatorString("(")
	for index, fieldInType := range recordType.Params() {
		if index > 0 {
			colorer.OneSpace()
		}
		ColorType(fieldInType, 0, false, colorer)
	}
	colorer.OperatorString(")")
}

func colorPrimitive(primitive *dectype.PrimitiveAtom, indentation int, inside bool, colorer coloring.Colorer) {
	fakeSymbol := token.NewTypeSymbolToken(primitive.HumanReadable(), token.SourceFileReference{}, 0)
	colorer.PrimitiveType(fakeSymbol)
}

func colorLocalType(primitive *dectype.LocalType, indentation int, inside bool, colorer coloring.Colorer) {
	colorer.LocalType(primitive.Identifier().Identifier().Symbol())
}

func colorAny(indentation int, inside bool, colorer coloring.Colorer) {
	fakeSymbol := token.NewTypeSymbolToken("ANY", token.SourceFileReference{}, 0)
	colorer.PrimitiveType(fakeSymbol)
}

func ColorAtom(atomType dtype.Atom, indentation int, inside bool, colorer coloring.Colorer) {
	switch t := atomType.(type) {
	case *dectype.FunctionAtom:
		colorFunctionType(t, indentation, inside, colorer)
	case *dectype.RecordAtom:
		colorRecordType(t, indentation, inside, colorer)
	case *dectype.CustomTypeAtom:
		colorCustomType(t, indentation, inside, colorer)
	case *dectype.PrimitiveAtom:
		colorPrimitive(t, indentation, inside, colorer)
	case *dectype.LocalType:
		colorLocalType(t, indentation, inside, colorer)
	default:
		panic(fmt.Sprintf("ColorAtom: unknown type %T", atomType))
	}
}

func ColorType(dType dtype.Type, indentation int, inside bool, colorer coloring.Colorer) {
	color.NoColor = false
	switch t := dType.(type) {
	case *dectype.Alias:
		colorAlias(t, colorer)
	case *dectype.InvokerType:
		colorTypeEmbed(t, colorer)
	default:
		{
			atom, wasAtom := dType.(dtype.Atom)
			if wasAtom {
				ColorAtom(atom, indentation, inside, colorer)
			} else {
				panic(fmt.Sprintf("ColorType: unknown type %v %T", dType, dType))
			}
		}
	}
}

func ColorTypeWithAtom(dType dtype.Type, indentation int, inside bool, colorer coloring.Colorer) {
	color.NoColor = false
	switch t := dType.(type) {
	case *dectype.Alias:
		{
			colorAlias(t, colorer)
			atom, _ := t.Resolve()
			colorer.UserInstruction(" which resolves to ")
			newDoubleLine(indentation+1, colorer)
			ColorAtom(atom, indentation+1, inside, colorer)
		}
	default:
		ColorType(dType, indentation+1, inside, colorer)
	}
}

func ColorTypesWithAtom(dTypes []dtype.Type, indentation int, inside bool, colorer coloring.Colorer) {
	color.NoColor = false
	for index, foundType := range dTypes {
		if index > 0 {
			colorer.NewLine(indentation)
		}
		ColorTypeWithAtom(foundType, indentation, inside, colorer)
	}
}
