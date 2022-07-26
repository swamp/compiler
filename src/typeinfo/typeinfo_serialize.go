/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package typeinfo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"io"
)

type SwtiType uint8

const (
	SwtiTypeCustom SwtiType = iota
	SwtiTypeFunction
	SwtiTypeAlias
	SwtiTypeRecord
	SwtiTypeArray
	SwtiTypeList
	SwtiTypeString
	SwtiTypeInt
	SwtiTypeFixed
	SwtiTypeBoolean
	SwtiTypeBlob
	SwtiTypeResourceName
	SwtiTypeChar
	SwtiTypeTuple
	SwtiTypeRefId
	SwtiTypeAny
	SwtiTypeAnyMatchingTypes
	SwtiTypeUnmanaged
)

func writeUint8(writer io.Writer, v byte) error {
	_, err := writer.Write([]byte{v})
	return err
}

func writeUint16(writer io.Writer, v int) error {
	if v > 65535 {
		panic("not allowed to write bigger than uint16")
	}
	buf := []byte{byte(0), byte(0)}
	binary.BigEndian.PutUint16(buf, uint16(v))
	_, err := writer.Write(buf)
	return err
}

func writeTypeID(writer io.Writer, v SwtiType) error {
	return writeUint8(writer, byte(v))
}

func writePrimitive(writer io.Writer, v SwtiType) error {
	return writeTypeID(writer, v)
}

func writeTypeRef(writer io.Writer, infoType InfoType) error {
	return writeUint16(writer, infoType.Index())
}

func writeList(writer io.Writer, list *ListType) error {
	if err := writeTypeID(writer, SwtiTypeList); err != nil {
		return err
	}
	if err := writeTypeRef(writer, list.itemType); err != nil {
		return err
	}

	if err := writeMemoryInfo(writer, list.itemSize, list.itemAlign); err != nil {
		return err
	}

	return nil
}

func writeUnmanaged(writer io.Writer, unmanaged *UnmanagedType) error {
	if err := writeTypeID(writer, SwtiTypeUnmanaged); err != nil {
		return err
	}

	if err := writeName(writer, unmanaged.name); err != nil {
		return err
	}

	hash := fnv.New32a()
	hash.Write([]byte(unmanaged.name))
	lowHash := int(hash.Sum32() & 0xffff)
	return writeUint16(writer, lowHash)
}

func writeTuple(writer io.Writer, tuple *TupleType) error {
	if err := writeTypeID(writer, SwtiTypeTuple); err != nil {
		return err
	}

	if err := writeMemoryInfo(writer, tuple.memoryInfo.MemorySize, tuple.memoryInfo.MemoryAlign); err != nil {
		return err
	}

	if err := writeCount(writer, len(tuple.fields)); err != nil {
		return err
	}

	for _, field := range tuple.fields {
		if err := writeMemoryOffsetInfo(writer, field.memoryOffsetInfo); err != nil {
			return err
		}
		if err := writeTypeRef(writer, field.fieldType); err != nil {
			return err
		}
	}

	return nil
}

func writeArray(writer io.Writer, array *ArrayType) error {
	if err := writeTypeID(writer, SwtiTypeArray); err != nil {
		return err
	}
	if err := writeTypeRef(writer, array.itemType); err != nil {
		return err
	}

	if err := writeMemoryInfo(writer, array.itemSize, array.itemAlign); err != nil {
		return err
	}

	return nil
}

func write16BitCount(writer io.Writer, count int) error {
	return writeUint16(writer, count)
}

func writeCount(writer io.Writer, count int) error {
	_, err := writer.Write([]byte{byte(count)})
	return err
}

func writeName(writer io.Writer, name string) error {
	if err := writeCount(writer, len(name)); err != nil {
		return err
	}

	_, err := writer.Write([]byte(name))

	return err
}

func writeMemoryOffset(writer io.Writer, offset MemoryOffset) error {
	return writeUint16(writer, int(offset))
}

func writeMemorySize(writer io.Writer, size MemorySize) error {
	return writeUint16(writer, int(size))
}

func writeMemoryAlign(writer io.Writer, align MemoryAlign) error {
	if align > 8 {
		panic(fmt.Errorf("illegal memory align"))
	}
	return writeUint8(writer, uint8(align))
}

func writeMemoryInfo(writer io.Writer, size MemorySize, align MemoryAlign) error {
	if err := writeMemorySize(writer, size); err != nil {
		return err
	}
	if err := writeMemoryAlign(writer, align); err != nil {
		return err
	}

	return nil
}

func writeMemoryOffsetInfo(writer io.Writer, info MemoryOffsetInfo) error {
	if err := writeMemoryOffset(writer, info.MemoryOffset); err != nil {
		return err
	}
	if err := writeMemoryInfo(writer, info.MemoryInfo.MemorySize, info.MemoryInfo.MemoryAlign); err != nil {
		return err
	}

	return nil
}

func writeRecord(writer io.Writer, record *RecordType) error {
	if err := writeTypeID(writer, SwtiTypeRecord); err != nil {
		return err
	}

	if err := writeMemoryInfo(writer, record.memoryInfo.MemorySize, record.memoryInfo.MemoryAlign); err != nil {
		return err
	}
	if err := writeCount(writer, len(record.fields)); err != nil {
		return err
	}
	for _, field := range record.fields {
		if err := writeName(writer, field.name); err != nil {
			return err
		}

		if err := writeMemoryOffsetInfo(writer, field.memoryOffsetInfo); err != nil {
			return err
		}

		if err := writeTypeRef(writer, field.fieldType); err != nil {
			return err
		}
	}

	return nil
}

func writeAlias(writer io.Writer, alias *AliasType) error {
	if err := writeTypeID(writer, SwtiTypeAlias); err != nil {
		return err
	}

	if err := writeName(writer, alias.name); err != nil {
		return err
	}

	if err := writeTypeRef(writer, alias.realType); err != nil {
		return err
	}

	return nil
}

func writeTypeRefId(writer io.Writer, typeRefId *TypeRefIdType) error {
	if err := writeTypeID(writer, SwtiTypeRefId); err != nil {
		return err
	}

	if err := writeTypeRef(writer, typeRefId.originalType); err != nil {
		return err
	}

	return nil
}

func writeFunction(writer io.Writer, fn *FunctionType) error {
	if err := writeTypeID(writer, SwtiTypeFunction); err != nil {
		return err
	}
	if err := writeCount(writer, len(fn.parameterTypes)); err != nil {
		return err
	}

	for _, parameterType := range fn.parameterTypes {
		if err := writeTypeRef(writer, parameterType); err != nil {
			return err
		}
	}

	return nil
}

func writeCustomTypeVariantField(writer io.Writer, variantField VariantField) error {
	if err := writeTypeRef(writer, variantField.fieldType); err != nil {
		return err
	}

	if err := writeMemoryOffsetInfo(writer, variantField.memoryOffset); err != nil {
		return err
	}

	return nil
}

func writeCustomTypeVariant(writer io.Writer, variant *Variant) error {
	if err := writeName(writer, variant.name); err != nil {
		return err
	}

	if err := writeMemoryInfo(writer, variant.memoryInfo.MemorySize, variant.memoryInfo.MemoryAlign); err != nil {
		return err
	}

	if err := writeCount(writer, len(variant.fields)); err != nil {
		return err
	}
	for _, variantField := range variant.fields {
		if err := writeCustomTypeVariantField(writer, variantField); err != nil {
			return err
		}
	}

	return nil
}

func writeCustom(writer io.Writer, custom *CustomType) error {
	if err := writeTypeID(writer, SwtiTypeCustom); err != nil {
		return err
	}

	if len(custom.name) == 0 {
		panic("Name must be set")
	}

	if err := writeName(writer, custom.name); err != nil {
		return err
	}

	if err := writeMemoryInfo(writer, custom.memoryInfo.MemorySize, custom.memoryInfo.MemoryAlign); err != nil {
		return err
	}

	if err := writeCount(writer, len(custom.variants)); err != nil {
		return err
	}

	for _, variant := range custom.variants {
		if err := writeCustomTypeVariant(writer, variant); err != nil {
			return err
		}
	}

	return nil
}

func writeInfoType(writer io.Writer, entry InfoType) error {
	switch t := entry.(type) {
	case *ListType:
		return writeList(writer, t)
	case *ArrayType:
		return writeArray(writer, t)
	case *IntType:
		return writePrimitive(writer, SwtiTypeInt)
	case *StringType:
		return writePrimitive(writer, SwtiTypeString)
	case *CharacterType:
		return writePrimitive(writer, SwtiTypeChar)
	case *ResourceNameType:
		return writePrimitive(writer, SwtiTypeResourceName)
	case *FixedType:
		return writePrimitive(writer, SwtiTypeFixed)
	case *BoolType:
		return writePrimitive(writer, SwtiTypeBoolean)
	case *BlobType:
		return writePrimitive(writer, SwtiTypeBlob)
	case *RecordType:
		return writeRecord(writer, t)
	case *AliasType:
		return writeAlias(writer, t)
	case *FunctionType:
		return writeFunction(writer, t)
	case *CustomType:
		return writeCustom(writer, t)
	case *TypeRefIdType:
		return writeTypeRefId(writer, t)
	case *TupleType:
		return writeTuple(writer, t)
	case *AnyType:
		// TODO:
		return writePrimitive(writer, SwtiTypeAny)
	case *UnmanagedType:
		return writeUnmanaged(writer, t)
	case *AnyMatchingTypes:
		// TODO:
		return writePrimitive(writer, SwtiTypeAnyMatchingTypes)
	}

	return fmt.Errorf("strange, unknown info type %v %T", entry, entry)
}

func writeVersion(writer io.Writer) error {
	const (
		major byte = 0
		minor byte = 1
		patch byte = 8
	)

	if err := writeUint8(writer, major); err != nil {
		return err
	}

	if err := writeUint8(writer, minor); err != nil {
		return err
	}

	return writeUint8(writer, patch)
}

func Serialize(c *Chunk, writer io.Writer) error {
	if err := writeVersion(writer); err != nil {
		return err
	}

	if err := write16BitCount(writer, len(c.infoTypes)); err != nil {
		return err
	}

	for _, entry := range c.infoTypes {
		if err := writeInfoType(writer, entry); err != nil {
			return err
		}
	}

	return nil
}

func SerializeToOctets(c *Chunk) ([]byte, error) {
	var buf bytes.Buffer

	if err := Serialize(c, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
