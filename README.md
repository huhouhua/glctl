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

1. Login to your GitLab instance:
```bash
glctl login https://gitlab.example.com --username myname
```

2. Create a new project:
```bash
glctl create project my-new-project
```

3. List your projects:
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

## Development

### Prerequisites
- Go 1.16 or later
- Make
- GitLab instance for testing

### Building from Source
```bash
make build
```

### Running Tests
```bash
make test
```

### Code Style
```bash
make lint
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the [GitHub repository](https://github.com/huhouhua/glctl/issues).
