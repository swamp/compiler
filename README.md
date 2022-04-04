# Swamp Compiler

[![Go Reference](https://pkg.go.dev/badge/github.com/swamp/compiler.svg)](https://pkg.go.dev/github.com/swamp/compiler)
[![Go Report Card](https://goreportcard.com/badge/github.com/swamp/compiler)](https://goreportcard.com/report/github.com/swamp/compiler)
[![Actions Status](https://github.com/swamp/compiler/workflows/Go/badge.svg)](https://github.com/swamp/compiler/actions)

Compiles `.swamp` files and produces `.swamp-pack` binaries.

Swamp is designed to be embedded into other applications, mainly for game engines.


## Syntax

### Functions

You must specify an annotation for all functions. All types separated by '->' (right arrow). The last type is normally the return type (see Currying).

```haskell

doubleInt : Int -> Int
doubleInt a =
    a * 2

```

### Primitive types

* `Int` (always signed 32-bit)
* `Fixed` (fixedpoint, always signed 32-bit)
* `String`
* `Bool` (`True` or `False`)
* `Blob` (binary data)
* `List`. Literal `[]`
* `Array` Literal `[| |]`

### Types

* CustomType (similar to enums or unions in other languages, but with associated data).

```
type Custom =
    Idle
    | Running Int
    | Sleeping
```

* Records (structs in other languages).

```haskell
{ name : String, x : Int, y : Int }
```

* TupleType (ordered types, minimum of two entries)

```haskell
( String, Int, Fixed )
```

### Collection Types

* List. Fast to add to.

```haskell
List Int
```

* Array. Fast to access using index.

```haskell
Array String
```

### Expressions

Most things in Swamp is an expression.

#### if

```haskell

sample : Int -> String
sample a =
    if a > 10 then "high" else "low"

```

#### let

```haskell

sample : Int -> Int -> Int
sample a b =
    let
       x = a + 10

       y = b - 1
    in
    x + y

```

Let features:

##### destructuring

Extract fields from a Record or Tuple during let expression.

```haskell

type alias Position = { x : Int, y : Int }

addXY : Position -> Int
addXY pos =
    let
        { x, y } = pos
    in
    x + y


addXYTuple : (Int, Int) -> (Int, Int)
addXYTuple pos =
    let
        x, y = pos
    in
    x + y    

```


#### case


```haskell

type CharacterState =
    Moving Int
    | Jumping
    | Idle


speed : CharacterState -> Int
speed state =
    case state of
       
        Moving speed -> speed

        Jumping -> 10

        _ -> 0
```

#### function call

Add terms after a function value to call it and return the result. If fewer arguments are passed, it creates a function value that saves the specified arguments for future calls (Currying).

```haskell

double : Int -> Int
double a =
    a * 2

main : Int -> Int
main _ =
    double 42
```

#### record lookup

Fetches the field from a record value. Use a `.` on a record value and then specify the field name.

```haskell

lookupX : { x : Int, y : Int } -> Int
lookupX position =
    position.x

```


#### guard

A list of if statements to be evaluated. Evaluated from top to bottom, uses `_` if no match is found.

```haskell
temperature : Int -> String
temperature x =
    | x > 15 -> "Warm"
    | x < -10 > -> "Cold"
    | _ -> "Neither hot nor cold"
```


#### construction

Easier way to fill in the values for a record.

```haskell

type alias Position = { x : Int, y : Int }


moveLeft : Position -> Position
moveLeft pos =
    Position (pos.x - 10) pos.y

```

#### cast

Explicitly cast from one type to another. Usually this conversion happen automatically in Swamp (if the types are strictly the same).

```haskell

type alias Velocity = Int

let
    velocity = 42 : Velocity
in
velocity

```

#### pipe

Send a result to the left or to the right.


```haskell

double : Int -> Int
double a =
    a * 2


sample : Int -> Int
sample a =
    double <| AnotherFile.Abs -10

```



### Statements

Statements must be top level in a swamp file and can not be part of expressions.

#### type alias

Use a type under another name.

```haskell

type alias Position = { x : Int, y : Int }


moveRight : Position -> Position
moveRight pos =
    { x = pos.x + 10, y = pos.y }

```

#### custom type

Define own type with variants that can have parameters.

```haskell

type alias Velocity = Int

type Custom =
    Idle
    | Running Velocity
    | Sleeping

```

#### import

Import from another `.swamp` file:

```haskell

import AnotherFile

sample : Int -> Int
sample a =
    AnotherFile.Abs -10

```


### Literals

#### List

```haskell
[ 2, 42, 99 ]
```


#### Array


```haskell
[| 2, 42, 99 |]

```

#### Bool

```haskell
True
False
```

#### Int

```haskell
12
```

#### Fixed

```haskell
42.9
```

#### String

```haskell
"Hello, world"
```

#### Type ID

References a type.

```haskell
$Position
```

#### Resource Name

```haskell
@directory/name.png
```
