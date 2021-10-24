package generate_sp

import (
	"encoding/binary"

	"github.com/swamp/compiler/src/assembler_sp"
	swamppacksp "github.com/swamp/compiler/src/swamp_pack_sp"
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

	return swamppacksp.Pack(ledgerOctets, dynamicMemory, typeInfoPayload)
}
