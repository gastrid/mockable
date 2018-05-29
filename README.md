# Mockable #

Mockable takes a bunch of functions in a file you specify and turns them into methods to a struct that implements an interface that's also created.

## Parameters ##

Mockable takes two compulsory parameters:
```
filename: the name of the file you want to transform
interface: the name of the interface that will be created
```

A variable name and a struct name will be generated from the interface name. If your interface name is `Foo`, the following
will be generated:
```go
type FooImpl struct {}
var defaultFoo Foo = &FooImpl{}
```

## What it does ##

Mockable will turn this:

```go
// Package ast declares the types used to represent syntax trees for Go
// packages.
//
package childir

import (
	"fmt"
	"time"
)

//go:generate mockable -interface=InterTest -filename=mocktest.go

// this is a comment
func SomethingWithTime() *time.Time {
	tt := time.Now()
	fmt.Println(tt)
	return &tt
}

// SomethingWithStrings does something with strings
func SomethingWithStrings(word string) string {
	fmt.Println(word)
	return word
}

func SomethingWithFloats(numb int64) int64 {
	fmt.Println(numb)
	return numb
}
```

into this:
```go
// Package ast declares the types used to represent syntax trees for Go
// packages.
//
package childir

import (
	"fmt"
	"time"
)

//go:generate genmock -interface=InterTest -mock-package=. -package=.

// this is a comment
func (mck *InterTestImpl) SomethingWithTime() *time.Time {
	tt := time.Now()
	fmt.Println(tt)
	return &tt
}

// SomethingWithStrings does something with strings
func (mck *InterTestImpl) SomethingWithStrings(word string) string {
	fmt.Println(word)
	return word
}

func (mck *InterTestImpl) SomethingWithFloats(numb int64) int64 {
	fmt.Println(numb)
	return numb
}

type InterTest interface {
	SomethingWithTime() *time.Time
	SomethingWithStrings(word string) string
	SomethingWithFloats(numb int64) int64
}
type InterTestImpl struct {
}

var defaultInterTest InterTest = &InterTestImpl{}

func SetInterTest(newS InterTest) InterTest {
	old := defaultInterTest
	defaultInterTest = newS
	return old
}
func SomethingWithTime() *time.Time { return defaultInterTest.SomethingWithTime() }

func SomethingWithStrings(word string) string { return defaultInterTest.SomethingWithStrings(word) }

func SomethingWithFloats(numb int64) int64 { return defaultInterTest.SomethingWithFloats(numb) }

```
