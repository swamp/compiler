/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

type LocalFileSystemPath string

type LocalFileSystemRoot string

type DocumentProvider interface {
	ReadDocument(uri LocalFileSystemPath) (string, error)
}
