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
	cssClassName := classNameFromType(typeToWrite)
	fmt.Fprintf(writer, "<a href='#%v'>", fullyQualifiedName.String())
	fmt.Fprintf(writer, span(cssClassName, value))
	fmt.Fprintf(writer, "</a>")
}

func (c *HtmlColorer) CustomType(t *dectype.CustomTypeAtom) {
	spanWrite(c.writer, "customtypename", t.Name())
}

func (c *HtmlColorer) CustomTypeVariant(t *dectype.CustomTypeVariant) {
	spanWrite(c.writer, "customtypevariant", t.Name().Name())
}

func (c *HtmlColorer) InvokerType(t *dectype.InvokerType) {
	spanWriteLink(c.writer, t.TypeGenerator(), t.TypeGenerator().HumanReadable())
}

func (c *HtmlColorer) RecordTypeField(t *dectype.RecordField) {
	spanWrite(c.writer, "recordtypefield", t.Name())
}

func (c *HtmlColorer) AliasName(t *dectype.Alias) {
	aliasName := classNameFromType(t)
	spanWrite(c.writer, aliasName, t.TypeIdentifier().Name())
}

func (c *HtmlColorer) TypeName(t *ast.TypeIdentifier) {
	spanWrite(c.writer, "typename", t.Name())
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
	// fmt.Fprintf(c.writer, "<a href='#%v'>", t.Raw())
	spanWrite(c.writer, "typesymbol", t.Raw())
	// fmt.Fprintf(c.writer, "</a>")
}

func (c *HtmlColorer) TypeGeneratorName(t token.TypeSymbolToken) {
	spanWrite(c.writer, "typegenerator", t.Raw())
}

func (c *HtmlColorer) ModuleReference(t token.TypeSymbolToken) {
	spanWrite(c.writer, "modulereference", t.Raw())
}

func (c *HtmlColorer) PrimitiveType(t token.TypeSymbolToken) {
	spanWrite(c.writer, "primitive", t.Raw())
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

func (c *HtmlColorer) OneSpace() {
	fmt.Fprintf(c.writer, " ")
}

func (c *HtmlColorer) RightArrow() {
	fmt.Fprintf(c.writer, "->")
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
