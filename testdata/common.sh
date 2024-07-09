#!/bin/bash
#docker version: 20.10.10+
#docker-compose version: 1.18.0+

set -e
set -o noglob

#
# Set Colors
#

bold='\033[1m'
underline='\033[4m'
reset='\033[0m'

red='\033[31m'
green='\033[32m'
white='\033[37m'
tan='\033[33m'
blue='\033[34m'

#
# 日志函数定义
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

function check_docker {
	if ! docker --version &> /dev/null
	then
		error "Need to install docker(20.10.10+) first and run this script again."
		exit 1
	fi

	# docker has been installed and check its version
	if [[ $(docker --version) =~ (([0-9]+)\.([0-9]+)([\.0-9]*)) ]]
	then
		docker_version=${BASH_REMATCH[1]}
		docker_version_part1=${BASH_REMATCH[2]}
		docker_version_part2=${BASH_REMATCH[3]}

		note "docker version: $docker_version"
		# the version of docker does not meet the requirement
		if [ "$docker_version_part1" -lt 17 ] || ([ "$docker_version_part1" -eq 17 ] && [ "$docker_version_part2" -lt 6 ])
		then
			error "Need to upgrade docker package to 20.10.10+."
			exit 1
		fi
	else
		error "Failed to parse docker version."
		exit 1
	fi
}

function check_dockercompose {
	if [! docker compose version] &> /dev/null || [! docker-compose --version] &> /dev/null
	then
		error "Need to install docker-compose(1.18.0+) or a docker-compose-plugin (https://docs.docker.com/compose/)by yourself first and run this script again."
		exit 1
	fi

	# either docker compose plugin has been installed
	if docker compose version &> /dev/null
	then
		note "$(docker compose version)"
		DOCKER_COMPOSE="docker compose"

	# or docker-compose has been installed, check its version
	elif [[ $(docker-compose --version) =~ (([0-9]+)\.([0-9]+)([\.0-9]*)) ]]
	then
		docker_compose_version=${BASH_REMATCH[1]}
		docker_compose_version_part1=${BASH_REMATCH[2]}
		docker_compose_version_part2=${BASH_REMATCH[3]}

		note "docker-compose version: $docker_compose_version"
		# the version of docker-compose does not meet the requirement
		if [ "$docker_compose_version_part1" -lt 1 ] || ([ "$docker_compose_version_part1" -eq 1 ] && [ "$docker_compose_version_part2" -lt 18 ])
		then
			error "Need to upgrade docker-compose package to 1.18.0+."
			exit 1
		fi
	else
		error "Failed to parse docker-compose version."
		exit 1
	fi
}
