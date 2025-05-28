
<div align="center">
	<h1>Gitlab CLI</h1>
	<p>glctl is a powerful GitLab command line tool. It provides a declarative API to manage GitLab resources, making it easier for you to perform common GitLab operations from the terminal.!</p>
</div>

<p align="center">
	<a href="#-installation">Installation</a> ❘
	<a href="#-features">Features</a> ❘
	<a href="#-quick-start">Quick-Start</a> ❘
	<a href="#-examples">Examples</a> ❘
	<a href="#-license">License</a>
</p>


![Workflow ci](https://github.com/huhouhua/glctl/actions/workflows/glctl.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/huhouhua/glctl)](https://goreportcard.com/report/github.com/huhouhua/glctl)
[![release](https://img.shields.io/github/release-pre/huhouhua/glctl.svg)](https://github.com/huhouhua/glctl/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/huhouhua/glctl?status.svg)](https://godoc.org/github.com/huhouhua/glctl)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/huhouhua/glctl?logo=go)
[![Test Coverage](https://codecov.io/gh/huhouhua/glctl/branch/main/graph/badge.svg)](https://codecov.io/gh/huhouhua/glctl)

```
Basic Commands:
  get         Display one or many resources
  edit        Edit a resource on the server
  delete      Delete resources by file names, stdin, resources and names, or by resources
  create      Create a resource from a file or from stdin

Authorization Commands:
  login       Login to gitlab
  logout      logout current gitlab

Advanced Commands:
  replace     Replace a repository file by file name or stdin.

 all formats are accepted. If replacing an existing repository file, the
complete repository file spec must be provided. This can be obtained by

  $ glctl get files PROJECT --path=my.yml --ref=BRANCH --raw

Settings Commands:
  completion  Output shell completion code for the specified shell (bash, zsh,
fish, or powershell)

Other Commands:
  version     Print the client and server version information

```

## 🤘&nbsp; Features
- Manage GitLab projects, issues, merge requests, and more
- Authenticate and manage GitLab sessions
- Create, get, edit, and delete GitLab resources
- Operation project branch file
- Shell completion support


## 📦&nbsp; Installation

### 📁 From Binary

Download the appropriate version for your platform from the [releases page](https://github.com/huhouhua/glctl/releases).

### 🛠️ From Source
 - compile glctl and place it in _output/
```bash
git clone https://github.com/huhouhua/glctl.git
cd glctl
make build
```

## 🚀&nbsp;Quick Start

### 📄&nbsp;Usage
  ```bash
  glctl <command> <subcommand> [flags]
  ```

### 🔐 Authentication

- start interactive login
```bash
glctl login https://gitlab.example.com
```
- start interactive Login by username
```bash
glctl login https://gitlab.example.com --username myname
```

- Login by specifying username and password
```bash
glctl login https://gitlab.example.com --username myname --password mypassword
```

- authenticate with private token and hostname
```bash
export GITLAB_URL=https://gitlab.example.com
export GITLAB_PRIVATE_TOKEN=305e146a4aa23fb4021a4f162102251e85f651a058a34fb2c27d633617cf8877
```

- authenticate with oauth token and hostname
```bash
export GITLAB_URL=https://gitlab.example.com
export GITLAB_OAUTH_TOKEN=aefb8b4e0895799aa60cf50eb8bcd9ae1fecf08fb6cc8249238219067e5aa926
```

- Loggin in using environment variables (Not recommended for shared environments)
```bash
export GITLAB_URL=https://gitlab.example.com
export GITLAB_USERNAME=myname
export GITLAB_PASSWORD=mypassword
```

## 🥙&nbsp;Example
- List your groups
```bash
glctl get groups
```

- List your projects:
```bash
glctl get projects
```

- List the group1/project1 branch:
```bash
glctl get branchs group1/project1
```

- create a develop branch from master branch for project group1/project1
```bash
glctl create branch develop --project=group1/project1 --ref=master
```

### 🥪 Available Commands

- `login` - Authenticate with GitLab
- `logout` - Log out from GitLab
- `create` - Create new GitLab resources (projects, issues, merge requests, etc.)
- `get` - Get information about GitLab resources
- `edit` - Edit existing GitLab resources
- `delete` - Delete GitLab resources
- `replace` - Replace existing GitLab resources
- `version` - Display version information
- `completion` - Generate shell completion scripts

### 🗒️&nbsp;Logged in user authorization file
Files are stored in `$HOME/.glctl.yaml` example:
```yaml
access_token: 305e146a4aa23fb4021a4f162102251e85f651a058a34fb2c27d633617cf8877
created_at: 1.748339041e+09
host_url: https://gitlab.example.com
refresh_token: aefb8b4e0895799aa60cf50eb8bcd9ae1fecf08fb6cc8249238219067e5aa926
scope: api
token_type: Bearer
user_name: root
```

## 🧠&nbsp;TODOs

- This cli tool is still in the development stage, and most of the resources are not completed. Your help is very much needed. 🙋‍♂️
- Declarative resources are still in the design stage

## 🤝&nbsp;Issues

If you have an issue: report it on the [issue tracker](https://github.com/huhouhua/glctl/issues)

## 👤&nbsp;Author

Kevin Berger (<huhouhuam@outlook.com>)

## 🧑‍💻&nbsp;Contributing

Contributions are always welcome. For more information, check out the [contributing guide](CONTRIBUTING.md)

## 📘&nbsp;License

Licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.