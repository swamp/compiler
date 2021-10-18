package assembler_sp

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	dectype "github.com/swamp/compiler/src/decorated/types"
)

type PackageConstants struct {
	constants             []*Constant
	functions             []*Constant
	externalFunctions     []*Constant
	strings               []*Constant
	resourceNames         []*Constant
	dynamicMapper         *DynamicMemoryMapper
	someConstantIDCounter uint
}

func NewPackageConstants() *PackageConstants {
	return &PackageConstants{
		dynamicMapper: DynamicMemoryMapperNew(128 * 1024),
	}
}

func (c *PackageConstants) String() string {
	s := "\n"
	for _, constant := range c.constants {
		if constant == nil {
			panic("swamp assembler: nil constant")
		}
		s += fmt.Sprintf("%v\n", constant)
	}
	return strings.TrimSpace(s)
}

func (c *PackageConstants) Constants() []*Constant {
	return c.constants
}

func (c *PackageConstants) Finalize() {
	if len(c.resourceNames) == 0 {
		return
	}
	pointerArea := c.dynamicMapper.Allocate(uint(int(dectype.Sizeof64BitPointer)*len(c.resourceNames)), uint32(dectype.Alignof64BitPointer), "Resource name chunk")
	for index, resourceName := range c.resourceNames {
		writePosition := pointerArea.Position + SourceDynamicMemoryPos(index*int(dectype.Sizeof64BitPointer))
		binary.LittleEndian.PutUint64(c.dynamicMapper.memory[writePosition:writePosition+SourceDynamicMemoryPos(dectype.Sizeof64BitPointer)], uint64(resourceName.PosRange().Position))
	}

	var resourceNameChunkOctets [16]byte
	binary.LittleEndian.PutUint64(resourceNameChunkOctets[0:8], uint64(pointerArea.Position))
	binary.LittleEndian.PutUint64(resourceNameChunkOctets[8:16], uint64(len(c.resourceNames)))
	resourceNameChunkPointer := c.dynamicMapper.Write(resourceNameChunkOctets[:], "ResourceNameChunk struct (character-pointer-pointer, resourceNameCount)")

	c.constants = append(c.constants, &Constant{
		constantType:   ConstantTypeResourceNameChunk,
		str:            "",
		source:         resourceNameChunkPointer,
		debugString:    "resource name chunk",
		resourceNameId: 0,
	})
}

func (c *PackageConstants) DynamicMemory() *DynamicMemoryMapper {
	return c.dynamicMapper
}

func (c *PackageConstants) AllocateStringOctets(s string) SourceDynamicMemoryPosRange {
	stringOctets := []byte(s)
	stringOctets = append(stringOctets, byte(0))
	stringOctetsPointer := c.dynamicMapper.Write(stringOctets, "string:"+s)

	return stringOctetsPointer
}

const SizeofSwampString = 16

func (c *PackageConstants) AllocateStringConstant(s string) *Constant {
	for _, constant := range c.strings {
		if constant.str == s {
			return constant
		}
	}

	stringOctetsPointer := c.AllocateStringOctets(s)

	var swampStringOctets [SizeofSwampString]byte
	binary.LittleEndian.PutUint64(swampStringOctets[0:8], uint64(stringOctetsPointer.Position))
	binary.LittleEndian.PutUint64(swampStringOctets[8:16], uint64(len(s)))

	swampStringPointer := c.dynamicMapper.Write(swampStringOctets[:], "SwampString struct (character-pointer, characterCount) for:"+s)

	newConstant := NewStringConstant("string", s, swampStringPointer)
	c.constants = append(c.constants, newConstant)
	c.strings = append(c.strings, newConstant)

	return newConstant
}

func (c *PackageConstants) AllocateResourceNameConstant(s string) *Constant {
	for _, resourceNameConstant := range c.resourceNames {
		if resourceNameConstant.str == s {
			return resourceNameConstant
		}
	}

	stringOctetsPointer := c.AllocateStringOctets(s)

	newConstant := NewResourceNameConstant(c.someConstantIDCounter, s, stringOctetsPointer)
	c.someConstantIDCounter++
	c.constants = append(c.constants, newConstant)
	c.resourceNames = append(c.resourceNames, newConstant)

	return newConstant
}

const (
	SizeofSwampFunc         = 9 * 8
	SizeofSwampExternalFunc = 18 * 8
)

func (c *PackageConstants) AllocateFunctionStruct(uniqueFullyQualifiedFunctionName string,
	opcodesPointer SourceDynamicMemoryPosRange, returnOctetSize dectype.MemorySize,
	returnAlignSize dectype.MemoryAlign, parameterCount uint, parameterOctetSize dectype.MemorySize, typeIndex uint) (*Constant, error) {
	var swampFuncStruct [SizeofSwampFunc]byte

	fullyQualifiedStringPointer := c.AllocateStringOctets(uniqueFullyQualifiedFunctionName)

	binary.LittleEndian.PutUint32(swampFuncStruct[0:4], uint32(0))
	binary.LittleEndian.PutUint64(swampFuncStruct[8:16], uint64(parameterCount))      // parameterCount
	binary.LittleEndian.PutUint64(swampFuncStruct[16:24], uint64(parameterOctetSize)) // parameters octet size

	binary.LittleEndian.PutUint64(swampFuncStruct[24:32], uint64(opcodesPointer.Position))
	binary.LittleEndian.PutUint64(swampFuncStruct[32:40], uint64(opcodesPointer.Size))

	binary.LittleEndian.PutUint64(swampFuncStruct[40:48], uint64(returnOctetSize)) // returnOctetSize
	binary.LittleEndian.PutUint64(swampFuncStruct[48:56], uint64(returnAlignSize)) // returnAlign

	binary.LittleEndian.PutUint64(swampFuncStruct[56:64], uint64(fullyQualifiedStringPointer.Position)) // debugName
	binary.LittleEndian.PutUint64(swampFuncStruct[64:72], uint64(typeIndex))                            // typeIndex

	funcPointer := c.dynamicMapper.Write(swampFuncStruct[:], "function Struct for:"+uniqueFullyQualifiedFunctionName)

	newConstant := NewFunctionReferenceConstantWithDebug("fn", uniqueFullyQualifiedFunctionName, funcPointer)
	c.constants = append(c.constants, newConstant)
	c.functions = append(c.functions, newConstant)

	return newConstant, nil
}

func (c *PackageConstants) AllocateExternalFunctionStruct(uniqueFullyQualifiedFunctionName string, returnValue SourceStackPosRange, parameters []SourceStackPosRange) (*Constant, error) {
	var swampFuncStruct [SizeofSwampExternalFunc]byte

	fullyQualifiedStringPointer := c.AllocateStringOctets(uniqueFullyQualifiedFunctionName)
	if len(parameters) == 0 {
		// panic(fmt.Errorf("not allowed to have zero paramters for %v", uniqueFullyQualifiedFunctionName))
	}

	binary.LittleEndian.PutUint32(swampFuncStruct[0:4], uint32(1))                  // external type
	binary.LittleEndian.PutUint64(swampFuncStruct[8:16], uint64(len(parameters)))   // parameterCount
	binary.LittleEndian.PutUint32(swampFuncStruct[16:20], uint32(returnValue.Pos))  // return pos
	binary.LittleEndian.PutUint32(swampFuncStruct[20:24], uint32(returnValue.Size)) // return size

	for index, param := range parameters {
		first := 24 + index*8
		firstEnd := first + 8
		second := 28 + index*8
		secondEnd := second + 8
		binary.LittleEndian.PutUint32(swampFuncStruct[first:firstEnd], uint32(param.Pos))    // params pos
		binary.LittleEndian.PutUint32(swampFuncStruct[second:secondEnd], uint32(param.Size)) // params size
	}

	binary.LittleEndian.PutUint64(swampFuncStruct[120:128], uint64(fullyQualifiedStringPointer.Position)) // debugName

	funcPointer := c.dynamicMapper.Write(swampFuncStruct[:], fmt.Sprintf("external function Struct for: '%s' param Count: %d", uniqueFullyQualifiedFunctionName, len(parameters)))

	newConstant := NewExternalFunctionReferenceConstantWithDebug("fn", uniqueFullyQualifiedFunctionName, funcPointer)
	c.constants = append(c.constants, newConstant)
	c.externalFunctions = append(c.externalFunctions, newConstant)

	return newConstant, nil
}

const SwampFuncOpcodeOffset = 24

func (c *PackageConstants) FetchOpcodes(functionConstant *Constant) []byte {
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

func (c *PackageConstants) AllocatePrepareFunctionConstant(uniqueFullyQualifiedFunctionName string,
	returnSize dectype.MemorySize, returnAlign dectype.MemoryAlign,
	parameterCount uint, parameterOctetSize dectype.MemorySize, typeId uint) (*Constant, error) {
	pointer := SourceDynamicMemoryPosRange{
		Position: 0,
		Size:     0,
	}

	return c.AllocateFunctionStruct(uniqueFullyQualifiedFunctionName, pointer, returnSize, returnAlign,
		parameterCount, parameterOctetSize, typeId)
}

func (c *PackageConstants) AllocatePrepareExternalFunctionConstant(uniqueFullyQualifiedFunctionName string, returnValue SourceStackPosRange, parameters []SourceStackPosRange) (*Constant, error) {
	return c.AllocateExternalFunctionStruct(uniqueFullyQualifiedFunctionName, returnValue, parameters)
}

func (c *PackageConstants) DefineFunctionOpcodes(funcConstant *Constant, opcodes []byte) error {
	opcodesPointer := c.dynamicMapper.Write(opcodes, "opcodes for:"+funcConstant.str)

	overwritePointer := SourceDynamicMemoryPos(uint(funcConstant.PosRange().Position) + SwampFuncOpcodeOffset)

	var opcodePointerOctets [16]byte

	binary.LittleEndian.PutUint64(opcodePointerOctets[0:8], uint64(opcodesPointer.Position))
	binary.LittleEndian.PutUint64(opcodePointerOctets[8:16], uint64(opcodesPointer.Size))

	c.dynamicMapper.Overwrite(overwritePointer, opcodePointerOctets[:], "opcodepointer"+funcConstant.str)

	return nil
}

func (c *PackageConstants) FindFunction(identifier VariableName) *Constant {
	for _, constant := range c.functions {
		if constant.str == string(identifier) {
			return constant
		}
	}

	return c.FindExternalFunction(identifier)
}

func (c *PackageConstants) FindExternalFunction(identifier VariableName) *Constant {
	for _, constant := range c.externalFunctions {
		if constant.str == string(identifier) {
			return constant
		}
	}

	log.Printf("couldn't find constant external function %v", identifier)
	c.DebugOutput()

	return nil
}

func (c *PackageConstants) FindStringConstant(s string) *Constant {
	for _, constant := range c.strings {
		if constant.str == s {
			return constant
		}
	}
	return nil
}

func (c *PackageConstants) DebugOutput() {
	log.Printf("functions:\n")
	for _, function := range c.functions {
		log.Printf("%v %v\n", function.str, function.debugString)
	}
}
