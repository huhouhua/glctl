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


export GITLAB_USERNAME=root
export GITLAB_PASSWORD=123qwe123
export GITLAB_URL=${GITLAB_URL:-http://localhost:8080}
./testdata/get_token.sh
export GITLAB_PRIVATE_TOKEN=$(cat token.txt)
export GITLAB_OAUTH_TOKEN=$(./testdata/get_oauth_token.sh)