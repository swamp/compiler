package doc

import (
	"fmt"

	"github.com/gomarkdown/markdown"
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

func span(className string, value string) string {
	return fmt.Sprintf("<span class='%v'>%v</span>", className, value)
}

func code(className string, value string) string {
	return fmt.Sprintf("<code class='%v'>%v</code>", className, value)
}

func typeToHtml(typeToConvert dtype.Type) string {
	switch t := typeToConvert.(type) {
	case *dectype.Alias:
		return span("alias", t.AstAlias().Name())
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
		s := span("invoker", t.TypeGenerator().HumanReadable())
		for _, param := range t.Params() {
			s += " "
			s += typeToHtml(param)
		}
		return s
	case *dectype.FunctionAtom:
		s := "("
		for index, parameterType := range t.FunctionParameterTypes() {
			if index > 0 {
				s += span("arrow", " -> ")
			}
			s += typeToHtml(parameterType)
		}
		s += ")"
		return s
	case *dectype.FunctionTypeReference:
		return typeToHtml(t.Next())
	case *dectype.PrimitiveAtom:
		return span("primitive", t.PrimitiveName().Name())
	case *dectype.PrimitiveTypeReference:
		return typeToHtml(t.Next())
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

func ModuleToHtml(module *decorated.Module) string {
	html := ""

	var markdownString string

	for _, def := range module.LocalDefinitions().Definitions() {
		maybeFunction, _ := def.Expression().(*decorated.FunctionValue)
		if maybeFunction != nil {
			commentBlock := maybeFunction.CommentBlock()
			if commentBlock != nil && commentBlock.Token().ForDocumentation {
				html += fmt.Sprintf("\n\n<h3>%v</h3>\n", def.FullyQualifiedVariableName())

				html += fmt.Sprintf("<div class='prototype'><code>%v</code></div>\n", typeToHtmlBlock(maybeFunction.Type()))

				params := ""
				for index, arg := range maybeFunction.Parameters() {
					if index > 0 {
						params += " "
					}
					params += span("argument", arg.Identifier().Name())
				}

				html += code("params", params) + "\n"

				token := commentBlock.Token()
				markdownString = token.Value()

				html += commentDiv(markdownString)
			}
		}
	}

	for _, localType := range module.LocalTypes().AllTypes() {
		alias, isAlias := localType.(*dectype.Alias)
		if isAlias {
			comment := alias.AstAlias().Comment()
			if comment != nil && comment.Token().ForDocumentation {
				markdownString = comment.Token().Value()
				html += fmt.Sprintf("\n\n<h3>%v</h3>\n", alias.AstAlias().Name())
				html += commentDiv(markdownString)
			}
		}
	}

	return html
}
