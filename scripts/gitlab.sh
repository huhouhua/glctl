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


set -e
set -o errexit
set -o nounset
set -o pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"
source $DIR/lib/loggin.sh
source $DIR/lib/environment.sh
source $DIR/lib/check.sh

INSTALL="false"
UNINSTALL="false"
UNPAUSE="false"
PAUSE="false"
TEST_DATA="false"

usage=$'Provide Gitlab installation,uninstall,start, stop, and add test data
If you need to install gitlab, please set --install
If you need to un install gitlab, please set --uninstall
If you need to un pause gitlab, please set --unpause
If you need to pause gitlab, please set --pause
If you need to add test data to gitlab, please set --data
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
  --pause)
    PAUSE="true"
    ;;
  --unpause)
    UNPAUSE="true"
    ;;
  --data)
    TEST_DATA="true"
    ;;
  *)
    note "$usage"
    exit 1
    ;;
  esac
  shift || true
done

export HOST_IP=$(hostname -I | awk '{print $1}')

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

  ${DOCKER_COMPOSE} --project-name gitlab -f ${DOCKER_COMPOSE_TEST_FILE} up -d
  info "Waiting for GitLab to become healthy..."
  until [ "$(${DOCKER} inspect --format='{{.State.Health.Status}}' gitlab)"  == "healthy" ]; do
    sleep 5
  done
  success "GitLab is up and running!"
  exit 0
fi

if [ "$UNINSTALL" == "true" ]; then
h2 "[Step $item]: uninstall..."
let item+=1

  ${DOCKER_COMPOSE} -f ${DOCKER_COMPOSE_TEST_FILE} down -v
  success $"---- uninstall successfully.----"
  exit 0
fi

if [ "$PAUSE" == "true" ]; then
h2 "[Step $item]: pause..."
let item+=1

  ${DOCKER_COMPOSE} -f ${DOCKER_COMPOSE_TEST_FILE} pause
  success $"---- pause successfully.----"
  exit 0
fi

if [ "$UNPAUSE" == "true" ]; then
h2 "[Step $item]: un pause..."
let item+=1

  ${DOCKER_COMPOSE} -f ${DOCKER_COMPOSE_TEST_FILE} unpause
  success $"---- un pause successfully.----"
  exit 0
fi