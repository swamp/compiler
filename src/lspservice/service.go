package lspservice

import (
	"fmt"
	"log"
	"os"

	"github.com/piot/go-lsp"
	"github.com/piot/lsp-server/lspserv"

	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type DecoratedTokenScanner interface {
	FindToken(documentURI token.DocumentURI, position token.Position) decorated.TypeOrToken
	RootTokens(documentURI token.DocumentURI) []decorated.TypeOrToken
}

type Compiler interface {
	Compile(filename string) (*decorated.Module, error)
}

type DocumentCacher interface {
	GetDocument(filename LocalFileSystemPath, version DocumentVersion) (*InMemoryDocument, error)
}

type Workspacer interface {
	AllModules() []*decorated.Module
}

type Service struct {
	scanner     DecoratedTokenScanner
	compiler    Compiler
	documents   DocumentCacher
	workspacer  Workspacer
	diagnostics *DiagnosticsForDocuments
}

func NewService(compiler Compiler, scanner DecoratedTokenScanner, documents DocumentCacher, workspacer Workspacer) *Service {
	diagnostics := NewDiagnosticsForDocuments()
	return &Service{scanner: scanner, compiler: compiler, documents: documents, workspacer: workspacer, diagnostics: diagnostics}
}

func (s *Service) Reset() error {
	return nil
}

func (s *Service) ResetCaches(lock bool) {
}

func (s *Service) ShutDown() {
}

func toDocumentURI(URI lsp.DocumentURI) token.DocumentURI {
	return token.DocumentURI(URI)
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

func lspToTokenRange(lspRange lsp.Range) token.Range {
	return token.MakeRange(
		token.MakePosition(lspRange.Start.Line, lspRange.Start.Character),
		token.MakePosition(lspRange.End.Line, lspRange.End.Character))
}

func (s *Service) HandleHover(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Hover, error) {
	tokenPosition := lspToTokenPosition(params.Position)
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	decoratedToken := s.scanner.FindToken(sourceFileURI, tokenPosition)

	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	tokenType, isItAType := decoratedToken.(dtype.Type)

	log.Printf("the token is %T and type:%v\n", decoratedToken, isItAType)

	var codeSignature string
	var name string
	var documentation string
	if isItAType {
		codeSignature = tokenType.HumanReadable()
		name = "type"
	} else {
		normalToken, _ := decoratedToken.(decorated.Token)
		if normalToken == nil {
			onlyHumanReadable, _ := decoratedToken.(decorated.HumanReadEnabler)
			if onlyHumanReadable != nil {
				name = onlyHumanReadable.HumanReadable()
			} else {
				log.Printf("not sure what this token is: %T\n", decoratedToken)
				return nil, nil
			}
		} else {
			codeSignature = normalToken.Type().HumanReadable()
			name = normalToken.HumanReadable()
			switch t := normalToken.(type) {
			case *decorated.FunctionReference:
				if t.FunctionValue().CommentBlock() != nil {
					documentation = t.FunctionValue().CommentBlock().Value()
				}
			}
		}
	}

	showString := fmt.Sprintf("`%v`\n\n%v\n\n ```swamp\n%v\n```\n", name, documentation, codeSignature)

	hover := &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.MUKMarkdown,
			Value: showString,
		},
		Range: tokenToLspRange(decoratedToken.FetchPositionLength().Range),
	}

	return hover, nil
}

func tokenToDefinition(decoratedToken decorated.TypeOrToken) (token.SourceFileReference, error) {
	switch t := decoratedToken.(type) {
	case *decorated.ImportStatement:
		return t.Module().FetchPositionLength(), nil
	case *decorated.FunctionParameterReference:
		return t.ParameterRef().FetchPositionLength(), nil
	case *decorated.LetVariableReference:
		return t.LetVariable().FetchPositionLength(), nil
	case *decorated.FunctionReference:
		return t.FunctionValue().FetchPositionLength(), nil
	case *decorated.ModuleReference:
		return t.Module().FetchPositionLength(), nil
	case *decorated.RecordFieldReference:
		return t.RecordTypeField().VariableIdentifier().FetchPositionLength(), nil
	case *decorated.CustomTypeVariantReference:
		return t.CustomTypeVariant().FetchPositionLength(), nil
	case *decorated.RecordConstructor:
		return t.Type().FetchPositionLength(), nil
	case *decorated.CustomTypeVariantConstructor:
		return tokenToDefinition(t.Reference())
	case *decorated.FunctionValue:
		return t.FetchPositionLength(), nil
	case *decorated.FunctionParameterDefinition:
		return t.FetchPositionLength(), nil
	case *decorated.Let:
		return t.FetchPositionLength(), nil
	case *decorated.LetVariable:
		return t.FetchPositionLength(), nil
	case *decorated.FunctionCall:
		return tokenToDefinition(t.FunctionExpression())
	case *decorated.CurryFunction:
		return tokenToDefinition(t.FunctionValue())
	case *decorated.Constant:
		return t.FunctionValue().AstFunctionValue().DebugFunctionIdentifier().FetchPositionLength(), nil
	// TYPES
	case *dectype.FunctionTypeReference:
		return t.FunctionAtom().FetchPositionLength(), nil
	case *dectype.TypeReference:
		return t.Next().FetchPositionLength(), nil
	case *dectype.TypeReferenceScoped:
		return t.Next().FetchPositionLength(), nil
	}

	return token.SourceFileReference{}, fmt.Errorf("couldn't find anything for %T", decoratedToken)
}

func (s *Service) HandleGotoDefinition(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	tokenPosition := lspToTokenPosition(params.Position)
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	decoratedToken := s.scanner.FindToken(sourceFileURI, tokenPosition)
	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	sourceFileReference, lookupErr := tokenToDefinition(decoratedToken)
	if lookupErr != nil {
		log.Printf("couldn't find any definition for %T\n", decoratedToken)
		return nil, nil
	}

	if sourceFileReference.Document == nil {
		log.Printf("couldn't go to definition for %T, no document \n", decoratedToken)
		return nil, nil
	}

	location := sourceFileReferenceToLocation(sourceFileReference)
	log.Printf("definition for %T resulted in %v \n", decoratedToken, sourceFileReference)

	return location, nil
}

func (s *Service) HandleLinkedEditingRange(params lsp.LinkedEditingRangeParams, conn lspserv.Connection) (*lsp.LinkedEditingRanges, error) {
	tokenPosition := lspToTokenPosition(params.Position)
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	decoratedToken := s.scanner.FindToken(sourceFileURI, tokenPosition)
	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return &lsp.LinkedEditingRanges{
			Ranges:      []lsp.Range{},
			WordPattern: nil,
		}, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	documentURI := token.MakeDocumentURI(string(params.TextDocument.URI))
	sourceFileReferences := findAllLinkedSymbolsInDocument(decoratedToken, documentURI)

	var renameRanges []lsp.Range

	for _, ref := range sourceFileReferences {
		renameRanges = append(renameRanges, *tokenToLspRange(ref.Range))
	}

	return &lsp.LinkedEditingRanges{
		Ranges:      renameRanges,
		WordPattern: nil,
	}, nil
}

func (s *Service) HandleGotoDeclaration(params lsp.DeclarationOptions, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, fmt.Errorf("concept of go to declaration not in the Swamp language")
}

func (s *Service) HandleGotoTypeDefinition(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, nil
}

func (s *Service) HandleGotoImplementation(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, fmt.Errorf("concept of go to implementation not in the Swamp language")
}

func findReferences(uri lsp.DocumentURI, position lsp.Position, scanner DecoratedTokenScanner) ([]token.SourceFileReference, error) {
	tokenPosition := lspToTokenPosition(position)
	sourceFileURI := toDocumentURI(uri)
	decoratedToken := scanner.FindToken(sourceFileURI, tokenPosition)
	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	var sourceFileReferences []token.SourceFileReference

	switch t := decoratedToken.(type) {
	case *decorated.ImportStatement:
		// sourceFileReference = token.MakeSourceFileReference(token.MakeSourceFileDocumentFromURI(t.Module().DocumentURI()), token.MakeRange(token.MakePosition(0, 0), token.MakePosition(0, 0)))
	case *decorated.FunctionParameterDefinition:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.LetVariable:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.LetAssignment:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.FunctionValue:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.FunctionName:
		for _, ref := range t.FunctionValue().References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	}

	return sourceFileReferences, nil
}

func (s *Service) HandleFindReferences(params lsp.ReferenceParams, conn lspserv.Connection) ([]*lsp.Location, error) {
	sourceFileReferences, err := findReferences(params.TextDocument.URI, params.Position, s.scanner)
	if err != nil {
		return nil, err
	}

	var locations []*lsp.Location

	for _, reference := range sourceFileReferences {
		location := sourceFileReferenceToLocation(reference)
		locations = append(locations, location)
	}

	return locations, nil
}

func functionParametersToDocumentSymbols(parameters []*decorated.FunctionParameterDefinition) []lsp.DocumentSymbol {
	var symbols []lsp.DocumentSymbol

	for _, param := range parameters {
		symbol := lsp.DocumentSymbol{
			Name:           param.Identifier().Name(),
			Detail:         param.Type().HumanReadable(),
			Kind:           lsp.SKVariable,
			Tags:           nil,
			Range:          *tokenToLspRange(param.FetchPositionLength().Range),
			SelectionRange: *tokenToLspRange(param.FetchPositionLength().Range),
			Children:       nil,
		}
		symbols = append(symbols, symbol)
	}

	return symbols
}

func convertRootTokenToOutlineSymbol(rootToken decorated.TypeOrToken) *lsp.DocumentSymbol {
	switch t := rootToken.(type) {
	case *decorated.NamedFunctionValue:
		return &lsp.DocumentSymbol{
			Name:           t.Value().AstFunctionValue().DebugFunctionIdentifier().Name(),
			Detail:         t.Value().Type().HumanReadable(),
			Kind:           lsp.SKFunction,
			Tags:           nil,
			Range:          *tokenToLspRange(t.FetchPositionLength().Range),
			SelectionRange: *tokenToLspRange(t.FunctionName().FetchPositionLength().Range),
			Children:       functionParametersToDocumentSymbols(t.Value().Parameters()),
		}
	}
	return nil
}

func (s *Service) HandleSymbol(params lsp.DocumentSymbolParams, conn lspserv.Connection) ([]*lsp.DocumentSymbol, error) {
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	rootTokens := s.scanner.RootTokens(sourceFileURI)

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

func findAllLinkedSymbolsInDocument(decoratedToken decorated.TypeOrToken, filterDocument token.DocumentURI) []token.SourceFileReference {
	var sourceFileReferences []token.SourceFileReference

	switch t := decoratedToken.(type) {
	case *decorated.ImportStatement:
		// sourceFileReference = token.MakeSourceFileReference(token.MakeSourceFileDocumentFromURI(t.Module().DocumentURI()), token.MakeRange(token.MakePosition(0, 0), token.MakePosition(0, 0)))
	case *decorated.FunctionParameterDefinition:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.FetchPositionLength())
		}
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.FunctionParameterReference:
		return findAllLinkedSymbolsInDocument(t.ParameterRef(), filterDocument)

	case *decorated.LetVariable:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.FetchPositionLength())
		}
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.LetVariableReference:
		return findAllLinkedSymbolsInDocument(t.LetVariable(), filterDocument)

		/*
			case *decorated.LetAssignment:
				for _, ref := range t.References() {
					if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
						sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
					}
				}

		*/

	case *decorated.FunctionValue:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.AstFunctionValue().DebugFunctionIdentifier().FetchPositionLength())
		}
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.FunctionReference:
		return findAllLinkedSymbolsInDocument(t.FunctionValue(), filterDocument)
	}

	return sourceFileReferences
}

func findLinkedSymbolsInDocument(decoratedToken decorated.TypeOrToken, filterDocument token.DocumentURI) []token.SourceFileReference {
	var sourceFileReferences []token.SourceFileReference

	switch t := decoratedToken.(type) {
	case *decorated.ImportStatement:
		// sourceFileReference = token.MakeSourceFileReference(token.MakeSourceFileDocumentFromURI(t.Module().DocumentURI()), token.MakeRange(token.MakePosition(0, 0), token.MakePosition(0, 0)))
	case *decorated.FunctionParameterDefinition:
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.FunctionParameterReference:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.ParameterRef().FetchPositionLength())
		}
	case *decorated.LetVariable:
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.LetVariableReference:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.LetVariable().FetchPositionLength())
		}
	case *decorated.LetAssignment:
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.FunctionValue:
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.FunctionReference:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.FunctionValue().FetchPositionLength())
		}
	}

	return sourceFileReferences
}

// HandleHighlights :
func (s *Service) HandleHighlights(params lsp.DocumentHighlightParams,
	conn lspserv.Connection) ([]*lsp.DocumentHighlight, error) {
	tokenPosition := lspToTokenPosition(params.Position)
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	decoratedToken := s.scanner.FindToken(sourceFileURI, tokenPosition)
	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	documentURI := token.MakeDocumentURI(string(params.TextDocument.URI))
	sourceFileReferences := findLinkedSymbolsInDocument(decoratedToken, documentURI)

	var highlights []*lsp.DocumentHighlight

	for _, reference := range sourceFileReferences {
		highlight := &lsp.DocumentHighlight{
			Range: *tokenToLspRange(reference.Range),
			Kind:  1, // Read only
		}

		highlights = append(highlights, highlight)
	}

	return highlights, nil
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

func (s *Service) HandleSemanticTokensFull(params lsp.SemanticTokensParams, conn lspserv.Connection) (*lsp.SemanticTokens, error) {
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	fmt.Fprintf(os.Stderr, "get root tokens\n")
	allTokens := s.scanner.RootTokens(sourceFileURI)
	fmt.Fprintf(os.Stderr, "root tokens done %v\n", len(allTokens))
	builder := NewSemanticBuilder()
	for _, foundToken := range allTokens {
		if err := addSemanticToken(foundToken, builder); err != nil {
			log.Printf("was problem with token %v %v", foundToken, err)
			return nil, err
		}
	}
	fmt.Fprintf(os.Stderr, "added all tokens done\n")
	return &lsp.SemanticTokens{
		ResultId: "",
		Data:     builder.EncodedValues(),
	}, nil
}

func (s *Service) HandleCodeLens(params lsp.CodeLensParams, conn lspserv.Connection) ([]*lsp.CodeLens, error) {
	var codeLenses []*lsp.CodeLens

	for _, rootToken := range s.scanner.RootTokens(toDocumentURI(params.TextDocument.URI)) {
		switch t := rootToken.(type) {
		case *decorated.NamedFunctionValue:
			value := t.Value()
			count := len(value.References())
			var textToDisplay string
			if count == 0 {
				textToDisplay = "no references"
			} else if count == 1 {
				textToDisplay = "one reference"
			} else {
				textToDisplay = fmt.Sprintf("%d references", count)
			}
			codeLens := &lsp.CodeLens{
				Range: *tokenToLspRange(value.Annotation().FetchPositionLength().Range),
				Command: lsp.Command{
					Title:     textToDisplay,
					Command:   "",
					Arguments: nil,
				},
				Data: nil,
			}
			codeLenses = append(codeLenses, codeLens)
		}
	}

	return codeLenses, nil
}

func (s *Service) HandleCodeLensResolve(params lsp.CodeLens, conn lspserv.Connection) (*lsp.CodeLens, error) {
	return nil, nil
}

func (s *Service) HandleDidChangeWatchedFiles(params lsp.DidChangeWatchedFilesParams, conn lspserv.Connection) error {
	return nil
}

func (s *Service) HandleDidOpen(params lsp.DidOpenTextDocumentParams, conn lspserv.Connection) error {
	return s.CompileAndReportErrors(params.TextDocument.URI, uint(params.TextDocument.Version), conn)
}

func (s *Service) getDocumentHelper(documentIdentifier lsp.VersionedTextDocumentIdentifier) (*InMemoryDocument, error) {
	localPath, pathErr := toDocumentURI(documentIdentifier.URI).ToLocalFilePath()
	if pathErr != nil {
		return nil, pathErr
	}

	return s.documents.GetDocument(LocalFileSystemPath(localPath), DocumentVersion(documentIdentifier.Version))
}

type DiagnosticsForDocument struct {
	diagnostics []lsp.Diagnostic
}

func (d *DiagnosticsForDocument) Add(lspDiagnostic lsp.Diagnostic) {
	d.diagnostics = append(d.diagnostics, lspDiagnostic)
}

func (d *DiagnosticsForDocument) All() []lsp.Diagnostic {
	return d.diagnostics
}

func (d *DiagnosticsForDocument) IsEmpty() bool {
	return len(d.diagnostics) == 0
}

func (d *DiagnosticsForDocument) Clear() {
	d.diagnostics = []lsp.Diagnostic{}
}

type DiagnosticsForDocuments struct {
	allDiagnostics map[string]*DiagnosticsForDocument
}

func NewDiagnosticsForDocuments() *DiagnosticsForDocuments {
	return &DiagnosticsForDocuments{make(map[string]*DiagnosticsForDocument)}
}

func (d *DiagnosticsForDocuments) All() map[string]*DiagnosticsForDocument {
	return d.allDiagnostics
}

func (d *DiagnosticsForDocuments) Clear() {
	for _, document := range d.allDiagnostics {
		document.Clear()
	}
}

func (d *DiagnosticsForDocuments) Tidy() {
	var toBeDeleted []string
	for key, document := range d.allDiagnostics {
		if document.IsEmpty() {
			toBeDeleted = append(toBeDeleted, key)
		}
	}

	for _, keyToDelete := range toBeDeleted {
		delete(d.allDiagnostics, keyToDelete)
	}
}

func (d *DiagnosticsForDocuments) Add(localPath LocalFileSystemPath, lspDiagnostic lsp.Diagnostic) {
	foundLocalPath := string(localPath)
	existingDiagDocument := d.allDiagnostics[foundLocalPath]
	if existingDiagDocument == nil {
		existingDiagDocument = &DiagnosticsForDocument{}
		d.allDiagnostics[foundLocalPath] = existingDiagDocument
	}

	existingDiagDocument.Add(lspDiagnostic)
}

func (s *Service) CompileAndReportErrors(uri lsp.DocumentURI, version uint, conn lspserv.Connection) error {
	localPath, localPathErr := toDocumentURI(uri).ToLocalFilePath()
	if localPathErr != nil {
		return localPathErr
	}

	s.diagnostics.Clear()
	_, compileErr := s.compiler.Compile(localPath)
	allDiagnostics := s.diagnostics
	if compileErr != nil {
		log.Printf("received err %T", compileErr)
		moduleErr, wasModuleErr := compileErr.(*decorated.ModuleError)
		if wasModuleErr {
			compileErr = moduleErr.WrappedError()
		}
		multiErr, wasMultiErr := compileErr.(*decorated.MultiErrors)
		if !wasMultiErr {
			return nil
		}
		for _, foundErr := range multiErr.Errors() {
			sourcePosition := foundErr.FetchPositionLength()
			if sourcePosition.Document == nil {
				continue
			}
			foundLocalPath, errLocalPath := sourcePosition.Document.Uri.ToLocalFilePath()
			if errLocalPath != nil {
				return errLocalPath
			}

			lspDiagnostic := lsp.Diagnostic{
				Range:           *tokenToLspRange(sourcePosition.Range),
				Severity:        lsp.Error,
				Code:            fmt.Sprintf("%T", foundErr),
				CodeDescription: nil,
				Source:          "swamp",
				Message:         foundErr.Error(),
			}
			allDiagnostics.Add(LocalFileSystemPath(foundLocalPath), lspDiagnostic)
		}
	}

	for _, workspaceModule := range s.workspacer.AllModules() {
		for _, warning := range workspaceModule.Warnings() {
			sourcePosition := warning.FetchPositionLength()
			if sourcePosition.Document == nil {
				continue
			}
			foundLocalPath, errLocalPath := sourcePosition.Document.Uri.ToLocalFilePath()
			if errLocalPath != nil {
				return errLocalPath
			}

			lspDiagnostic := lsp.Diagnostic{
				Range:           *tokenToLspRange(sourcePosition.Range),
				Severity:        lsp.Warning,
				Code:            fmt.Sprintf("%T", warning),
				CodeDescription: nil,
				Source:          "swamp",
				Message:         warning.Warning(),
			}
			allDiagnostics.Add(LocalFileSystemPath(foundLocalPath), lspDiagnostic)
		}
	}

	for uri, diagDocument := range allDiagnostics.All() {
		params := lsp.PublishDiagnosticsParams{
			URI:         lsp.DocumentURI(uri),
			Version:     0,
			Diagnostics: diagDocument.All(),
		}
		conn.PublishDiagnostics(params)
	}

	s.diagnostics.Tidy()

	return nil
}

func (s *Service) HandleDidChange(params lsp.DidChangeTextDocumentParams, conn lspserv.Connection) error {
	foundDocument, err := s.getDocumentHelper(params.TextDocument)
	if err != nil {
		return err
	}

	for _, contentChange := range params.ContentChanges {
		editRange := lspToTokenRange(contentChange.Range)
		if changeErr := foundDocument.MakeChange(editRange, contentChange.Text); changeErr != nil {
			return changeErr
		}
	}
	foundDocument.UpdateVersion(DocumentVersion(params.TextDocument.Version))

	return s.CompileAndReportErrors(params.TextDocument.URI, uint(params.TextDocument.Version), conn)
}

func (s *Service) HandleDidClose(params lsp.DidCloseTextDocumentParams, conn lspserv.Connection) error {
	return nil
}

func (s *Service) HandleWillSave(params lsp.WillSaveTextDocumentParams, conn lspserv.Connection) error {
	return nil
}

func (s *Service) HandleDidSave(params lsp.DidSaveTextDocumentParams, conn lspserv.Connection) error {
	return nil
}
