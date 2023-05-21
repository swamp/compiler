/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package doc

import (
	"fmt"
	"io"
	"sort"

	"github.com/gomarkdown/markdown"
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/ast/codewriter"
	"github.com/swamp/compiler/src/coloring"
	"github.com/swamp/compiler/src/decorated/decoratedcodewriter"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func div(className string, value string) string {
	return fmt.Sprintf("<div class=\"%v\">%v</div>", className, value)
}

func descriptionDiv(value string) string {
	return div("description", value)
}

func commentDiv(markdownString string) (string, error) {
	markdownConverted, err := markdownToHtml(markdownString)
	if err != nil {
		return "", err
	}
	return descriptionDiv(markdownConverted), nil
}

func commentDivWrite(writer io.Writer, markdownString string) error {
	output, err := commentDiv(markdownString)
	if err != nil {
		return err
	}
	fmt.Fprint(writer, output)

	return nil
}

func span(className string, value string) string {
	return fmt.Sprintf("<span class=\"%v\">%v</span>", className, value)
}

func code(className string, value string) string {
	return fmt.Sprintf("<pre><code class=\"%v\">%v</code></pre>", className, value)
}

func codeWrite(writer io.Writer, className string, value string) {
	fmt.Fprintf(writer, code(className, value))
}

func classNameFromType(typeToConvert dtype.Type) string {
	switch t := typeToConvert.(type) {
	case *dectype.Alias:
		return classNameFromType(t.Next())
	case *dectype.AliasReference:
		return classNameFromType(t.Next())
	case *dectype.CustomTypeReference:
		return classNameFromType(t.Next())
	case *dectype.CustomTypeAtom:
		return "customtype"
	case *dectype.TupleTypeAtom:
		return "tuple"
	case *dectype.FunctionAtom:
		return "functiontype"
	case *dectype.FunctionTypeReference:
		return classNameFromType(t.Next())
	case *dectype.PrimitiveAtom:
		if t.PrimitiveName().Name() == "Any" {
			return "anytype"
		}

		return "primitivetype"
	case *dectype.PrimitiveTypeReference:
		return classNameFromType(t.Next())
	case *dectype.ResolvedLocalType:
		return "localtype"
	case *dectype.UnmanagedType:
		return "unmanagedtype"
	case *dectype.RecordAtom:
		return "recordtype"
	default:
		panic(fmt.Sprintf("can not get css class name from %T", typeToConvert))
	}
}

func typeToHtml(typeToConvert dtype.Type, colorer coloring.DecoratedColorer) {
	decoratedcodewriter.WriteType(typeToConvert, colorer, 0)
}

func typeToHtmlBlock(typeToConvert dtype.Type, colorer coloring.DecoratedColorer, writer io.Writer) {
	typeToHtml(typeToConvert, colorer)
	fmt.Fprintf(writer, "\n")
}

func expressionToHtmlBlock(expression ast.Expression, colorer coloring.Colorer) {
	codewriter.WriteExpression(expression, colorer, 0)
}

func markdownToHtml(rawMarkdownString string) (string, error) {
	markdownString, err := ConvertSwampMarkdown(rawMarkdownString)
	if err != nil {
		return "", err
	}
	raw := []byte(markdownString)
	outputRaw := markdown.ToHTML(raw, nil, nil)

	return string(outputRaw), nil
}

func shouldIncludeCommentBlock(commentBlock *ast.MultilineComment) bool {
	return commentBlock != nil && commentBlock.Token().ForDocumentation
}

func sortFunctionKeys(functions map[*decorated.FullyQualifiedPackageVariableName]*decorated.FunctionValue) []*decorated.FullyQualifiedPackageVariableName {
	keys := make([]*decorated.FullyQualifiedPackageVariableName, 0, len(functions))
	for k := range functions {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	return keys
}

func sortConstantKeys(constants map[*decorated.FullyQualifiedPackageVariableName]*decorated.Constant) []*decorated.FullyQualifiedPackageVariableName {
	keys := make([]*decorated.FullyQualifiedPackageVariableName, 0, len(constants))
	for k := range constants {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	return keys
}

func filterTypes(types []decorated.NamedType) []dtype.Type {
	var filteredTypes []dtype.Type
	for _, localType := range types {
		switch t := localType.RealType().(type) {
		case *dectype.Alias:
			comment := t.AstAlias().Comment()
			if comment != nil && comment.Token().ForDocumentation {
				filteredTypes = append(filteredTypes, t)
			}
		case *dectype.CustomTypeAtom:
			comment := t.AstCustomType().Comment()
			if comment != nil && comment.Token().ForDocumentation {
				filteredTypes = append(filteredTypes, t)
			}
		case *dectype.RecordAtom:
			{
				comment := t.AstRecord().Comment()
				if comment != nil && comment.Token().ForDocumentation {
					filteredTypes = append(filteredTypes, t)
				}
			}
		}
	}

	return filteredTypes
}

func filterDefinitions(definitions []decorated.ModuleDef) (map[*decorated.FullyQualifiedPackageVariableName]*decorated.FunctionValue, map[*decorated.FullyQualifiedPackageVariableName]*decorated.Constant) {
	filteredFunctions := make(map[*decorated.FullyQualifiedPackageVariableName]*decorated.FunctionValue)
	filteredConstants := make(map[*decorated.FullyQualifiedPackageVariableName]*decorated.Constant)

	for _, def := range definitions {
		switch t := def.Expression().(type) {
		case *decorated.FunctionValue:
			{
				if shouldIncludeCommentBlock(t.CommentBlock()) {
					filteredFunctions[def.FullyQualifiedVariableName()] = t
				}
			}
		case *decorated.Constant:
			{
				if shouldIncludeCommentBlock(t.CommentBlock()) {
					filteredConstants[def.FullyQualifiedVariableName()] = t
				}
			}
		}
	}

	return filteredFunctions, filteredConstants
}

func writeHeaderForType(writer io.Writer, colorer coloring.DecoratedColorer, fullyQualifiedName *decorated.FullyQualifiedPackageVariableName, p dtype.Type) {
	fmt.Fprintf(writer, "\n\n<h3 id='%v'>%v</h3>\n", fullyQualifiedName, fullyQualifiedName.Identifier().Name())

	fmt.Fprintf(writer, "<div class='swamp-function-prototype'><code>\n")

	typeToHtmlBlock(p, colorer, writer)

	fmt.Fprintf(writer, "\n</code></div>")
}

func writeHeaderForConstant(writer io.Writer, colorer coloring.Colorer, fullyQualifiedName *decorated.FullyQualifiedPackageVariableName, expression ast.Expression) {
	fmt.Fprintf(writer, "\n\n<h3 id='%v'>%v</h3>\n", fullyQualifiedName, fullyQualifiedName.Identifier().Name())

	fmt.Fprintf(writer, "<pre class='swamp-value'>\n")
	expressionToHtmlBlock(expression, colorer)
	fmt.Fprintf(writer, "</pre>\n")
}

func documentType(astType ast.Type, comment token.MultiLineCommentToken, colorer coloring.Colorer, writer io.Writer) error {
	fmt.Fprintf(writer, "\n\n<h3 id='%v'>%v</h3>\n", astType.Name(), astType.Name())
	fmt.Fprintf(writer, "<pre class='swamp'>")
	codewriter.WriteType(astType, colorer, false, 0)
	fmt.Fprintf(writer, "</pre>")
	markdownString := comment.Value()
	return commentDivWrite(writer, markdownString)
}

func findTypeNames(d dtype.Type) (string, dectype.ArtifactFullyQualifiedTypeName) {
	switch t := d.(type) {
	case *dectype.CustomTypeAtom:
		return t.Name(), t.ArtifactTypeName()
	case *dectype.Alias:
		return t.TypeIdentifier().Name(), t.ArtifactTypeName()
	case *dectype.PrimitiveTypeReference:
		moduleName := dectype.MakeModuleNameFromString(t.PrimitiveAtom().PrimitiveName().Name(), t.FetchPositionLength().Document)
		return t.AstIdentifier().SomeTypeIdentifier().Name(), dectype.ArtifactFullyQualifiedTypeName{ModuleName: moduleName}
	case *dectype.UnmanagedType:
		moduleName := dectype.MakeModuleNameFromString("builtin", t.FetchPositionLength().Document)
		return "Unmanaged", dectype.ArtifactFullyQualifiedTypeName{ModuleName: moduleName}
	case *dectype.RecordAtom:
		return "Record", dectype.ArtifactFullyQualifiedTypeName{}
	case *dectype.CustomTypeReference:
		return findTypeNames(t.CustomTypeAtom())
	default:
		panic(fmt.Errorf("unknown type to get name from %T", t))
	}
}

func documentDecoratedType(decoratedType dtype.Type, comment token.MultiLineCommentToken, colorer coloring.DecoratedColorer, writer io.Writer) error {
	shortName, fullyQualifiedName := findTypeNames(decoratedType)
	fmt.Fprintf(writer, "\n\n<h3 id='%v'>%v</h3>\n", fullyQualifiedName.String(), shortName)
	fmt.Fprintf(writer, "<pre><code class=\"swamp\">")
	decoratedcodewriter.WriteType(decoratedType, colorer, 0)
	fmt.Fprintf(writer, "</code></pre>")
	markdownString := comment.Value()
	return commentDivWrite(writer, markdownString)
}

func documentConstant(name *decorated.FullyQualifiedPackageVariableName, constant *decorated.Constant, colorer coloring.Colorer, writer io.Writer) error {
	expression := constant.AstConstant().Expression()
	writeHeaderForConstant(writer, colorer, name, expression)
	comment := constant.CommentBlock()
	markdownString := comment.Value()
	return commentDivWrite(writer, markdownString)
}

func ModuleToHtml(writer io.Writer, module *decorated.Module) error {
	var markdownString string

	colorer := &HtmlColorer{writer: writer}

	filteredTypes := filterTypes(module.LocalTypes().AllInOrderTypes())

	filteredFunctions, filteredConstants := filterDefinitions(module.LocalDefinitions().Definitions())
	if len(filteredFunctions) == 0 && len(filteredConstants) == 0 && len(filteredTypes) == 0 {
		return nil
	}

	fmt.Fprintf(writer, "\n\n<h3>%v</h3>\n", module.FullyQualifiedModuleName())

	sortedConstantKeys := sortConstantKeys(filteredConstants)
	for _, constantName := range sortedConstantKeys {
		filteredConstant := filteredConstants[constantName]
		if err := documentConstant(constantName, filteredConstant, colorer, writer); err != nil {
			return err
		}
	}

	//sortedTypeKeys := filteredTypes //sortTypeKeys(filteredTypes)
	for _, localType := range filteredTypes {
		switch t := localType.(type) {
		case *dectype.Alias:
			{
				comment := t.AstAlias().Comment()
				if err := documentDecoratedType(t, comment.Token(), colorer, writer); err != nil {
					return err
				}
			}
		case *dectype.CustomTypeAtom:
			{
				comment := t.AstCustomType().Comment()
				if err := documentDecoratedType(t, comment.Token(), colorer, writer); err != nil {
					return err
				}
			}
		case *dectype.RecordAtom:
			{
				comment := t.AstRecord().Comment()
				if err := documentDecoratedType(t, comment.Token(), colorer, writer); err != nil {
					return err
				}
			}
		default:
			{
				panic(fmt.Errorf("type %T is not handled", t))
			}
		}
	}

	sortedFunctionKeys := sortFunctionKeys(filteredFunctions)
	for _, functionName := range sortedFunctionKeys {
		filteredFunction := filteredFunctions[functionName]

		commentBlock := filteredFunction.CommentBlock()
		writeHeaderForType(writer, colorer, functionName, filteredFunction.Type())

		params := ""
		for index, arg := range filteredFunction.Parameters() {
			if index > 0 {
				params += " "
			}
			params += span("argument", arg.Parameter().Name())
		}

		codeWrite(writer, "params", params)
		fmt.Fprintln(writer, "")

		token := commentBlock.Token()
		markdownString = token.Value()

		if err := commentDivWrite(writer, markdownString); err != nil {
			return err
		}
	}

	return nil
}
