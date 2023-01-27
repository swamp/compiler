/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package typeinfo

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

type InfoType interface {
	Index() int
	Ref() string
	HumanReadable() string
}

type (
	MemoryOffset uint16
	MemorySize   uint16
	MemoryAlign  uint8
)

type MemoryInfo struct {
	MemorySize  MemorySize
	MemoryAlign MemoryAlign
}

type MemoryOffsetInfo struct {
	MemoryOffset MemoryOffset
	MemoryInfo   MemoryInfo
}

type Type struct {
	index int
}

func (t *Type) Index() int {
	return t.index
}

func (t *Type) String() string {
	return fmt.Sprintf("index: %d", t.index)
}

func (t *Type) Ref() string {
	return fmt.Sprintf("%d", t.index)
}

type BoolType struct {
	Type
}

func (t *BoolType) String() string {
	return "Bool"
}

func (t *BoolType) HumanReadable() string {
	return "Bool"
}

type LocalType struct {
	Type
	name string
}

func (t *LocalType) String() string {
	return t.name
}

func (t *LocalType) HumanReadable() string {
	return t.name
}

type AnyMatchingTypes struct {
	Type
}

func (t *AnyMatchingTypes) String() string {
	return "AnyMatching"
}

func (t *AnyMatchingTypes) HumanReadable() string {
	return "*"
}

type BlobType struct {
	Type
}

func (t *BlobType) String() string {
	return "Blob"
}

func (t *BlobType) HumanReadable() string {
	return "Blob"
}

type AnyType struct {
	Type
}

func (t *AnyType) String() string {
	return "Any"
}

func (t *AnyType) HumanReadable() string {
	return "Any"
}

type UnmanagedType struct {
	Type
	name string
}

func (t *UnmanagedType) String() string {
	return fmt.Sprintf("Unmanaged<%v>", t.name)
}

func (t *UnmanagedType) HumanReadable() string {
	return fmt.Sprintf("Unmanaged<%v>", t.name)
}

type IntType struct {
	Type
}

func (t *IntType) String() string {
	return "Int"
}

func (t *IntType) HumanReadable() string {
	return "Int"
}

type FixedType struct {
	Type
}

func (t *FixedType) String() string {
	return "Fixed"
}

func (t *FixedType) HumanReadable() string {
	return "Fixed"
}

type StringType struct {
	Type
}

func (t *StringType) String() string {
	return "String"
}

func (t *StringType) HumanReadable() string {
	return "String"
}

type CharacterType struct {
	Type
}

func (t *CharacterType) String() string {
	return "Char"
}

func (t *CharacterType) HumanReadable() string {
	return "Char"
}

type ResourceNameType struct {
	Type
}

func (t *ResourceNameType) String() string {
	return "ResourceName"
}

func (t *ResourceNameType) HumanReadable() string {
	return "ResourceName"
}

type TypeRefIdType struct {
	Type
	originalType InfoType
}

func (t *TypeRefIdType) String() string {
	return fmt.Sprintf("TypeRefId <%v>", t.originalType)
}

func (t *TypeRefIdType) HumanReadable() string {
	return fmt.Sprintf("TypeRefId <%v>", t.originalType.HumanReadable())
}

type ListType struct {
	Type
	itemType  InfoType
	itemSize  MemorySize
	itemAlign MemoryAlign
}

func (t *ListType) String() string {
	return fmt.Sprintf("List <%v>", t.itemType.Ref())
}

func (t *ListType) HumanReadable() string {
	return fmt.Sprintf("List<%v>", t.itemType.HumanReadable())
}

type ArrayType struct {
	Type
	itemType  InfoType
	itemSize  MemorySize
	itemAlign MemoryAlign
}

func (t *ArrayType) String() string {
	return fmt.Sprintf("Array <%v>", t.itemType.Ref())
}

func (t *ArrayType) HumanReadable() string {
	return fmt.Sprintf("Array<%v>", t.itemType.HumanReadable())
}

type AliasType struct {
	Type
	name     string
	realType InfoType
}

func (t *AliasType) String() string {
	return fmt.Sprintf("alias %s -> %v", t.name, t.realType.Ref())
}

func (t *AliasType) HumanReadable() string {
	return t.name
}

func (t *AliasType) HumanReadableExpanded() string {
	return fmt.Sprintf("%s => %v", t.name, t.realType.HumanReadable())
}

type RecordField struct {
	name             string
	fieldType        InfoType
	memoryOffsetInfo MemoryOffsetInfo
}

func (t RecordField) String() string {
	return fmt.Sprintf("field %s %v", t.name, t.fieldType.Ref())
}

func (t RecordField) HumanReadable() string {
	return fmt.Sprintf("%s : %v", t.name, t.fieldType.HumanReadable())
}

type RecordType struct {
	Type
	fields     []RecordField
	memoryInfo MemoryInfo
}

func (t *RecordType) String() string {
	return fmt.Sprintf("record fields: %v", t.fields)
}

func fieldsToHumanReadable(fields []RecordField) string {
	s := ""
	for index, field := range fields {
		if index > 0 {
			s += ", "
		}
		s += field.HumanReadable()
	}

	return s
}

func (t *RecordType) HumanReadable() string {
	return fmt.Sprintf("{ %v }", fieldsToHumanReadable(t.fields))
}

type VariantField struct {
	index        uint
	memoryOffset MemoryOffsetInfo
	fieldType    InfoType
}

func (t VariantField) String() string {
	return fmt.Sprintf("variantfield %d: %v (%v)", t.index, t.fieldType.Ref(), t.fieldType)
}

func (t VariantField) HumanReadable() string {
	return fmt.Sprintf("%d: %v", t.index, t.fieldType.HumanReadable())
}

type Variant struct {
	Type
	inCustomType InfoType
	name         string
	fields       []VariantField
	memoryInfo   MemoryInfo
}

func Refs(types []InfoType) string {
	if len(types) == 0 {
		return ""
	}
	s := ""
	for index, infoType := range types {
		if index > 0 {
			s += ", "
		}
		if infoType == nil {
			s += "NIL!!!"
		} else {
			s += infoType.Ref()
		}
	}
	return s
}

func (t *Variant) String() string {
	return fmt.Sprintf("variant index:%d %s  (inside:%d) parameters:%v", t.Type.index, t.name, t.inCustomType.Index(), t.fields)
}

func (t *Variant) HumanReadable() string {
	return fmt.Sprintf("%s%s", t.name, t.fields)
}

type CustomType struct {
	Type
	name       string
	generics   []InfoType
	variants   []*Variant
	memoryInfo MemoryInfo
}

func (t *CustomType) String() string {
	return fmt.Sprintf("custom %d variants:%v", t.index, t.variants)
}

func variantsToHumanReadable(variants []*Variant) string {
	s := ""
	for index, variant := range variants {
		if index > 0 {
			s += " | "
		}
		s += variant.HumanReadable()
	}

	return s
}

func (t *CustomType) HumanReadable() string {
	return fmt.Sprintf("%s(%v)", t.name, variantsToHumanReadable(t.variants))
}

type FunctionType struct {
	Type
	parameterTypes []InfoType
}

func (t *FunctionType) String() string {
	return fmt.Sprintf("function (params:%v)", Refs(t.parameterTypes))
}

func typesToHumanReadable(types []InfoType, separator string) string {
	s := ""
	for index, infoType := range types {
		if index > 0 {
			s += separator
		}

		s += infoType.HumanReadable()
	}

	return s
}

func (t *FunctionType) HumanReadable() string {
	return fmt.Sprintf("( %v )", typesToHumanReadable(t.parameterTypes, " -> "))
}

type TupleTypeField struct {
	memoryOffsetInfo MemoryOffsetInfo
	fieldType        InfoType
}

func (t TupleTypeField) String() string {
	return fmt.Sprintf("tuplefield %s", t.fieldType.Ref())
}

func (t TupleTypeField) HumanReadable() string {
	return fmt.Sprintf("%v", t.fieldType.HumanReadable())
}

type TupleType struct {
	Type
	fields     []TupleTypeField
	memoryInfo MemoryInfo
}

func (t *TupleType) String() string {
	return fmt.Sprintf("tuple params:%v", t.fields)
}

func (t *TupleType) HumanReadable() string {
	return fmt.Sprintf("( %v )", t.fields)
}

type Chunk struct {
	infoTypes []InfoType
}

func customsAreSame(custom *CustomType, other *CustomType) bool {
	if custom.name != other.name {
		return false
	}

	for genericIndex, generic := range custom.generics {
		otherGeneric := other.generics[genericIndex]
		if generic.Index() != otherGeneric.Index() {
			return false
		}
	}

	return true
}

func (c *Chunk) doWeHaveCustom(custom *CustomType) int {
	for index, infoType := range c.infoTypes {
		other, isCustom := infoType.(*CustomType)
		if isCustom {
			if customsAreSame(custom, other) {
				return index
			}
		}
	}

	return -1
}

func recordsAreSame(record *RecordType, other *RecordType) bool {
	if len(other.fields) != len(record.fields) {
		return false
	}

	for fieldIndex, field := range record.fields {
		otherField := other.fields[fieldIndex]
		if field.name != otherField.name {
			return false
		}

		if otherField.fieldType == nil {
			panic(fmt.Errorf("problem with other field %v", otherField.name))
		}

		if field.fieldType == nil {
			panic("problem with this field")
		}

		if field.fieldType.Index() != otherField.fieldType.Index() {
			return false
		}
	}

	return true
}

func (c *Chunk) doWeHaveRecord(record *RecordType) int {
	for index, infoType := range c.infoTypes {
		other, isRecord := infoType.(*RecordType)
		if isRecord {
			if recordsAreSame(record, other) {
				return index
			}
		}
	}

	return -1
}

func functionsAreSame(fn *FunctionType, other *FunctionType) bool {
	if len(other.parameterTypes) != len(fn.parameterTypes) {
		return false
	}

	for paramIndex, parameterType := range fn.parameterTypes {
		otherParamType := other.parameterTypes[paramIndex]

		if parameterType.Index() != otherParamType.Index() {
			return false
		}
	}

	return true
}

func tuplesAreSame(fn *TupleType, other *TupleType) bool {
	if len(other.fields) != len(fn.fields) {
		return false
	}

	for paramIndex, parameterType := range fn.fields {
		otherParamType := other.fields[paramIndex].fieldType

		if parameterType.fieldType.Index() != otherParamType.Index() {
			return false
		}
	}

	return true
}

func (c *Chunk) doWeHaveFunction(fn *FunctionType) int {
	for index, infoType := range c.infoTypes {
		other, isFunction := infoType.(*FunctionType)
		if isFunction {
			if functionsAreSame(fn, other) {
				return index
			}
		}
	}

	return -1
}

func (c *Chunk) doWeHaveTuple(fn *TupleType) int {
	for index, infoType := range c.infoTypes {
		other, isFunction := infoType.(*TupleType)
		if isFunction {
			if tuplesAreSame(fn, other) {
				return index
			}
		}
	}

	return -1
}

func (c *Chunk) doWeHaveArray(array *ArrayType) int {
	for index, infoType := range c.infoTypes {
		other, isArray := infoType.(*ArrayType)
		if isArray {
			if other.itemType.Index() == array.itemType.Index() {
				return index
			}
		}
	}

	return -1
}

func (c *Chunk) doWeHaveList(list *ListType) int {
	for index, infoType := range c.infoTypes {
		other, isList := infoType.(*ListType)
		if isList {
			if other.itemType.Index() == list.itemType.Index() {
				return index
			}
		}
	}

	return -1
}

func (c *Chunk) doWeHaveTypeRef(proposedRefIdType *TypeRefIdType) int {
	for index, infoType := range c.infoTypes {
		other, isTypeRef := infoType.(*TypeRefIdType)
		if isTypeRef {
			if other.originalType.Index() == proposedRefIdType.originalType.Index() {
				return index
			}
		}
	}

	return -1
}

func (c *Chunk) doWeHaveAlias(alias *AliasType) int {
	for index, infoType := range c.infoTypes {
		other, isAlias := infoType.(*AliasType)
		if isAlias {
			if alias.name == other.name && other.realType.Index() == alias.realType.Index() {
				return index
			}
		}
	}

	return -1
}

func (c *Chunk) doWeHaveInt() int {
	for index, infoType := range c.infoTypes {
		_, isInt := infoType.(*IntType)
		if isInt {
			return index
		}
	}

	return -1
}

func (c *Chunk) doWeHaveFixed() int {
	for index, infoType := range c.infoTypes {
		_, isInt := infoType.(*FixedType)
		if isInt {
			return index
		}
	}

	return -1
}

func (c *Chunk) doWeHaveBool() int {
	for index, infoType := range c.infoTypes {
		_, isInt := infoType.(*BoolType)
		if isInt {
			return index
		}
	}

	return -1
}

func (c *Chunk) doWeHaveBlob() int {
	for index, infoType := range c.infoTypes {
		_, isBlob := infoType.(*BlobType)
		if isBlob {
			return index
		}
	}

	return -1
}

func (c *Chunk) doWeHaveAny() int {
	for index, infoType := range c.infoTypes {
		_, isAny := infoType.(*AnyType)
		if isAny {
			return index
		}
	}

	return -1
}

func (c *Chunk) doWeHaveUnmanaged(unmanaged *UnmanagedType) int {
	for index, infoType := range c.infoTypes {
		other, isUnmanaged := infoType.(*UnmanagedType)
		if isUnmanaged {
			if other.name == unmanaged.name {
				return index
			}
		}
	}

	return -1
}

func (c *Chunk) doWeHaveAnyMatchingTypes() int {
	for index, infoType := range c.infoTypes {
		_, isAny := infoType.(*AnyMatchingTypes)
		if isAny {
			return index
		}
	}

	return -1
}

func (c *Chunk) doWeHaveString() int {
	for index, infoType := range c.infoTypes {
		_, isInt := infoType.(*StringType)
		if isInt {
			return index
		}
	}

	return -1
}

func (c *Chunk) doWeHaveCharacter() int {
	for index, infoType := range c.infoTypes {
		_, isInt := infoType.(*CharacterType)
		if isInt {
			return index
		}
	}

	return -1
}

func (c *Chunk) doWeHaveResourceName() int {
	for index, infoType := range c.infoTypes {
		_, isResourceName := infoType.(*ResourceNameType)
		if isResourceName {
			return index
		}
	}

	return -1
}

func (c *Chunk) consumeCustomVariant(variant *dectype.CustomTypeVariantAtom) (*Variant, error) {
	consumedTypes, consumeErr := c.ConsumeTypes(variant.ParameterTypes())
	if consumeErr != nil {
		return nil, consumeErr
	}

	if len(variant.ParameterTypes()) != len(variant.Fields()) {
		panic("illegal custom type parametercount field count")
	}

	var fields []VariantField

	for index, field := range variant.Fields() {
		var memInfo MemoryOffsetInfo
		if !dectype.TypeIsTemplateHasLocalTypes(field.Type()) {
			memInfo = memoryOffsetInfo(field.MemoryOffset(), field.Type())
		}

		newFields := VariantField{
			index:        uint(index),
			fieldType:    consumedTypes[index],
			memoryOffset: memInfo,
		}
		fields = append(fields, newFields)
	}

	memorySize, memoryAlign := dectype.GetMemorySizeAndAlignmentInternal(variant)

	proposedNewVariant := &Variant{
		inCustomType: &CustomType{},
		Type:         Type{},
		name:         variant.Name().Name(),
		fields:       fields,
		memoryInfo: MemoryInfo{
			MemorySize:  MemorySize(memorySize),
			MemoryAlign: MemoryAlign(memoryAlign),
		},
	}

	proposedNewVariant.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewVariant)

	return proposedNewVariant, nil
}

func (c *Chunk) consumeCustom(custom *dectype.CustomTypeAtom) (*CustomType, error) {
	customName := custom.ArtifactTypeName().String()
	if len(customName) == 0 {
		panic("custom name must be set here")
	}

	var consumedParameters []InfoType
	for _, generic := range custom.Parameters() {
		consumedParameter, err := c.Consume(generic)
		if err != nil {
			return nil, err
		}

		consumedParameters = append(consumedParameters, consumedParameter)
	}

	proposedNewCustom := &CustomType{
		Type:     Type{},
		name:     customName,
		variants: nil,
		generics: consumedParameters,
		memoryInfo: MemoryInfo{
			MemorySize:  MemorySize(0),
			MemoryAlign: MemoryAlign(0),
		},
	}

	indexCustom := c.doWeHaveCustom(proposedNewCustom)
	if indexCustom != -1 {
		return c.infoTypes[indexCustom].(*CustomType), nil
	}

	var consumedVariants []*Variant
	for _, variant := range custom.Variants() {
		newVariant, err := c.consumeCustomVariant(variant)
		if err != nil {
			return nil, err
		}

		consumedVariants = append(consumedVariants, newVariant)
	}

	memorySize, memoryAlign := dectype.GetMemorySizeAndAlignment(custom)
	proposedNewCustom.memoryInfo = MemoryInfo{
		MemorySize:  MemorySize(memorySize),
		MemoryAlign: MemoryAlign(memoryAlign),
	}
	proposedNewCustom.variants = consumedVariants

	proposedNewCustom.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewCustom)

	for _, variant := range consumedVariants {
		variant.inCustomType = proposedNewCustom
	}

	return proposedNewCustom, nil
}

func memoryOffsetInfo(offset dectype.MemoryOffset, p dtype.Type) MemoryOffsetInfo {
	size, align := dectype.GetMemorySizeAndAlignment(p)
	return MemoryOffsetInfo{
		MemoryOffset: MemoryOffset(offset),
		MemoryInfo: MemoryInfo{
			MemorySize:  MemorySize(size),
			MemoryAlign: MemoryAlign(align),
		},
	}
}

func (c *Chunk) consumeRecord(record *dectype.RecordAtom) (*RecordType, error) {
	var fields []RecordField

	for _, field := range record.SortedFields() {
		consumeFieldType, err := c.Consume(field.Type())
		if err != nil {
			return nil, err
		}
		recordField := RecordField{
			name:             field.Name(),
			memoryOffsetInfo: memoryOffsetInfo(field.MemoryOffset(), field.Type()),
			fieldType:        consumeFieldType,
		}
		fields = append(fields, recordField)
	}

	memorySize, memoryAlign := dectype.GetMemorySizeAndAlignment(record)

	proposedNewRecord := &RecordType{
		Type:   Type{},
		fields: fields,
		memoryInfo: MemoryInfo{
			MemorySize:  MemorySize(memorySize),
			MemoryAlign: MemoryAlign(memoryAlign),
		},
	}

	indexRecord := c.doWeHaveRecord(proposedNewRecord)
	if indexRecord != -1 {
		return c.infoTypes[indexRecord].(*RecordType), nil
	}

	proposedNewRecord.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewRecord)

	return proposedNewRecord, nil
}

func (c *Chunk) consumeFunction(fn *dectype.FunctionAtom) (*FunctionType, error) {
	var types []InfoType

	for _, paramType := range fn.FunctionParameterTypes() {
		consumeType, err := c.Consume(paramType)
		if err != nil {
			return nil, err
		}
		if consumeType == nil {
			return nil, fmt.Errorf("this should not be needed")
		}
		types = append(types, consumeType)
	}

	proposedNewFunction := &FunctionType{
		Type:           Type{},
		parameterTypes: types,
	}

	indexRecord := c.doWeHaveFunction(proposedNewFunction)
	if indexRecord != -1 {
		return c.infoTypes[indexRecord].(*FunctionType), nil
	}

	proposedNewFunction.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewFunction)

	return proposedNewFunction, nil
}

func (c *Chunk) consumeTuple(fn *dectype.TupleTypeAtom) (*TupleType, error) {
	var tupleFields []TupleTypeField

	for _, field := range fn.Fields() {
		consumeType, err := c.Consume(field.Type())
		if err != nil {
			return nil, err
		}
		if consumeType == nil {
			return nil, fmt.Errorf("this should not be needed")
		}
		tupleField := TupleTypeField{
			memoryOffsetInfo: memoryOffsetInfo(field.MemoryOffset(), field.Type()),
			fieldType:        consumeType,
		}
		tupleFields = append(tupleFields, tupleField)
	}

	memorySize, memoryAlign := dectype.GetMemorySizeAndAlignment(fn)

	proposedNewTuple := &TupleType{
		Type:   Type{},
		fields: tupleFields,
		memoryInfo: MemoryInfo{
			MemorySize:  MemorySize(memorySize),
			MemoryAlign: MemoryAlign(memoryAlign),
		},
	}

	indexRecord := c.doWeHaveTuple(proposedNewTuple)
	if indexRecord != -1 {
		return c.infoTypes[indexRecord].(*TupleType), nil
	}

	proposedNewTuple.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewTuple)

	return proposedNewTuple, nil
}

func (c *Chunk) consumeArray(genericTypes []dtype.Type) (InfoType, error) {
	if len(genericTypes) != 1 {
		return nil, fmt.Errorf("can only have one parameter to array")
	}

	itemType := genericTypes[0]

	consumedType, err := c.Consume(itemType)
	if err != nil {
		return nil, err
	}

	memorySize, memoryAlign := dectype.GetMemorySizeAndAlignmentInternal(itemType)

	proposedNewArray := &ArrayType{
		Type:      Type{},
		itemType:  consumedType,
		itemSize:  MemorySize(memorySize),
		itemAlign: MemoryAlign(memoryAlign),
	}

	indexArray := c.doWeHaveArray(proposedNewArray)
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*ArrayType), nil
	}

	proposedNewArray.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewArray)

	return proposedNewArray, nil
}

func (c *Chunk) consumeList(genericTypes []dtype.Type) (InfoType, error) {
	if len(genericTypes) != 1 {
		return nil, fmt.Errorf("can only have one parameter to lists")
	}

	itemType := genericTypes[0]
	consumedType, err := c.Consume(itemType)
	if err != nil {
		return nil, err
	}
	if consumedType == nil {
		return nil, nil
	}

	memorySize, memoryAlign := dectype.GetMemorySizeAndAlignmentInternal(itemType)
	if !dectype.IsLocalType(itemType) && memorySize == 0 {
		panic(fmt.Errorf("can not have zero item size"))
	}

	proposedNewList := &ListType{
		Type:      Type{},
		itemType:  consumedType,
		itemSize:  MemorySize(memorySize),
		itemAlign: MemoryAlign(memoryAlign),
	}

	indexArray := c.doWeHaveList(proposedNewList)
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*ListType), nil
	}

	proposedNewList.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewList)

	return proposedNewList, nil
}

func (c *Chunk) consumeInt() (InfoType, error) {
	indexArray := c.doWeHaveInt()
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*IntType), nil
	}

	proposedNewInt := &IntType{}

	proposedNewInt.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewInt)

	return proposedNewInt, nil
}

func (c *Chunk) consumeFixed() (InfoType, error) {
	indexArray := c.doWeHaveFixed()
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*FixedType), nil
	}

	proposedNewInt := &FixedType{}

	proposedNewInt.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewInt)

	return proposedNewInt, nil
}

func (c *Chunk) consumeString() (InfoType, error) {
	indexArray := c.doWeHaveString()
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*StringType), nil
	}

	proposedNewInt := &StringType{}

	proposedNewInt.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewInt)

	return proposedNewInt, nil
}

func (c *Chunk) consumeCharacter() (InfoType, error) {
	indexArray := c.doWeHaveCharacter()
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*CharacterType), nil
	}

	proposedNewInt := &CharacterType{}

	proposedNewInt.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewInt)

	return proposedNewInt, nil
}

func (c *Chunk) consumeResourceName() (InfoType, error) {
	indexArray := c.doWeHaveResourceName()
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*ResourceNameType), nil
	}

	proposedNewResourceName := &ResourceNameType{}

	proposedNewResourceName.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewResourceName)

	return proposedNewResourceName, nil
}

func (c *Chunk) consumeTypeRef(referencedType dtype.Type) (InfoType, error) {
	consumedReferencedType, err := c.Consume(referencedType)
	if err != nil {
		return nil, err
	}
	if consumedReferencedType == nil {
		return nil, nil
	}

	proposedNewRefIdType := &TypeRefIdType{
		Type:         Type{},
		originalType: consumedReferencedType,
	}

	indexArray := c.doWeHaveTypeRef(proposedNewRefIdType)
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*TypeRefIdType), nil
	}

	proposedNewRefIdType.index = len(c.infoTypes)

	c.infoTypes = append(c.infoTypes, proposedNewRefIdType)

	return proposedNewRefIdType, nil
}

func (c *Chunk) consumeBool() (InfoType, error) {
	indexArray := c.doWeHaveBool()
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*BoolType), nil
	}

	proposedNewInt := &BoolType{}

	proposedNewInt.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewInt)

	return proposedNewInt, nil
}

func (c *Chunk) consumeBlob() (InfoType, error) {
	indexArray := c.doWeHaveBlob()
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*BlobType), nil
	}

	proposedNewInt := &BlobType{}

	proposedNewInt.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewInt)

	return proposedNewInt, nil
}

func (c *Chunk) consumeAny() (InfoType, error) {
	indexArray := c.doWeHaveAny()
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*AnyType), nil
	}

	proposedNewAnyType := &AnyType{}

	proposedNewAnyType.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewAnyType)

	return proposedNewAnyType, nil
}

func (c *Chunk) consumeUnmanaged(unmanagedType *dectype.UnmanagedType) (InfoType, error) {
	if len(unmanagedType.Identifier().NativeLanguageTypeName().Name()) == 0 {
		return nil, fmt.Errorf("must have a nativeLanguageTypeName")
	}

	propsedNewUnmanaged := &UnmanagedType{
		Type: Type{},
		name: unmanagedType.Identifier().NativeLanguageTypeName().Name(),
	}

	indexArray := c.doWeHaveUnmanaged(propsedNewUnmanaged)
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*UnmanagedType), nil
	}

	propsedNewUnmanaged.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, propsedNewUnmanaged)

	return propsedNewUnmanaged, nil
}

func (c *Chunk) consumeAnyMatchingTypes() (InfoType, error) {
	indexArray := c.doWeHaveAnyMatchingTypes()
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*AnyMatchingTypes), nil
	}

	proposedNewAnyType := &AnyMatchingTypes{}
	proposedNewAnyType.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewAnyType)

	return proposedNewAnyType, nil
}

func (c *Chunk) consumePrimitive(primitive *dectype.PrimitiveAtom) (InfoType, error) {
	name := primitive.PrimitiveName().Name()
	if name == "List" {
		return c.consumeList(primitive.GenericTypes())
	} else if name == "Array" {
		return c.consumeArray(primitive.GenericTypes())
	} else if name == "Int" {
		return c.consumeInt()
	} else if name == "Fixed" {
		return c.consumeFixed()
	} else if name == "Bool" {
		return c.consumeBool()
	} else if name == "String" {
		return c.consumeString()
	} else if name == "Char" {
		return c.consumeCharacter()
	} else if name == "ResourceName" {
		return c.consumeResourceName()
	} else if name == "Blob" {
		return c.consumeBlob()
	} else if name == "TypeRef" {
		return c.consumeTypeRef(primitive.GenericTypes()[0])
	} else if name == "Any" {
		return c.consumeAny()
	}

	return nil, fmt.Errorf("chunk: consume: unknown primitive %v", primitive)
}

func (c *Chunk) consumeLocalType(localType *dectype.LocalTypeDefinition) (InfoType, error) {
	return &LocalType{
		Type: Type{},
		name: localType.Identifier().Name(),
	}, nil
}

func (c *Chunk) consumeAlias(alias *dectype.Alias) (InfoType, error) {
	consumedType, err := c.Consume(alias.Next())
	if err != nil {
		return nil, err
	}

	if consumedType == nil {
		return nil, nil
	}

	fullyQualifiedName := alias.ArtifactTypeName()

	proposedNewAlias := &AliasType{
		Type:     Type{},
		name:     fullyQualifiedName.String(),
		realType: consumedType,
	}

	indexArray := c.doWeHaveAlias(proposedNewAlias)
	if indexArray != -1 {
		return c.infoTypes[indexArray].(*AliasType), nil
	}

	proposedNewAlias.index = len(c.infoTypes)
	c.infoTypes = append(c.infoTypes, proposedNewAlias)

	return proposedNewAlias, nil
}

func (c *Chunk) ConsumeAtom(a dtype.Atom) (InfoType, error) {
	switch t := a.(type) {
	case *dectype.CustomTypeAtom:
		return c.consumeCustom(t)
	case *dectype.CustomTypeVariantAtom:
		customType, err := c.consumeCustom(t.InCustomType())
		if err != nil {
			return nil, err
		}
		createdVariant := customType.variants[t.Index()]
		return createdVariant, nil
	case *dectype.RecordAtom:
		return c.consumeRecord(t)
	case *dectype.FunctionAtom:
		return c.consumeFunction(t)
	case *dectype.PrimitiveAtom:
		return c.consumePrimitive(t)
	case *dectype.TupleTypeAtom:
		return c.consumeTuple(t)
	case *dectype.UnmanagedType:
		return c.consumeUnmanaged(t)
	case *dectype.AnyMatchingTypes:
		return c.consumeAnyMatchingTypes()
	}

	return nil, fmt.Errorf("unknown atom %T", a)
}

func (c *Chunk) ConsumeType(d dtype.Type) (InfoType, error) {
	return c.Consume(d)
}

func (c *Chunk) Lookup(d dtype.Type) (int, error) {
	infoType, err := c.ConsumeType(d)
	if err != nil {
		return -1, err
	}

	return infoType.Index(), nil
}

func (c *Chunk) Consume(p dtype.Type) (InfoType, error) {
	atom, isAtom := p.(dtype.Atom)
	if isAtom {
		return c.ConsumeAtom(atom)
	}
	switch t := p.(type) {
	case *dectype.Alias:
		return c.consumeAlias(t)
	case *dectype.FunctionTypeReference:
		return c.Consume(t.Next())
	case *dectype.CustomTypeVariantReference:
		// intentionally ignore
		return c.Consume(t.Next())
	case *dectype.LocalTypeDefinition:
		// intentionally ignore
		return c.Consume(t.Next())
	case *dectype.InvokerType:
		invokerAtom, resolveErr := t.Resolve()
		if resolveErr != nil {
			return nil, resolveErr
		}
		if invokerAtom == nil {
			return nil, fmt.Errorf("wrong atom invoke")
		}
		return c.ConsumeAtom(invokerAtom)
	case *dectype.PrimitiveTypeReference:
		return c.Consume(t.Next())
	case *dectype.AliasReference:
		return c.Consume(t.Next())
	case *dectype.CustomTypeReference:
		return c.Consume(t.Next())
	}

	err := fmt.Errorf("chunk: consume: unknown thing %T", p)
	log.Print(err.Error())
	return nil, err
}

func (c *Chunk) ConsumeTypes(types []dtype.Type) ([]InfoType, error) {
	var consumedTypes []InfoType

	for _, t := range types {
		consumed, err := c.Consume(t)
		if err != nil {
			return nil, err
		}
		consumedTypes = append(consumedTypes, consumed)
	}

	return consumedTypes, nil
}

func (c *Chunk) DebugOutputHumanReadable() {
	for index, t := range c.infoTypes {
		s := t.HumanReadable()
		alias, isAlias := t.(*AliasType)
		if isAlias {
			s = alias.HumanReadableExpanded()
		}
		log.Printf("%d : %v\n", index, s)
	}
}

func (c *Chunk) DebugOutputStrict() {
	for index, t := range c.infoTypes {
		log.Printf("%d : %v\n", index, t)
	}
}

func (c *Chunk) DebugOutput() {
	c.DebugOutputHumanReadable()
	c.DebugOutputStrict()
}
