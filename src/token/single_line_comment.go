/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

// SingleLineCommentToken :
type SingleLineCommentToken struct {
	CommentToken
}

func NewSingleLineCommentToken(raw string, text string, forDocumentation bool, position PositionLength) SingleLineCommentToken {
	return SingleLineCommentToken{CommentToken: CommentToken{RawString: raw, CommentString: text, PositionLength: position, ForDocumentation: forDocumentation}}
}
