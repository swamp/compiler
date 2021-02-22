/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

// SingleLineCommentToken :
type SingleLineCommentToken struct {
	CommentToken
}

func NewSingleLineCommentToken(raw string, text string, forDocumentation bool, sourceFileReference SourceFileReference) SingleLineCommentToken {
	return SingleLineCommentToken{CommentToken: CommentToken{RawString: raw, CommentString: text, SourceFileReference: sourceFileReference, ForDocumentation: forDocumentation}}
}

func (s SingleLineCommentToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}

func (s SingleLineCommentToken) String() string {
	return s.CommentToken.Raw()
}
