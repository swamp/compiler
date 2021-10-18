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

func (c *ScopeVariables) DefineVariable(name string, posRange SourceStackPosRange) error {
	if uint(posRange.Size) == 0 {
		return fmt.Errorf("octet size zero is not allowed for allocate stack memory")
	}
	_, alreadyHas := c.nameToVariable[name]
	if alreadyHas {
		return fmt.Errorf("cannot define variable again '%s'", name)
	}

	v := NewVariable(VariableName(name), posRange)

	c.nameToVariable[name] = v

	return nil
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
