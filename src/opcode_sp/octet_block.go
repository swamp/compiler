/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package opcode_sp

import (
	"encoding/binary"
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type OctetBlock struct {
	octets []uint8
}

func NewOctetBlock(octets []uint8) *OctetBlock {
	return &OctetBlock{octets: octets}
}

func (o *OctetBlock) Octets() []uint8 {
	return o.octets
}

func (o *OctetBlock) Replace(pos opcode_sp_type.ProgramCounter, v opcode_sp_type.DeltaPC) {
	binary.LittleEndian.PutUint16(o.octets[pos.Value():pos.Value()+2], uint16(v))
}

func (o *OctetBlock) FixUpLabelInject(r *LabelInject) error {
	referencedLabel := r.ReferencedLabel()
	if !referencedLabel.IsDefined() {
		return fmt.Errorf("label %v is not defined", referencedLabel)
	}

	delta, deltaErr := r.ForwardDeltaPC()
	if deltaErr != nil {
		return deltaErr
	}

	o.Replace(r.LocatedAtPosition(), delta)

	return nil
}

func (o *OctetBlock) FixUpLabelInjects(injects []*LabelInject) error {
	for _, inject := range injects {
		fixupErr := o.FixUpLabelInject(inject)
		if fixupErr != nil {
			return fixupErr
		}
	}

	return nil
}
