/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

// Package swamppack_sp packs constants into a .swamp-pack.
package swamppack_sp

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"

	raff "github.com/piot/raff-go/src"
)

type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

func writeChunkHeader(writer io.Writer, icon raff.FourOctets, name raff.FourOctets, payload []byte) error {
	if err := raff.WriteChunk(writer, icon, name, payload); err != nil {
		return err
	}

	return nil
}

func writePackHeader(writer io.Writer) error {
	name := raff.FourOctets{'s', 'p', 'k', '5'}
	packetIcon := raff.FourOctets{0xF0, 0x9F, 0x93, 0xA6}
	return writeChunkHeader(writer, packetIcon, name, nil)
}

func writeVersion(writer io.Writer, version Version) error {
	header := []byte{version.Major, version.Minor, version.Patch}

	if _, writeErr := writer.Write(header); writeErr != nil {
		return writeErr
	}

	return nil
}

func writeTypeInfo(writer io.Writer, payload []byte) error {
	name := raff.FourOctets{'s', 't', 'i', '0'}
	packetIcon := raff.FourOctets{0xF0, 0x9F, 0x93, 0x9C}
	return writeChunkHeader(writer, packetIcon, name, payload)
}

func writeDynamicMemory(writer io.Writer, payload []byte) error {
	name := raff.FourOctets{'d', 'm', 'e', '0'}
	packetIcon := raff.FourOctets{0xF0, 0x9F, 0x92, 0xBB}
	log.Printf("dynamicMemory: %s", hex.Dump(payload))
	return writeChunkHeader(writer, packetIcon, name, payload)
}

func writeLedger(writer io.Writer, payload []byte) error {
	name := raff.FourOctets{'l', 'd', 'g', '0'}
	packetIcon := raff.FourOctets{0xF0, 0x9F, 0x97, 0x92}
	return writeChunkHeader(writer, packetIcon, name, payload)
}

func Pack(ledger []byte, dynamicMemory []byte, typeInfo []byte) ([]byte, error) {
	var buf bytes.Buffer

	if err := raff.WriteHeader(&buf); err != nil {
		return nil, err
	}

	if writeErr := writePackHeader(&buf); writeErr != nil {
		return nil, writeErr
	}

	if writeErr := writeTypeInfo(&buf, typeInfo); writeErr != nil {
		return nil, writeErr
	}

	if writeErr := writeDynamicMemory(&buf, dynamicMemory); writeErr != nil {
		return nil, writeErr
	}

	if writeErr := writeLedger(&buf, ledger); writeErr != nil {
		return nil, writeErr
	}

	return buf.Bytes(), nil
}
