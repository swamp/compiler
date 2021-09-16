package generate_sp

import (
	"fmt"
	"log"
	"strings"

	"github.com/swamp/compiler/src/assembler_sp"
)

type ConstantType uint

const (
	ConstantTypeString ConstantType = iota
	ConstantTypeBoolean
	ConstantTypeInteger
	ConstantTypeResourceName
	ConstantTypeFunction
	ConstantTypeFunctionExternal
)

type Constant struct {
	posRange     assembler_sp.ConstantPosAndRange
	constantType ConstantType
	str          string
	b            bool
	integer      int32
}

func (v *Constant) ConstantType() ConstantType {
	return v.constantType
}

func (v *Constant) IntegerValue() int32 {
	return v.integer
}

func (v *Constant) StringValue() string {
	return v.str
}

func (v *Constant) BooleanValue() bool {
	return v.b
}

func (v *Constant) FunctionReferenceFullyQualifiedName() string {
	return v.str
}

func NewStringConstant(posRange assembler_sp.ConstantPosAndRange, str string) *Constant {
	return &Constant{posRange: posRange, constantType: ConstantTypeString, str: str}
}

func NewIntegerConstant(posRange assembler_sp.ConstantPosAndRange, i int32) *Constant {
	return &Constant{posRange: posRange, constantType: ConstantTypeInteger, integer: i}
}

func NewResourceNameConstant(posRange assembler_sp.ConstantPosAndRange, str string) *Constant {
	return &Constant{posRange: posRange, constantType: ConstantTypeResourceName, str: str}
}

func NewFunctionReferenceConstantWithDebug(posRange assembler_sp.ConstantPosAndRange, uniqueFullyQualifiedName string) *Constant {
	return &Constant{posRange: posRange, constantType: ConstantTypeFunction, str: uniqueFullyQualifiedName}
}

func NewExternalFunctionReferenceConstantWithDebug(posRange assembler_sp.ConstantPosAndRange, uniqueFullyQualifiedName string) *Constant {
	return &Constant{posRange: posRange, constantType: ConstantTypeFunctionExternal, str: uniqueFullyQualifiedName}
}

func NewBooleanConstant(posRange assembler_sp.ConstantPosAndRange, b bool) *Constant {
	return &Constant{posRange: posRange, constantType: ConstantTypeBoolean, b: b}
}

type StartMemoryLayout struct {
	pointer   uint32
	posRanges []assembler_sp.ConstantPosAndRange
}

func NewStartMemoryLayout() *StartMemoryLayout {
	return &StartMemoryLayout{}
}

func (c *StartMemoryLayout) AllocateSpace(octetSize uint32, debugString string) assembler_sp.ConstantPosAndRange {
	posRange := assembler_sp.ConstantPosAndRange{Pos: c.pointer, Size: octetSize, DebugString: debugString}
	c.pointer += octetSize
	if (c.pointer % 8) != 0 {
		c.pointer += 8 - (c.pointer % 8)
	}

	c.posRanges = append(c.posRanges, posRange)

	return posRange
}

func (c *StartMemoryLayout) AllRanges() []assembler_sp.ConstantPosAndRange {
	return c.posRanges
}

func (c *StartMemoryLayout) DebugOutput() {
	lastPointer := uint32(0)
	log.Printf("constants memory layout\n")
	for _, posRange := range c.posRanges {
		if posRange.Pos < lastPointer {
			panic("problem with order")
		}
		lastPointer = posRange.Pos + posRange.Size

		log.Printf("%04d: %d > '%s'\n", posRange.Pos, posRange.Size, posRange.DebugString)
	}
}

type StartMemoryConstants struct {
	constants    []*Constant
	memoryLayout *StartMemoryLayout
}

func NewStartMemoryConstants() *StartMemoryConstants {
	return &StartMemoryConstants{
		memoryLayout: NewStartMemoryLayout(),
	}
}

func (c *StartMemoryConstants) String() string {
	s := "\n"
	for _, constant := range c.constants {
		if constant == nil {
			panic("swamp assembler: nil constant")
		}
		s += fmt.Sprintf("%v\n", constant)
	}
	return strings.TrimSpace(s)
}

func (c *StartMemoryConstants) Constants() []*Constant {
	return c.constants
}

func (c *StartMemoryConstants) AllocateSpace(size uint32, debugString string) assembler_sp.ConstantPosAndRange {
	return c.memoryLayout.AllocateSpace(size, debugString)
}

func (c *StartMemoryConstants) CopyConstants(constants []*Constant) {
	for _, constantToCopy := range constants {
		c.constants = append(c.constants, constantToCopy)
	}
}

func (c *StartMemoryConstants) AllocateStringConstant(s string) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeString {
			if constant.str == s {
				return constant
			}
		}
	}
	posRange := c.AllocateSpace(uint32(len(s)+1), "string")
	newConstant := NewStringConstant(posRange, s)
	c.constants = append(c.constants, newConstant)

	return newConstant
}

func (c *StartMemoryConstants) AllocateIntegerConstant(i int32) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeInteger {
			if constant.integer == i {
				return constant
			}
		}
	}
	posRange := c.AllocateSpace(4, "int32")
	newConstant := NewIntegerConstant(posRange, i)
	c.constants = append(c.constants, newConstant)

	return newConstant
}

func (c *StartMemoryConstants) AllocateResourceNameConstant(name string) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeResourceName {
			if constant.str == name {
				return constant
			}
		}
	}
	posRange := c.AllocateSpace(uint32(len(name)+1), "resource name")
	newConstant := NewResourceNameConstant(posRange, name)
	c.constants = append(c.constants, newConstant)

	return newConstant
}

func (c *StartMemoryConstants) AllocateBooleanConstant(t bool) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeBoolean {
			if constant.b == t {
				return constant
			}
		}
	}
	posRange := c.AllocateSpace(1, "boolean")
	newConstant := NewBooleanConstant(posRange, t)
	c.constants = append(c.constants, newConstant)
	return newConstant
}

func (c *StartMemoryConstants) AllocateFunctionReferenceConstant(uniqueFullyQualifiedFunctionName string) (*Constant, error) {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeFunction {
			if constant.str == uniqueFullyQualifiedFunctionName {
				return constant, nil
			}
		}
	}
	posRange := c.AllocateSpace(uint32(len(uniqueFullyQualifiedFunctionName)+1), "function ref constant")
	newConstant := NewFunctionReferenceConstantWithDebug(posRange, uniqueFullyQualifiedFunctionName)
	c.constants = append(c.constants, newConstant)

	return newConstant, nil
}

func (c *StartMemoryConstants) AllocateExternalFunctionReferenceConstant(uniqueFullyQualifiedFunctionName string) (*Constant, error) {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeFunctionExternal {
			if constant.str == uniqueFullyQualifiedFunctionName {
				return constant, nil
			}
		}
	}
	posRange := c.AllocateSpace(uint32(len(uniqueFullyQualifiedFunctionName)+1), "external func")
	newConstant := NewExternalFunctionReferenceConstantWithDebug(posRange, uniqueFullyQualifiedFunctionName)
	c.constants = append(c.constants, newConstant)

	return newConstant, nil
}

/*
func (c *StartMemoryConstants) findFunc(identifier *VariableName) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeFunction {
			if constant.str == identifier.Name() {
				return constant
			}
		}
	}

	return nil
}
*/

func (c *StartMemoryConstants) FindStringConstant(s string) *Constant {
	for _, constant := range c.constants {
		if constant.constantType == ConstantTypeString {
			if constant.str == s {
				return constant
			}
		}
	}
	return nil
}
