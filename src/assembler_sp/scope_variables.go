package assembler_sp

import (
	"fmt"
)

type ScopeVariables struct {
	parent         *ScopeVariables
	nameToVariable map[string]*VariableImpl
}

func NewFunctionVariables() *ScopeVariables {
	return &ScopeVariables{nameToVariable: make(map[string]*VariableImpl)}
}

func NewFunctionVariablesWithParent(parent *ScopeVariables) *ScopeVariables {
	return &ScopeVariables{nameToVariable: make(map[string]*VariableImpl), parent: parent}
}

func (c *ScopeVariables) DefineVariable(name string, posRange SourceStackPosRange) {
	if uint(posRange.Size) == 0 {
		panic(fmt.Errorf("octet size zero is not allowed for allocate stack memory"))
	}
	_, alreadyHas := c.nameToVariable[name]
	if alreadyHas {
		panic("cannot define variable again")
	}

	v := NewVariable(VariableName(name), posRange)

	c.nameToVariable[name] = v
}

func (c *ScopeVariables) FindVariable(name string) (SourceStackPosRange, error) {
	existingVariable, alreadyHas := c.nameToVariable[name]
	if !alreadyHas {
		if c.parent != nil {
			return c.parent.FindVariable(name)
		}
		return SourceStackPosRange{}, fmt.Errorf("could not find variable %v", name)
	}

	return existingVariable.source, nil
}
