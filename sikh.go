package sikh

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/kyleraywed/sikh/keymaps"
	"golang.org/x/term"
)

type Sikh struct {
	isRunning atomic.Bool
}

// Start reading keypresses. Define logic via a handler function.
func (sikh *Sikh) Start(handler func(string)) {
	if sikh.isRunning.Load() {
		return
	}

	sikh.isRunning.Store(true)

	for {
		if !sikh.isRunning.Load() {
			return
		}

		byteRep, err := sikh.getRawKeystroke()
		if err != nil {
			log.Println(err)
		}

		if strRep, ok := keymaps.StandardMap[byteRep]; ok {
			handler(strRep)
		}
	}
}

func (sikh *Sikh) Halt() {
	sikh.isRunning.Store(false)
}

// Present the user with a prompt to easier map out keys, ESC to quit.
func (sikh *Sikh) ReadBytes() {
	h := func(br [4]byte) {
		switch br {
		case [4]byte{27, 0, 0, 0}: // esc
			sikh.Halt()
		case [4]byte{3, 0, 0, 0}: // ctrl+c
			fmt.Println(br, "Press [Esc] to exit.")
		default:
			fmt.Println(br)
		}
	}

	if sikh.isRunning.Load() {
		return
	}

	sikh.isRunning.Store(true)
	fmt.Println("Press [Esc] [27, 0, 0, 0] to quit")

	for {
		if !sikh.isRunning.Load() {
			return
		}

		byteRep, err := sikh.getRawKeystroke()
		if err != nil {
			log.Println(err)
		}

		h(byteRep)
	}
}

func (sikh *Sikh) getRawKeystroke() ([4]byte, error) {
	var rep [4]byte

	// initially, i had Start() put the term in raw and Halt() take it out, but
	// that results in some pretty strange formatting issues.

	// the solution was found in the great luck that computers are really fast,
	// such that going in and out of raw mode for every keypress creates no noticable overhead.

	// "It works on my machine"

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return rep, err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 4)
	_, err = os.Stdin.Read(b)

	if err != nil {
		return rep, err
	}

	copy(rep[:], b[:4])
	return rep, nil
}
