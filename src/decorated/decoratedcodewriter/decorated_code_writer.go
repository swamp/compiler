package decoratedcodewriter

import (
	"fmt"

	"github.com/swamp/compiler/src/coloring"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func WriteCustomType(customType *dectype.CustomTypeAtom, colorer coloring.DecoratedColorer, indentation int) {
	colorer.KeywordString("type")
	colorer.OneSpace()
	colorer.CustomType(customType)
	indentation++
	colorer.NewLine(indentation)
	colorer.KeywordString("=")
	colorer.OneSpace()

	for index, variant := range customType.Variants() {
		if index > 0 {
			colorer.NewLine(indentation)
			colorer.OperatorString("|")
			colorer.OneSpace()
		}

		colorer.CustomTypeVariant(variant)

		hasParams := variant.ParameterCount() > 0
		if !hasParams {
			continue
		}

		colorer.OneSpace()

		for paramIndex, variantParam := range variant.ParameterTypes() {
			if paramIndex > 0 {
				colorer.OneSpace()
			}
			WriteType(variantParam, colorer, indentation)
		}
	}
}

func WriteRecordType(recordType *dectype.RecordAtom, colorer coloring.DecoratedColorer, indentation int) {
	colorer.OperatorString("{")
	colorer.OneSpace()

	for index, field := range recordType.ParseOrderedFields() {
		if index > 0 {
			colorer.NewLine(indentation)
			colorer.OperatorString(",")
			colorer.OneSpace()
		}

		colorer.RecordTypeField(field)
		colorer.OneSpace()
		colorer.OperatorString(":")
		colorer.OneSpace()
		WriteType(field.Type(), colorer, indentation)
	}

	colorer.NewLine(indentation)
	colorer.OperatorString("}")
}

func WriteAliasType(alias *dectype.Alias, colorer coloring.DecoratedColorer, indentation int) {
	colorer.KeywordString("type alias")
	colorer.OneSpace()
	colorer.AliasName(alias)
	colorer.OneSpace()
	colorer.KeywordString("=")
	colorer.NewLine(indentation + 1)
	WriteType(alias.Next(), colorer, indentation+1)
}

func WriteAliasReference(aliasReference *dectype.AliasReference, colorer coloring.DecoratedColorer) {
	colorer.AliasName(aliasReference.Alias())
}

func WriteInvokerType(invokerType *dectype.InvokerType, colorer coloring.DecoratedColorer, indentation int) {
	colorer.InvokerType(invokerType)

	for _, parameterType := range invokerType.Params() {
		colorer.OneSpace()
		WriteType(parameterType, colorer, indentation)
	}
}

func WriteUnmanagedType(unmanaged *dectype.UnmanagedType, colorer coloring.DecoratedColorer) {
	colorer.KeywordString("Unmanaged")
	colorer.OperatorString("<")
	colorer.UnmanagedName(unmanaged.Identifier())
	colorer.OperatorString(">")
}

func WritePrimitiveType(primitive *dectype.PrimitiveAtom, colorer coloring.DecoratedColorer, indentation int) {
}

func WritePrimitiveTypeReference(primitiveReference *dectype.PrimitiveTypeReference, colorer coloring.DecoratedColorer, indentation int) {
	colorer.PrimitiveTypeName(primitiveReference.PrimitiveAtom().PrimitiveName())
}

func WriteType(decoratedType dtype.Type, colorer coloring.DecoratedColorer, indentation int) {
	switch t := decoratedType.(type) {
	case *dectype.CustomTypeAtom:
		WriteCustomType(t, colorer, indentation)
	case *dectype.InvokerType:
		WriteInvokerType(t, colorer, indentation)
	case *dectype.RecordAtom:
		WriteRecordType(t, colorer, indentation)
	case *dectype.Alias:
		WriteAliasType(t, colorer, indentation)
	case *dectype.AliasReference:
		WriteAliasReference(t, colorer)
	case *dectype.UnmanagedType:
		WriteUnmanagedType(t, colorer)
	case *dectype.PrimitiveTypeReference:
		WritePrimitiveTypeReference(t, colorer, indentation)
	case *dectype.PrimitiveAtom:
		WritePrimitiveType(t, colorer, indentation)
	default:
		panic(fmt.Errorf("couldn't write type %T in decorated writer", decoratedType))
	}
}
