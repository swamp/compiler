package generate_sp

import (
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/opcodes/opcode_sp"
)

func toFilePosition(cache *FileUrlCache, source token.SourceFileReference) opcode_sp.FilePosition {
	var fileID FileUrlID
	if source.Document == nil {
		fileID = 0xffff
	} else {
		fileID = cache.GetID(string(source.Document.Uri))
	}
	return opcode_sp.FilePosition{
		SourceFileID: uint(fileID),
		Start: opcode_sp.LineCol{
			Line:   uint(source.Range.Start().Line()),
			Column: uint(source.Range.Start().Column()),
		},
		End: opcode_sp.LineCol{
			Line:   uint(source.Range.End().Line()),
			Column: uint(source.Range.End().Column()),
		},
	}
}
