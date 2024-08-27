#!/usr/bin/env bash

# Copyright 2024 The huhouhua Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http:www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Function install jq on Debian/Ubuntu
install::jq::debian()
{
    info "Installing jq on Debian/Ubuntu..."
    sudo apt-get update
    sudo apt-get install -y jq
}

# Function to install jq on CentOS/RHEL/Fedora
install::jq::centos_fedora()
{
    info "Installing jq on CentOS/RHEL/Fedora..."
    if [ -f /etc/yum.repos.d/epel.repo ]; then
        sudo yum install -y jq
    else
        sudo yum install -y epel-release
        sudo yum install -y jq
    fi
}

# Function to install jq on macOS
install::jq::macos()
{
    info "Installing jq on macOS..."
    if command -v brew >/dev/null 2>&1; then
        brew install jq
    else
        note "Homebrew is not installed. Please install Homebrew first."
        exit 1
    fi
}

# Detect OS and install jq
install::jq()
{
  if [ -f /etc/os-release ]; then
      . /etc/os-release
      case "$ID" in
          ubuntu|debian)
              install::jq::debian
              ;;
          centos|rhel|fedora)
              install::jq::centos_fedora
              ;;
          *)
              error "Unsupported Linux distribution. Please install jq manually."
              exit 1
              ;;
      esac
  elif [ "$(uname)" == "Darwin" ]; then
      install::jq::macos
  else
      error "Unsupported operating system. Please install jq manually."
      exit 1
  fi
}
