#!/usr/bin/env bash

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
