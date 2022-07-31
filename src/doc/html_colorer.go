/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package doc

import (
	"fmt"
	"io"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type HtmlColorer struct {
	writer io.Writer
}

func spanWrite(writer io.Writer, className string, value string) {
	fmt.Fprintf(writer, span(className, value))
}

func spanWriteLink(writer io.Writer, typeToWrite dtype.Type, value string) {
	_, fullyQualifiedName := findTypeNames(typeToWrite)
	if fullyQualifiedName.String() == "" {
		panic(fmt.Errorf("must have fully qualified name for %T", typeToWrite))
	}
	fullyQualifiedNameString := fullyQualifiedName.String()
	cssClassName := classNameFromType(typeToWrite)

	fmt.Fprintf(writer, "<a href=\"#%v\" title=\"%v\">", fullyQualifiedNameString, fullyQualifiedNameString)

	if len(fullyQualifiedName.ModuleName.Path().Prefix()) > 0 {
		for index, part := range fullyQualifiedName.ModuleName.Path().Prefix() {
			if index > 0 {
				spanWrite(writer, "operator", ".")
			}

			spanWrite(writer, "modulereferenceprefix", part.Name())
		}

		spanWrite(writer, "operator", ".")
	}

	fmt.Fprintf(writer, span(cssClassName, fullyQualifiedName.Last()))

	fmt.Fprintf(writer, "</a>")
}

func (c *HtmlColorer) CustomType(t *dectype.CustomTypeAtom) {
	spanWrite(c.writer, "customtypename", t.Name())
}

func (c *HtmlColorer) CustomTypeName(t *dectype.CustomTypeReference) {
	spanWriteLink(c.writer, t.CustomTypeAtom(), t.CustomTypeAtom().Name())
}

func (c *HtmlColorer) CustomTypeVariant(t *dectype.CustomTypeVariantAtom) {
	spanWrite(c.writer, "customtypevariant", t.Name().Name())
}

func (c *HtmlColorer) InvokerType(t *dectype.InvokerType) {
	spanWriteLink(c.writer, t.TypeGenerator(), t.TypeGenerator().HumanReadable())
}

func (c *HtmlColorer) RecordTypeField(t *dectype.RecordField) {
	spanWrite(c.writer, "recordtypefield", t.Name())
}

func (c *HtmlColorer) AliasName(t *dectype.Alias) {
	spanWriteLink(c.writer, t, t.TypeIdentifier().Name())
}

func (c *HtmlColorer) RightArrow() {
	spanWrite(c.writer, "arrow", "➞")
}

func (c *HtmlColorer) PrimitiveTypeName(t *ast.TypeIdentifier) {
	spanWrite(c.writer, "primitivetype", t.Name())
}

func (c *HtmlColorer) UnmanagedName(t *ast.UnmanagedType) {
	spanWrite(c.writer, "unmanagedname", t.Name())
}

func (c *HtmlColorer) Operator(t token.OperatorToken) {
	spanWrite(c.writer, "operator", t.Raw())
}

func (c *HtmlColorer) VariableSymbol(t token.VariableSymbolToken) {
	spanWrite(c.writer, "variable", t.Raw())
}

func (c *HtmlColorer) Definition(t token.VariableSymbolToken) {
	spanWrite(c.writer, "definition", t.Raw())
}

func (c *HtmlColorer) LocalType(t token.VariableSymbolToken) {
	spanWrite(c.writer, "localtype", t.Raw())
}

func (c *HtmlColorer) Parameter(t token.VariableSymbolToken) {
	spanWrite(c.writer, "parameter", t.Raw())
}

func (c *HtmlColorer) RecordField(t token.VariableSymbolToken) {
	spanWrite(c.writer, "recordfield", t.Raw())
}

func (c *HtmlColorer) TypeSymbol(t token.TypeSymbolToken) {
	spanWrite(c.writer, "typesymbol", t.Raw())
}

func (c *HtmlColorer) TypeGeneratorName(t token.TypeSymbolToken) {
	spanWrite(c.writer, "typegenerator", t.Raw())
}

func (c *HtmlColorer) ModuleReference(t token.TypeSymbolToken) {
	spanWrite(c.writer, "modulereference", t.Raw())
}

func (c *HtmlColorer) PrimitiveType(t token.TypeSymbolToken) {
	spanWrite(c.writer, "primitivetype", t.Raw())
}

func (c *HtmlColorer) AliasNameSymbol(t token.TypeSymbolToken) {
	spanWrite(c.writer, "alias", t.Raw())
}

func (c *HtmlColorer) NumberLiteral(t token.NumberToken) {
	spanWrite(c.writer, "number", t.Raw())
}

func (c *HtmlColorer) BooleanLiteral(t token.BooleanToken) {
	spanWrite(c.writer, "boolean", t.Raw())
}

func (c *HtmlColorer) KeywordString(t string) {
	spanWrite(c.writer, "keyword", t)
}

func (c *HtmlColorer) NewLine(indentation int) {
	fmt.Fprintf(c.writer, "\n")

	for i := 0; i < indentation; i++ {
		fmt.Fprintf(c.writer, "  ")
	}
}

func (c *HtmlColorer) OperatorString(t string) {
	spanWrite(c.writer, "operator", t)
}

func (c *HtmlColorer) StringLiteral(s token.StringToken) {
	spanWrite(c.writer, "string", s.Raw())
}

func (c *HtmlColorer) LocalTypeName(s *dectype.LocalType) {
	spanWrite(c.writer, "localtype", s.Identifier().Name())
}

func (c *HtmlColorer) OneSpace() {
	fmt.Fprintf(c.writer, " ")
}

func (c *HtmlColorer) RightPipe() {
	fmt.Fprintf(c.writer, "|>")
}

func (c *HtmlColorer) LeftPipe() {
	fmt.Fprintf(c.writer, "<|")
}

func (c *HtmlColorer) String() string {
	return "htmlcolorer"
}

func (c *HtmlColorer) UserInstruction(t string) {
}
