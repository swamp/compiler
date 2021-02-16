/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package executetest

import (
	"testing"
)

func TestSimple(t *testing.T) {
	executeTest(t,
		`
main : Bool -> Int
main ignored =
    42 + 8
`, "int: 50 refcount:1")
}

func TestAdvanced(t *testing.T) {
	executeTest(t, `
type alias Sprite =
    { x : Int
    , y : Int
    }


type alias World =
    { drawTasks : List (List Sprite)
    }


drawSprite : Sprite -> Bool
drawSprite sprite =
    True


drawSprites : List Sprite -> List Bool
drawSprites sprites =
    List.map drawSprite sprites


drawWorld : World -> List (List Bool)
drawWorld world =
    List.map drawSprites world.drawTasks


main : Bool -> List (List Bool)
main ignore =
    drawWorld { drawTasks = [ [ { x = 10, y = 20 }, { x = 44, y = 98 } ], [ { x = 99, y = 98 } ] ] }
`, `
list: refcount:1
..list: refcount:2
....bool: True
..list: refcount:2
....bool: True
....bool: True
`)
}

func TestTrickyEmptyList(t *testing.T) {
	executeTest(t,
		`
main : Bool -> Int
main ignored =
    let
        x = List.head []
    in
    case x of
        Just a -> a

        _ -> -1
`,
		`int: -1 refcount:3`)
}

func TestSimple2(t *testing.T) {
	executeTest(t,
		`
third : Int -> Int
third a =
    a - 2


another : Int -> Int
another a =
    third 18 + a


main : Bool -> Int
main ignored =
    another 42 + 8
`, "int: 66 refcount:1")
}

func TestCustomType(t *testing.T) {
	executeTest(t,
		`
type SomeEnum =
    First String
    | Anon
    | Second Int


a : Bool -> SomeEnum
a dummy =
    First "Hello"


main : Bool -> SomeEnum
main x =
    a x
`, `
enum: 0 refcount:2
..string: 'Hello' refcount:3
`)
}

func TestCustomTypeMultipleParameters(t *testing.T) {
	executeTest(t,
		`
type SomeEnum =
    Anon
    | First String Int
    | Second Int


a : Bool -> SomeEnum
a dummy =
    First "Hello" 42


display : SomeEnum -> String
display v =
    case v of
        Anon -> "Anonymous"

        First str num -> str

        _ -> "Dont know"


main : Bool -> String
main x =
    let
        e = a x
    in
    display e
`, `
string: 'Hello' refcount:4
`)
}

func TestPosition(t *testing.T) {
	executeTest(t,
		`
type alias Position =
    { x : Int
    , y : Int
    }


type alias Sprite =
    { pos : Position
    }


move : Position -> Position -> Position
move pos delta =
    let
        newX = pos.x + delta.x

        newY = pos.y - delta.y
    in
    { x = newX, y = newY }


main : Bool -> Position
main a =
    move { x = 10, y = 20 } { x = 1, y = 2 }
`, `
struct: field_count: 2 refcount:2
..int: 11 refcount:1
..int: 18 refcount:1
`)
}

func TestOwnAppender(t *testing.T) {
	executeTest(t,
		`
ownAppender : List Int -> List Int -> List Int
ownAppender lista listb =
    42 :: lista ++ listb ++ [ 9 ]


main : Bool -> List Int
main ignored =
    ownAppender [ 1, 2 ] (11 :: [ 3, 4 ])
`, `
list: refcount:2
..int: 42 refcount:5
..int: 1 refcount:5
..int: 2 refcount:5
..int: 11 refcount:4
..int: 3 refcount:4
..int: 4 refcount:4
..int: 9 refcount:3
`)
}

func TestCase(t *testing.T) {
	executeTest(t,
		`
type alias Tinkering =
    { solder : Bool
    }


type alias Studying =
    { subject : String
    }


type alias Work =
    { web : Int
    }


type alias Amplitude =
    { something : String
    }


type Child =
    Aron Tinkering Amplitude
    | Alexandra
    | Alma Studying
    | Isabelle Work


some : Child -> String
some child =
    case child of
        Aron x amp ->
            if x.solder then
                "Aron"
            else
                amp.something

        Alexandra -> "Alexandris"

        _ -> "Unknown"


main : Bool -> String
main ignore =
    some ( Aron { solder = False } { something = "return it" } )
`, `
    string: 'return it' refcount:4
  `)
}

// -- just a comment
func TestListOfRecordWithLists(t *testing.T) {
	executeTest(t,
		`
type alias Enemy =
    { values : List Int
    }


type alias World =
    { enemies : List Enemy
    }


updateEnemy : List Enemy -> Bool
updateEnemy enemies =
    True


updateWorld : World -> Bool
updateWorld w =
    updateEnemy w.enemies


main : Bool -> Bool
main ignore =
    updateWorld { enemies = [ { values = [ 1, 3 ] } ] }

    `, `
    bool: True
  `)
}

func TestSimpleListOfLists(t *testing.T) {
	executeTest(t,
		`
type alias Sprite =
    { x : Int
    , y : Int
    }


drawSprite : Sprite -> Bool
drawSprite sprite =
    if sprite.x > 10 then
        True
    else
        False


main : Bool -> List Bool
main ignore =
    List.map drawSprite [ { x = 10, y = 20 }, { x = 44, y = 98 } ]
    `, `
list: refcount:1
..bool: True
..bool: False
  `)
}

func TestSimpleListConCatMap(t *testing.T) {
	executeTest(t,
		`
makeItLower : Int -> Int
makeItLower v =
    v - 10


makeThemLower : List Int -> List Int
makeThemLower lst =
    List.map makeItLower lst


main : Bool -> List Int
main ignore =
    List.concatMap makeThemLower [ [ 191, 23222, 310, 8000 ] ]
    `, `
list: refcount:1
..int: 181 refcount:3
..int: 23212 refcount:3
..int: 300 refcount:3
..int: 7990 refcount:3
  `)
}

func TestListOfLists(t *testing.T) {
	executeTest(t,
		`
type alias Sprite =
    { x : Int
    , y : Int
    }


type alias World =
    { drawTasks : List (List Sprite)
    }


drawSprite : Sprite -> Bool
drawSprite sprite =
    True


drawSprites : List Sprite -> List Bool
drawSprites sprites =
    List.map drawSprite sprites


drawWorld : World -> List (List Bool)
drawWorld world =
    List.map drawSprites world.drawTasks


main : Bool -> List (List Bool)
main ignore =
    drawWorld { drawTasks = [ [ { x = 10, y = 20 }, { x = 44, y = 98 } ], [ { x = 99, y = 98 } ] ] }
    `, `
list: refcount:1
..list: refcount:2
....bool: True
..list: refcount:2
....bool: True
....bool: True
  `)
}

func TestSimpleListAny(t *testing.T) {
	executeTest(t,
		`
isMoreThanTen : Int -> Bool
isMoreThanTen v =
    v > 10


main : Bool -> Bool
main ignore =
    List.any isMoreThanTen [ 1, 2, 34, 499 ]
    `, `
bool: True
  `)
}

func TestSimpleListAnyFalse(t *testing.T) {
	executeTest(t,
		`
isMoreThanTen : Int -> Bool
isMoreThanTen v =
    v < 1000


main : Bool -> Bool
main ignore =
    List.any isMoreThanTen [ 9191, 23222, 31000, 8000 ]
    `, `
bool: False
  `)
}

func TestSimpleListFilter(t *testing.T) {
	executeTest(t,
		`
isLessThanThousand : Int -> Bool
isLessThanThousand v =
    v < 1000


main : Bool -> List Int
main ignore =
    List.filter isLessThanThousand [ 191, 23222, 310, 8000 ]
    `, `
list: refcount:1
..int: 310 refcount:4
..int: 191 refcount:4
  `)
}

func TestSimpleListRemove(t *testing.T) {
	executeTest(t,
		`
isLessThanThousand : Int -> Bool
isLessThanThousand v =
    v < 1000


main : Bool -> List Int
main ignore =
    List.remove isLessThanThousand [ 191, 23222, 310, 8000 ]
    `, `
list: refcount:1
..int: 8000 refcount:4
..int: 23222 refcount:4
  `)
}

func TestSimpleListHead(t *testing.T) {
	executeTest(t,
		`
main : Bool -> Maybe Int
main ignore =
    List.head [ 191, 23222, 310, 8000 ]
    `, `
enum: 1 refcount:1
..int: 191 refcount:4
  `)
}

func TestMaybe(t *testing.T) {
	executeTest(t,
		`type alias Thing =
    { something : String
    }


main : Bool -> Maybe Thing
main ignore =
    Just { something = "hello" }
`, `
enum: 1 refcount:2
..struct: field_count: 1 refcount:2
....string: 'hello' refcount:3
`)
}

func TestSimpleListHeadNothing(t *testing.T) {
	executeTest(t,
		`
main : Bool -> Maybe Int
main ignore =
    List.head []
    `, `
enum: 0 refcount:1
  `)
}

func TestSimpleListConcat(t *testing.T) {
	executeTest(t,
		`
main : Bool -> List Int
main ignore =
    List.concat [ [ 10 ], [ 191, 23222, 310, 8000 ] ]
    `, `
list: refcount:1
..int: 8000 refcount:4
..int: 310 refcount:4
..int: 23222 refcount:4
..int: 191 refcount:4
..int: 10 refcount:4
  `)
}

func TestSimpleListEmpty(t *testing.T) {
	executeTest(t,
		`
type alias Answer =
    { first : Bool
    , second : Bool
    }


main : Bool -> Answer
main ignore =
    let
        first = List.isEmpty []

        second = List.isEmpty [ 2 ]
    in
    { first = first, second = second }
`, `
struct: field_count: 2 refcount:2
..bool: True
..bool: False
  `)
}

func TestSimpleListRange(t *testing.T) {
	executeTest(t,
		`
type alias Answer =
    { first : List Int
    , empty : List Int
    , negative : List Int
    }


main : Bool -> Answer
main ignore =
    let
        first = List.range 2 10

        empty = List.range 2 1

        negative = List.range -2 7
    in
    { first = first, empty = empty, negative = negative }
`, `
struct: field_count: 3 refcount:2
..emptylist
..list: refcount:1
....int: 2 refcount:2
....int: 3 refcount:2
....int: 4 refcount:2
....int: 5 refcount:2
....int: 6 refcount:2
....int: 7 refcount:2
....int: 8 refcount:2
....int: 9 refcount:2
....int: 10 refcount:2
..list: refcount:1
....int: -2 refcount:2
....int: -1 refcount:2
....int: 0 refcount:2
....int: 1 refcount:2
....int: 2 refcount:2
....int: 3 refcount:2
....int: 4 refcount:2
....int: 5 refcount:2
....int: 6 refcount:2
....int: 7 refcount:2
  `)
}

func TestSimpleListLength(t *testing.T) {
	executeTest(t,
		`
type alias Answer =
    { first : Int
    , empty : Int
    }


main : Bool -> Answer
main ignore =
    let
        first = List.length [ 2, 3, 99 ]

        empty = List.length []
    in
    { first = first, empty = empty }
`, `
struct: field_count: 2 refcount:2
..int: 0 refcount:1
..int: 3 refcount:1
  `)
}

func TestArrayFromList(t *testing.T) {
	executeTest(t,
		`
type alias Scale2 =
    { scaleX : Int
    , scaleY : Int
    }


type alias Sprite =
    { dummy : Int
    , scale : Scale2
    }


spritesToDraw : Bool -> Array Sprite
spritesToDraw doIt =
    Array.fromList [ { dummy = 0, scale = { scaleX = 10, scaleY = 10 } } ]


main : Bool -> Sprite
main ignore =
    let
        mightBeSprite = Array.get 0 (spritesToDraw True)
    in
    case mightBeSprite of
        Just sprite -> sprite

        Nothing -> { dummy = -1, scale = { scaleX = -1, scaleY = -1 } }
`, `
struct: field_count: 2 refcount:4
..int: 0 refcount:4
..struct: field_count: 2 refcount:2
....int: 10 refcount:4
....int: 10 refcount:4
  `)
}
