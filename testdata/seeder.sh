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

source ../scripts/lib/loggin.sh

Header="-HContent-Type: application/json"
CCURL="curl -f -s -XPOST" # Create
UCURL="curl -f -s -XPUT" # Update
RCURL="curl -f -s -XGET" # Retrieve
DCURL="curl -f -s -XDELETE" # Delete

gitlab::credentials()
{
  export GITLAB_USERNAME=root
  export GITLAB_PASSWORD=$(gitlab::get::password)
  export GITLAB_URL=${GITLAB_URL:-http://localhost:8080}
  export GITLAB_PRIVATE_TOKEN=$(gitlab:get::token)
  export GITLAB_OAUTH_TOKEN=$(gitlab::get::oath_token)

}

# get access token
gitlab::login(){
  ${CCURL} "${GITLAB_URL}/oauth/token?grant_type=password&username=${GITLAB_USERNAME}&password=${GITLAB_PASSWORD}" | jq '.access_token' | tr -d '"'
}

gitlab::get::password()
{


}

gitlab::get:oath_token()
{
  ${CCURL}  "${GITLAB_URL}/oauth/token?grant_type=password&username=${GITLAB_USERNAME}&password=${GITLAB_PASSWORD}" | jq '.access_token' | tr -d '"'
}

  #
  # From: https://stackoverflow.com/questions/47948887/login-to-gitlab-using-curl
  #
gitlab::get:token()
{
  # curl for the login page to get a session cookie and the sources with the auth tokens
  body_header=$(${RCURL} -c cookies.txt -i "${GITLAB_URL}/users/sign_in" -s)

  # grep the auth token for the user login for
  #   not sure whether another token on the page will work, too - there are 3 of them
  csrf_token=$(echo $body_header | perl -ne 'print "$1\n" if /new_user.*?authenticity_token"[[:blank:]]value="(.+?)"/' | sed -n 1p)

  # send login credentials with curl, using cookies and token from previous request
  ${RCURL} -b cookies.txt -c cookies.txt -i "${GITLAB_URL}/users/sign_in" \
      --data "user[login]=${GITLAB_USERNAME}&user[password]=${GITLAB_PASSWORD}" \
      --data-urlencode "authenticity_token=${csrf_token}"

  # send curl GET request to personal access token page to get auth token
  body_header=$(${RCURL} -H 'user-agent: curl' -b cookies.txt -i "${GITLAB_URL}/profile/personal_access_tokens" -s)
  csrf_token=$(echo $body_header | perl -ne 'print "$1\n" if /authenticity_token"[[:blank:]]value="(.+?)"/' | sed -n 1p)

  # curl POST request to send the "generate personal access token form"
  # the response will be a redirect, so we have to follow using `-L`
  body_header=$(${RCURL} -L -b cookies.txt "${GITLAB_URL}/profile/personal_access_tokens" \
      --data-urlencode "authenticity_token=${csrf_token}" \
      --data 'personal_access_token[name]=golab-generated&personal_access_token[expires_at]=&personal_access_token[scopes][]=api')

  # Scrape the personal access token from the response HTML
  personal_access_token=$(echo $body_header | perl -ne 'print "$1\n" if /created-personal-access-token"[[:blank:]]value="(.+?)"/' | sed -n 1p)

  return  $personal_access_token
}

echo "=============================================="
echo "#              SEEDER SCRIPT                 #"
echo "=============================================="


# create user
echo "Creating users"
user1_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=John+Doe&email=john.doe@gmail.com&username=john.doe&password=123qwe123&skip_confirmation=true" | jq '.id')
user2_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=John+Smith&email=john.smith@gmail.com&username=john.smith&password=123qwe123&skip_confirmation=true" | jq '.id')
user3_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=Matt+Hunter&email=matt.hunter@gmail.com&username=matt.hunter&password=123qwe123&skip_confirmation=true" | jq '.id')
user4_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=Amelia+Walsh&email=amelia.walsh@gmail.com&username=amelia.walsh&password=123qwe123&skip_confirmation=true" | jq '.id')
user5_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=Kevin+McLean&email=kevin.mclean@gmail.com&username=kevin.mclean&password=123qwe123&skip_confirmation=true" | jq '.id')
user6_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=Kylie+Morrison&email=kylie.morrison@gmail.com&username=kylie.morrison&password=123qwe123&skip_confirmation=true" | jq '.id')
user7_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=Rebecca+Gray&email=rebecca.gray@gmail.com&username=rebecca.gray&password=123qwe123&skip_confirmation=true" | jq '.id')
user8_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=Simon+Turner&email=simon.turner@gmail.com&username=simon.turner&password=123qwe123&skip_confirmation=true" | jq '.id')
user9_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=Olivia+White&email=olivia.white@gmail.com&username=olivia.white&password=123qwe123&skip_confirmation=true" | jq '.id')
user10_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=Frank+Watson&email=frank.watson@gmail.com&username=frank.watson&password=123qwe123&skip_confirmation=true" | jq '.id')
user11_id=$(curl -X POST "${GITLAB_URL}/api/v4/users?access_token=${access_token}&name=Paul+Lyman&email=paul.lyman@gmail.com&username=paul.lyman&password=123qwe123&skip_confirmation=true" | jq '.id')

# create group
echo "Creating groups"
pgroup1_id=$(curl -X POST "${GITLAB_URL}/api/v4/groups?access_token=${access_token}&name=Group1&path=Group1" | jq '.id')
pgroup2_id=$(curl -X POST "${GITLAB_URL}/api/v4/groups?access_token=${access_token}&name=Group2&path=Group2" | jq '.id')

# add user to group
echo "Adding users to group"
curl -X POST "${GITLAB_URL}/api/v4/groups/${pgroup1_id}/members?access_token=${access_token}&user_id=${user1_id}&access_level=30"
curl -X POST "${GITLAB_URL}/api/v4/groups/${pgroup1_id}/members?access_token=${access_token}&user_id=${user2_id}&access_level=40"
curl -X POST "${GITLAB_URL}/api/v4/groups/${pgroup1_id}/members?access_token=${access_token}&user_id=${user3_id}&access_level=50"
curl -X POST "${GITLAB_URL}/api/v4/groups/${pgroup2_id}/members?access_token=${access_token}&user_id=${user4_id}&access_level=30"
curl -X POST "${GITLAB_URL}/api/v4/groups/${pgroup2_id}/members?access_token=${access_token}&user_id=${user5_id}&access_level=40"
curl -X POST "${GITLAB_URL}/api/v4/groups/${pgroup2_id}/members?access_token=${access_token}&user_id=${user6_id}&access_level=50"

# create subgroup
echo "Creating subgroup"
sgroup1_id=$(curl -X POST "${GITLAB_URL}/api/v4/groups?access_token=${access_token}&name=SubGroup1&path=SubGroup1&parent_id=${pgroup1_id}" | jq '.id')
sgroup2_id=$(curl -X POST "${GITLAB_URL}/api/v4/groups?access_token=${access_token}&name=SubGroup2&path=SubGroup2&parent_id=${pgroup1_id}" | jq '.id')
sgroup3_id=$(curl -X POST "${GITLAB_URL}/api/v4/groups?access_token=${access_token}&name=SubGroup3&path=SubGroup3&parent_id=${pgroup2_id}" | jq '.id')
sgroup4_id=$(curl -X POST "${GITLAB_URL}/api/v4/groups?access_token=${access_token}&name=SubGroup4&path=SubGroup4&parent_id=${pgroup2_id}" | jq '.id')

echo sleeping for 5 seconds..
sleep 5

# create group project
echo "Creating a project in group/subgroup"
groupproject1_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project1&namespace_id=${pgroup1_id}" | jq '.id')
groupproject2_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project2&namespace_id=${pgroup1_id}" | jq '.id')
groupproject3_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project3&namespace_id=${pgroup1_id}" | jq '.id')
groupproject4_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project4&namespace_id=${sgroup1_id}" | jq '.id')
groupproject5_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project5&namespace_id=${sgroup1_id}" | jq '.id')
groupproject6_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project6&namespace_id=${sgroup1_id}" | jq '.id')
groupproject7_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project7&namespace_id=${sgroup2_id}" | jq '.id')
groupproject8_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project8&namespace_id=${sgroup2_id}" | jq '.id')
groupproject9_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project9&namespace_id=${sgroup2_id}" | jq '.id')
groupproject10_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project10&namespace_id=${pgroup2_id}" | jq '.id')
groupproject11_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project11&namespace_id=${pgroup2_id}" | jq '.id')
groupproject12_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project12&namespace_id=${pgroup2_id}" | jq '.id')
groupproject13_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project13&namespace_id=${sgroup3_id}" | jq '.id')
groupproject14_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project14&namespace_id=${sgroup3_id}" | jq '.id')
groupproject15_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project15&namespace_id=${sgroup3_id}" | jq '.id')
groupproject16_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project16&namespace_id=${sgroup4_id}" | jq '.id')
groupproject17_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project17&namespace_id=${sgroup4_id}" | jq '.id')
groupproject18_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects?access_token=${access_token}&name=Project18&namespace_id=${sgroup4_id}" | jq '.id')

# add user to project
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject1_id}/members?access_token=${access_token}&user_id=${user1_id}&access_level=30"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject1_id}/members?access_token=${access_token}&user_id=${user2_id}&access_level=40"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject1_id}/members?access_token=${access_token}&user_id=${user3_id}&access_level=50"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject1_id}/members?access_token=${access_token}&user_id=${user4_id}&access_level=30"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject2_id}/members?access_token=${access_token}&user_id=${user5_id}&access_level=40"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject2_id}/members?access_token=${access_token}&user_id=${user6_id}&access_level=50"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject2_id}/members?access_token=${access_token}&user_id=${user7_id}&access_level=40"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject2_id}/members?access_token=${access_token}&user_id=${user8_id}&access_level=40"

# create user project
echo "Creating users project"
project1_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects/user/${user7_id}?access_token=${access_token}&name=Project19" | jq '.id')
project2_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects/user/${user8_id}?access_token=${access_token}&name=Project20" | jq '.id')
project3_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects/user/${user9_id}?access_token=${access_token}&name=Project21" | jq '.id')
project4_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects/user/${user10_id}?access_token=${access_token}&name=Project22" | jq '.id')
project5_id=$(curl -X POST "${GITLAB_URL}/api/v4/projects/user/${user11_id}?access_token=${access_token}&name=Project23" | jq '.id')

# create hooks for projects
echo "Creating hooks for projects"
curl -X POST "${GITLAB_URL}/api/v4/projects/${project1_id}/hooks?access_token=${access_token}&url=http%3A%2F%2Fwww.sample1.com%2F"
curl -X POST "${GITLAB_URL}/api/v4/projects/${project2_id}/hooks?access_token=${access_token}&url=http%3A%2F%2Fwww.sample2.com%2F"
curl -X POST "${GITLAB_URL}/api/v4/projects/${project3_id}/hooks?access_token=${access_token}&url=http%3A%2F%2Fwww.sample3.com%2F"
curl -X POST "${GITLAB_URL}/api/v4/projects/${project4_id}/hooks?access_token=${access_token}&url=http%3A%2F%2Fwww.sample4.com%2F"
curl -X POST "${GITLAB_URL}/api/v4/projects/${project5_id}/hooks?access_token=${access_token}&url=http%3A%2F%2Fwww.sample5.com%2F"

# push some commits in preparation for creation of git tags
echo "Push commits for projects"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject1_id}/repository/files/README.md?access_token=${access_token}&branch=master&content=Test&commit_message=First+commit"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject2_id}/repository/files/README.md?access_token=${access_token}&branch=master&content=Test&commit_message=First+commit"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject3_id}/repository/files/README.md?access_token=${access_token}&branch=master&content=Test&commit_message=First+commit"

# create tags
echo "Creating tags for projects"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject1_id}/repository/tags?access_token=${access_token}&tag_name=sample_1.0&ref=master"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject2_id}/repository/tags?access_token=${access_token}&tag_name=sample_1.0&ref=master"
curl -X POST "${GITLAB_URL}/api/v4/projects/${groupproject3_id}/repository/tags?access_token=${access_token}&tag_name=sample_1.0&ref=master"