# Swamp Compiler

[![Go Reference](https://pkg.go.dev/badge/github.com/swamp/compiler.svg)](https://pkg.go.dev/github.com/swamp/compiler)
[![Go Report Card](https://goreportcard.com/badge/github.com/swamp/compiler)](https://goreportcard.com/report/github.com/swamp/compiler)
[![Actions Status](https://github.com/swamp/compiler/workflows/Go/badge.svg)](https://github.com/swamp/compiler/actions)

Compiles `.swamp` files and produces `.swamp-pack` binaries.

Swamp is designed to be embedded into other applications, mainly for game engines.


## Syntax

### Functions

You must specify an annotation for all functions. All types separated by '->' (right arrow). The last type is the return type.

```haskell

doubleInt : Int -> Int
doubleInt a =
    a * 2

```

### Primitives

* Int (always 32-bit)
* FixedPoint (always 32-bit)
* String
* Bool (True or False)
* Blob (binary data)
* List []
* Array

### Types

* CustomType (similar to enums in other languages, but with associated data).
* Records (structs in other languages).

### Collections

* List. Fast to add to.
* Array. Fast to access using index.


### Statements

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

#### alias

Use something under another name.

```haskell

type alias Position = { x : Int, y : Int }


moveRight : Position -> Position
moveRight pos =
    { x = pos.x + 10, y = pos.y }

```

#### import

Import from another `.swamp` file:

```haskell

import AnotherFile

sample : Int -> Int
sample a =
    AnotherFile.Abs -10

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
