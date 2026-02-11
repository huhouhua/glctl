#!/usr/bin/env bash

# Copyright 2024 The Kevin Berger <huhouhuam@outlook.com> Authors
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


set -e
set -o errexit
set -o nounset
set -o pipefail

GLCTL_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source ${GLCTL_ROOT}/scripts/lib/loggin.sh
source ${GLCTL_ROOT}/scripts/lib/environment.sh
source ${GLCTL_ROOT}/scripts/lib/check.sh

INSTALL="false"
UNINSTALL="false"
START="false"
STOP="false"

usage=$'Provide Gitlab installation,uninstall,start, stop, and add test data
If you need to install gitlab, please set --install
If you need to un install gitlab, please set --uninstall
If you need to run gitlab service, please set --start
If you need to stop gitlab service, please set --stop
To view help information, please set --help'
item=0
while [ $# -gt 0 ]; do
  case $1 in
  --help)
    note "$usage"
    exit 0
    ;;
  --install)
    INSTALL="true"
    ;;
  --uninstall)
    UNINSTALL="true"
    ;;
  --stop)
    STOP="true"
    ;;
  --start)
    START="true"
    ;;
  *)
    note "$usage"
    exit 1
    ;;
  esac
  shift || true
done

h2 "[Step $item]: checking if docker is installed ..."; let item+=1
check::docker

h2 "[Step $item]: checking docker compose is installed ..."; let item+=1
check::docker::compose

h2 "[Step $item]: checking $DOCKER_COMPOSE_TEST_FILE ..."
let item+=1
if [ ! -f ${DOCKER_COMPOSE_TEST_FILE} ]; then
    error "${DOCKER_COMPOSE_TEST_FILE} docker-compose file not exist!"
    exit 1
fi
echo ""

if [ "$INSTALL" == "true" ]; then
h2 "[Step $item]: install..."
let item+=1

  "${DOCKER_COMPOSE[@]}" --project-name gitlab-test -f ${DOCKER_COMPOSE_TEST_FILE} up -d
  wait::gitlab
  exit 0
fi

if [ "$UNINSTALL" == "true" ]; then
h2 "[Step $item]: uninstall..."
let item+=1

  "${DOCKER_COMPOSE[@]}" -f ${DOCKER_COMPOSE_TEST_FILE} down -v
  success $"---- uninstall successfully.----"
  exit 0
fi

if [ "$STOP" == "true" ]; then
h2 "[Step $item]: stop..."
let item+=1

  ${DOCKER} stop gitlab-test
  success $"---- stop successfully.----"
  exit 0
fi

if [ "$START" == "true" ]; then
h2 "[Step $item]: start..."
let item+=1

  ${DOCKER} start gitlab-test
  wait::gitlab
  exit 0
fi