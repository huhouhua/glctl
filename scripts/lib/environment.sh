#!/usr/bin/env bash

DIR="$(cd "$(dirname "$0")" && pwd)"

readonly DOCKER_COMPOSE_TEST_FILE=$DIR/docker-compose-test.yaml
readonly DOCKER_COMPOSE=docker-compose
readonly DOCKER=docker
