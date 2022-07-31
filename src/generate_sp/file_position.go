/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/opcodes/opcode_sp"
)

func toFilePosition(cache *assembler_sp.FileUrlCache, source token.SourceFileReference) opcode_sp.FilePosition {
	var fileID assembler_sp.FileUrlID
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
