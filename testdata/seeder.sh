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

DIR="$(cd "$(dirname "$0")" && pwd)"
source $DIR/../scripts/lib/loggin.sh
source $DIR/credentials.sh

gitlab::credentials

info "GITLAB_USERNAME:${GITLAB_USERNAME}"
info "GITLAB_PASSWORD:${GITLAB_PASSWORD}"
info "GITLAB_URL:${GITLAB_URL}"
info "GITLAB_PRIVATE_TOKEN:${GITLAB_PRIVATE_TOKEN}"
info "GITLAB_OAUTH_TOKEN:${GITLAB_OAUTH_TOKEN}"

echo "=============================================="
echo "#              SEEDER SCRIPT                 #"
echo "=============================================="

set -e

echo ""
# create user
info "Creating users......"
user1_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=John+Doe&email=john.doe@gmail.com&username=john.doe&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user1_id} john.doe created"
user2_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=John+Smith&email=john.smith@gmail.com&username=john.smith&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user2_id} john.smith created"
user3_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Matt+Hunter&email=matt.hunter@gmail.com&username=matt.hunter&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user3_id} matt.hunter created"
user4_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Amelia+Walsh&email=amelia.walsh@gmail.com&username=amelia.walsh&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user4_id} amelia.walsh created"
user5_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Kevin+McLean&email=kevin.mclean@gmail.com&username=kevin.mclean&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user5_id} kevin.mclean created"
user6_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Kylie+Morrison&email=kylie.morrison@gmail.com&username=kylie.morrison&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user6_id} kylie.morrison created"
user7_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Rebecca+Gray&email=rebecca.gray@gmail.com&username=rebecca.gray&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user7_id} rebecca.gray created"
user8_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Simon+Turner&email=simon.turner@gmail.com&username=simon.turner&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user8_id} simon.turner created"
user9_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Olivia+White&email=olivia.white@gmail.com&username=olivia.white&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user9_id} olivia.white created"
user10_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Frank+Watson&email=frank.watson@gmail.com&username=frank.watson&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user10_id} frank.watson created"
user11_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Paul+Lyman&email=paul.lyman@gmail.com&username=paul.lyman&password=123qwe123&skip_confirmation=true" | jq '.id')
success "${user11_id} paul.lyman created"

echo ""
# create group
info "Creating groups"
pgroup1_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Group1&path=Group1" | jq '.id')
success "${pgroup1_id} Group1 created"
pgroup2_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Group2&path=Group2" | jq '.id')
success "${pgroup2_id} Group2 created"

echo ""
# add user to group
info "Adding users to group"
${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user1_id}&access_level=30"
success "add ${user1_id} to Group1"
${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user2_id}&access_level=40"
success "add ${user2_id} to Group1"
${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user3_id}&access_level=50"
success "add ${user3_id} to Group1"
${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user4_id}&access_level=30"
success "add ${user4_id} to Group2"
${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user5_id}&access_level=40"
success "add ${user5_id} to Group2"
${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user6_id}&access_level=50"
success "add ${user6_id} to Group2"

echo ""
# create subgroup
info "Creating subgroup"
sgroup1_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup1&path=SubGroup1&parent_id=${pgroup1_id}" | jq '.id')
success "${sgroup1_id} SubGroup1 created"
sgroup2_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup2&path=SubGroup2&parent_id=${pgroup1_id}" | jq '.id')
success "${sgroup2_id} SubGroup2 created"
sgroup3_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup3&path=SubGroup3&parent_id=${pgroup2_id}" | jq '.id')
success "${sgroup3_id} SubGroup3 created"
sgroup4_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup4&path=SubGroup4&parent_id=${pgroup2_id}" | jq '.id')
success "${sgroup4_id} SubGroup4 created"

echo ""
info sleeping for 5 seconds..
sleep 5

echo ""
# create group project
info "Creating a project in group/subgroup"
group_project1_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project1&namespace_id=${pgroup1_id}" | jq '.id')
success "${group_project1_id} Group1/Project1 created"
group_project2_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project2&namespace_id=${pgroup1_id}" | jq '.id')
success "${group_project2_id} Group1/Project2 created"
group_project3_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project3&namespace_id=${pgroup1_id}" | jq '.id')
success "${group_project3_id} Group1/Project3 created"
group_project4_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project4&namespace_id=${sgroup1_id}" | jq '.id')
success "${group_project4_id} Group1/SubGroup1/Project4 created"
group_project5_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project5&namespace_id=${sgroup1_id}" | jq '.id')
success "${group_project5_id} Group1/SubGroup1/Project5 created"
group_project6_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project6&namespace_id=${sgroup1_id}" | jq '.id')
success "${group_project6_id} Group1/SubGroup1/Project6 created"
group_project7_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project7&namespace_id=${sgroup2_id}" | jq '.id')
success "${group_project7_id} Group1/SubGroup2/Project7 created"
group_project8_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project8&namespace_id=${sgroup2_id}" | jq '.id')
success "${group_project8_id} Group1/SubGroup2/Project8 created"
group_project9_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project9&namespace_id=${sgroup2_id}" | jq '.id')
success "${group_project9_id} Group1/SubGroup2/Project9 created"
group_project10_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project10&namespace_id=${pgroup2_id}" | jq '.id')
success "${group_project10_id} Group2/Project10 created"
group_project11_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project11&namespace_id=${pgroup2_id}" | jq '.id')
success "${group_project11_id} Group2/Project11 created"
group_project12_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project12&namespace_id=${pgroup2_id}" | jq '.id')
success "${group_project12_id} Group2/Project12 created"
group_project13_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project13&namespace_id=${sgroup3_id}" | jq '.id')
success "${group_project13_id} Group2/SubGroup3/Project13 created"
group_project14_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project14&namespace_id=${sgroup3_id}" | jq '.id')
success "${group_project14_id} Group2/SubGroup3/Project14 created"
group_project15_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project15&namespace_id=${sgroup3_id}" | jq '.id')
success "${group_project15_id} Group2/SubGroup3/Project15 created"
group_project16_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project16&namespace_id=${sgroup4_id}" | jq '.id')
success "${group_project16_id} Group2/SubGroup4/Project16 created"
group_project17_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project17&namespace_id=${sgroup4_id}" | jq '.id')
success "${group_project17_id} Group2/SubGroup4/Project17 created"
group_project18_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project18&namespace_id=${sgroup4_id}" | jq '.id')
success "${group_project18_id} Group2/SubGroup4/Project18 created"

echo ""
# add user to project
info "Add user to project"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user1_id}&access_level=30"
success "add ${user1_id} to Group1/Project1"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user2_id}&access_level=40"
success "add ${user2_id} to Group1/Project1"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user3_id}&access_level=50"
success "add ${user3_id} to Group1/Project1"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user4_id}&access_level=30"
success "add ${user4_id} to Group1/Project1"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user5_id}&access_level=40"
success "add ${user5_id} to Group1/Project2"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user6_id}&access_level=50"
success "add ${user6_id} to Group1/Project2"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user7_id}&access_level=40"
success "add ${user7_id} to Group1/Project2"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user8_id}&access_level=40"
success "add ${user8_id} to Group1/Project2"

echo ""
# create user project
info "Creating users project"
project1_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects/user/${user7_id}" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project19" | jq '.id')
success "${user7_id} in Project19 created"
project2_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects/user/${user8_id}" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project20" | jq '.id')
success "${user8_id} in Project20 created"
project3_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects/user/${user9_id}" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project21" | jq '.id')
success "${user9_id} in Project21 created"
project4_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects/user/${user10_id}" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project22" | jq '.id')
success "${user10_id} in Project22 created"
project5_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects/user/${user11_id}" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project23" | jq '.id')
success "${user11_id} in Project23 created"

echo ""
# create hooks for projects
info "Creating hooks for projects"
${CCURL} "${GITLAB_URL}/api/v4/projects/${project1_id}/hooks" -d "access_token=${GITLAB_PRIVATE_TOKEN}&url=http%3A%2F%2Fwww.sample1.com%2F"
success "hooks in Project1 created"
${CCURL} "${GITLAB_URL}/api/v4/projects/${project2_id}/hooks" -d "access_token=${GITLAB_PRIVATE_TOKEN}&url=http%3A%2F%2Fwww.sample2.com%2F"
success "hooks in Project2 created"
${CCURL} "${GITLAB_URL}/api/v4/projects/${project3_id}/hooks" -d "access_token=${GITLAB_PRIVATE_TOKEN}&url=http%3A%2F%2Fwww.sample3.com%2F"
success "hooks in Project3 created"
${CCURL} "${GITLAB_URL}/api/v4/projects/${project4_id}/hooks" -d "access_token=${GITLAB_PRIVATE_TOKEN}&url=http%3A%2F%2Fwww.sample4.com%2F"
success "hooks in Project4 created"
${CCURL} "${GITLAB_URL}/api/v4/projects/${project5_id}/hooks" -d "access_token=${GITLAB_PRIVATE_TOKEN}&url=http%3A%2F%2Fwww.sample5.com%2F"
success "hooks in Project5 created"

echo ""
# push some commits in preparation for creation of git tags
info "Push commits for projects"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/repository/files/README.md" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=master&content=Test&commit_message=First+commit"
success "commit master branch to Group1/Project1"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/repository/files/README.md" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=master&content=Test&commit_message=First+commit"
success "commit master branch to Group1/Project2"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project3_id}/repository/files/README.md" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=master&content=Test&commit_message=First+commit"
success "commit master branch to Group1/Project3"

echo ""
# create tags
info "Creating tags for projects"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/repository/tags" -d "access_token=${GITLAB_PRIVATE_TOKEN}&tag_name=sample_1.0&ref=master"
success "Group1/Project1 sample_1.0 master branch tag created"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/repository/tags" -d "access_token=${GITLAB_PRIVATE_TOKEN}&tag_name=sample_1.0&ref=master"
success "Group1/Project2 sample_1.0 master branch tag created"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project3_id}/repository/tags" -d "access_token=${GITLAB_PRIVATE_TOKEN}&tag_name=sample_1.0&ref=master"
success "Group1/Project3 sample_1.0 master branch tag created"