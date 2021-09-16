/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/instruction_sp"
	"github.com/swamp/compiler/src/opcode_sp"
	"github.com/swamp/compiler/src/opcode_sp_type"
	swampdisasm "github.com/swamp/disassembler/lib"
)

type VariableType uint

type CodeCommand interface {
	String() string
}

type Code struct {
	statements []CodeCommand
	labels     []*Label
}

func (c *Code) Commands() []CodeCommand {
	return c.statements
}

func (c *Code) PrintOut() {
	for _, cmd := range c.statements {
		fmt.Printf("%v\n", cmd)
	}
}

func NewCode() *Code {
	return &Code{}
}

func (c *Code) addStatement(cmd CodeCommand) {
	c.statements = append(c.statements, cmd)
}

func (c *Code) Copy(other *Code) {
	for _, cmd := range other.statements {
		lbl, _ := cmd.(*Label)
		if lbl != nil {
			c.labels = append(c.labels, lbl)
		}
		c.addStatement(cmd)
	}
}

func (c *Code) Label(identifier VariableName, debugString string) *Label {
	o := &Label{identifier: identifier, debugString: debugString}
	c.labels = append(c.labels, o)
	c.addStatement(o)
	return o
}

func (c *Code) ListAppend(target TargetStackPos, a SourceStackPos, b SourceStackPos) {
	o := &ListAppend{target: target, a: a, b: b}
	c.addStatement(o)
}

func (c *Code) StringAppend(target TargetStackPos, a SourceStackPos, b SourceStackPos) {
	o := &StringAppend{target: target, a: a, b: b}
	c.addStatement(o)
}

func (c *Code) ListConj(target TargetStackPos, item SourceStackPos, list SourceStackPos) {
	o := &ListConj{target: target, item: item, list: list}
	c.addStatement(o)
}

func (c *Code) IntBinaryOperator(target TargetStackPos, a SourceStackPos, b SourceStackPos, operator instruction_sp.BinaryOperatorType) {
	o := &IntBinaryOperator{target: target, a: a, b: b, operator: operator}
	c.addStatement(o)
}

func (c *Code) UnaryOperator(target TargetStackPos, a SourceStackPos, operator instruction_sp.UnaryOperatorType) {
	o := &UnaryOperator{target: target, a: a, operator: operator}
	c.addStatement(o)
}

func (c *Code) ListLiteral(target TargetStackPos, values []SourceStackPos, itemSize StackRange) {
	o := &ListLiteral{target: target, values: values, itemSize: itemSize}
	c.addStatement(o)
}

func (c *Code) ArrayLiteral(target TargetStackPos, values []SourceStackPos, itemSize StackRange) {
	o := &ArrayLiteral{target: target, values: values, itemSize: itemSize}
	c.addStatement(o)
}

func (c *Code) CaseEnum(test SourceStackPos, consequences []*CaseConsequence, defaultConsequence *CaseConsequence) {
	o := &Case{test: test, consequences: consequences, defaultConsequence: defaultConsequence}
	c.addStatement(o)
}

func (c *Code) CopyConstant(target TargetStackPos, source SourceDynamicMemoryPos) {
	o := &CopyConstant{target: target, source: source}
	c.addStatement(o)
}

func (c *Code) LoadInteger(target TargetStackPos, intValue int32) {
	o := &LoadInteger{target: target, intValue: intValue}
	c.addStatement(o)
}

func (c *Code) LoadRune(target TargetStackPos, runeValue uint8) {
	o := &LoadRune{target: target, rune: runeValue}
	c.addStatement(o)
}

func (c *Code) LoadBool(target TargetStackPos, boolValue bool) {
	o := &LoadBool{target: target, boolean: boolValue}
	c.addStatement(o)
}

func (c *Code) LoadZeroMemoryPointer(target TargetStackPos, zeroMemoryPointer SourceDynamicMemoryPos) {
	o := &LoadZeroMemoryPointer{target: target, sourceZeroMemory: zeroMemoryPointer}
	c.addStatement(o)
}

func (c *Code) CopyMemory(target TargetStackPos, source SourceStackPosRange) {
	o := &CopyMemory{target: target, source: source}
	c.addStatement(o)
}

func (c *Code) CasePatternMatching(test SourceStackPosRange, consequences []*CaseConsequencePatternMatching, defaultConsequence *CaseConsequencePatternMatching) {
	o := &CasePatternMatching{test: test, consequences: consequences, defaultConsequence: defaultConsequence}
	c.addStatement(o)
}

func (c *Code) BranchFalse(condition SourceStackPos, jump *Label) {
	if jump == nil {
		panic("swamp assembler: null jump")
	}
	o := &BranchFalse{condition: condition, jump: jump}
	c.addStatement(o)
}

func (c *Code) BranchTrue(condition SourceStackPos, jump *Label) {
	if jump == nil {
		panic("swamp assembler: null jump")
	}
	o := &BranchTrue{condition: condition, jump: jump}
	c.addStatement(o)
}

func (c *Code) Jump(jump *Label) {
	if jump == nil {
		panic("swamp assembler: null jump")
	}
	o := &Jump{jump: jump}
	c.addStatement(o)
}

func (c *Code) Return() {
	o := &Return{}
	c.addStatement(o)
}

func (c *Code) Call(function SourceStackPos, newBasePointer TargetStackPos) {
	o := &Call{function: function, newBasePointer: newBasePointer}
	c.addStatement(o)
}

func (c *Code) Recur() {
	o := &Recur{}
	c.addStatement(o)
}

func (c *Code) CallExternal(target TargetStackPosRange, function SourceStackPos, arguments SourceStackPos) {
	o := &CallExternal{target: target, function: function, newBasePointer: arguments}
	c.addStatement(o)
}

func (c *Code) Curry(target TargetStackPos, typeIDConstant uint16, function SourceStackPos, arguments []SourceStackPos) {
	o := &Curry{target: target, typeIDConstant: typeIDConstant, function: function, arguments: arguments}
	c.addStatement(o)
}

func targetStackPosition(pos TargetStackPos) opcode_sp_type.TargetStackPosition {
	return opcode_sp_type.TargetStackPosition(pos)
}

func sourceStackPosition(pos SourceStackPos) opcode_sp_type.SourceStackPosition {
	return opcode_sp_type.SourceStackPosition(pos)
}

func sourceDynamicMemoryPos(pos SourceDynamicMemoryPos) opcode_sp_type.SourceDynamicMemoryPosition {
	return opcode_sp_type.SourceDynamicMemoryPosition(pos)
}

func sourceStackRange(size SourceStackRange) opcode_sp_type.SourceStackRange {
	return opcode_sp_type.SourceStackRange(size)
}

func stackRange(size StackRange) opcode_sp_type.StackRange {
	return opcode_sp_type.StackRange(size)
}

func sourceStackPositionRange(pos SourceStackPosRange) opcode_sp_type.SourceStackPositionRange {
	return opcode_sp_type.SourceStackPositionRange{Position: sourceStackPosition(pos.Pos), Range: sourceStackRange(pos.Size)}
}

func writeUnaryOperator(stream *opcode_sp.Stream, operator *UnaryOperator) {
	stream.IntUnaryOperator(targetStackPosition(operator.target), operator.operator, sourceStackPosition(operator.a))
}

func writeListAppend(stream *opcode_sp.Stream, operator *ListAppend) {
	stream.ListAppend(targetStackPosition(operator.target), sourceStackPosition(operator.a), sourceStackPosition(operator.b))
}

func writeStringAppend(stream *opcode_sp.Stream, operator *StringAppend) {
	stream.StringAppend(targetStackPosition(operator.target), sourceStackPosition(operator.a), sourceStackPosition(operator.b))
}

func writeListConj(stream *opcode_sp.Stream, operator *ListConj) {
	stream.ListConj(targetStackPosition(operator.target), sourceStackPosition(operator.list), sourceStackPosition(operator.item))
}

func writeBinaryOperator(stream *opcode_sp.Stream, operator *BinaryOperator) {
	stream.BinaryOperator(targetStackPosition(operator.target), operator.operator, sourceStackPosition(operator.a), sourceStackPosition(operator.b))
}

func writeIntBinaryOperator(stream *opcode_sp.Stream, operator *IntBinaryOperator) {
	stream.IntBinaryOperator(targetStackPosition(operator.target), operator.operator, sourceStackPosition(operator.a), sourceStackPosition(operator.b))
}

func writeCopyMemory(stream *opcode_sp.Stream, operator *CopyMemory) {
	stream.MemoryCopy(targetStackPosition(operator.target), sourceStackPositionRange(operator.source))
}

func writeZeroMemoryPointer(stream *opcode_sp.Stream, operator *LoadZeroMemoryPointer) {
	stream.LoadZeroMemoryPointer(targetStackPosition(operator.target), sourceDynamicMemoryPos(operator.sourceZeroMemory))
}

func writeBranchFalse(stream *opcode_sp.Stream, branch *BranchFalse) {
	stream.BranchFalse(sourceStackPosition(branch.Condition()), branch.Jump().OpLabel())
}

func writeBranchTrue(stream *opcode_sp.Stream, branch *BranchTrue) {
	stream.BranchTrue(sourceStackPosition(branch.Condition()), branch.Jump().OpLabel())
}

func writeJump(stream *opcode_sp.Stream, jump *Jump) {
	stream.Jump(jump.Jump().OpLabel())
}

func writeCase(stream *opcode_sp.Stream, caseExpr *Case) {
	var opLabels []instruction_sp.EnumCaseJump

	for _, consequence := range caseExpr.consequences {
		label := consequence.label.OpLabel()

		caseJump := instruction_sp.NewEnumCaseJump(consequence.InternalEnumIndex(), label)
		opLabels = append(opLabels, caseJump)
	}

	defaultCons := caseExpr.defaultConsequence

	if caseExpr.defaultConsequence != nil {
		label := defaultCons.label.OpLabel()
		caseJump := instruction_sp.NewEnumCaseJump(0xff, label)
		opLabels = append(opLabels, caseJump)
	}

	stream.EnumCase(sourceStackPosition(caseExpr.test), opLabels)
}

func writeCasePatternMatching(stream *opcode_sp.Stream, caseExpr *CasePatternMatching) {
	var opLabels []instruction_sp.CasePatternMatchingJump

	for _, consequence := range caseExpr.consequences {
		label := consequence.label.OpLabel()

		caseJump := instruction_sp.NewCasePatternMatchingJump(sourceStackPosition(consequence.LiteralVariable()), label)
		opLabels = append(opLabels, caseJump)
	}

	defaultCons := caseExpr.defaultConsequence

	if caseExpr.defaultConsequence != nil {
		label := defaultCons.label.OpLabel()
		caseJump := instruction_sp.NewCasePatternMatchingJump(sourceStackPosition(caseExpr.test.Pos), label)
		opLabels = append(opLabels, caseJump)
	}

	stream.CasePatternMatching(sourceStackPositionRange(caseExpr.test), opLabels)
}

func writeList(stream *opcode_sp.Stream, listLiteral *ListLiteral) {
	var registers []opcode_sp_type.SourceStackPosition

	for _, argument := range listLiteral.values {
		registers = append(registers, sourceStackPosition(argument))
	}

	stream.CreateList(targetStackPosition(listLiteral.target), stackRange(listLiteral.itemSize), registers)
}

func writeLoadInteger(stream *opcode_sp.Stream, loadInteger *LoadInteger) {
	stream.LoadInteger(targetStackPosition(loadInteger.target), loadInteger.intValue)
}

func writeLoadBool(stream *opcode_sp.Stream, loadBool *LoadBool) {
	stream.LoadBool(targetStackPosition(loadBool.target), loadBool.boolean)
}

func writeCreateArray(stream *opcode_sp.Stream, arrayLiteral *ArrayLiteral) {
	var registers []opcode_sp_type.SourceStackPosition

	for _, argument := range arrayLiteral.values {
		registers = append(registers, sourceStackPosition(argument))
	}

	stream.CreateArray(targetStackPosition(arrayLiteral.target), stackRange(arrayLiteral.itemSize), registers)
}

func writeCallExternal(stream *opcode_sp.Stream, call *CallExternal) {
}

func writeCall(stream *opcode_sp.Stream, call *Call) {
	stream.Call(targetStackPosition(call.newBasePointer), sourceStackPosition(call.function))
}

func writeRecur(stream *opcode_sp.Stream, call *Recur) {
	stream.TailCall()
}

func writeCurry(stream *opcode_sp.Stream, call *Curry) {
	var arguments []opcode_sp_type.SourceStackPosition

	for _, argument := range call.arguments {
		arguments = append(arguments, sourceStackPosition(argument))
	}

	stream.Curry(targetStackPosition(call.target), call.typeIDConstant, sourceStackPosition(call.function), arguments)
}

func writeReturn(stream *opcode_sp.Stream) {
	stream.Return()
}

func handleStatement(cmd CodeCommand, opStream *opcode_sp.Stream) {
	switch t := cmd.(type) {
	case *BinaryOperator:
		writeBinaryOperator(opStream, t)
	case *UnaryOperator:
		writeUnaryOperator(opStream, t)
	case *BranchFalse:
		writeBranchFalse(opStream, t)
	case *BranchTrue:
		writeBranchTrue(opStream, t)
	case *Jump:
		writeJump(opStream, t)
	case *Case:
		writeCase(opStream, t)
	case *CasePatternMatching:
		writeCasePatternMatching(opStream, t)
	case *ListLiteral:
		writeList(opStream, t)
	case *LoadInteger:
		writeLoadInteger(opStream, t)
	case *LoadBool:
		writeLoadBool(opStream, t)
	case *ArrayLiteral:
		writeCreateArray(opStream, t)
	case *ListAppend:
		writeListAppend(opStream, t)
	case *ListConj:
		writeListConj(opStream, t)
	case *StringAppend:
		writeStringAppend(opStream, t)
	case *Return:
		writeReturn(opStream)
	case *Label:
		opStream.Label(t.OpLabel())
	case *Call:
		writeCall(opStream, t)
	case *Recur:
		writeRecur(opStream, t)
	case *CallExternal:
		writeCallExternal(opStream, t)
	case *Curry:
		writeCurry(opStream, t)
	case *IntBinaryOperator:
		writeIntBinaryOperator(opStream, t)
	case *CopyMemory:
		writeCopyMemory(opStream, t)
	case *LoadZeroMemoryPointer:
		writeZeroMemoryPointer(opStream, t)
	default:
		panic(fmt.Sprintf("swamp assembler: unknown cmd %v", cmd))
	}
}

func (c *Code) Resolve(verboseFlag bool) ([]byte, error) {
	if verboseFlag {
		// context.ShowSummary()
	}

	opStream := opcode_sp.NewStream()

	for _, label := range c.labels {
		opLabel := opStream.CreateLabel(label.Name())
		label.SetOpLabel(opLabel)
	}

	for _, cmd := range c.statements {
		handleStatement(cmd, opStream)
	}

	octets, err := opStream.Serialize()

	if verboseFlag {
		fmt.Println("--- disassembly ---")

		stringLines := swampdisasm.Disassemble(octets)
		for _, line := range stringLines {
			fmt.Printf("%s\n", line)
		}
	}

	return octets, err
}
