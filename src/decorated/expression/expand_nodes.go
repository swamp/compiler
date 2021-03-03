package decorated

import (
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type Statement interface {
	TypeOrToken
	StatementString() string
}

type TypeOrToken interface {
	String() string
	FetchPositionLength() token.SourceFileReference
}

type HumanReadEnabler interface {
	HumanReadable() string
}

type Token interface {
	TypeOrToken
	HumanReadEnabler
	Type() dtype.Type
}

func expandChildNodesFunctionValue(fn *FunctionValue) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.Expression())...)
	for _, parameter := range fn.Parameters() {
		tokens = append(tokens, expandChildNodes(parameter)...)
	}
	return tokens
}

func expandChildNodesFunctionReference(fn *FunctionReference) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.ident)...)
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

func expandChildNodesAnnotation(fn *AnnotationStatement) []TypeOrToken {
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

func expandChildNodesRecordLiteral(recordLiteral *RecordLiteral) []TypeOrToken {
	var tokens []TypeOrToken

	if recordLiteral.RecordTemplate() != nil {
		tokens = append(tokens, expandChildNodes(recordLiteral.RecordTemplate())...)
	}

	for _, assignment := range recordLiteral.ParseOrderedAssignments() {
		tokens = append(tokens, expandChildNodes(assignment.FieldName())...)
		tokens = append(tokens, expandChildNodes(assignment.Expression())...)
	}

	return tokens
}

func expandChildNodesNamedFunctionValue(namedFunctionValue *NamedFunctionValue) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(namedFunctionValue.FunctionName())...)
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

func expandChildNodesRecordConstructor(constructor *RecordConstructor) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(constructor.typeIdentifier)...)

	for _, arg := range constructor.arguments {
		tokens = append(tokens, expandChildNodes(arg.FieldName())...)
		tokens = append(tokens, expandChildNodes(arg.Expression())...)
	}

	return tokens
}

func expandChildNodesGuard(guard *Guard) []TypeOrToken {
	var tokens []TypeOrToken
	for _, item := range guard.Items() {
		tokens = append(tokens, expandChildNodes(item.Condition())...)
		tokens = append(tokens, expandChildNodes(item.Expression())...)
	}

	if guard.DefaultGuard() != nil {
		tokens = append(tokens, expandChildNodes(guard.DefaultGuard().Expression())...)
	}

	return tokens
}

func expandChildNodesCustomTypeVariantReference(constructor *CustomTypeVariantReference) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(constructor.typeIdentifier)...)
	return tokens
}

func expandChildNodesCaseForCustomType(caseForCustomType *CaseCustomType) []TypeOrToken {
	var tokens []TypeOrToken

	tokens = append(tokens, expandChildNodes(caseForCustomType.Test())...)

	for _, consequence := range caseForCustomType.Consequences() {
		tokens = append(tokens, expandChildNodes(consequence.Identifier())...)
		for _, param := range consequence.Parameters() {
			tokens = append(tokens, expandChildNodes(param)...)
		}
		tokens = append(tokens, expandChildNodes(consequence.Expression())...)
	}

	tokens = append(tokens, expandChildNodes(caseForCustomType.DefaultCase())...)

	return tokens
}

func expandChildNodesCaseForPatternMatching(caseForCustomType *CasePatternMatching) []TypeOrToken {
	var tokens []TypeOrToken

	tokens = append(tokens, expandChildNodes(caseForCustomType.Test())...)

	for _, consequence := range caseForCustomType.Consequences() {
		tokens = append(tokens, expandChildNodes(consequence.Literal())...)
		tokens = append(tokens, expandChildNodes(consequence.Expression())...)
	}

	tokens = append(tokens, expandChildNodes(caseForCustomType.DefaultCase())...)

	return tokens
}

func expandChildNodesBinaryOperator(namedFunctionValue *BinaryOperator) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(namedFunctionValue.Left())...)
	tokens = append(tokens, expandChildNodes(namedFunctionValue.Right())...)
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

func expandChildNodesIf(ifExpression *If) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(ifExpression.Condition())...)
	tokens = append(tokens, expandChildNodes(ifExpression.Consequence())...)
	tokens = append(tokens, expandChildNodes(ifExpression.Alternative())...)

	return tokens
}

func expandChildNodes(node Node) []TypeOrToken {
	tokens := []TypeOrToken{node}
	switch t := node.(type) {
	case *ast.TypeIdentifier:
		if t.ModuleReference() != nil {
			tokens = append(tokens, t.ModuleReference())
		}
		return tokens
	case *ast.VariableIdentifier:
		if t.ModuleReference() != nil {
			tokens = append(tokens, t.ModuleReference())
		}
		return tokens
	case *ast.ModuleReference:
		for _, part := range t.Parts() {
			tokens = append(tokens, part)
		}
		return tokens
	case *ast.ModuleNamePart:
		return tokens
	case *AnnotationStatement:
		return append(tokens, expandChildNodesAnnotation(t)...)
	case *ImportStatement: // TODO:
		return tokens
	case *FunctionValue:
		return append(tokens, expandChildNodesFunctionValue(t)...)
	case *FunctionReference:
		return append(tokens, expandChildNodesFunctionReference(t)...)
	case *FunctionCall:
		return append(tokens, expandChildNodesFunctionCall(t)...)
	case *CurryFunction:
		return append(tokens, expandChildNodesCurryFunction(t)...)
	case *Let:
		return append(tokens, expandChildNodesLet(t)...)
	case *If:
		return append(tokens, expandChildNodesIf(t)...)
	case *LetAssignment:
		return append(tokens, expandChildNodesLetAssignment(t)...)
	case *ListLiteral:
		return append(tokens, expandChildNodesListLiteral(t)...)
	case *RecordLiteral:
		return append(tokens, expandChildNodesRecordLiteral(t)...)
	case *FunctionParameterDefinition:
		return append(tokens, expandChildNodes(t.identifier)...)
	case *NamedFunctionValue:
		return append(tokens, expandChildNodesNamedFunctionValue(t)...)
	case *CustomTypeVariantConstructor:
		return append(tokens, expandChildNodesCustomTypeVariantConstructor(t)...)
	case *RecordConstructor:
		return append(tokens, expandChildNodesRecordConstructor(t)...)
	case *Guard:
		return append(tokens, expandChildNodesGuard(t)...)
	case *CustomTypeVariantReference:
		return append(tokens, expandChildNodesCustomTypeVariantReference(t)...)
	case *CaseCustomType:
		return append(tokens, expandChildNodesCaseForCustomType(t)...)
	case *CasePatternMatching:
		return append(tokens, expandChildNodesCaseForPatternMatching(t)...)
	case *PipeRightOperator:
		return expandChildNodes(&t.BinaryOperator)
	case *PipeLeftOperator:
		return expandChildNodes(&t.BinaryOperator)
	case *ArithmeticOperator:
		return expandChildNodes(&t.BinaryOperator)
	case *LogicalOperator:
		return expandChildNodes(&t.BinaryOperator)
	case *ConsOperator:
		return expandChildNodes(&t.BinaryOperator)
	case *BooleanOperator:
		return expandChildNodes(&t.BinaryOperator)
	case *BitwiseOperator:
		return expandChildNodes(&t.BinaryOperator)

	case *ArithmeticUnaryOperator:
		return expandChildNodes(&t.UnaryOperator)
	case *FunctionName: // Should not be expanded
		return tokens
	case *ExternalFunctionDeclaration: // Should not be expanded
		return tokens
	case *LetVariableReference: // Should not be expanded
		return tokens
	case *LetVariable: // Should not be expanded
		return tokens
	case *RecordFieldReference: // Should not be expanded
		return tokens
	case *FunctionParameterReference: // Should not be expanded
		return tokens
	case *IntegerLiteral: // Should not be expanded
		return tokens
	case *CharacterLiteral: // Should not be expanded
		return tokens
	case *TypeIdLiteral: // Should not be expanded
		return tokens
	case *StringInterpolation: // Should not be expanded
		return tokens
	case *BooleanLiteral: // Should not be expanded
		return tokens
	case *StringLiteral: // Should not be expanded
		return tokens
	case *RecordLiteralField: // Should not be expanded
		return tokens
	case *BitwiseUnaryOperator:
		return expandChildNodes(&t.UnaryOperator)
	case *LogicalUnaryOperator:
		return expandChildNodes(&t.UnaryOperator)

	case *UnaryOperator:
		return expandChildNodes(t.Left())
	case *BinaryOperator:
		return expandChildNodesBinaryOperator(t)
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
