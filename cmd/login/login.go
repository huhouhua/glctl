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

package login

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/huhouhua/gitlab-repo-operator/cmd/require"
	"github.com/huhouhua/gitlab-repo-operator/cmd/types"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"strings"
)

var loginDesc = "This command authenticates you to a Gitlab server, retrieves your OAuth Token and then save it to $HOME/.grepo.yaml file."

type loginOptions struct {
	ServerAddress string
	User          string
	Password      string
}

func NewLoginCmd() *cobra.Command {
	var o loginOptions
	cmd := &cobra.Command{
		Use:     "login [OPTIONS] [SERVER]",
		Short:   "Login to gitlab",
		Long:    loginDesc,
		Example: `grepo login http://localhost:8080`,
		Args:    require.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&o.User, "username", "u", "", "Username")
	flags.StringVarP(&o.Password, "password", "p", "", "Password")
	//cmdutil.VerifyMarkFlagRequired(cmd, "project")

	return cmd
}

// Complete completes all the required options.
func (o *loginOptions) Complete(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		o.ServerAddress = args[0]
	}
	if strings.TrimSpace(o.User) == "" {
		o.User = promptUserNameInput()
	}
	if strings.TrimSpace(o.Password) == "" {
		o.Password = promptPasswordInput()
	}
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *loginOptions) Validate(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(o.ServerAddress) == "" {
		return fmt.Errorf("please enter the gitlab url")
	}
	if strings.TrimSpace(o.User) == "" {
		return fmt.Errorf("please enter the username ")
	}
	if strings.TrimSpace(o.Password) == "" {
		return fmt.Errorf("please enter the password ")
	}
	return nil
}

// Run executes a create subcommand using the specified options.
func (o *loginOptions) Run(args []string) error {
	uri := fmt.Sprintf("%s/oauth/token?grant_type=password&username=%s&password=%s", o.ServerAddress, o.User, o.Password)
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
	var cfg = types.GitLabOauthInfo{}
	err = mapstructure.Decode(cfgMap, &cfg)
	if err != nil {
		return err
	}
	cfgFile := fmt.Sprintf("%s/.grepo.yaml", home)
	// add host_url and user to config file
	cfg.HostUrl = &o.ServerAddress
	cfg.UserName = &o.User
	b, err = yaml.Marshal(cfg)
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

func promptPasswordInput() string {
	for i := 0; i < 3; i++ {
		fmt.Print("Password: ")
		input, err := gopass.GetPasswd()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		password := strings.TrimSpace(string(input))
		if password != "" {
			return password
		}
	}
	return ""
}

func promptUserNameInput() string {
	for i := 0; i < 3; i++ {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s: ", "Username")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		username := strings.Replace(input, "\n", "", -1)
		if strings.TrimSpace(username) != "" {
			return username
		}
	}
	return ""
}
