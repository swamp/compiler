/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"encoding/binary"

	"github.com/swamp/assembler/lib/assembler_sp"
	swamppack "github.com/swamp/pack/lib"
)

func PackLedger(constants []*assembler_sp.Constant) ([]byte, error) {
	var octets []byte

	for _, constant := range constants {
		var entryOctets [8]byte
		binary.LittleEndian.PutUint32(entryOctets[0:4], uint32(constant.ConstantType()))
		binary.LittleEndian.PutUint32(entryOctets[4:8], uint32(constant.PosRange().Position))
		octets = append(octets, entryOctets[:]...)
	}

	var entryOctets [8]byte

	binary.LittleEndian.PutUint32(entryOctets[0:4], uint32(0))
	binary.LittleEndian.PutUint32(entryOctets[4:8], uint32(0))
	octets = append(octets, entryOctets[:]...)

	return octets, nil
}

func Pack(constants []*assembler_sp.Constant, dynamicMemory []byte, typeInfoPayload []byte) ([]byte, error) {
	ledgerOctets, ledgerErr := PackLedger(constants)
	if ledgerErr != nil {
		return nil, ledgerErr
	}

	return swamppack.Pack(ledgerOctets, dynamicMemory, typeInfoPayload)
}
