/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package coloring

import (
	"fmt"
	"strings"
	"testing"
)

func TestSomething(t *testing.T) {
	code := `
type alias GameplayState =
{ playerX: Int
, time : Int
}


defaultCamera : Bool -> Camera
defaultCamera ignore =
    normalCamera (FrameBuffer.default True) 2


zoomRenderToTextureSprite : Int -> Camera -> Sprite
zoomRenderToTextureSprite time camera =
    let
        percentage =
            Math.cos (time*8) / 7

        scaleFactor =
            percentage + 300
    in
    spriteFromRenderToTexture camera |> scaleSprite scaleFactor
`
	output, err := SyntaxColor(strings.TrimSpace(code))

	if err != nil {
		t.Fatal(err)
	}

	fmt.Print(output)
}
