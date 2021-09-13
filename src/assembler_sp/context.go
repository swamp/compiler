/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"
	"strings"

	swampopcodetype "github.com/swamp/opcodes/type"
)

type FunctionContextConstants struct {
	constants             []*Constant
	someConstantIDCounter int
}

func (c *FunctionContextConstants) String() string {
	s := "\n"
	for _, constant := range c.constants {
		if constant == nil {
			panic("swamp assembler: nil constant")
		}
		s += fmt.Sprintf("%v\n", constant)
	}
	return strings.TrimSpace(s)
}

func (c *FunctionContextConstants) Constants() []*Constant {
	return c.constants
}

func (c *FunctionContextConstants) CopyConstants(constants []*Constant) {
	for _, constantToCopy := range constants {
		c.constants = append(c.constants, constantToCopy)
	}
}

func (c *FunctionContextConstants) AllocateStringConstant(s string) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeString {
			if constant.str == s {
				return constant
			}
		}
	}
	c.someConstantIDCounter++
	newConstant := NewStringConstant(c.someConstantIDCounter, s)
	c.constants = append(c.constants, newConstant)

	return newConstant
}

func (c *FunctionContextConstants) AllocateIntegerConstant(i int32) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeInteger {
			if constant.integer == i {
				return constant
			}
		}
	}
	c.someConstantIDCounter++
	newConstant := NewIntegerConstant(c.someConstantIDCounter, i)
	c.constants = append(c.constants, newConstant)

	return newConstant
}

func (c *FunctionContextConstants) AllocateResourceNameConstant(name string) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeResourceName {
			if constant.str == name {
				return constant
			}
		}
	}
	c.someConstantIDCounter++
	newConstant := NewResourceNameConstant(c.someConstantIDCounter, name)
	c.constants = append(c.constants, newConstant)

	return newConstant
}

func (c *FunctionContextConstants) AllocateBooleanConstant(t bool) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeBoolean {
			if constant.b == t {
				return constant
			}
		}
	}
	c.someConstantIDCounter++
	newConstant := NewBooleanConstant(c.someConstantIDCounter, t)
	c.constants = append(c.constants, newConstant)
	return newConstant
}

func (c *FunctionContextConstants) AllocateFunctionReferenceConstant(uniqueFullyQualifiedFunctionName string) (*Constant, error) {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeFunction {
			if constant.str == uniqueFullyQualifiedFunctionName {
				return constant, nil
			}
		}
	}
	c.someConstantIDCounter++
	newConstant := NewFunctionReferenceConstantWithDebug(c.someConstantIDCounter, uniqueFullyQualifiedFunctionName)
	c.constants = append(c.constants, newConstant)

	return newConstant, nil
}

func (c *FunctionContextConstants) AllocateExternalFunctionReferenceConstant(uniqueFullyQualifiedFunctionName string) (*Constant, error) {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeFunctionExternal {
			if constant.str == uniqueFullyQualifiedFunctionName {
				return constant, nil
			}
		}
	}
	c.someConstantIDCounter++
	newConstant := NewExternalFunctionReferenceConstantWithDebug(c.someConstantIDCounter, uniqueFullyQualifiedFunctionName)
	c.constants = append(c.constants, newConstant)

	return newConstant, nil
}

func (c *FunctionContextConstants) findFunc(identifier *VariableName) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeFunction {
			if constant.str == identifier.Name() {
				return constant
			}
		}
	}
	/*
		if c.parent != nil {
			return c.parent.findFunc(identifier)
		}
	*/
	return nil
}

func (c *FunctionContextConstants) FindStringConstant(s string) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeString {
			if constant.str == s {
				return constant
			}
		}
	}
	return nil
}

type FunctionRegisterLayout struct {
	mappedVariables          []*VariableImpl
	returnVariable           *VariableImpl
	someIDCounter            int
	constants                *FunctionContextConstants
	highestUsedRegisterValue uint8
}

func NewFunctionRegisterLayout() *FunctionRegisterLayout {
	return &FunctionRegisterLayout{mappedVariables: make([]*VariableImpl, 256)}
}

func (c *FunctionRegisterLayout) findFirstFreeIndex() uint8 {
	for index, v := range c.mappedVariables {
		if v == nil {
			return uint8(index)
		}
	}
	panic("swamp assembler: out of variable space")
}

func (c *FunctionRegisterLayout) findHighestUsedIndex() uint8 {
	rememberIndex := uint8(0)
	for index, v := range c.mappedVariables {
		if v != nil {
			rememberIndex = uint8(index)
		}
	}
	return rememberIndex
}

func (c *FunctionRegisterLayout) HighestUsedRegisterValue() uint8 {
	return c.highestUsedRegisterValue
}

func (c *FunctionRegisterLayout) RegisterCountUsedMax() uint8 {
	return c.highestUsedRegisterValue + 1
}

func (c *FunctionRegisterLayout) addVariable(v *VariableImpl) {
	registerValue := c.findFirstFreeIndex()
	if registerValue > c.highestUsedRegisterValue {
		c.highestUsedRegisterValue = registerValue
	}
	register := swampopcodetype.NewRegister(registerValue)
	v.SetRegister(register)
	c.mappedVariables[registerValue] = v
}

func (c *FunctionRegisterLayout) AddVariable(context *Context, identifier *VariableName) *VariableImpl {
	c.someIDCounter++
	v := NewVariable(context, c.someIDCounter, identifier)
	c.addVariable(v)
	return v
}

func (c *FunctionRegisterLayout) AddTempVariable(context *Context, debugName string) *VariableImpl {
	c.someIDCounter++
	v := NewTempVariable(context, c.someIDCounter, debugName)
	c.addVariable(v)
	return v
}

func (c *FunctionRegisterLayout) String() string {
	s := "\n"
	for index, variable := range c.mappedVariables {
		if variable == nil {
			break
		}
		registerIndex := uint8(index)
		s += fmt.Sprintf("%v: %v\n", registerIndex, variable)
	}

	if c.constants != nil {
		s += c.constants.String()
	}

	return strings.TrimSpace(s)
}

func (c *FunctionRegisterLayout) ShowSummary() {
	fmt.Printf("-- register summary --\n%v\n", c)
}

type Context struct {
	nameToVariable map[string]*VariableImpl
	variables      []*VariableImpl
	tempVariables  []*VariableImpl
	parent         *Context
	root           *FunctionRootContext
	constants      *FunctionContextConstants
	layouter       *FunctionRegisterLayout
}

func (c *Context) MakeScopeContext() *Context {
	newContext := &Context{
		nameToVariable: make(map[string]*VariableImpl),
		root:           c.root, parent: c, constants: c.constants, layouter: c.layouter,
	}
	newContext.parent = c
	return newContext
}

func (c *Context) Parent() *Context {
	return c.parent
}

func (r *Context) Constants() *FunctionContextConstants {
	return r.constants
}

func (c *Context) Free() {
}

func (c *Context) AllocateVariable(identifier *VariableName) *VariableImpl {
	v := c.layouter.AddVariable(c, identifier)

	found := c.nameToVariable[identifier.Name()]
	if found != nil {
		panic(fmt.Sprintf("swamp assembler: tried to add multiple variable of name %v", identifier.Name()))
	}

	c.nameToVariable[identifier.Name()] = v
	c.variables = append(c.variables, v)
	return v
}

func (c *Context) AllocateTempVariable(tempName string) *VariableImpl {
	v := c.layouter.AddTempVariable(c, tempName)
	c.tempVariables = append(c.tempVariables, v)
	return v
}

func (c *Context) FindVariable(identifier *VariableName) *VariableImpl {
	found := c.nameToVariable[identifier.Name()]
	if found == nil && c.parent != nil {
		return c.parent.FindVariable(identifier)
	}
	return found
}

func (c *Context) String() string {
	s := "\n"
	for _, variable := range c.variables {
		s += fmt.Sprintf("%v\n", variable)
	}
	for _, tempVariable := range c.tempVariables {
		s += fmt.Sprintf("%v\n", tempVariable)
	}
	return strings.TrimSpace(s)
}

func (c *Context) ShowSummary() {
	fmt.Printf("---------- Variables ------------\n")
	for _, variable := range c.variables {
		fmt.Printf("%v\n", variable)
	}
	fmt.Printf("---------- Constants ------------\n")
	fmt.Printf("%v\n", c.constants)
	fmt.Printf("---------------------------------\n")
}

type FunctionRootContext struct {
	returnVariable *VariableImpl
	constants      *FunctionContextConstants
	layouter       *FunctionRegisterLayout
	scopeContext   *Context
}

func NewFunctionRootContext() *FunctionRootContext {
	c := &FunctionRootContext{constants: &FunctionContextConstants{}, layouter: NewFunctionRegisterLayout()}
	bootstrap := &Context{root: c, constants: c.constants, layouter: c.layouter}
	c.scopeContext = bootstrap.MakeScopeContext()
	c.allocateReturnVariable("return")
	return c
}

func (c *FunctionRootContext) allocateReturnVariable(tempName string) *VariableImpl {
	c.returnVariable = c.layouter.AddTempVariable(c.scopeContext, tempName)
	if c.returnVariable.Register().Value() != 0 {
		panic("swamp assembler: layouter register value is zero")
	}
	return c.returnVariable
}

func (r *FunctionRootContext) ReturnVariable() *VariableImpl {
	return r.returnVariable
}

func (r *FunctionRootContext) ScopeContext() *Context {
	return r.scopeContext
}

func (r *FunctionRootContext) Layouter() *FunctionRegisterLayout {
	return r.layouter
}

func (r *FunctionRootContext) Constants() *FunctionContextConstants {
	return r.constants
}

func (r *FunctionRootContext) ShowSummary() {
	r.layouter.ShowSummary()
}

func (r *FunctionRootContext) String() string {
	return r.constants.String() + r.scopeContext.String()
}
