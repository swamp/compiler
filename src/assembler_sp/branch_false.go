/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type BranchFalse struct {
	condition SourceStackPos
	jump      *Label
}

func (o *BranchFalse) String() string {
	return fmt.Sprintf("[brfa %v jump:%v]", o.condition, o.jump)
}

func (o *BranchFalse) Condition() SourceStackPos {
	return o.condition
}

func (o *BranchFalse) Jump() *Label {
	return o.jump
}

type BranchTrue struct {
	condition SourceStackPos
	jump      *Label
}

func (o *BranchTrue) String() string {
	return fmt.Sprintf("[breq %v jump:%v]", o.condition, o.jump)
}

func (o *BranchTrue) Condition() SourceStackPos {
	return o.condition
}

func (o *BranchTrue) Jump() *Label {
	return o.jump
}
