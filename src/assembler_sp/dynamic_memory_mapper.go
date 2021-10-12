package assembler_sp

import (
	"fmt"
	"log"
)

type DynamicMemoryInfo struct {
	DebugString string
	Allocation  SourceDynamicMemoryPosRange
}

type DynamicMemoryMapper struct {
	position        uint32
	maxOctetSize    uint
	memory          []byte
	maxIndexWritten uint32
	infos           []DynamicMemoryInfo
}

func DynamicMemoryMapperNew(maxOctetSize uint) *DynamicMemoryMapper {
	return &DynamicMemoryMapper{maxOctetSize: maxOctetSize, memory: make([]byte, maxOctetSize)}
}

func (m *DynamicMemoryMapper) Octets() []byte {
	return m.memory[0:m.position]
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

	posRange := SourceDynamicMemoryPosRange{Position: pos, Size: DynamicMemoryRange(octetSize)}
	info := DynamicMemoryInfo{
		DebugString: debugString,
		Allocation:  posRange,
	}
	m.infos = append(m.infos, info)

	return posRange
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

func (m *DynamicMemoryMapper) Overwrite(position SourceDynamicMemoryPos, data []byte, debugString string) error {
	endPos := uint32(position) + uint32(len(data))
	copy(m.memory[position:endPos], data)
	return nil
}

func (m *DynamicMemoryMapper) DebugOutput() {
	log.Printf("Dynamic Memory (PackageConstants):\n")
	for _, info := range m.infos {
		log.Printf("  %v %v\n", info.Allocation, info.DebugString)
	}
}
