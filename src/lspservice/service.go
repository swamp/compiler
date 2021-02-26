package lspservice

import (
	"fmt"
	"log"
	"net/url"

	"github.com/piot/go-lsp"
	"github.com/piot/lsp-server/lspserv"

	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

/*
   "tokenTypes": [
       "namespace",
       "type",
       "class",
       "enum",
       "interface",
       "struct",
       "typeParameter",
       "parameter",
       "variable",
       "property",
       "enumMember",
       "event",
       "function",
       "method",
       "macro",
       "keyword",
       "modifier",
       "comment",
       "string",
       "number",
       "regexp",
       "operator"
   ],
   "tokenModifiers": [

*/
type SemanticBuilder struct {
	tokenTypes      []string
	tokenModifiers  []string
	lastRange       token.Range
	lastDebug       string
	encodedIntegers []uint
}

func NewSemanticBuilder() *SemanticBuilder {
	self := &SemanticBuilder{
		tokenTypes: []string{
			"namespace",
			"type",
			"class",
			"enum",
			"interface",
			"struct",
			"typeParameter",
			"parameter",
			"variable",
			"property",
			"enumMember",
			"event",
			"function",
			"method",
			"macro",
			"keyword",
			"modifier",
			"comment",
			"string",
			"number",
			"regexp",
			"operator",
		},
		tokenModifiers: []string{
			"declaration",
			"definition",
			"readonly",
			"static",
			"deprecated",
			"abstract",
			"async",
			"modification",
			"documentation",
			"defaultLibrary",
		},
		lastRange: token.NewPositionLength(token.MakePosition(0, 0), 0, 0),
	}
	return self
}

func FindInStrings(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}

	return -1
}

func (s *SemanticBuilder) EncodedValues() []uint {
	return s.encodedIntegers
}

func (s *SemanticBuilder) EncodeSymbol(debugString string, tokenRange token.Range, tokenType string, modifiers []string) error {
	if !tokenRange.IsAfter(s.lastRange) {
		return fmt.Errorf("they must be in order! %v %v and now %v %v", s.lastRange, s.lastDebug, tokenRange, debugString)
	}
	// log.Printf("adding symbol %v '%v'\n", tokenRange, debugString)

	tokenTypeId := FindInStrings(s.tokenTypes, tokenType)
	if tokenTypeId < 0 {
		return fmt.Errorf("unknown token type %v", tokenType)
	}

	var modifierBitMask uint
	for _, modifier := range modifiers {
		modifierId := FindInStrings(s.tokenModifiers, modifier)
		if modifierId < 0 {
			return fmt.Errorf("unknown token type %v", tokenType)
		}
		modifierBitMask |= 1 << modifierId
	}

	lastLine := s.lastRange.Position().Line()
	lastStartColumn := s.lastRange.Position().Column()

	deltaLine := uint(tokenRange.Position().Line() - lastLine)
	if deltaLine != 0 {
		lastStartColumn = 0
	}
	deltaColumnFromLastStartColumn := uint(tokenRange.Position().Column() - lastStartColumn)

	tokenLength := tokenRange.SingleLineLength()
	if tokenLength == -1 {
		return fmt.Errorf("token spans multiple lines %v %s", tokenType, debugString)
	}

	encodedIntegers := [5]uint{deltaLine, deltaColumnFromLastStartColumn, uint(tokenLength), uint(tokenTypeId), modifierBitMask}

	s.encodedIntegers = append(s.encodedIntegers, encodedIntegers[:]...)
	s.lastRange = tokenRange
	s.lastDebug = debugString

	return nil
}

type DecoratedTokenScanner interface {
	FindToken(position token.Position) decorated.TypeOrToken
	RootTokens() []decorated.TypeOrToken
}

type Compiler interface {
	Compile(filename string) error
}

type Service struct {
	scanner  DecoratedTokenScanner
	compiler Compiler
}

func NewService(compiler Compiler, scanner DecoratedTokenScanner) *Service {
	return &Service{scanner: scanner, compiler: compiler}
}

func (s *Service) Reset() error {
	return nil
}

func (s *Service) ResetCaches(lock bool) {
}

func (s *Service) ShutDown() {
}

func lspToTokenPosition(position lsp.Position) token.Position {
	return token.MakePosition(position.Line, position.Character)
}

func sourceFileReferenceToLocation(reference token.SourceFileReference) *lsp.Location {
	return &lsp.Location{
		URI:   lsp.DocumentURI(reference.Document.Uri),
		Range: *tokenToLspRange(reference.Range),
	}
}

func tokenToLspPosition(position token.Position) lsp.Position {
	return lsp.Position{
		Line:      position.Line(),
		Character: position.Column(),
	}
}

func tokenToLspRange(rangeToken token.Range) *lsp.Range {
	exclusiveEndPosition := lsp.Position{
		Line:      rangeToken.End().Line(),
		Character: rangeToken.End().Column() + 1,
	}

	return &lsp.Range{
		Start: tokenToLspPosition(rangeToken.Start()),
		End:   exclusiveEndPosition,
	}
}

func (s *Service) HandleHover(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Hover, error) {
	tokenPosition := lspToTokenPosition(params.Position)

	decoratedToken := s.scanner.FindToken(tokenPosition)

	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	_, isItAType := decoratedToken.(dtype.Type)

	log.Printf("the token is %T and type:%v\n", decoratedToken, isItAType)

	showString := decoratedToken.String()
	if !isItAType {
		normalToken, _ := decoratedToken.(decorated.Token)
		showString += fmt.Sprintf(" : %v", normalToken.Type().HumanReadable())
	}

	hover := &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.MUKMarkdown,
			Value: showString,
		},
		Range: tokenToLspRange(decoratedToken.FetchPositionLength().Range),
	}

	return hover, nil
}

func (s *Service) HandleGotoDefinition(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	tokenPosition := lspToTokenPosition(params.Position)

	decoratedToken := s.scanner.FindToken(tokenPosition)
	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	log.Printf("found: %T %v\n", decoratedToken, decoratedToken)
	var sourceFileReference token.SourceFileReference

	switch t := decoratedToken.(type) {
	case *decorated.Import:
		sourceFileReference = token.MakeSourceFileReference(token.MakeSourceFileDocumentFromURI(t.Module().DocumentURI()), token.MakeRange(token.MakePosition(0, 0), token.MakePosition(0, 0)))
	case *decorated.FunctionParameterReference:
		sourceFileReference = t.ParameterRef().FetchPositionLength()
	case *decorated.LetVariableReference:
		sourceFileReference = t.LetVariable().FetchPositionLength()
	case *decorated.FunctionReference:
		sourceFileReference = t.FunctionValue().FetchPositionLength()
	}

	if sourceFileReference.Document == nil {
		return nil, fmt.Errorf("couldn't get a reference for %T\n", decoratedToken)
	}

	location := sourceFileReferenceToLocation(sourceFileReference)

	return location, nil
}

func (s *Service) HandleGotoDeclaration(params lsp.DeclarationOptions, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, nil
}

func (s *Service) HandleGotoTypeDefinition(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, nil
}

func (s *Service) HandleGotoImplementation(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, nil
}

func (s *Service) HandleFindReferences(params lsp.ReferenceParams, conn lspserv.Connection) ([]*lsp.Location, error) {
	return nil, nil
}

func convertRootTokenToOutlineSymbol(rootToken decorated.TypeOrToken) *lsp.DocumentSymbol {
	switch t := rootToken.(type) {
	case *decorated.FunctionValue:
		return &lsp.DocumentSymbol{
			Name:           t.AstFunctionValue().DebugFunctionIdentifier().Name(),
			Detail:         t.Type().HumanReadable(),
			Kind:           lsp.SKFunction,
			Tags:           nil,
			Range:          *tokenToLspRange(t.FetchPositionLength().Range),
			SelectionRange: *tokenToLspRange(t.FetchPositionLength().Range),
			Children:       nil,
		}
	}
	return nil
}

func (s *Service) HandleSymbol(params lsp.DocumentSymbolParams, conn lspserv.Connection) ([]*lsp.DocumentSymbol, error) {
	rootTokens := s.scanner.RootTokens()

	var symbols []*lsp.DocumentSymbol

	for _, rootToken := range rootTokens {
		documentSymbol := convertRootTokenToOutlineSymbol(rootToken)
		if documentSymbol == nil {
			continue
		}
		symbols = append(symbols, documentSymbol)
	}
	return symbols, nil
} // Used for outline

func (s *Service) HandleCompletion(params lsp.CompletionParams, conn lspserv.Connection) (*lsp.CompletionList, error) {
	return nil, nil
} // Intellisense when pressing '.'.

func (s *Service) HandleCompletionItemResolve(params lsp.CompletionItem, conn lspserv.Connection) (*lsp.CompletionItem, error) {
	return nil, nil
}

func (s *Service) HandleSignatureHelp(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.SignatureHelp, error) {
	return nil, nil
}

func (s *Service) HandleFormatting(params lsp.DocumentFormattingParams, conn lspserv.Connection) ([]*lsp.TextEdit, error) {
	return nil, nil
}

// HandleHighlights :
func (s *Service) HandleHighlights(params lsp.DocumentHighlightParams,
	conn lspserv.Connection) ([]*lsp.DocumentHighlight, error) {
	tokenPosition := lspToTokenPosition(params.Position)

	decoratedToken := s.scanner.FindToken(tokenPosition)
	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	log.Printf("found: %T %v\n", decoratedToken, decoratedToken)
	var sourceFileReferences []token.SourceFileReference

	switch t := decoratedToken.(type) {
	case *decorated.Import:
		// sourceFileReference = token.MakeSourceFileReference(token.MakeSourceFileDocumentFromURI(t.Module().DocumentURI()), token.MakeRange(token.MakePosition(0, 0), token.MakePosition(0, 0)))
	case *decorated.FunctionParameterDefinition:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.FunctionParameterReference:
		sourceFileReferences = append(sourceFileReferences, t.ParameterRef().FetchPositionLength())
	case *decorated.LetVariable:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.LetVariableReference:
		sourceFileReferences = append(sourceFileReferences, t.LetVariable().FetchPositionLength())
	case *decorated.LetAssignment:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.FunctionValue:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.FunctionReference:
		sourceFileReferences = append(sourceFileReferences, t.FunctionValue().FetchPositionLength())
	}

	var hilights []*lsp.DocumentHighlight

	for _, reference := range sourceFileReferences {
		highlight := &lsp.DocumentHighlight{
			Range: *tokenToLspRange(reference.Range),
			Kind:  1, // Read only
		}

		hilights = append(hilights, highlight)
	}

	return hilights, nil
}

func (s *Service) HandleCodeAction(params lsp.CodeActionParams, conn lspserv.Connection) (*lsp.CodeAction, error) {
	return nil, nil
}

func (s *Service) HandleCodeActionResolve(params lsp.CodeAction, conn lspserv.Connection) (*lsp.CodeAction, error) {
	return nil, nil
}

func (s *Service) HandleRename(params lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	return nil, nil
}

/*
"namespace",
"type",
"class",
"enum",
"interface",
"struct",
"typeParameter",
"parameter",
"variable",
"property",
"enumMember",
"event",
"function",
"method",
"macro",
"keyword",
"modifier",
"comment",
"string",
"number",
"regexp",
"operator",
},
tokenModifiers: []string{
"declaration",
"definition",
"readonly",
"static",
"deprecated",
"abstract",
"async",
"modification",
"documentation",
"defaultLibrary",

*/

func addSemanticTokenFunctionValue(f *decorated.FunctionValue, builder *SemanticBuilder) error {
	functionNameIdentifier := f.AstFunctionValue().DebugFunctionIdentifier()
	functionNameRange := functionNameIdentifier.FetchPositionLength().Range

	if err := builder.EncodeSymbol(functionNameIdentifier.Name(), functionNameRange, "function", []string{"definition"}); err != nil {
		return err
	}

	for _, parameter := range f.Parameters() {
		if err := builder.EncodeSymbol(parameter.String(), parameter.FetchPositionLength().Range, "parameter", []string{}); err != nil {
			return err
		}
	}

	return addSemanticToken(f.Expression(), builder)
}

func addSemanticTokenAnnotation(f *decorated.Annotation, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(f.String(), f.Identifier().FetchPositionLength().Range, "function", []string{"declaration"}); err != nil {
		return err
	}
	if err := addSemanticToken(f.Type(), builder); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenFunctionType(f *dectype.FunctionAtom, builder *SemanticBuilder) error {
	for _, paramType := range f.FunctionParameterTypes() {
		if err := addSemanticToken(paramType, builder); err != nil {
			return err
		}
	}

	return nil
}

func addTypeReferencePrimitive(referenceRange token.Range, primitive *dectype.PrimitiveAtom, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(primitive.PrimitiveName().Name(), referenceRange, "type", []string{"declaration", "defaultLibrary"}); err != nil {
		return err
	}

	return nil
}

func IsBuiltInType(typeToCheckUnaliased dtype.Type) bool {
	typeToCheck := dectype.Unalias(typeToCheckUnaliased)
	switch t := typeToCheck.(type) {
	case *dectype.TypeReference:
		return IsBuiltInType(t.Next())
	case *dectype.InvokerType:
		typeToInvoke := t.TypeGenerator()
		typeRef, _ := typeToInvoke.(*dectype.TypeReference)
		if typeRef != nil {
			typeToInvoke = typeRef.Next()
		}
		typeToInvokeName := typeToInvoke.DecoratedName()
		return typeToInvokeName == "List" || typeToInvokeName == "Array"
	case *dectype.PrimitiveAtom:
		typeName := t.AtomName()
		return typeName == "Int" ||
			typeName == "Fixed" || typeName == "Bool" || typeName == "ResourceName" ||
			typeName == "TypeId" || typeName == "Blob"
	}

	return false
}

func addTypeReferenceInvoker(referenceRange token.Range, invoker *dectype.InvokerType, builder *SemanticBuilder) error {
	tokenModifiers := []string{"declaration"}
	if IsBuiltInType(invoker) {
		tokenModifiers = append(tokenModifiers, "defaultLibrary")
	}

	if err := builder.EncodeSymbol(invoker.TypeGenerator().HumanReadable(), referenceRange, "type", tokenModifiers); err != nil {
		return err
	}

	for _, param := range invoker.Params() {
		var tokenModifiersForParam []string
		if IsBuiltInType(param) {
			tokenModifiersForParam = append(tokenModifiersForParam, "defaultLibrary")
		}
		if err := builder.EncodeSymbol(param.HumanReadable(), param.FetchPositionLength().Range, "typeParameter", tokenModifiersForParam); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenImport(decoratedImport *decorated.Import, builder *SemanticBuilder) error {
	keyword := decoratedImport.AstImport().Keyword()
	if err := builder.EncodeSymbol(keyword.Name(), keyword.FetchPositionLength().Range, "keyword", nil); err != nil {
		return err
	}

	for _, segment := range decoratedImport.AstImport().Path() {
		if err := builder.EncodeSymbol(segment.Name(), segment.FetchPositionLength().Range, "namespace", nil); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenLet(decoratedLet *decorated.Let, builder *SemanticBuilder) error {
	keyword := decoratedLet.AstLet().Keyword()
	if err := builder.EncodeSymbol(keyword.Raw(), keyword.FetchPositionLength().Range, "keyword", nil); err != nil {
		return err
	}

	for _, assignment := range decoratedLet.Assignments() {
		if err := builder.EncodeSymbol(assignment.LetVariable().Name().Name(), assignment.LetVariable().FetchPositionLength().Range, "variable", nil); err != nil {
			return err
		}

		addSemanticToken(assignment.Expression(), builder)
	}

	if err := builder.EncodeSymbol(decoratedLet.AstLet().InKeyword().String(), decoratedLet.AstLet().InKeyword().FetchPositionLength().Range, "keyword", nil); err != nil {
		return err
	}

	return addSemanticToken(decoratedLet.Consequence(), builder)
}

func addSemanticTokenLetVariableReference(letVarReference *decorated.LetVariableReference, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(letVarReference.String(), letVarReference.FetchPositionLength().Range, "variable", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenString(stringLiteral *decorated.StringLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(stringLiteral.Value(), stringLiteral.FetchPositionLength().Range, "string", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenFixed(fixedLiteral *decorated.FixedLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(fixedLiteral.String(), fixedLiteral.FetchPositionLength().Range, "number", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenInteger(integerLiteral *decorated.IntegerLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(integerLiteral.String(), integerLiteral.FetchPositionLength().Range, "number", nil); err != nil {
		return err
	}
	return nil
}

func addSemanticTokenListLiteral(listLiteral *decorated.ListLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(listLiteral.String(), listLiteral.AstListLiteral().StartParenToken().Range, "operator", nil); err != nil {
		return err
	}

	for _, expression := range listLiteral.Expressions() {
		if err := addSemanticToken(expression, builder); err != nil {
			return err
		}
	}

	if err := builder.EncodeSymbol(listLiteral.String(), listLiteral.AstListLiteral().EndParenToken().Range, "operator", nil); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenArrayLiteral(arrayLiteral *decorated.ArrayLiteral, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(arrayLiteral.String(), arrayLiteral.AstArrayLiteral().StartParenToken().Range, "operator", nil); err != nil {
		return err
	}

	for _, expression := range arrayLiteral.Expressions() {
		if err := addSemanticToken(expression, builder); err != nil {
			return err
		}
	}

	if err := builder.EncodeSymbol(arrayLiteral.String(), arrayLiteral.AstArrayLiteral().EndParenToken().Range, "operator", nil); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenFunctionCall(funcCall *decorated.FunctionCall, builder *SemanticBuilder) error {
	if err := addSemanticToken(funcCall.FunctionValue(), builder); err != nil {
		return err
	}

	for _, argument := range funcCall.Arguments() {
		if err := addSemanticToken(argument, builder); err != nil {
			return err
		}
	}

	return nil
}

func addSemanticTokenFunctionReference(functionReference *decorated.FunctionReference, builder *SemanticBuilder) error {
	isScoped := functionReference.Identifier().ModuleReference() != nil
	if isScoped {
		for _, namespaceParts := range functionReference.Identifier().ModuleReference().Parts() {
			if err := builder.EncodeSymbol(namespaceParts.TypeIdentifier().Name(), namespaceParts.TypeIdentifier().FetchPositionLength().Range, "namespace", nil); err != nil {
				return err
			}
		}
	}

	if err := builder.EncodeSymbol(functionReference.String(), functionReference.Identifier().FetchPositionLength().Range, "function", nil); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenFunctionParameterReference(parameter *decorated.FunctionParameterReference, builder *SemanticBuilder) error {
	if err := builder.EncodeSymbol(parameter.String(), parameter.Identifier().FetchPositionLength().Range, "parameter", nil); err != nil {
		return err
	}

	return nil
}

func addSemanticTokenTypeReference(typeReference *dectype.TypeReference, builder *SemanticBuilder) error {
	next := typeReference.Next()

	log.Printf("typeReference %T %v\n", next, next)
	referenceRange := typeReference.FetchPositionLength().Range
	switch t := next.(type) {
	case *dectype.PrimitiveAtom:
		return addTypeReferencePrimitive(referenceRange, t, builder)
	case *dectype.InvokerType:
		return addTypeReferenceInvoker(referenceRange, t, builder)
	}

	return addSemanticToken(next, builder)
}

func addSemanticToken(typeOrToken decorated.TypeOrToken, builder *SemanticBuilder) error {
	switch t := typeOrToken.(type) {
	case *decorated.FunctionValue:
		return addSemanticTokenFunctionValue(t, builder)
	case *decorated.Annotation:
		return addSemanticTokenAnnotation(t, builder)
	case *dectype.TypeReference:
		return addSemanticTokenTypeReference(t, builder)
	case *dectype.FunctionAtom:
		return addSemanticTokenFunctionType(t, builder)
	case *dectype.InvokerType:
		return addTypeReferenceInvoker(t.FetchPositionLength().Range, t, builder)
	case *decorated.Import:
		return addSemanticTokenImport(t, builder)
	case *decorated.Let:
		return addSemanticTokenLet(t, builder)
	case *decorated.LetVariableReference:
		return addSemanticTokenLetVariableReference(t, builder)
	case *decorated.StringLiteral:
		return addSemanticTokenString(t, builder)
	case *decorated.FixedLiteral:
		return addSemanticTokenFixed(t, builder)
	case *decorated.IntegerLiteral:
		return addSemanticTokenInteger(t, builder)
	case *decorated.ListLiteral:
		return addSemanticTokenListLiteral(t, builder)
	case *decorated.ArrayLiteral:
		return addSemanticTokenArrayLiteral(t, builder)
	case *decorated.FunctionCall:
		return addSemanticTokenFunctionCall(t, builder)
	case *decorated.FunctionReference:
		return addSemanticTokenFunctionReference(t, builder)
	case *decorated.FunctionParameterReference:
		return addSemanticTokenFunctionParameterReference(t, builder)
	default:
		log.Printf("unknown %T %v\n", t, t)
	}

	return nil
}

func (s *Service) HandleSemanticTokensFull(params lsp.SemanticTokensParams, conn lspserv.Connection) (*lsp.SemanticTokens, error) {
	allTokens := s.scanner.RootTokens()
	builder := NewSemanticBuilder()
	for _, foundToken := range allTokens {
		if err := addSemanticToken(foundToken, builder); err != nil {
			return nil, err
		}
	}
	return &lsp.SemanticTokens{
		ResultId: "",
		Data:     builder.EncodedValues(),
	}, nil
}

func (s *Service) HandleCodeLens(params lsp.CodeLensParams, conn lspserv.Connection) ([]*lsp.CodeLens, error) {
	return []*lsp.CodeLens{
		/*{
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      4,
					Character: 0,
				},
				End: lsp.Position{
					Line:      4,
					Character: 4,
				},
			},
			Command: lsp.Command{
				Title:     "Some Command here",
				Command:   "swamp.somecommand",
				Arguments: nil,
			},
			Data: nil,
		},*/
	}, nil
}

func (s *Service) HandleCodeLensResolve(params lsp.CodeLens, conn lspserv.Connection) (*lsp.CodeLens, error) {
	return nil, nil
}

func (s *Service) HandleDidChangeWatchedFiles(params lsp.DidChangeWatchedFilesParams, conn lspserv.Connection) error {
	return nil
}

func (s *Service) HandleDidOpen(params lsp.DidOpenTextDocumentParams, conn lspserv.Connection) error {
	fullUrl, urlErr := url.Parse(string(params.TextDocument.URI))
	if urlErr != nil {
		return urlErr
	}
	s.compiler.Compile(fullUrl.Path)
	return nil
}
