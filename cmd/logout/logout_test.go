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

package logout

import (
	"fmt"
	"os"
	"testing"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"

	"github.com/stretchr/testify/assert"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
)

func TestLogout(t *testing.T) {
	tests := []struct {
		name           string
		optionsFunc    func(opt *Options) error
		completeFunc   func(streams genericiooptions.IOStreams) error
		expectedOutput string
	}{
		{
			name: "logout Succeeded",
			optionsFunc: func(opt *Options) error {
				temp, err := os.CreateTemp("", ".glctl_logout_succeeded")
				if err != nil {
					return err
				}
				opt.path = temp.Name()
				return nil
			},
			completeFunc: func(streams genericiooptions.IOStreams) error {
				return nil
			},
			expectedOutput: "logout Succeeded",
		},
		{
			name: "logout Fail",
			optionsFunc: func(opt *Options) error {
				opt.path = "/tmp/glctl/glctl_logout_fail"
				return nil
			},
			expectedOutput: "not exist",
		},
	}
	for _, tc := range tests {
		streams := genericiooptions.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewLogoutCmd(streams)
			var cmdOptions = NewOptions(streams)
			if tc.optionsFunc != nil {
				err := tc.optionsFunc(cmdOptions)
				if err != nil {
					t.Errorf("%s: Invalid test case. Specify expected result.\n", tc.name)
					return
				}
			}
			out := cmdtesting.RunForStdout(streams, func() {
				var err error
				if tc.completeFunc != nil {
					if err = tc.completeFunc(streams); err != nil {
						_, _ = fmt.Fprint(streams.Out, err)
						return
					}
				} else {
					if err = cmdOptions.Complete(cmd, nil); err != nil {
						_, _ = fmt.Fprint(streams.Out, err)
						return
					}
				}
				defer func() {
					_ = os.Remove(cmdOptions.path)
				}()
				if err = cmdOptions.Validate(cmd, nil); err != nil {
					_, _ = fmt.Fprint(streams.Out, err)
					return
				}
				if err = cmdOptions.Run(nil); err != nil {
					_, _ = fmt.Fprint(streams.Out, err)
					return
				}
			})
			cmdtesting.TInfo(out)
			if tc.expectedOutput == "" {
				t.Errorf("%s: Invalid test case. Specify expected result.\n", tc.name)
			}
			assert.Containsf(
				t,
				out,
				tc.expectedOutput,
				"%s : Unexpected output! Expected\n%s\ngot\n%s",
				tc.name,
				tc.expectedOutput,
				out,
			)
		})
	}
}
