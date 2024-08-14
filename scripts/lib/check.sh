#!/usr/bin/env bash
#docker version: 20.10.10+
#docker-compose version: 1.18.0+

set -e

function check::docker()
{
	if ! ${DOCKER}  --version &> /dev/null
	then
		error "Need to install docker(20.10.10+) first and run this script again."
		exit 1
	fi

	# docker has been installed and check its version
	if [[ $(${DOCKER} --version) =~ (([0-9]+)\.([0-9]+)([\.0-9]*)) ]]
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

function check::docker::compose()
{
	if [! docker compose version] &> /dev/null || [! ${DOCKER_COMPOSE}  --version] &> /dev/null
	then
		error "Need to install docker-compose(1.18.0+) or a docker-compose-plugin (https://docs.docker.com/compose/)by yourself first and run this script again."
		exit 1
	fi

	# either docker compose plugin has been installed
	if docker compose version &> /dev/null
	then
		note "$(docker compose version)"

	# or docker-compose has been installed, check its version
	elif [[ $(${DOCKER_COMPOSE} --version) =~ (([0-9]+)\.([0-9]+)([\.0-9]*)) ]]
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