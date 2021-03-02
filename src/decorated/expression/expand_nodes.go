package decorated

import (
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type TypeOrToken interface {
	String() string
	FetchPositionLength() token.SourceFileReference
}

type Token interface {
	TypeOrToken
	Type() dtype.Type
	HumanReadable() string
}

func expandChildNodesFunctionValue(fn *FunctionValue) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.Expression())...)
	for _, parameter := range fn.Parameters() {
		tokens = append(tokens, expandChildNodes(parameter)...)
	}
	return tokens
}

func expandChildNodesFunctionCall(fn *FunctionCall) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.FunctionValue())...)
	for _, argument := range fn.Arguments() {
		tokens = append(tokens, expandChildNodes(argument)...)
	}
	return tokens
}

func expandChildNodesCurryFunction(fn *CurryFunction) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.FunctionValue())...)
	for _, argument := range fn.ArgumentsToSave() {
		tokens = append(tokens, expandChildNodes(argument)...)
	}
	return tokens
}

func expandChildNodesAnnotation(fn *Annotation) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.Type())...)
	return tokens
}

func expandChildNodesFunctionType(fn *dectype.FunctionAtom) []TypeOrToken {
	var tokens []TypeOrToken
	for _, parameter := range fn.FunctionParameterTypes() {
		tokens = append(tokens, expandChildNodes(parameter)...)
	}
	return tokens
}

func expandChildNodesCustomType(fn *dectype.CustomTypeAtom) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.TypeIdentifier())...)
	for _, variant := range fn.Variants() {
		tokens = append(tokens, expandChildNodes(variant.Name())...)
		for _, param := range variant.ParameterTypes() {
			tokens = append(tokens, expandChildNodes(param)...)
		}
	}
	return tokens
}

func expandChildNodesRecordType(fn *dectype.RecordAtom) []TypeOrToken {
	var tokens []TypeOrToken
	for _, field := range fn.ParseOrderedFields() {
		tokens = append(tokens, expandChildNodes(field.VariableIdentifier())...)
		tokens = append(tokens, expandChildNodes(field.Type())...)
	}
	return tokens
}

func expandChildNodesFunctionTypeReference(fn *dectype.FunctionTypeReference) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.FunctionAtom())...)
	return tokens
}

func expandChildNodesTypeReference(fn *dectype.TypeReference) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.Next())...)
	return tokens
}

func expandChildNodesPrimitive(fn *dectype.PrimitiveAtom) []TypeOrToken {
	var tokens []TypeOrToken
	for _, parameter := range fn.GenericTypes() {
		tokens = append(tokens, expandChildNodes(parameter)...)
	}
	return tokens
}

func expandChildNodesInvokerType(fn *dectype.InvokerType) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.TypeGenerator())...)
	for _, param := range fn.Params() {
		tokens = append(tokens, expandChildNodes(param)...)
	}
	return tokens
}

func expandChildNodesLetAssignment(assignment *LetAssignment) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(assignment.LetVariable())...)
	tokens = append(tokens, expandChildNodes(assignment.Expression())...)

	return tokens
}

func expandChildNodesListLiteral(listLiteral *ListLiteral) []TypeOrToken {
	var tokens []TypeOrToken
	for _, expression := range listLiteral.Expressions() {
		tokens = append(tokens, expandChildNodes(expression)...)
	}

	return tokens
}

func expandChildNodesNamedFunctionValue(namedFunctionValue *NamedFunctionValue) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(namedFunctionValue.Identifier())...)
	tokens = append(tokens, expandChildNodes(namedFunctionValue.Value())...)

	return tokens
}

func expandChildNodesCustomTypeVariantConstructor(constructor *CustomTypeVariantConstructor) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(constructor.Reference())...)

	for _, arg := range constructor.arguments {
		tokens = append(tokens, expandChildNodes(arg)...)
	}

	return tokens
}

func expandChildNodesCustomTypeVariantReference(constructor *CustomTypeVariantReference) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(constructor.typeIdentifier)...)
	return tokens
}

func expandChildNodesCaseForCustomType(namedFunctionValue *CaseCustomType) []TypeOrToken {
	var tokens []TypeOrToken
	for _, consequence := range namedFunctionValue.Consequences() {
		tokens = append(tokens, expandChildNodes(consequence.Identifier())...)
		for _, param := range consequence.Parameters() {
			tokens = append(tokens, expandChildNodes(param)...)
		}
		tokens = append(tokens, expandChildNodes(consequence.Expression())...)
	}

	tokens = append(tokens, expandChildNodes(namedFunctionValue.DefaultCase())...)

	return tokens
}

func expandChildNodesRecordLookups(lookup *RecordLookups) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(lookup.Expression())...)
	for _, lookupField := range lookup.LookupFields() {
		tokens = append(tokens, expandChildNodes(lookupField.reference)...)
	}

	return tokens
}

func expandChildNodesLet(let *Let) []TypeOrToken {
	var tokens []TypeOrToken
	for _, assignment := range let.Assignments() {
		tokens = append(tokens, expandChildNodes(assignment)...)
	}

	tokens = append(tokens, expandChildNodes(let.Consequence())...)

	return tokens
}

/*
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.CaseCustomType
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LogicalOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.BooleanOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 file:///home/peter/own/hackman/swamp/gameplay/Projectile.swamp:24:1: warning: 'projectile' not used in function checkForCollide
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.If
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.BooleanOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.BooleanOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.BooleanOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LogicalOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LogicalOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.Guard
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.CasePatternMatching
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeLeftOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.ArithmeticOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.If
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.TypeIdLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.StringInterpolation
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.CaseCustomType
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.If
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.CaseCustomType
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.StringLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.CaseCustomType
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordConstructor
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.ArithmeticOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.IntegerLiteral
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.ArithmeticOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionParameterReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordFieldReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.FunctionReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariableReference
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.CaseCustomType
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.If
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.LetVariable
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.PipeRightOperator
2021/03/02 18:50:02 expand_nodes: could not expand: *decorated.RecordLiteral
*/

func expandChildNodes(node Node) []TypeOrToken {
	tokens := []TypeOrToken{node}
	switch t := node.(type) {
	case *ast.TypeIdentifier:
		return tokens
	case *ast.VariableIdentifier:
		return tokens
	case *Annotation:
		return append(tokens, expandChildNodesAnnotation(t)...)
	case *FunctionValue:
		return append(tokens, expandChildNodesFunctionValue(t)...)
	case *FunctionCall:
		return append(tokens, expandChildNodesFunctionCall(t)...)
	case *CurryFunction:
		return append(tokens, expandChildNodesCurryFunction(t)...)
	case *Let:
		return append(tokens, expandChildNodesLet(t)...)
	case *LetAssignment:
		return append(tokens, expandChildNodesLetAssignment(t)...)
	case *ListLiteral:
		return append(tokens, expandChildNodesListLiteral(t)...)
	case *FunctionParameterDefinition:
		return append(tokens, expandChildNodes(t.identifier)...)
	case *NamedFunctionValue:
		return append(tokens, expandChildNodesNamedFunctionValue(t)...)
	case *CustomTypeVariantConstructor:
		return append(tokens, expandChildNodesCustomTypeVariantConstructor(t)...)
	case *CustomTypeVariantReference:
		return append(tokens, expandChildNodesCustomTypeVariantReference(t)...)
	case *CaseCustomType:
		return append(tokens, expandChildNodesCaseForCustomType(t)...)
	case *RecordLookups:
		return append(tokens, expandChildNodesRecordLookups(t)...)
	case *AsmConstant:
		return tokens
	case *dectype.LocalType:
		return tokens
	case *dectype.Alias:
		return append(tokens, expandChildNodes(t.Next())...)
	case *dectype.PrimitiveAtom:
		return append(tokens, expandChildNodesPrimitive(t)...)
	case *dectype.InvokerType:
		return append(tokens, expandChildNodesInvokerType(t)...)
	case *dectype.FunctionAtom:
		return append(tokens, expandChildNodesFunctionType(t)...)
	case *dectype.CustomTypeAtom:
		return append(tokens, expandChildNodesCustomType(t)...)
	case *dectype.RecordAtom:
		return append(tokens, expandChildNodesRecordType(t)...)
	case *dectype.FunctionTypeReference:
		return append(tokens, expandChildNodesFunctionTypeReference(t)...)
	case *dectype.TypeReference:
		return append(tokens, expandChildNodesTypeReference(t)...)
	default:
		log.Printf("expand_nodes: could not expand: %T\n", t)
		return tokens
	}
}

func ExpandAllChildNodes(nodes []Node) []TypeOrToken {
	var tokens []TypeOrToken
	for _, node := range nodes {
		tokens = append(tokens, expandChildNodes(node)...)
	}

	return tokens
}
