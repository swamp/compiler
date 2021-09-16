package assembler_sp

import (
	"fmt"
)

type FunctionVariables struct {
	parent         *FunctionVariables
	nameToVariable map[string]*VariableImpl
}

func NewFunctionVariables() *FunctionVariables {
	return &FunctionVariables{nameToVariable: make(map[string]*VariableImpl)}
}

func NewFunctionVariablesWithParent(parent *FunctionVariables) *FunctionVariables {
	return &FunctionVariables{nameToVariable: make(map[string]*VariableImpl), parent: parent}
}

func (c *FunctionVariables) DefineVariable(name string, posRange SourceStackPosRange) {
	_, alreadyHas := c.nameToVariable[name]
	if alreadyHas {
		panic("cannot define variable again")
	}

	v := NewVariable(VariableName(name), posRange)

	c.nameToVariable[name] = v
}

func (c *FunctionVariables) FindVariable(name string) (SourceStackPosRange, error) {
	existingVariable, alreadyHas := c.nameToVariable[name]
	if !alreadyHas {
		if c.parent != nil {
			return c.parent.FindVariable(name)
		}
		return SourceStackPosRange{}, fmt.Errorf("could not find variable %v", name)
	}

	return existingVariable.source, nil
}
