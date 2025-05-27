# glctl

![Workflow ci](https://github.com/huhouhua/glctl/actions/workflows/glctl.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/huhouhua/glctl)](https://goreportcard.com/report/github.com/huhouhua/glctl)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/huhouhua/glctl?status.svg)](https://godoc.org/github.com/huhouhua/glctl)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/huhouhua/glctl?logo=go)
[![Test Coverage](https://codecov.io/gh/huhouhua/glctl/branch/main/graph/badge.svg)](https://codecov.io/gh/huhouhua/glctl)

glctl is a powerful GitLab command line tool. It provides a declarative API to manage GitLab resources, making it easier for you to perform common GitLab operations from the terminal.

## Features
- Manage GitLab projects, issues, merge requests, and more
- Authenticate and manage GitLab sessions
- Create, read, update, and delete GitLab resources
- Validate GitLab configurations
- Shell completion support

### Usage
  ```bash
  glctl <command> <subcommand> [flags]
  ```

## Installation

### From Binary

Download the appropriate version for your platform from the [releases page](https://github.com/huhouhua/glctl/releases).

### From Source
 - compile glctl and place it in _output/
```bash
git clone https://github.com/huhouhua/glctl.git
cd glctl
make build
```

## Quick Start

### Authentication

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

### Create a new project:
```bash
glctl create project my-new-project
```

### List your projects:
```bash
glctl get projects
```

## Available Commands

- `login` - Authenticate with GitLab
- `logout` - Log out from GitLab
- `create` - Create new GitLab resources (projects, issues, merge requests, etc.)
- `get` - Get information about GitLab resources
- `edit` - Edit existing GitLab resources
- `delete` - Delete GitLab resources
- `replace` - Replace existing GitLab resources
- `validate` - Validate GitLab configurations
- `version` - Display version information
- `completion` - Generate shell completion scripts


## Logged in user authorization file
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

## Issues

If you have an issue: report it on the [issue tracker](https://github.com/huhouhua/glctl/issues)

## Author

Kevin Berger (<huhouhuam@outlook.com>)

## Contributing

Contributions are always welcome. For more information, check out the [contributing guide](CONTRIBUTING.md)

## License

Licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.