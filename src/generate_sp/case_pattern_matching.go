package generate_sp

import (
	"fmt"

	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

type PatternMatchingType uint8

const (
	PatternMatchingTypeInt PatternMatchingType = iota
	PatternMatchingTypeString
)

/*
	var instructionMatchingType instruction_sp.PatternMatchingType
	switch matchingType {
	case PatternMatchingTypeInt:
		instructionMatchingType = instruction_sp.PatternMatchingTypeInt
	case PatternMatchingTypeString:
		instructionMatchingType = instruction_sp.PatternMatchingTypeString
	}
*/

func generateCasePatternMatchingInt(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, caseExpr *decorated.CaseForPatternMatching, matchingType PatternMatchingType, genContext *generateContext) error {
	testVar, testErr := generateExpressionWithSourceVar(code, caseExpr.Test(), genContext, "cast-test")
	if testErr != nil {
		return testErr
	}

	var consequences []*assembler_sp.CaseConsequencePatternMatchingInt
	var consequencesCodes []*assembler_sp.Code

	for _, consequence := range caseExpr.Consequences() {
		consequenceContext := genContext.MakeScopeContext()
		consequencesCode := assembler_sp.NewCode()

		intValue := int32(0)
		intLiteral, wasIntLiteral := consequence.Literal().(*decorated.IntegerLiteral)
		if wasIntLiteral {
			intValue = intLiteral.Value()
		} else {
			characterLiteral, wasCharacterLiteral := consequence.Literal().(*decorated.CharacterLiteral)
			if !wasCharacterLiteral {
				panic(fmt.Errorf("unsupported literal"))
			}
			intValue = characterLiteral.Value()
		}

		labelVariableName := assembler_sp.VariableName("a1")
		caseLabel := consequencesCode.Label(labelVariableName, "case")

		caseExprErr := generateExpression(consequencesCode, target, consequence.Expression(), false, consequenceContext)
		if caseExprErr != nil {
			return caseExprErr
		}

		asmConsequence := assembler_sp.NewCaseConsequencePatternMatchingInt(intValue, caseLabel)
		consequences = append(consequences, asmConsequence)

		consequencesCodes = append(consequencesCodes, consequencesCode)
	}

	defaultCode := assembler_sp.NewCode()
	defaultContext := genContext.MakeScopeContext()

	decoratedDefault := caseExpr.DefaultCase()
	defaultLabel := defaultCode.Label("default", "default")
	caseExprErr := generateExpression(defaultCode, target, decoratedDefault, false, defaultContext)
	if caseExprErr != nil {
		return caseExprErr
	}
	consequencesCodes = append(consequencesCodes, defaultCode)
	//		endLabel := consequencesBlockCode.Label(nil, "if-end")

	consequencesBlockCode := assembler_sp.NewCode()

	lastConsequnce := consequencesCodes[len(consequencesCodes)-1]

	labelVariableEndName := assembler_sp.VariableName("case end")
	endLabel := lastConsequnce.Label(labelVariableEndName, "caseend")

	for index, consequenceCode := range consequencesCodes {
		if index != len(consequencesCodes)-1 {
			consequenceCode.Jump(endLabel)
		}
	}

	for _, consequenceCode := range consequencesCodes {
		consequencesBlockCode.Copy(consequenceCode)
	}

	code.CasePatternMatchingInt(testVar.Pos, consequences, defaultLabel)

	code.Copy(consequencesBlockCode)

	return nil
}

func generateCasePatternMatchingMultiple(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	caseExpr *decorated.CaseForPatternMatching, genContext *generateContext) error {
	matchType := dectype.UnaliasWithResolveInvoker(caseExpr.ComparisonType())
	primitiveAtom, wasPrimitiveAtom := matchType.(*dectype.PrimitiveAtom)
	if !wasPrimitiveAtom {
		panic(fmt.Errorf("must have primitive atom"))
	}

	var matchingType PatternMatchingType
	switch primitiveAtom.PrimitiveName().Name() {
	case "Int":
		matchingType = PatternMatchingTypeInt
	case "Char":
		matchingType = PatternMatchingTypeInt
	case "String":
		matchingType = PatternMatchingTypeString
	default:
		panic(fmt.Errorf("not supported matching type %v", primitiveAtom.PrimitiveName()))
	}

	switch matchingType {
	case PatternMatchingTypeInt:
		return generateCasePatternMatchingInt(code, target, caseExpr, matchingType, genContext)
	}

	panic(fmt.Errorf("not supported pattern matching type %v", caseExpr.Type()))
}

func handleCasePatternMatchingMultiple(code *assembler_sp.Code,
	caseExpr *decorated.CaseForPatternMatching, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := allocMemoryForType(genContext.context.stackMemory, caseExpr.Type(), "casePatternMatchingResult")
	if err := generateCasePatternMatchingMultiple(code, posRange, caseExpr, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
