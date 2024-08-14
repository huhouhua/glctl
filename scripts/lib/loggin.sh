#!/usr/bin/env bash

set -e
set -o noglob

#
# Set Colors
#

readonly bold='\033[1m'
readonly underline='\033[4m'
readonly reset='\033[0m'

readonly red='\033[31m'
readonly green='\033[32m'
readonly white='\033[37m'
readonly tan='\033[33m'
readonly blue='\033[34m'

#
# Log function definition
#

underline() { printf "${underline}${bold}%s${reset}\n" "$@"
}
h1() { printf "\n${underline}${bold}${blue}%s${reset}\n" "$@"
}
h2() { printf "\n${underline}${bold}${white}%s${reset}\n" "$@"
}
debug() { printf "${white}%s${reset}\n" "$@"
}
info() { printf "${white}➜ %s${reset}\n" "$@"
}
success() { printf "${green}✔ %s${reset}\n" "$@"
}
error() { printf "${red}✖ %s${reset}\n" "$@"
}
warn() { printf "${tan}➜ %s${reset}\n" "$@"
}
bold() { printf "${bold}%s${reset}\n" "$@"
}
note() { printf "\n${underline}${bold}${blue}Note:${reset} ${blue}%s${reset}\n" "$@"
}
