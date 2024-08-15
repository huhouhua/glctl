#!/usr/bin/env bash

export GITLAB_USERNAME=root
export GITLAB_PASSWORD=123qwe123
export GITLAB_URL=${GITLAB_URL:-http://localhost:8080}
./testdata/get_token.sh
export GITLAB_PRIVATE_TOKEN=$(cat token.txt)
export GITLAB_OAUTH_TOKEN=$(./testdata/get_oauth_token.sh)