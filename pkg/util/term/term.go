// Copyright 2024 The Kevin Berger <huhouhuam@outlook.com> Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package term provides structures and helper functions to work with
// terminal (state, sizes).
package term

import (
	"io"
	"os"

	"github.com/huhouhua/glctl/pkg/util/interrupt"

	"github.com/moby/term"
)

// SafeFunc is a function to be invoked by TTY.
type SafeFunc func() error

// TTY helps invoke a function and preserve the state of the terminal, even if the process is
// terminated during execution. It also provides support for terminal resizing for remote command
// execution/attachment.
type TTY struct {
	// In is a reader representing stdin. It is a required field.
	In io.Reader
	// Out is a writer representing stdout. It must be set to support terminal resizing. It is an
	// optional field.
	Out io.Writer
	// Raw is true if the terminal should be set raw.
	Raw bool
	// TryDev indicates the TTY should try to open /dev/tty if the provided input
	// is not a file descriptor.
	TryDev bool

	// Parent is an optional interrupt handler provided to this function - if provided
	// it will be invoked after the terminal state is restored. If it is not provided,
	// a signal received during the TTY will result in os.Exit(0) being invoked.
	Parent *interrupt.Handler

	// sizeQueue is set after a call to MonitorSize() and is used to monitor SIGWINCH signals when the
	// user's terminal resizes.
	sizeQueue *sizeQueue
}

// Safe invokes the provided function and will attempt to ensure that when the
// function returns (or a termination signal is sent) that the terminal state
// is reset to the condition it was in prior to the function being invoked. If
// t.Raw is true the terminal will be put into raw mode prior to calling the function.
// If the input file descriptor is not a TTY and TryDev is true, the /dev/tty file
// will be opened (if available).
func (t TTY) Safe(fn SafeFunc) error {
	inFd, isTerminal := term.GetFdInfo(t.In)

	if !isTerminal && t.TryDev {
		if f, err := os.Open("/dev/tty"); err == nil {
			defer f.Close()
			inFd = f.Fd()
			isTerminal = term.IsTerminal(inFd)
		}
	}
	if !isTerminal {
		return fn()
	}

	var state *term.State
	var err error
	if t.Raw {
		state, err = term.MakeRaw(inFd)
	} else {
		state, err = term.SaveState(inFd)
	}
	if err != nil {
		return err
	}
	return interrupt.Chain(t.Parent, func() {
		if t.sizeQueue != nil {
			t.sizeQueue.stop()
		}

		_ = term.RestoreTerminal(inFd, state)
	}).Run(fn)
}
