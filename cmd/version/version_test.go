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

package version

import (
	"fmt"
	"testing"

	"github.com/huhouhua/glctl/pkg/cli/genericiooptions"

	"github.com/stretchr/testify/assert"

	cmdtesting "github.com/huhouhua/glctl/cmd/testing"
	cmdutil "github.com/huhouhua/glctl/cmd/util"
)

func TestRunVersion(t *testing.T) {
	tests := []struct {
		name           string
		optionsFunc    func(opt *Options)
		args           []string
		expectedOutput string
	}{{
		name:           "lit all",
		args:           []string{},
		expectedOutput: "gitlab.Version{Version:\"14.7.2-ee\", Revision:\"39a169b2f25\"}",
	}, {
		name: "lit client",
		args: []string{},
		optionsFunc: func(opt *Options) {
			opt.ClientOnly = true
		},
		expectedOutput: "Client Version",
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		streams := genericiooptions.NewTestIOStreamsForPipe()
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCmdVersion(factory, streams)
			cmdOptions := NewOptions(streams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			out := cmdtesting.RunForStdout(streams, func() {
				var err error
				if err = cmdOptions.Complete(factory, cmd); err != nil {
					_, _ = fmt.Fprint(streams.Out, err.Error())
					return
				}
				if err = cmdOptions.Validate(); err != nil {
					_, _ = fmt.Fprint(streams.Out, err.Error())
					return
				}
				if err = cmdOptions.Run(); err != nil {
					_, _ = fmt.Fprint(streams.Out, err.Error())
					return
				}
			})
			cmdtesting.TInfo("\n" + out)
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
