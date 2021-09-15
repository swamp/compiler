package assembler_sp

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Constants struct {
	constants         []*Constant
	functions         []*Constant
	externalFunctions []*Constant
	strings           []*Constant
	zeroMapper        *ZeroMemoryMapper
	dynamicMapper     *DynamicMemoryMapper
}

func ConstantsNew() *Constants {
	return &Constants{
		zeroMapper:    ZeroMemoryMapperNew(32 * 1024),
		dynamicMapper: DynamicMemoryMapperNew(128 * 1024),
	}
}

func (c *Constants) String() string {
	s := "\n"
	for _, constant := range c.constants {
		if constant == nil {
			panic("swamp assembler: nil constant")
		}
		s += fmt.Sprintf("%v\n", constant)
	}
	return strings.TrimSpace(s)
}

func (c *Constants) Constants() []*Constant {
	return c.constants
}

func (c *Constants) CopyConstants(constants []*Constant) {
	for _, constantToCopy := range constants {
		c.constants = append(c.constants, constantToCopy)
	}
}

// typedef struct SwampString {
//    const char* characters; // 8 octets
//    size_t characterCount; // 8 octets
// } SwampString;

const SizeofSwampString = 16

func (c *Constants) AllocateStringConstant(s string) *Constant {
	for _, constant := range c.strings {
		if constant.str == s {
			return constant
		}
	}

	stringOctets := []byte(s)
	stringOctets = append(stringOctets, byte(0))
	stringOctetsPointer := c.dynamicMapper.Write(stringOctets, "string")

	var swampStringOctets [SizeofSwampString]byte
	binary.LittleEndian.PutUint64(swampStringOctets[0:8], uint64(stringOctetsPointer.Position))
	binary.LittleEndian.PutUint64(swampStringOctets[8:16], uint64(len(s)))

	swampStringPointer := c.zeroMapper.Write(swampStringOctets[:], "swampStringOctets")

	newConstant := NewStringConstant("string", s, swampStringPointer)
	c.constants = append(c.constants, newConstant)
	c.strings = append(c.strings, newConstant)

	return newConstant
}

func intValue(memory *ZeroMemoryMapper, pos SourceZeroMemoryPos) int32 {
	posRange := SourceZeroMemoryPosRange{
		Position: pos,
		Size:     4,
	}
	return int32(binary.LittleEndian.Uint32(memory.Read(posRange)))
}

/*
func (c *Constants) AllocateResourceNameConstant(name string) *Constant {
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

*/

/*

typedef struct SwampFunc {
    size_t curryOctetSize;
    const uint8_t* curryOctets;
    const struct SwampFunc* curryFunction;

    size_t parameterCount;
    size_t parametersOctetSize;
    const uint8_t* opcodes;
    size_t opcodeCount;

    size_t totalStackUsed;
    size_t returnOctetSize;

//    const swamp_value** constants; // or frozen variables in closure
  //  size_t constant_count;
    const char* debugName;
    uint16_t typeIndex;
} SwampFunc;
*/

const SizeofSwampFunc = 11 * 8

func (c *Constants) AllocateFunctionConstant(uniqueFullyQualifiedFunctionName string, opcodes []byte) (*Constant, error) {
	for _, constant := range c.functions {
		if constant.str == uniqueFullyQualifiedFunctionName {
			return constant, nil
		}
	}

	opcodesPointer := c.dynamicMapper.Write(opcodes, "opcodes")

	var swampStringOctets [SizeofSwampFunc]byte
	binary.LittleEndian.PutUint64(swampStringOctets[0:8], uint64(0))
	binary.LittleEndian.PutUint64(swampStringOctets[8:16], uint64(0))  // Curry Octets
	binary.LittleEndian.PutUint64(swampStringOctets[16:24], uint64(0)) // Curry fn
	binary.LittleEndian.PutUint64(swampStringOctets[24:32], uint64(0)) // parameterCount
	binary.LittleEndian.PutUint64(swampStringOctets[24:32], uint64(0)) // parameters octet size

	binary.LittleEndian.PutUint64(swampStringOctets[32:40], uint64(opcodesPointer.Position))
	binary.LittleEndian.PutUint64(swampStringOctets[40:48], uint64(opcodesPointer.Size))

	binary.LittleEndian.PutUint64(swampStringOctets[48:56], uint64(0)) // totalStackUsed
	binary.LittleEndian.PutUint64(swampStringOctets[56:64], uint64(0)) // returnOctetSize
	binary.LittleEndian.PutUint64(swampStringOctets[64:72], uint64(0)) // debugName
	binary.LittleEndian.PutUint64(swampStringOctets[72:76], uint64(0)) // typeIndex

	funcPointer := c.zeroMapper.Write(swampStringOctets[:], "fn:"+uniqueFullyQualifiedFunctionName)

	newConstant := NewFunctionReferenceConstantWithDebug("fn", uniqueFullyQualifiedFunctionName, funcPointer)
	c.constants = append(c.constants, newConstant)
	c.functions = append(c.functions, newConstant)

	return newConstant, nil
}

/*
func (c *Constants) AllocateExternalFunctionConstant(uniqueFullyQualifiedFunctionName string) (*Constant, error) {
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


*/
func (c *Constants) FindFunction(identifier VariableName) *Constant {
	for _, constant := range c.functions {
		if constant.str == string(identifier) {
			return constant
		}
	}
	/*
		if c.parent != nil {
			return c.parent.findFunc(identifier)
		}
	*/
	return nil
}

func (c *Constants) FindStringConstant(s string) *Constant {
	for _, constant := range c.strings {
		if constant.str == s {
			return constant
		}
	}
	return nil
}
