// Copyright 2024 The huhouhua Authors
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

package util

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// cmdTestCase describes a test case that works with releases.
type CmdTestCase struct {
	Name      string
	Cmd       string
	WantError bool
}

func TInfo(msg interface{}) {
	fmt.Println("--- INFO:", msg)
}
func TOut(msg interface{}) {
	fmt.Println("--- OUTPUT:", msg)
}

type TestCmdFunc = func(buffer *bytes.Buffer) (*cobra.Command, error)

// A helper to ignore os.Exit(1) errors when running a cobra Command
func ExecuteCommand(cmdFunc TestCmdFunc, cmd string) (stdout string, err error) {
	testCh := make(chan func() (string, error), 1)
	defer func() {
		close(testCh)
	}()
	//go func(ch chan func() (string, error)) {
	//	addSignal(ch)
	//}(testCh)
	addSignal(testCh, cmdFunc, cmd, func(f TestCmdFunc, c string, ch chan func() (string, error)) {
		args, err := shellwords.Parse(c)
		if err != nil {
			ch <- func() (string, error) {
				return "", err
			}
			return
		}
		std, err := ExecuteCommandOfArgs(f, args...)
		ch <- func() (string, error) {
			return std, err
		}
	})
	resultFunc := <-testCh
	return resultFunc()
}

func ExecuteCommandOfArgs(cmdFunc TestCmdFunc, args ...string) (stdout string, err error) {
	buf := new(bytes.Buffer)
	//root, err := NewRootCmd(buf)
	cmd, err := cmdFunc(buf)
	if err != nil {
		return "", err
	}
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

	// see https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	_, err = root.ExecuteC()

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

	return stdout, buf.String(), err
}

func addSignal(testCh chan func() (string, error), cmdFunc TestCmdFunc, cmd string,
	test func(cmdFunc TestCmdFunc, cmd string, testCh chan func() (string, error))) {
	// 注册信号处理器
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		close(sigCh)
	}()
	go test(cmdFunc, cmd, testCh)
	for {
		select {
		case sig := <-sigCh:
			logrus.Infof("Received signal: %s\n", sig.String())
			testCh <- func() (string, error) {
				return "", errors.New(fmt.Sprintf("exit:%s", sig.String()))
			}
			//case <-time.After(time.Minute * 3): //等待ret返回结果，超时为3分钟
			//	testCh <- func() (string, error) {
			//		return "", errors.New("time out!")
			//	}
		}
	}
}
