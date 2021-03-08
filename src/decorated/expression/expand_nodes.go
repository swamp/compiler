package decorated

import (
	"log"
	"reflect"

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
	optionalModuleRef := fn.NameReference().ModuleReference()
	if optionalModuleRef != nil {
		tokens = append(tokens, optionalModuleRef)
	}
	return tokens
}

func expandChildNodesFunctionCall(fn *FunctionCall) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.FunctionExpression())...)
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

func expandChildNodesImportStatement(importStatement *ImportStatement) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(importStatement.ModuleReference())...)
	if importStatement.Alias() != nil {
		tokens = append(tokens, expandChildNodes(importStatement.Alias())...)
	}
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
	// tokens = append(tokens, expandChildNodes(fn.TypeReference())...) Can not expand type identifiers, need meaning
	for _, variant := range fn.Variants() {
		tokens = append(tokens, expandChildNodes(variant)...)
		for _, param := range variant.ParameterTypes() {
			tokens = append(tokens, expandChildNodes(param)...)
		}
	}
	return tokens
}

func expandChildNodesRecordType(fn *dectype.RecordAtom) []TypeOrToken {
	var tokens []TypeOrToken
	for _, field := range fn.ParseOrderedFields() {
		tokens = append(tokens, expandChildNodes(field.FieldName())...)
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

func expandChildNodesTypeReferenceScoped(fn *dectype.TypeReferenceScoped) []TypeOrToken {
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

func expandChildNodesArrayLiteral(arrayLiteral *ArrayLiteral) []TypeOrToken {
	var tokens []TypeOrToken
	for _, expression := range arrayLiteral.Expressions() {
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
	optionalModuleRef := constructor.NamedTypeReference().ModuleReference()
	if optionalModuleRef != nil {
		tokens = append(tokens, optionalModuleRef)
	}

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
	// tokens = append(tokens, expandChildNodes(constructor.typeIdentifier)...) // TODO: Need meaning
	return tokens
}

func expandChildNodesCaseForTypeAlias(typeAlias *dectype.Alias) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(typeAlias.TypeIdentifier())...)
	tokens = append(tokens, expandChildNodes(typeAlias.Next())...)

	return tokens
}

func expandChildNodesCaseForCustomType(caseForCustomType *CaseCustomType) []TypeOrToken {
	var tokens []TypeOrToken

	tokens = append(tokens, expandChildNodes(caseForCustomType.Test())...)

	for _, consequence := range caseForCustomType.Consequences() {
		tokens = append(tokens, expandChildNodes(consequence.VariantReference())...)
		for _, param := range consequence.Parameters() {
			tokens = append(tokens, expandChildNodes(param)...)
		}
		tokens = append(tokens, expandChildNodes(consequence.Expression())...)
	}

	if caseForCustomType.DefaultCase() != nil {
		tokens = append(tokens, expandChildNodes(caseForCustomType.DefaultCase())...)
	}

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
	if node == nil || reflect.ValueOf(node).IsNil() {
		panic("can not be nil")
	}
	tokens := []TypeOrToken{node}
	switch t := node.(type) {
	case *ModuleReference:
		return tokens
	case *AnnotationStatement:
		return append(tokens, expandChildNodesAnnotation(t)...)
	case *ImportStatement:
		return append(tokens, expandChildNodesImportStatement(t)...)
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
	case *ArrayLiteral:
		return append(tokens, expandChildNodesArrayLiteral(t)...)
	case *RecordLiteral:
		return append(tokens, expandChildNodesRecordLiteral(t)...)
	case *FunctionParameterDefinition:
		return tokens
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
	case *FixedLiteral: // Should not be expanded
		return tokens
	case *CharacterLiteral: // Should not be expanded
		return tokens
	case *TypeIdLiteral: // Should not be expanded
		return tokens
	case *ResourceNameLiteral: // Should not be expanded
		return tokens
	case *StringInterpolation: // Should not be expanded
		return tokens
	case *BooleanLiteral: // Should not be expanded
		return tokens
	case *StringLiteral: // Should not be expanded
		return tokens
	case *MultilineComment: // Should not be expanded
		return tokens
	case *RecordLiteralField: // Should not be expanded
		return tokens
	case *BitwiseUnaryOperator:
		return expandChildNodes(&t.UnaryOperator)
	case *LogicalUnaryOperator:
		return expandChildNodes(&t.UnaryOperator)
	case *UnaryOperator:
		return expandChildNodes(t.Left())
	case *Constant:
		return tokens
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
	case *dectype.CustomTypeVariant:
		return tokens
	case *dectype.RecordAtom:
		return append(tokens, expandChildNodesRecordType(t)...)
	case *dectype.RecordFieldName:
		return tokens
	case *dectype.FunctionTypeReference:
		return append(tokens, expandChildNodesFunctionTypeReference(t)...)
	case *dectype.TypeReference:
		return append(tokens, expandChildNodesTypeReference(t)...)
	case *dectype.TypeReferenceScoped:
		return append(tokens, expandChildNodesTypeReferenceScoped(t)...)
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