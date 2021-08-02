package doc

import (
	"fmt"
	"io"
	"sort"

	"github.com/gomarkdown/markdown"
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
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
	return fmt.Sprintf("<code class='%v'>%v</code>", className, value)
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

func markdownToHtml(markdownString string) string {
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
		alias, isAlias := localType.(*dectype.Alias)
		if isAlias {
			comment := alias.AstAlias().Comment()
			if comment != nil && comment.Token().ForDocumentation {
				filteredTypes[name] = alias
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

func writeHeader(writer io.Writer, fullyQualifiedName *decorated.FullyQualifiedPackageVariableName, p dtype.Type) {
	fmt.Fprintf(writer, "\n\n<h3>%v</h3>\n", fullyQualifiedName.Identifier().Name())

	fmt.Fprintf(writer, "<div class='prototype'><code>%v</code></div>\n", typeToHtmlBlock(p))
}

func ModuleToHtml(writer io.Writer, module *decorated.Module) {
	var markdownString string

	filteredTypes := filterTypes(module.LocalTypes().AllTypes())

	filteredFunctions, filteredConstants := filterDefinitions(module.LocalDefinitions().Definitions())
	if len(filteredFunctions) == 0 && len(filteredConstants) == 0 && len(filteredTypes) == 0 {
		return
	}

	fmt.Fprintf(writer, "\n\n<h2>Module %v</h2>\n", module.FullyQualifiedModuleName().Last())

	sortedConstantKeys := sortConstantKeys(filteredConstants)
	for _, constantName := range sortedConstantKeys {
		filteredConstant := filteredConstants[constantName]
		writeHeader(writer, constantName, filteredConstant.Type())
	}

	sortedTypeKeys := sortTypeKeys(filteredTypes)
	for _, localTypeName := range sortedTypeKeys {
		localType := filteredTypes[localTypeName]
		alias, isAlias := localType.(*dectype.Alias)
		if isAlias {
			comment := alias.AstAlias().Comment()
			if comment != nil && comment.Token().ForDocumentation {
				markdownString = comment.Token().Value()
				fmt.Fprintf(writer, "\n\n<h3>%v</h3>\n", alias.AstAlias().Name())
				commentDivWrite(writer, markdownString)
			}
		}
	}

	sortedFunctionKeys := sortFunctionKeys(filteredFunctions)
	for _, functionName := range sortedFunctionKeys {
		filteredFunction := filteredFunctions[functionName]

		commentBlock := filteredFunction.CommentBlock()
		writeHeader(writer, functionName, filteredFunction.Type())

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
