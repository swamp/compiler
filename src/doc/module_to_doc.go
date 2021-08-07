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
	return fmt.Sprintf("<div class='%v'>%v</div>", className, value)
}

func descriptionDiv(value string) string {
	return div("description", value)
}

func commentDiv(markdownString string) string {
	return descriptionDiv(markdownToHtml(markdownString))
}

func commentDivWrite(writer io.Writer, markdownString string) {
	fmt.Fprint(writer, commentDiv(markdownString))
}

func span(className string, value string) string {
	return fmt.Sprintf("<span class='%v'>%v</span>", className, value)
}

func code(className string, value string) string {
	return fmt.Sprintf("<pre><code class='%v'>%v</code></pre>", className, value)
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
	case *dectype.InvokerType:
		return "invokertype"
	case *dectype.FunctionAtom:
		return "functiontype"
	case *dectype.FunctionTypeReference:
		return classNameFromType(t.Next())
	case *dectype.PrimitiveAtom:
		return "primitive"
	case *dectype.PrimitiveTypeReference:
		return classNameFromType(t.Next())
	case *dectype.LocalType:
		return "localtype"
	case *dectype.UnmanagedType:
		return "unmanagedtype"
	case *dectype.RecordAtom:
		return "recordtype"
	default:
		panic(fmt.Sprintf("can not get css class name from %T", typeToConvert))
	}
}

func typeToHtml(typeToConvert dtype.Type) string {
	switch t := typeToConvert.(type) {
	case *dectype.Alias:
		className := classNameFromType(t.Next())
		if className == "primitive" {
			className = "alias"
		}
		return span(className, t.AstAlias().Name())

	case *dectype.AliasReference:
		return typeToHtml(t.Next())
	case *dectype.CustomTypeReference:
		return typeToHtml(t.Next())
	case *dectype.CustomTypeAtom:
		return span("customtype", t.AstCustomType().Name())
	case *dectype.TupleTypeAtom:
		s := span("paren", "(")
		for index, parameterType := range t.ParameterTypes() {
			if index > 0 {
				s += span("comma", ", ")
			}
			s += typeToHtml(parameterType)
		}
		s += span("paren", ")")
		return s
	case *dectype.InvokerType:
		s := span("invokertype", t.TypeGenerator().HumanReadable())
		for _, param := range t.Params() {
			s += " "
			s += typeToHtml(param)
		}
		return s
	case *dectype.FunctionAtom:
		s := span("paren", "(")
		for index, parameterType := range t.FunctionParameterTypes() {
			if index > 0 {
				s += span("arrow", " &#8594; ")
			}
			s += typeToHtml(parameterType)
		}
		s += span("paren", ")")
		return s
	case *dectype.FunctionTypeReference:
		return typeToHtml(t.Next())
	case *dectype.PrimitiveAtom:
		return span("primitive", t.PrimitiveName().Name())
	case *dectype.PrimitiveTypeReference:
		return typeToHtml(t.Next())
	case *dectype.LocalType:
		return span("localtype", t.Identifier().Name())
	default:
		panic(fmt.Sprintf("can not understand %T", typeToConvert))
	}
}

func typeToHtmlBlock(typeToConvert dtype.Type) string {
	return typeToHtml(typeToConvert) + "\n"
}

func expressionToHtmlBlock(expression ast.Expression, colorer coloring.Colorer) {
	codewriter.WriteExpression(expression, colorer, 0)
}

func markdownToHtml(rawMarkdownString string) string {
	markdownString := ConvertSwampMarkdown(rawMarkdownString)
	raw := []byte(markdownString)
	outputRaw := markdown.ToHTML(raw, nil, nil)

	return string(outputRaw)
}

func shouldIncludeCommentBlock(commentBlock *ast.MultilineComment) bool {
	return commentBlock != nil && commentBlock.Token().ForDocumentation
}

func sortTypeKeys(types map[string]dtype.Type) []string {
	keys := make([]string, 0, len(types))
	for k := range types {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
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

func filterTypes(types map[string]dtype.Type) map[string]dtype.Type {
	filteredTypes := make(map[string]dtype.Type)
	for name, localType := range types {
		switch t := localType.(type) {
		case *dectype.Alias:
			comment := t.AstAlias().Comment()
			if comment != nil && comment.Token().ForDocumentation {
				filteredTypes[name] = t
			}
		case *dectype.CustomTypeAtom:
			comment := t.AstCustomType().Comment()
			if comment != nil && comment.Token().ForDocumentation {
				filteredTypes[name] = t
			}
		case *dectype.RecordAtom:
			{
				comment := t.AstRecord().Comment()
				if comment != nil && comment.Token().ForDocumentation {
					filteredTypes[name] = t
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

func writeHeaderForType(writer io.Writer, fullyQualifiedName *decorated.FullyQualifiedPackageVariableName, p dtype.Type) {
	fmt.Fprintf(writer, "\n\n<h3 id='%v'>%v</h3>\n", fullyQualifiedName, fullyQualifiedName.Identifier().Name())

	fmt.Fprintf(writer, "<div class='swamp-function-prototype'><code>%v</code></div>\n", typeToHtmlBlock(p))
}

func writeHeaderForConstant(writer io.Writer, colorer coloring.Colorer, fullyQualifiedName *decorated.FullyQualifiedPackageVariableName, expression ast.Expression) {
	fmt.Fprintf(writer, "\n\n<h3 id='%v'>%v</h3>\n", fullyQualifiedName, fullyQualifiedName.Identifier().Name())

	fmt.Fprintf(writer, "<pre class='swamp-value'>\n")
	expressionToHtmlBlock(expression, colorer)
	fmt.Fprintf(writer, "</pre>\n")
}

func documentType(astType ast.Type, comment token.MultiLineCommentToken, colorer coloring.Colorer, writer io.Writer) {
	fmt.Fprintf(writer, "\n\n<h3 id='%v'>%v</h3>\n", astType.Name(), astType.Name())
	fmt.Fprintf(writer, "<pre class='swamp'>")
	codewriter.WriteType(astType, colorer, false, 0)
	fmt.Fprintf(writer, "</pre>")
	markdownString := comment.Value()
	commentDivWrite(writer, markdownString)
}

func findTypeNames(d dtype.Type) (string, dectype.ArtifactFullyQualifiedTypeName) {
	switch t := d.(type) {
	case *dectype.CustomTypeAtom:
		return t.Name(), t.ArtifactTypeName()
	case *dectype.Alias:
		return t.TypeIdentifier().Name(), t.ArtifactTypeName()
	case *dectype.PrimitiveTypeReference:
		moduleName := dectype.MakeModuleNameFromString(t.PrimitiveAtom().PrimitiveName().Name())
		return t.AstIdentifier().SomeTypeIdentifier().Name(), dectype.ArtifactFullyQualifiedTypeName{ModuleName: moduleName}
	default:
		panic(fmt.Errorf("unknown type to get name from %T", t))
	}
}

func documentDecoratedType(decoratedType dtype.Type, comment token.MultiLineCommentToken, colorer coloring.DecoratedColorer, writer io.Writer) {
	shortName, fullyQualifiedName := findTypeNames(decoratedType)
	fmt.Fprintf(writer, "\n\n<h3 id='%v'>%v</h3>\n", fullyQualifiedName.String(), shortName)
	fmt.Fprintf(writer, "<pre><code class='swamp'>")
	decoratedcodewriter.WriteType(decoratedType, colorer, 0)
	fmt.Fprintf(writer, "</code></pre>")
	markdownString := comment.Value()
	commentDivWrite(writer, markdownString)
}

func documentConstant(name *decorated.FullyQualifiedPackageVariableName, constant *decorated.Constant, colorer coloring.Colorer, writer io.Writer) {
	expression := constant.AstConstant().Expression()
	writeHeaderForConstant(writer, colorer, name, expression)
	comment := constant.CommentBlock()
	markdownString := comment.Value()
	commentDivWrite(writer, markdownString)
}

func ModuleToHtml(writer io.Writer, module *decorated.Module) {
	var markdownString string

	colorer := &HtmlColorer{writer: writer}

	filteredTypes := filterTypes(module.LocalTypes().AllTypes())

	filteredFunctions, filteredConstants := filterDefinitions(module.LocalDefinitions().Definitions())
	if len(filteredFunctions) == 0 && len(filteredConstants) == 0 && len(filteredTypes) == 0 {
		return
	}

	fmt.Fprintf(writer, "\n\n<h2>Module %v</h2>\n", module.FullyQualifiedModuleName().Last())

	sortedConstantKeys := sortConstantKeys(filteredConstants)
	for _, constantName := range sortedConstantKeys {
		filteredConstant := filteredConstants[constantName]
		documentConstant(constantName, filteredConstant, colorer, writer)
	}

	sortedTypeKeys := sortTypeKeys(filteredTypes)
	for _, localTypeName := range sortedTypeKeys {
		localType := filteredTypes[localTypeName]
		switch t := localType.(type) {
		case *dectype.Alias:
			{
				comment := t.AstAlias().Comment()
				documentDecoratedType(t, comment.Token(), colorer, writer)
			}
		case *dectype.CustomTypeAtom:
			{
				comment := t.AstCustomType().Comment()
				documentDecoratedType(t, comment.Token(), colorer, writer)
			}
		case *dectype.RecordAtom:
			{
				comment := t.AstRecord().Comment()
				documentDecoratedType(t, comment.Token(), colorer, writer)
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
		writeHeaderForType(writer, functionName, filteredFunction.Type())

		params := ""
		for index, arg := range filteredFunction.Parameters() {
			if index > 0 {
				params += " "
			}
			params += span("argument", arg.Identifier().Name())
		}

		codeWrite(writer, "params", params)
		fmt.Fprintln(writer, "")

		token := commentBlock.Token()
		markdownString = token.Value()

		commentDivWrite(writer, markdownString)
	}
}
