# Sikh

A small **s**td**i**n **k**ey-**h**ooking library. It reads the raw input from a keypress, compares it against a map, and returns a string coresponding with the key pressed.

```go
type Sikh struct { // size=4, class=8
    isRunning atomic.Bool
}

func (sikh *Sikh) Halt()
func (sikh *Sikh) ReadBytes()
func (sikh *Sikh) Start(handler func(string))
func (sikh *Sikh) getRawKeystroke() ([4]byte, error)
```

Usage

```go
package main

import (
    "fmt"
    "github.com/kyleraywed/sikh"
)

func main() {
    // create the hook
	var sikh sikh.Sikh

    // switch over the string and implement logic
	sikh.Start(func(s string) {
		switch s {
		case "[Ctrl+c]":
			sikh.Halt() // do not forget this
		case "[Esc]":
			fmt.Println("Trying to escape? Try ctrl+c")
		default:
			fmt.Println(s)
		}
	})

	// read raw bytes from keypresses until Esc is pressed
	sikh.ReadBytes()
}
```

Notes and design

- Start blocks until Halt() is called; if you don't include some logic to call Halt(), you will have to kill the process yourself.