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


Header="-HContent-Type: application/json"
CCURL="curl -f -s -XPOST" # Create
UCURL="curl -f -s -XPUT" # Update
RCURL="curl -f -s -XGET" # Retrieve
DCURL="curl -f -s -XDELETE" # Delete

# get access token
gitlab::login(){
  ${CCURL} "${GITLAB_URL}/oauth/token" -d "grant_type=password&username=${GITLAB_USERNAME}&password=${GITLAB_PASSWORD}" | jq '.access_token' | tr -d '"'
}

  #
  # From: https://stackoverflow.com/questions/47948887/login-to-gitlab-using-curl
  #
gitlab::get::token()
{
# curl for the login page to get a session cookie and the sources with the auth tokens
body_header=$(curl  -c cookies.txt -i "${GITLAB_URL}/users/sign_in" -s)

# grep the auth token for the user login for
#   not sure whether another token on the page will work, too - there are 3 of them
csrf_token=$(echo $body_header | perl -ne 'print "$1\n" if /new_user.*?authenticity_token"[[:blank:]]value="(.+?)"/' | sed -n 1p)

# send login credentials with curl, using cookies and token from previous request
curl --location  -b cookies.txt -c cookies.txt -i "${GITLAB_URL}/users/sign_in" \
    --data "user[login]=${GITLAB_USERNAME}&user[password]=${GITLAB_PASSWORD}" \
    --data-urlencode "authenticity_token=${csrf_token}"

# send curl GET request to personal access token page to get auth token
body_header=$(curl  -H 'user-agent: curl' -b cookies.txt -i "${GITLAB_URL}/profile/personal_access_tokens" -s)

csrf_token=$(echo $body_header | perl -ne 'print "$1\n" if /authenticity_token"[[:blank:]]value="(.+?)"/' | sed -n 1p)

# curl POST request to send the "generate personal access token form"
# the response will be a redirect, so we have to follow using `-L`
body_header=$(curl -L -b cookies.txt "${GITLAB_URL}/profile/personal_access_tokens" \
    --data-urlencode "authenticity_token=${csrf_token}" \
    --data 'personal_access_token[name]=golab-generated&personal_access_token[expires_at]=&personal_access_token[scopes][]=api')

# Scrape the personal access token from the response HTML
personal_access_token=$(echo $body_header | perl -ne 'print "$1\n" if /created-personal-access-token"[[:blank:]]value="(.+?)"/' | sed -n 1p)
return  $personal_access_token
}

export GITLAB_USERNAME=${GITLAB_USERNAME:-root}
export GITLAB_PASSWORD=${GITLAB_PASSWORD:-password}
export GITLAB_URL=${GITLAB_URL:-http://localhost:8080}
export GITLAB_PRIVATE_TOKEN=${GITLAB_PRIVATE_TOKEN:-$(gitlab::login)}