package lspservice

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func addSemanticTokenFunctionValue(f *decorated.FunctionValue, builder *SemanticBuilder) error {
	for _, parameter := range f.Parameters() {
		if err := builder.EncodeSymbol(parameter.FetchPositionLength().Range, "parameter", []string{}); err != nil {
			return err
		}
	}

	return addSemanticToken(f.Expression(), builder)
}

func addSemanticTokenFunctionName(f *decorated.FunctionName, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(f.FetchPositionLength().Range, "function", []string{"definition"}); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenNamedFunctionValue(f *decorated.NamedFunctionValue, builder *SemanticBuilder) error {
	if err := addSemanticToken(f.FunctionName(), builder); err != nil {
		return err
	}

	return addSemanticToken(f.Value(), builder)
}

func addSemanticTokenAnnotation(f *decorated.AnnotationStatement, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(f.Identifier().FetchPositionLength().Range, "function", []string{"declaration"}); err != nil {
		return err
	}
	if err := addSemanticToken(f.Type(), builder); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenFunctionType(f *dectype.FunctionAtom, builder *SemanticBuilder) error {
	for _, paramType := range f.FunctionParameterTypes() {
		if err := addSemanticToken(paramType, builder); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenTypeAlias(f *dectype.Alias, builder *SemanticBuilder) error {
	if err := encodeKeyword(builder, f.AstAlias().KeywordType()); err != nil {
		return err
	}

	if err := encodeKeyword(builder, f.AstAlias().KeywordAlias()); err != nil {
		return err
	}

	if err := encodeEnum(builder, f.TypeIdentifier()); err != nil {
		return err
	}

	if err := addSemanticToken(f.Next(), builder); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenRecordType(f *dectype.RecordAtom, builder *SemanticBuilder) error {
	for _, paramType := range f.ParseOrderedFields() {
		if err := builder.EncodeSymbol(paramType.VariableIdentifier().FetchPositionLength().Range, "property", []string{"declaration"}); err != nil {
			return err
		}

		if err := addSemanticToken(paramType.Type(), builder); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenCustomType(f *dectype.CustomTypeAtom, builder *SemanticBuilder) error {
	if err := encodeKeyword(builder, f.AstCustomType().KeywordType()); err != nil {
		return err
	}
	if err := encodeEnum(builder, f.TypeIdentifier()); err != nil {
		return err
	}

	for _, variant := range f.Variants() {
		if err := encodeEnumMemberTypeIdentifier(builder, variant.Name()); err != nil {
			return err
		}

		for _, parameter := range variant.ParameterTypes() {
			if err := addSemanticToken(parameter, builder); err != nil {
				return err
			}
		}
	}

	return nil
}

func addSemanticTokenConstant(f *decorated.Constant, builder *SemanticBuilder) error {
	if err := encodeConstant(f.FetchPositionLength().Range, builder); err != nil {
		return err
	}
	return nil
}

func addTypeReferencePrimitive(referenceRange token.Range, primitive *dectype.PrimitiveAtom, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(referenceRange, "type", []string{"declaration", "defaultLibrary"}); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenCustomTypeVariantConstructor(constructor *decorated.CustomTypeVariantConstructor, builder *SemanticBuilder) error {
	if err := addSemanticToken(constructor.Reference(), builder); err != nil {
		return err
	}

	for _, arguments := range constructor.Arguments() {
		if err := addSemanticToken(arguments, builder); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenCustomTypeVariantReference(ref *decorated.CustomTypeVariantReference, builder *SemanticBuilder) error {
	encodeEnumMember(builder, ref.AstIdentifier().SomeTypeIdentifier())

	return nil
}

func addSemanticTokenRecordConstructor(constructor *decorated.RecordConstructor, builder *SemanticBuilder) error {
	if err := encodeStructReferenceWithModuleReference(builder, constructor.NamedTypeReference().AstIdentifier().SomeTypeIdentifier()); err != nil {
		return err
	}

	for _, argument := range constructor.ParseOrderArguments() {
		if err := addSemanticToken(argument, builder); err != nil {
			return err
		}
	}
	return nil
}

func encodeProperty(builder *SemanticBuilder, identifier *ast.VariableIdentifier) error {
	return builder.EncodeSymbol(identifier.FetchPositionLength().Range, "property", nil)
}

func addSemanticTokenRecordsLookup(lookups *decorated.RecordLookups, builder *SemanticBuilder) error {
	if err := addSemanticToken(lookups.Expression(), builder); err != nil {
		return err
	}

	for _, lookup := range lookups.LookupFields() {
		if err := encodeProperty(builder, lookup.Identifier()); err != nil {
			return err
		}
	}

	return nil
}

func encodeKeyword(builder *SemanticBuilder, keyword token.Keyword) error {
	return builder.EncodeSymbol(keyword.FetchPositionLength().Range, "keyword", nil)
}

func encodeComment(builder *SemanticBuilder, comment token.CommentToken) error {
	return builder.EncodeSymbol(comment.FetchPositionLength().Range, "comment", nil)
}

func encodeMultilineComment(builder *SemanticBuilder, comment token.MultiLineCommentToken) error {
	return builder.EncodeSymbol(comment.FetchPositionLength().Range, "comment", nil)
}

func encodeOperator(builder *SemanticBuilder, operator token.OperatorToken) error {
	return builder.EncodeSymbol(operator.FetchPositionLength().Range, "operator", nil)
}

func encodeEnumMember(builder *SemanticBuilder, identifier ast.TypeIdentifierNormalOrScoped) error {
	return builder.EncodeSymbol(identifier.FetchPositionLength().Range, "enumMember", nil)
}

func encodeEnumMemberTypeIdentifier(builder *SemanticBuilder, identifier *ast.TypeIdentifier) error {
	return builder.EncodeSymbol(identifier.FetchPositionLength().Range, "enumMember", nil)
}

func encodeEnum(builder *SemanticBuilder, identifier *ast.TypeIdentifier) error {
	return builder.EncodeSymbol(identifier.FetchPositionLength().Range, "enum", nil)
}

func encodeConstantTypeIdentifier(builder *SemanticBuilder, identifier *ast.TypeIdentifier) error {
	return builder.EncodeSymbol(identifier.FetchPositionLength().Range, "enum", nil)
}

func encodeConstant(rangeFound token.Range, builder *SemanticBuilder) error {
	return builder.EncodeSymbol(rangeFound, "macro", nil)
}

func encodeVariable(builder *SemanticBuilder, identifier *ast.VariableIdentifier) error {
	return builder.EncodeSymbol(identifier.FetchPositionLength().Range, "variable", []string{"readonly"})
}

func encodeModuleReference(builder *SemanticBuilder, astModuleReference *ast.ModuleReference) error {
	for _, namespacePart := range astModuleReference.Parts() {
		if err := builder.EncodeSymbol(namespacePart.TypeIdentifier().FetchPositionLength().Range, "namespace", nil); err != nil {
			return err
		}
	}
	return nil
}

func encodeModuleAlias(builder *SemanticBuilder, astModuleAlias *ast.TypeIdentifier) error {
	if err := builder.EncodeSymbol(astModuleAlias.FetchPositionLength().Range, "namespace", nil); err != nil {
		return err
	}
	return nil
}

func encodeStructReferenceWithModuleReference(builder *SemanticBuilder, identifier ast.TypeIdentifierNormalOrScoped) error {
	scoped, isScoped := identifier.(*ast.TypeIdentifierScoped)
	if isScoped {
		encodeModuleReference(builder, scoped.ModuleReference())
	}
	return builder.EncodeSymbol(identifier.FetchPositionLength().Range, "class", nil)
}

func addSemanticTokenCaseForCustomType(caseNode *decorated.CaseCustomType, builder *SemanticBuilder) error {
	keywordCase := caseNode.AstCaseCustomType().KeywordCase()
	keywordOf := caseNode.AstCaseCustomType().KeywordOf()
	if err := encodeKeyword(builder, keywordCase); err != nil {
		return err
	}

	if err := addSemanticToken(caseNode.Test(), builder); err != nil {
		return err
	}

	if err := encodeKeyword(builder, keywordOf); err != nil {
		return err
	}

	for _, consequence := range caseNode.Consequences() {
		if err := encodeEnumMember(builder, consequence.VariantReference().AstIdentifier().SomeTypeIdentifier()); err != nil {
			return err
		}

		for _, param := range consequence.Parameters() {
			if err := encodeVariable(builder, param.Identifier()); err != nil {
				return err
			}
		}

		if err := addSemanticToken(consequence.Expression(), builder); err != nil {
			return err
		}
	}

	if caseNode.DefaultCase() != nil {
		if err := addSemanticToken(caseNode.DefaultCase(), builder); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenCaseForPatternMatching(caseNode *decorated.CasePatternMatching, builder *SemanticBuilder) error {
	keywordCase := caseNode.AstCasePatternMatching().KeywordCase()
	keywordOf := caseNode.AstCasePatternMatching().KeywordOf()
	if err := encodeKeyword(builder, keywordCase); err != nil {
		return err
	}

	if err := addSemanticToken(caseNode.Test(), builder); err != nil {
		return err
	}

	if err := encodeKeyword(builder, keywordOf); err != nil {
		return err
	}

	for _, consequence := range caseNode.Consequences() {
		if err := addSemanticToken(consequence.Literal(), builder); err != nil {
			return err
		}
		if err := addSemanticToken(consequence.Expression(), builder); err != nil {
			return err
		}
	}

	if caseNode.DefaultCase() != nil {
		if err := addSemanticToken(caseNode.DefaultCase(), builder); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenGuardToken(basic ast.GuardItemBasic, builder *SemanticBuilder) error {
	guardToken := basic.GuardToken
	operatorToken := token.NewOperatorToken(guardToken.Type(), guardToken.SourceFileReference, guardToken.Raw(), guardToken.DebugString())
	return encodeOperator(builder, operatorToken)
}

func addSemanticTokenGuard(guard *decorated.Guard, builder *SemanticBuilder) error {
	for _, consequence := range guard.Items() {
		addSemanticTokenGuardToken(consequence.AstGuardItem().GuardItemBasic, builder)
		if err := addSemanticToken(consequence.Condition(), builder); err != nil {
			return err
		}
		if err := addSemanticToken(consequence.Expression(), builder); err != nil {
			return err
		}
	}

	if guard.DefaultGuard() != nil {
		addSemanticTokenGuardToken(guard.DefaultGuard().AstGuardDefault().GuardItemBasic, builder)
		if err := addSemanticToken(guard.DefaultGuard().Expression(), builder); err != nil {
			return err
		}
	}

	return nil
}

func IsBuiltInType(typeToCheckUnaliased dtype.Type) bool {
	typeToCheck := dectype.Unalias(typeToCheckUnaliased)
	switch t := typeToCheck.(type) {
	case *dectype.TypeReference:
		return IsBuiltInType(t.Next())
	case *dectype.InvokerType:
		typeToInvoke := t.TypeGenerator()
		typeRef, _ := typeToInvoke.(*dectype.TypeReference)
		if typeRef != nil {
			typeToInvoke = typeRef.Next()
		}
		typeToInvokeName := typeToInvoke.DecoratedName()
		return typeToInvokeName == "List" || typeToInvokeName == "Array"
	case *dectype.PrimitiveAtom:
		typeName := t.AtomName()
		return typeName == "Int" ||
			typeName == "Fixed" || typeName == "Bool" || typeName == "ResourceName" ||
			typeName == "TypeId" || typeName == "Blob"
	}

	return false
}

func addTypeReferenceInvoker(referenceRange token.Range, invoker *dectype.InvokerType, builder *SemanticBuilder) error {
	tokenModifiers := []string{"declaration"}
	if IsBuiltInType(invoker) {
		tokenModifiers = append(tokenModifiers, "defaultLibrary")
	}

	if err := builder.EncodeSymbol(referenceRange, "type", tokenModifiers); err != nil {
		return err
	}

	for _, param := range invoker.Params() {
		var tokenModifiersForParam []string
		if IsBuiltInType(param) {
			tokenModifiersForParam = append(tokenModifiersForParam, "defaultLibrary")
		}
		if err := builder.EncodeSymbol(param.FetchPositionLength().Range, "typeParameter", tokenModifiersForParam); err != nil {
			return err
		}
	}

	return nil
}

func addTypeReferenceCustomType(referenceRange token.Range, invoker *dectype.CustomTypeAtom, builder *SemanticBuilder) error {
	tokenModifiers := []string{"declaration"}
	if IsBuiltInType(invoker) {
		tokenModifiers = append(tokenModifiers, "defaultLibrary")
	}

	if err := builder.EncodeSymbol(referenceRange, "type", tokenModifiers); err != nil {
		return err
	}

	return nil
}

func addTypeReferenceAlias(referenceRange token.Range, alias *dectype.Alias, builder *SemanticBuilder) error {
	tokenModifiers := []string{"declaration"}

	if err := builder.EncodeSymbol(referenceRange, "type", tokenModifiers); err != nil {
		return err
	}

	return nil
}

func addTypeReferenceFunctionType(referenceRange token.Range, functionType *dectype.FunctionAtom, builder *SemanticBuilder) error {
	return nil
}

func addSemanticTokenImport(decoratedImport *decorated.ImportStatement, builder *SemanticBuilder) error {
	astImport := decoratedImport.AstImport()
	if err := encodeKeyword(builder, astImport.KeywordImport()); err != nil {
		return err
	}

	for _, segment := range decoratedImport.AstImport().ModuleName().Parts() {
		if err := builder.EncodeSymbol(segment.FetchPositionLength().Range, "namespace", nil); err != nil {
			return err
		}
	}

	if astImport.KeywordAs() != nil {
		if err := encodeKeyword(builder, *astImport.KeywordAs()); err != nil {
			return err
		}
		if err := encodeModuleAlias(builder, astImport.Alias()); err != nil {
			return err
		}

	}

	if astImport.KeywordExposing() != nil {
		if err := encodeKeyword(builder, *astImport.KeywordExposing()); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenLet(decoratedLet *decorated.Let, builder *SemanticBuilder) error {
	keyword := decoratedLet.AstLet().Keyword()
	if err := builder.EncodeSymbol(keyword.FetchPositionLength().Range, "keyword", nil); err != nil {
		return err
	}

	for _, assignment := range decoratedLet.Assignments() {
		if assignment.LetVariable().Comment() != nil {
			encodeComment(builder, assignment.LetVariable().Comment().Token().CommentToken)
		}
		if err := builder.EncodeSymbol(assignment.LetVariable().FetchPositionLength().Range, "variable", []string{"readonly"}); err != nil {
			return err
		}

		addSemanticToken(assignment.Expression(), builder)
	}

	if err := builder.EncodeSymbol(decoratedLet.AstLet().InKeyword().FetchPositionLength().Range, "keyword", nil); err != nil {
		return err
	}

	return addSemanticToken(decoratedLet.Consequence(), builder)
}

func addSemanticTokenIf(decoratedIf *decorated.If, builder *SemanticBuilder) error {
	encodeKeyword(builder, decoratedIf.AstIf().KeywordIf())
	if err := addSemanticToken(decoratedIf.Condition(), builder); err != nil {
		return err
	}
	encodeKeyword(builder, decoratedIf.AstIf().KeywordThen())
	if err := addSemanticToken(decoratedIf.Consequence(), builder); err != nil {
		return err
	}
	encodeKeyword(builder, decoratedIf.AstIf().KeywordElse())
	if err := addSemanticToken(decoratedIf.Alternative(), builder); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenMultilineComment(decoratedMultilineComment *decorated.MultilineComment, builder *SemanticBuilder) error {
	return encodeMultilineComment(builder, decoratedMultilineComment.AstComment().Token())
}

func addSemanticTokenBinaryOperator(operator *decorated.BinaryOperator, builder *SemanticBuilder) error {
	if err := addSemanticToken(operator.Left(), builder); err != nil {
		return err
	}

	if err := addSemanticToken(operator.Right(), builder); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenUnaryOperator(operator *decorated.UnaryOperator, builder *SemanticBuilder) error {
	if err := addSemanticToken(operator.Left(), builder); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenRecordLiteral(recordLiteral *decorated.RecordLiteral, builder *SemanticBuilder) error {
	if recordLiteral.RecordTemplate() != nil {
		if err := addSemanticToken(recordLiteral.RecordTemplate(), builder); err != nil {
			return err
		}
	}

	for _, assignment := range recordLiteral.ParseOrderedAssignments() {
		if err := builder.EncodeSymbol(assignment.FieldName().FetchPositionLength().Range, "property", nil); err != nil {
			return err
		}

		addSemanticToken(assignment.Expression(), builder)
	}

	return nil
}

func addSemanticTokenLetVariableReference(letVarReference *decorated.LetVariableReference, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(letVarReference.FetchPositionLength().Range, "variable", []string{"readonly"}); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenString(stringLiteral *decorated.StringLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(stringLiteral.FetchPositionLength().Range, "string", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenResourceNameLiteral(resourceNameLiteral *decorated.ResourceNameLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(resourceNameLiteral.FetchPositionLength().Range, "operator", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenBooleanLiteral(stringLiteral *decorated.BooleanLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(stringLiteral.FetchPositionLength().Range, "number", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenStringInterpolation(stringInterpolation *decorated.StringInterpolation, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(stringInterpolation.FetchPositionLength().Range, "string", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenTypeId(typeId *decorated.TypeIdLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(typeId.FetchPositionLength().Range, "macro", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenChar(stringLiteral *decorated.CharacterLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(stringLiteral.FetchPositionLength().Range, "string", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenFixed(fixedLiteral *decorated.FixedLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(fixedLiteral.FetchPositionLength().Range, "number", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenInteger(integerLiteral *decorated.IntegerLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(integerLiteral.FetchPositionLength().Range, "number", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenListLiteral(listLiteral *decorated.ListLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(listLiteral.AstListLiteral().StartParenToken().Range, "operator", nil); err != nil {
		return err
	}

	for _, expression := range listLiteral.Expressions() {
		if err := addSemanticToken(expression, builder); err != nil {
			return err
		}
	}

	if err := builder.EncodeSymbol(listLiteral.AstListLiteral().EndParenToken().Range, "operator", nil); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenArrayLiteral(arrayLiteral *decorated.ArrayLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(arrayLiteral.AstArrayLiteral().StartParenToken().Range, "operator", nil); err != nil {
		return err
	}

	for _, expression := range arrayLiteral.Expressions() {
		if err := addSemanticToken(expression, builder); err != nil {
			return err
		}
	}

	if err := builder.EncodeSymbol(arrayLiteral.AstArrayLiteral().EndParenToken().Range, "operator", nil); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenFunctionCall(funcCall *decorated.FunctionCall, builder *SemanticBuilder) error {
	if err := addSemanticToken(funcCall.FunctionExpression(), builder); err != nil {
		return err
	}

	for _, argument := range funcCall.Arguments() {
		if err := addSemanticToken(argument, builder); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenCurryFunction(funcCall *decorated.CurryFunction, builder *SemanticBuilder) error {
	if err := addSemanticToken(funcCall.FunctionValue(), builder); err != nil {
		return err
	}

	for _, argument := range funcCall.ArgumentsToSave() {
		if err := addSemanticToken(argument, builder); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenFunctionReference(functionReference *decorated.FunctionReference, builder *SemanticBuilder) error {
	isScoped := functionReference.NameReference().ModuleReference() != nil
	if isScoped {
		encodeModuleReference(builder, functionReference.NameReference().ModuleReference().AstModuleReference())
	}

	if err := builder.EncodeSymbol(functionReference.Identifier().FetchPositionLength().Range, "function", nil); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenFunctionParameterReference(parameter *decorated.FunctionParameterReference, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(parameter.Identifier().FetchPositionLength().Range, "parameter", nil); err != nil {
		return err
	}

	return nil
}

func typeReferenceHelper(next dtype.Type, referenceRange token.Range, builder *SemanticBuilder) error {
	switch t := next.(type) {
	case *dectype.PrimitiveAtom:
		return addTypeReferencePrimitive(referenceRange, t, builder)
	case *dectype.InvokerType:
		return addTypeReferenceInvoker(referenceRange, t, builder)
	case *dectype.Alias:
		return addTypeReferenceAlias(referenceRange, t, builder)
	case *dectype.CustomTypeAtom:
		return addTypeReferenceCustomType(referenceRange, t, builder)
	case *dectype.FunctionAtom:
		return addTypeReferenceFunctionType(referenceRange, t, builder)
	}
	log.Printf("unhandled typeReference %T %v\n", next, next)

	return fmt.Errorf("unhandled typeReference %T %v\n", next, next)
}

func addSemanticTokenTypeReference(typeReference *dectype.TypeReference, builder *SemanticBuilder) error {
	referenceRange := typeReference.FetchPositionLength().Range
	next := typeReference.Next()

	return typeReferenceHelper(next, referenceRange, builder)
}

func addSemanticTokenTypeReferenceScoped(typeReference *dectype.TypeReferenceScoped, builder *SemanticBuilder) error {
	if err := encodeModuleReference(builder, typeReference.TypeIdentifierScoped().ModuleReference()); err != nil {
		return err
	}
	referenceRange := typeReference.FetchPositionLength().Range
	return typeReferenceHelper(typeReference.Next(), referenceRange, builder)
}

func typeReferenceParamHelper(param dtype.Type, builder *SemanticBuilder) error {
	switch t := param.(type) {
	case *dectype.FunctionTypeReference:
		return addSemanticTokenFunctionTypeReference(t, builder)
	case *dectype.TypeReferenceScoped:
		return addSemanticTokenTypeReferenceScoped(t, builder)
	case *dectype.TypeReference:
		return addSemanticTokenTypeReference(t, builder)
	case *dectype.InvokerType:
		return addSemanticToken(t, builder)
	case *dectype.FunctionAtom:
		return addSemanticToken(t, builder)
	case *dectype.RecordAtom:
		return addSemanticToken(t, builder)
	}
	log.Printf("unknown param %T", param)
	return fmt.Errorf("unknown param %T", param)
}

func addSemanticTokenFunctionTypeReference(typeReference *dectype.FunctionTypeReference, builder *SemanticBuilder) error {
	for _, param := range typeReference.FunctionAtom().FunctionParameterTypes() {
		if err := typeReferenceParamHelper(param, builder); err != nil {
			return err
		}
	}
	/*
		referenceRange := typeReference.FetchPositionLength().Range
		switch t := next.(type) {
		case *dectype.PrimitiveAtom:
			return addTypeReferencePrimitive(referenceRange, t, builder)
		case *dectype.InvokerType:
			return addTypeReferenceInvoker(referenceRange, t, builder)
		case *dectype.FunctionAtom:
			return addTypeReferenceFunctionType(referenceRange, t, builder)
		}
	*/

	// log.Printf("unhandled function typeReference %T %v\n", next, next)

	return typeReferenceHelper(typeReference.Next(), typeReference.FetchPositionLength().Range, builder)
}

func addSemanticToken(typeOrToken decorated.TypeOrToken, builder *SemanticBuilder) error {
	// log.Printf("addSemantic for %T", typeOrToken)
	switch t := typeOrToken.(type) {
	case *decorated.NamedFunctionValue:
		return addSemanticTokenNamedFunctionValue(t, builder)
	case *decorated.FunctionValue:
		return addSemanticTokenFunctionValue(t, builder)
	case *decorated.AnnotationStatement:
		return addSemanticTokenAnnotation(t, builder)
	case *dectype.InvokerType:
		return addTypeReferenceInvoker(t.FetchPositionLength().Range, t, builder)
	case *decorated.ImportStatement:
		return addSemanticTokenImport(t, builder)
	case *decorated.Let:
		return addSemanticTokenLet(t, builder)
	case *decorated.If:
		return addSemanticTokenIf(t, builder)
	case *decorated.MultilineComment:
		return addSemanticTokenMultilineComment(t, builder)
	case *decorated.BinaryOperator:
		return addSemanticTokenBinaryOperator(t, builder)
	case *decorated.UnaryOperator:
		return addSemanticTokenUnaryOperator(t, builder)
	case *decorated.RecordLiteral:
		return addSemanticTokenRecordLiteral(t, builder)
	case *decorated.ResourceNameLiteral:
		return addSemanticTokenResourceNameLiteral(t, builder)
	case *decorated.LetVariableReference:
		return addSemanticTokenLetVariableReference(t, builder)
	case *decorated.StringLiteral:
		return addSemanticTokenString(t, builder)
	case *decorated.BooleanLiteral:
		return addSemanticTokenBooleanLiteral(t, builder)
	case *decorated.ArithmeticUnaryOperator:
		return addSemanticTokenUnaryOperator(&t.UnaryOperator, builder)
	case *decorated.StringInterpolation:
		return addSemanticTokenStringInterpolation(t, builder)
	case *decorated.TypeIdLiteral:
		return addSemanticTokenTypeId(t, builder)
	case *decorated.CharacterLiteral:
		return addSemanticTokenChar(t, builder)
	case *decorated.FixedLiteral:
		return addSemanticTokenFixed(t, builder)
	case *decorated.IntegerLiteral:
		return addSemanticTokenInteger(t, builder)
	case *decorated.ListLiteral:
		return addSemanticTokenListLiteral(t, builder)
	case *decorated.ArrayLiteral:
		return addSemanticTokenArrayLiteral(t, builder)
	case *decorated.FunctionCall:
		return addSemanticTokenFunctionCall(t, builder)
	case *decorated.CurryFunction:
		return addSemanticTokenCurryFunction(t, builder)
	case *decorated.FunctionName:
		return addSemanticTokenFunctionName(t, builder)
	case *decorated.FunctionReference:
		return addSemanticTokenFunctionReference(t, builder)
	case *decorated.FunctionParameterReference:
		return addSemanticTokenFunctionParameterReference(t, builder)
	case *decorated.CustomTypeVariantConstructor:
		return addSemanticTokenCustomTypeVariantConstructor(t, builder)
	case *decorated.CustomTypeVariantReference:
		return addSemanticTokenCustomTypeVariantReference(t, builder)
	case *decorated.RecordConstructor:
		return addSemanticTokenRecordConstructor(t, builder)
	case *decorated.RecordLookups:
		return addSemanticTokenRecordsLookup(t, builder)
	case *decorated.CaseCustomType:
		return addSemanticTokenCaseForCustomType(t, builder)
	case *decorated.CasePatternMatching:
		return addSemanticTokenCaseForPatternMatching(t, builder)
	case *decorated.Guard:
		return addSemanticTokenGuard(t, builder)
	case *decorated.ArithmeticOperator:
		return addSemanticToken(&t.BinaryOperator, builder)
	case *decorated.PipeLeftOperator:
		return addSemanticToken(&t.BinaryOperator, builder)
	case *decorated.PipeRightOperator:
		return addSemanticToken(&t.BinaryOperator, builder)
	case *decorated.BitwiseOperator:
		return addSemanticToken(&t.BinaryOperator, builder)
	case *decorated.BooleanOperator:
		return addSemanticToken(&t.BinaryOperator, builder)
	case *decorated.LogicalOperator:
		return addSemanticToken(&t.BinaryOperator, builder)
	case *decorated.LogicalUnaryOperator:
		return addSemanticToken(&t.UnaryOperator, builder)
	case *decorated.ConsOperator:
		return addSemanticToken(&t.BinaryOperator, builder)
	case *decorated.Constant:
		return addSemanticTokenConstant(t, builder)
	case *decorated.ExternalFunctionDeclaration:
		return nil // TODO: decorate external
	case *decorated.AsmConstant:
		return nil // TODO: decorate asm constant

		// TYPES
		//

	case *dectype.TypeReference:
		return addSemanticTokenTypeReference(t, builder)
	case *dectype.TypeReferenceScoped:
		return addSemanticTokenTypeReferenceScoped(t, builder)
	case *dectype.FunctionTypeReference:
		return addSemanticTokenFunctionTypeReference(t, builder)
	case *dectype.FunctionAtom:
		return addSemanticTokenFunctionType(t, builder)
	case *dectype.Alias:
		return addSemanticTokenTypeAlias(t, builder)
	case *dectype.RecordAtom:
		return addSemanticTokenRecordType(t, builder)
	case *dectype.CustomTypeAtom:
		return addSemanticTokenCustomType(t, builder)

	default:
		log.Printf("semantic unhandled %T\n", t)
	}

	return nil
}
