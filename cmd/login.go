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

package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"strings"
)

var loginDesc = "This command authenticates you to a Gitlab server, retrieves your OAuth Token and then save it to $HOME/.grepo.yaml file."

type loginOptions struct {
	serverAddress string
	user          string
	password      string
}

func newLoginCmd() *cobra.Command {
	var opts loginOptions
	cmd := &cobra.Command{
		Use:               "login [OPTIONS] [SERVER]",
		Short:             "Login to gitlab",
		Long:              loginDesc,
		Example:           `grepo login http://localhost:8080`,
		Args:              RequiresMaxArgs(1),
		SilenceErrors:     true,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.serverAddress = args[0]
			}
			return runLogin(cmd.Context(), opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.user, "username", "u", "", "Username")
	flags.StringVarP(&opts.password, "password", "p", "", "Password")
	return cmd
}

func runLogin(ctx context.Context, opts loginOptions) error {
	if opts.serverAddress == "" {
		return fmt.Errorf("please enter the gitlab url")
	}
	if opts.user == "" {
		u, err := promptStringInput("Username")
		if err != nil {
			return err
		}
		opts.user = u
	}
	if opts.password == "" {
		p, err := promptPasswordInput()
		if err != nil {
			return err
		}
		opts.password = p
	}
	uri := fmt.Sprintf("%s/oauth/token?grant_type=password&username=%s&password=%s", opts.serverAddress, opts.user, opts.password)
	resp, err := http.Post(uri, "application/json", nil)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var cfgMap map[string]interface{}
	if err = json.Unmarshal(b, &cfgMap); err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Login failed!\n%s", cfgMap["error_description"])
	}
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	cfgFile := fmt.Sprintf("%s/.grepo.yaml", home)
	// add host_url and user to config file
	cfgMap["host_url"] = opts.serverAddress
	cfgMap["user"] = opts.user
	b, err = yaml.Marshal(cfgMap)
	if err != nil {
		return err
	}
	if err := os.WriteFile(cfgFile, b, 0600); err != nil {
		return err
	}
	fmt.Printf("%s file has been created by login command \n", cfgFile)
	fmt.Printf("\nLogin Succeeded \n")
	return nil
}

func promptPasswordInput() (string, error) {
	fmt.Print("Password: ")
	password, err := gopass.GetPasswd()
	return strings.TrimSpace(string(password)), err
}

func promptStringInput(askFor string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", askFor)
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// convert CRLF to LF
	return strings.Replace(username, "\n", "", -1), nil
}
