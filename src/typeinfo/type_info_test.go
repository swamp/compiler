package typeinfo

import (
	"testing"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

func MakeFakeSourceFileReference() token.SourceFileReference {
	return token.NewInternalSourceFileReference()
}

func MakeFakeVariable(name string) *ast.VariableIdentifier {
	return ast.NewVariableIdentifier(token.NewVariableSymbolToken(name, MakeFakeSourceFileReference(), -1))
}

func MakeFakeLocalTypeName(name string) *ast.LocalTypeName {
	return ast.NewLocalTypeName(MakeFakeVariable(name))
}

func MakeFakeLocalTypeNameReference(name string) *ast.LocalTypeNameReference {
	return ast.NewLocalTypeNameReference(MakeFakeLocalTypeName(name),
		ast.NewLocalTypeNameDefinition(MakeFakeLocalTypeName(name)))
}

func MakeFakeTypeIdentifier(name string) *ast.TypeIdentifier {
	return ast.NewTypeIdentifier(token.NewTypeSymbolToken(name, MakeFakeSourceFileReference(), -1))
}

func MakeFakeAstTypeReferenceWithLocalTypeNames(name string, arguments []string) *ast.TypeReference {
	var astTypes []ast.Type
	for _, argumentTypeName := range arguments {
		astTypes = append(astTypes, MakeFakeLocalTypeName(argumentTypeName))
	}
	return ast.NewTypeReference(MakeFakeTypeIdentifier(name), astTypes)
}

func test(t *testing.T) {

}
