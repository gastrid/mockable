This mockable will turn this:

```
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
```
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
