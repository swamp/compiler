package assembler_sp

import (
	"fmt"
)

type FunctionVariables struct {
	nameToVariable map[string]*VariableImpl
}

func NewFunctionVariables() *FunctionVariables {
	return &FunctionVariables{nameToVariable: make(map[string]*VariableImpl)}
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
		return SourceStackPosRange{}, fmt.Errorf("could not find variable %v", name)
	}

	return existingVariable.source, nil
}
