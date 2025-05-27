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

package testing

import (
	"bytes"
	"fmt"
	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

// cmdTestCase describes a test case that works with releases.

func TInfo(msg interface{}) {
	fmt.Println("--- INFO:", msg)
}
func TOut(msg interface{}) {
	fmt.Println("--- OUTPUT:", msg)
}

func ErrorAssertionWithEqual(t *testing.T, wantErr, err error) bool {
	var fu assert.ErrorAssertionFunc = func(t assert.TestingT, err error, i ...interface{}) bool {
		if err == nil {
			return true
		}
		return assert.Equal(t, wantErr, err)
	}
	return fu(t, err)
}

type TestCmdFunc = func(buffer *bytes.Buffer) *cobra.Command

// A helper to ignore os.Exit(1) errors when running a cobra Command
func ExecuteCommand(cmdFunc TestCmdFunc, cmd string) (stdout string, err error) {
	args, err := shellwords.Parse(cmd)
	if err != nil {
		return "", err
	}
	return ExecuteCommandOfArgs(cmdFunc, args...)
}

func ExecuteCommandOfArgs(cmdFunc TestCmdFunc, args ...string) (stdout string, err error) {
	buf := new(bytes.Buffer)
	//nolint:gocritic
	//root, err := NewRootCmd(buf)
	cmd := cmdFunc(buf)
	for i, arg := range args {
		TInfo(fmt.Sprintf("(%d) %s", i, arg))
	}
	// programs can exit with error here..
	stdout, _, err = ExecuteCommandC(cmd, args...)
	TInfo("The command successfully returned values for assertion.")
	return stdout, err

}

// A helper to ignore os.Exit(1) errors when running a cobra Command
func ExecuteCommandC(root *cobra.Command, args ...string) (stdout string, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	stdout = Run(func() {
		_, err = root.ExecuteC()
	})
	return stdout, buf.String(), err
}

func Run(exec func()) (stdout string) {
	// see https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	exec()
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			panic(err)
		}
		outC <- buf.String()
	}()

	// back to normal state
	if err := w.Close(); err != nil {
		TInfo(err)
	}
	os.Stdout = old // restoring the real stdout
	stdout = <-outC
	return stdout
}

func RunForStdout(streams genericiooptions.IOStreams, exec func()) (stdout string) {
	// see https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := streams.Out // keep backup of the real stdout
	exec()
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, streams.In); err != nil {
			panic(err)
		}
		outC <- buf.String()
	}()

	// back to normal state
	if err := streams.Out.(*os.File).Close(); err != nil {
		TInfo(err)
	}
	streams.Out = old // restoring the real stdout
	stdout = <-outC
	return stdout
}
