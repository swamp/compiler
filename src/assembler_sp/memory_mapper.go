package assembler_sp

import (
	"fmt"
	"log"
)

type ZeroMemoryPointer uint32

type StackMemoryMapper struct {
	position     uint32
	maxOctetSize uint
	memory       []byte
}

func NewStackMemoryMapper(maxOctetSize uint) *StackMemoryMapper {
	return &StackMemoryMapper{maxOctetSize: maxOctetSize}
}

func (m *StackMemoryMapper) Allocate(octetSize uint, align uint32, debugString string) TargetStackPosRange {
	if octetSize == 0 {
		panic(fmt.Errorf("octet size zero is not allowed for allocate stack memory"))
	}
	if align == 0 {
		panic(fmt.Errorf("align zero size is not allowed for allocate stack memory"))
	}
	extra := m.position % align
	if extra != 0 {
		m.position += align - extra
	}
	pos := m.position
	log.Printf("stack allocate: %v [%v] ('%v') => %v\n", octetSize, align, debugString, pos)

	m.position += uint32(octetSize)

	posRange := TargetStackPosRange{
		Pos:  TargetStackPos(pos),
		Size: StackRange(octetSize),
	}
	return posRange
}

func (m *StackMemoryMapper) Set(pos TargetStackPos) {
	m.position = uint32(pos)
}

type DynamicMemoryMapper struct {
	position        uint32
	maxOctetSize    uint
	memory          []byte
	maxIndexWritten uint32
}

func DynamicMemoryMapperNew(maxOctetSize uint) *DynamicMemoryMapper {
	return &DynamicMemoryMapper{maxOctetSize: maxOctetSize, memory: make([]byte, maxOctetSize)}
}

func (m *DynamicMemoryMapper) Allocate(octetSize uint, align uint32, debugString string) SourceDynamicMemoryPosRange {
	if octetSize == 0 {
		panic(fmt.Errorf("octet size zero is not allowed for allocate DynamicMemoryMapper memory"))
	}
	if align == 0 {
		panic(fmt.Errorf("align zero size is not allowed for allocate DynamicMemoryMapper memory"))
	}
	extra := m.position % align
	if extra != 0 {
		m.position += align - extra
	}
	pos := SourceDynamicMemoryPos(m.position)

	m.position += uint32(octetSize)

	return SourceDynamicMemoryPosRange{Position: pos, Size: DynamicMemoryRange(octetSize)}
}

func (m *DynamicMemoryMapper) Read(pos SourceDynamicMemoryPosRange) []byte {
	start := pos.Position
	endPos := uint32(start) + uint32(pos.Size)
	return m.memory[start:endPos]
}

func (m *DynamicMemoryMapper) Write(data []byte, debugString string) SourceDynamicMemoryPosRange {
	posRange := m.Allocate(uint(len(data)), 1, debugString)
	position := posRange.Position
	endPos := uint32(position) + uint32(len(data))
	if endPos > m.maxIndexWritten {
		m.maxIndexWritten = endPos - 1
	}
	copy(m.memory[position:endPos], data)

	return posRange
}
