package assembler_sp

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"
)

type Constants struct {
	constants         []*Constant
	functions         []*Constant
	externalFunctions []*Constant
	strings           []*Constant
	dynamicMapper     *DynamicMemoryMapper
}

func NewConstants() *Constants {
	return &Constants{
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

func (c *Constants) DynamicMemory() *DynamicMemoryMapper {
	return c.dynamicMapper
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
	stringOctetsPointer := c.dynamicMapper.Write(stringOctets, "string:"+s)

	var swampStringOctets [SizeofSwampString]byte
	binary.LittleEndian.PutUint64(swampStringOctets[0:8], uint64(stringOctetsPointer.Position))
	binary.LittleEndian.PutUint64(swampStringOctets[8:16], uint64(len(s)))

	swampStringPointer := c.dynamicMapper.Write(swampStringOctets[:], "SwampString struct (character-pointer, characterCount) for:"+s)

	newConstant := NewStringConstant("string", s, swampStringPointer)
	c.constants = append(c.constants, newConstant)
	c.strings = append(c.strings, newConstant)

	return newConstant
}

func intValue(memory *DynamicMemoryMapper, pos SourceDynamicMemoryPos) int32 {
	posRange := SourceDynamicMemoryPosRange{
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

func (c *Constants) AllocateFunctionStruct(uniqueFullyQualifiedFunctionName string,
	opcodesPointer SourceDynamicMemoryPosRange) (*Constant, error) {
	var swampFuncStruct [SizeofSwampFunc]byte

	binary.LittleEndian.PutUint64(swampFuncStruct[0:8], uint64(0))
	binary.LittleEndian.PutUint64(swampFuncStruct[8:16], uint64(0))  // Curry Octets
	binary.LittleEndian.PutUint64(swampFuncStruct[16:24], uint64(0)) // Curry fn
	binary.LittleEndian.PutUint64(swampFuncStruct[24:32], uint64(0)) // parameterCount
	binary.LittleEndian.PutUint64(swampFuncStruct[24:32], uint64(0)) // parameters octet size

	binary.LittleEndian.PutUint64(swampFuncStruct[32:40], uint64(opcodesPointer.Position))
	binary.LittleEndian.PutUint64(swampFuncStruct[40:48], uint64(opcodesPointer.Size))

	binary.LittleEndian.PutUint64(swampFuncStruct[48:56], uint64(0)) // totalStackUsed
	binary.LittleEndian.PutUint64(swampFuncStruct[56:64], uint64(0)) // returnOctetSize
	binary.LittleEndian.PutUint64(swampFuncStruct[64:72], uint64(0)) // debugName
	binary.LittleEndian.PutUint64(swampFuncStruct[72:80], uint64(0)) // typeIndex

	funcPointer := c.dynamicMapper.Write(swampFuncStruct[:], "function Struct for:"+uniqueFullyQualifiedFunctionName)

	newConstant := NewFunctionReferenceConstantWithDebug("fn", uniqueFullyQualifiedFunctionName, funcPointer)
	c.constants = append(c.constants, newConstant)
	c.functions = append(c.functions, newConstant)

	return newConstant, nil
}

const SwampFuncOpcodeOffset = 40

func (c *Constants) FetchOpcodes(functionConstant *Constant) []byte {
	readSection := SourceDynamicMemoryPosRange{
		Position: SourceDynamicMemoryPos(uint(functionConstant.source.Position + SwampFuncOpcodeOffset)),
		Size:     DynamicMemoryRange(8 + 8),
	}
	opcodePointerAndSize := c.dynamicMapper.Read(readSection)
	opcodePosition := binary.LittleEndian.Uint64(opcodePointerAndSize[0:8])
	opcodeSize := binary.LittleEndian.Uint64(opcodePointerAndSize[8:16])

	readOpcodeSection := SourceDynamicMemoryPosRange{
		Position: SourceDynamicMemoryPos(opcodePosition),
		Size:     DynamicMemoryRange(opcodeSize),
	}

	return c.dynamicMapper.Read(readOpcodeSection)
}

func (c *Constants) AllocatePrepareFunctionConstant(uniqueFullyQualifiedFunctionName string) (*Constant, error) {
	pointer := SourceDynamicMemoryPosRange{
		Position: 0,
		Size:     0,
	}

	return c.AllocateFunctionStruct(uniqueFullyQualifiedFunctionName, pointer)
}

func (c *Constants) DefineFunctionOpcodes(funcConstant *Constant, opcodes []byte) error {
	opcodesPointer := c.dynamicMapper.Write(opcodes, "opcodes for:"+funcConstant.str)

	overwritePointer := SourceDynamicMemoryPos(uint(funcConstant.PosRange().Position) + SwampFuncOpcodeOffset)

	var opcodePointerOctets [16]byte

	binary.LittleEndian.PutUint64(opcodePointerOctets[0:8], uint64(opcodesPointer.Position))
	binary.LittleEndian.PutUint64(opcodePointerOctets[8:16], uint64(opcodesPointer.Size))

	c.dynamicMapper.Overwrite(overwritePointer, opcodePointerOctets[:], "opcodepointer"+funcConstant.str)

	return nil
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

	log.Printf("couldn't find %v", identifier)
	c.DebugOutput()

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

func (c *Constants) DebugOutput() {
	log.Printf("functions:\n")
	for _, function := range c.functions {
		log.Printf("%v %v\n", function.str, function.debugString)
	}
}
