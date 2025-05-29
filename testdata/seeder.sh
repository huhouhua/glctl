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

GLCTL_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source ${GLCTL_ROOT}/scripts/lib/loggin.sh
source ${GLCTL_ROOT}/scripts/lib/tools.sh
install::jq

source ${GLCTL_ROOT}/testdata/credentials.sh

info "GITLAB_USERNAME:${GITLAB_USERNAME}"
info "GITLAB_PASSWORD:${GITLAB_PASSWORD}"
info "GITLAB_URL:${GITLAB_URL}"
info "GITLAB_PRIVATE_TOKEN:${GITLAB_PRIVATE_TOKEN}"

echo "=============================================="
echo "#              SEEDER SCRIPT                 #"
echo "=============================================="

set -e

echo ""
# create user
info "Creating users......"

user1_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=John+Doe&email=john.doe@gmail.com&username=john.doe&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user1_id}" "${user1_id} john.doe created"

user2_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=John+Smith&email=john.smith@gmail.com&username=john.smith&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user2_id}" "${user2_id} john.smith created"

user3_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Matt+Hunter&email=matt.hunter@gmail.com&username=matt.hunter&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user3_id}"  "${user3_id} matt.hunter created"

user4_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Amelia+Walsh&email=amelia.walsh@gmail.com&username=amelia.walsh&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user4_id}" "${user4_id} amelia.walsh created"

user5_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Kevin+McLean&email=kevin.mclean@gmail.com&username=kevin.mclean&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user5_id}" "${user5_id} kevin.mclean created"

user6_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Kylie+Morrison&email=kylie.morrison@gmail.com&username=kylie.morrison&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user6_id}" "${user6_id} kylie.morrison created"

user7_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Rebecca+Gray&email=rebecca.gray@gmail.com&username=rebecca.gray&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user7_id}" "${user7_id} rebecca.gray created"

user8_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Simon+Turner&email=simon.turner@gmail.com&username=simon.turner&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user8_id}" "${user8_id} simon.turner created"

user9_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Olivia+White&email=olivia.white@gmail.com&username=olivia.white&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user9_id}" "${user9_id} olivia.white created"

user10_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Frank+Watson&email=frank.watson@gmail.com&username=frank.watson&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user10_id}" "${user10_id} frank.watson created"

user11_id=$(${CCURL} "${GITLAB_URL}/api/v4/users" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Paul+Lyman&email=paul.lyman@gmail.com&username=paul.lyman&password=123qwe123&skip_confirmation=true" | jq '.id')
print::result "${user11_id}" "${user11_id} paul.lyman created"

echo ""
# create group
info "Creating groups"
pgroup1_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Group1&path=Group1" | jq '.id')
print::result "${pgroup1_id}" "${pgroup1_id} Group1 created"

pgroup2_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Group2&path=Group2" | jq '.id')
print::result "${pgroup2_id}" "${pgroup2_id} Group2 created"

pgroup3_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Group3&path=Group3" | jq '.id')
print::result "${pgroup3_id}" "${pgroup3_id} Group3 created"

pgroup4_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Group4&path=Group4" | jq '.id')
print::result "${pgroup4_id}" "${pgroup4_id} Group4 created"

pgroup5_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Group5&path=Group5" | jq '.id')
print::result "${pgroup5_id}" "${pgroup5_id} Group5 created"

echo ""
# add user to group
info "Adding users to group"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user1_id}&access_level=30"; echo
success "add ${user1_id} to Group1"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user2_id}&access_level=40";
success "add ${user2_id} to Group1"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user3_id}&access_level=50"; echo
success "add ${user3_id} to Group1"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user4_id}&access_level=30"; echo
success "add ${user4_id} to Group2"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user5_id}&access_level=40"; echo
success "add ${user5_id} to Group2"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user6_id}&access_level=50"; echo
success "add ${user6_id} to Group2"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup3_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user7_id}&access_level=30"; echo
success "add ${user7_id} to Group3"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup3_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user8_id}&access_level=40"; echo
success "add ${user8_id} to Group3"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup3_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user9_id}&access_level=50"; echo
success "add ${user9_id} to Group3"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup4_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user10_id}&access_level=30"; echo
success "add ${user10_id} to Group4"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup4_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user11_id}&access_level=40"; echo
success "add ${user11_id} to Group4"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup4_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user1_id}&access_level=50"; echo
success "add ${user1_id} to Group4"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup5_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user2_id}&access_level=30"; echo
success "add ${user2_id} to Group5"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup5_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user3_id}&access_level=40"; echo
success "add ${user3_id} to Group5"

${CCURL} "${GITLAB_URL}/api/v4/groups/${pgroup5_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user4_id}&access_level=50"; echo
success "add ${user4_id} to Group5"

echo ""
# create subgroup
info "Creating subgroup"
sgroup1_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup1&path=SubGroup1&parent_id=${pgroup1_id}" | jq '.id')
print::result "${sgroup1_id}" "${sgroup1_id} SubGroup1 created"

sgroup2_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup2&path=SubGroup2&parent_id=${pgroup1_id}" | jq '.id')
print::result "${sgroup2_id}" "${sgroup2_id} SubGroup2 created"

sgroup3_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup3&path=SubGroup3&parent_id=${pgroup2_id}" | jq '.id')
print::result "${sgroup3_id}" "${sgroup3_id} SubGroup3 created"

sgroup4_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup4&path=SubGroup4&parent_id=${pgroup2_id}" | jq '.id')
print::result "${sgroup4_id}" "${sgroup4_id} SubGroup4 created"

sgroup5_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup5&path=SubGroup5&parent_id=${pgroup3_id}" | jq '.id')
print::result "${sgroup5_id}" "${sgroup5_id} SubGroup5 created"

sgroup6_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup6&path=SubGroup6&parent_id=${pgroup3_id}" | jq '.id')
print::result "${sgroup6_id}" "${sgroup6_id} SubGroup6 created"

sgroup7_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup7&path=SubGroup7&parent_id=${pgroup4_id}" | jq '.id')
print::result "${sgroup7_id}" "${sgroup7_id} SubGroup7 created"

sgroup8_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup8&path=SubGroup8&parent_id=${pgroup4_id}" | jq '.id')
print::result "${sgroup8_id}" "${sgroup8_id} SubGroup8 created"

sgroup9_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup9&path=SubGroup9&parent_id=${pgroup5_id}" | jq '.id')
print::result "${sgroup9_id}" "${sgroup9_id} SubGroup9 created"

sgroup10_id=$(${CCURL} "${GITLAB_URL}/api/v4/groups" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=SubGroup10&path=SubGroup10&parent_id=${pgroup5_id}" | jq '.id')
print::result "${sgroup10_id}" "${sgroup10_id} SubGroup10 created"

echo ""
info "sleeping for 5 seconds.."
sleep 5

echo ""
# create group project
info "Creating a project in group/subgroup"
group_project1_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project1&namespace_id=${pgroup1_id}" | jq '.id')
print::result "${group_project1_id}" "Group1/Project1 created"

group_project2_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project2&namespace_id=${pgroup1_id}" | jq '.id')
print::result "${group_project2_id}" "Group1/Project2 created"

group_project3_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project3&namespace_id=${pgroup1_id}" | jq '.id')
print::result "${group_project3_id}" "Group1/Project3 created"

group_project4_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project4&namespace_id=${sgroup1_id}" | jq '.id')
print::result "${group_project4_id}" "Group1/SubGroup1/Project4 created"

group_project5_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project5&namespace_id=${sgroup1_id}" | jq '.id')
print::result "${group_project5_id}" "Group1/SubGroup1/Project5 created"

group_project6_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project6&namespace_id=${sgroup1_id}" | jq '.id')
print::result "${group_project6_id}" "Group1/SubGroup1/Project6 created"

group_project7_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project7&namespace_id=${sgroup2_id}" | jq '.id')
print::result "${group_project7_id}" "Group1/SubGroup2/Project7 created"

group_project8_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project8&namespace_id=${sgroup2_id}" | jq '.id')
print::result "${group_project8_id}" "Group1/SubGroup2/Project8 created"

group_project9_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project9&namespace_id=${sgroup2_id}" | jq '.id')
print::result "${group_project9_id}" "Group1/SubGroup2/Project9 created"

group_project10_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project10&namespace_id=${pgroup2_id}" | jq '.id')
print::result "${group_project10_id}" "Group2/Project10 created"

group_project11_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project11&namespace_id=${pgroup2_id}" | jq '.id')
print::result "${group_project11_id}" "Group2/Project11 created"

group_project12_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project12&namespace_id=${pgroup2_id}" | jq '.id')
print::result "${group_project12_id}" "Group2/Project12 created"

group_project13_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project13&namespace_id=${sgroup3_id}" | jq '.id')
print::result "${group_project13_id}" "Group2/SubGroup3/Project13 created"

group_project14_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project14&namespace_id=${sgroup3_id}" | jq '.id')
print::result "${group_project14_id}" "Group2/SubGroup3/Project14 created"

group_project15_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project15&namespace_id=${sgroup3_id}" | jq '.id')
print::result "${group_project15_id}" "Group2/SubGroup3/Project15 created"

group_project16_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project16&namespace_id=${sgroup4_id}" | jq '.id')
print::result "${group_project16_id}" "Group2/SubGroup4/Project16 created"

group_project17_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project17&namespace_id=${sgroup4_id}" | jq '.id')
print::result "${group_project17_id}" "Group2/SubGroup4/Project17 created"

group_project18_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project18&namespace_id=${sgroup4_id}" | jq '.id')
print::result "${group_project18_id}" "Group2/SubGroup4/Project18 created"

echo ""
# add test branch to project
info "Add test branch to project"
project1_branch_unprotect=$(${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/repository/branches" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=unprotect&ref=main" | jq '.name')
print::result "${project1_branch_unprotect}" "Group1/Project1 ${project1_branch_unprotect} created"

project1_branch_protect=$(${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/repository/branches" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=protect&ref=main" | jq '.name')
print::result "${project1_branch_protect}" "Group1/Project1 ${project1_branch_protect} created"

project2_branch_unprotect=$(${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/repository/branches" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=unprotect&ref=main" | jq '.name')
print::result "${project2_branch_unprotect}" "Group1/Project2 ${project2_branch_unprotect} created"

project2_branch_protect=$(${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/repository/branches" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=protect&ref=main" | jq '.name')
print::result "${project2_branch_protect}" "Group1/Project2 ${project2_branch_protect} created"

project3_branch_unprotect=$(${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project3_id}/repository/branches" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=unprotect&ref=main" | jq '.name')
print::result "${project3_branch_unprotect}" "Group1/Project3 ${project3_branch_unprotect} created"

project3_branch_protect=$(${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project3_id}/repository/branches" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=protect&ref=main" | jq '.name')
print::result "${project3_branch_protect}" "Group1/Project3 ${project3_branch_protect} created"


echo ""
# add user to project
info "Add user to project"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user1_id}&access_level=30"; echo
success "add ${user1_id} to Group1/Project1"

${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user2_id}&access_level=40"; echo
success "add ${user2_id} to Group1/Project1"

${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user4_id}&access_level=30"; echo
success "add ${user4_id} to Group1/Project1"

${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user5_id}&access_level=40"; echo
success "add ${user5_id} to Group1/Project2"

${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user7_id}&access_level=40"; echo
success "add ${user7_id} to Group1/Project2"

${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/members" -d "access_token=${GITLAB_PRIVATE_TOKEN}&user_id=${user8_id}&access_level=40"; echo
success "add ${user8_id} to Group1/Project2"

echo ""
# create user project
info "Creating users project"
project1_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects/user/${user7_id}" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project19" | jq '.id')
print::result "${project1_id}" "${user7_id} in Project19 created"

project2_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects/user/${user8_id}" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project20" | jq '.id')
print::result "${project2_id}" "${user8_id} in Project20 created"

project3_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects/user/${user9_id}" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project21" | jq '.id')
print::result "${project3_id}" "${user9_id} in Project21 created"

project4_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects/user/${user10_id}" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project22" | jq '.id')
print::result "${project4_id}" "${user10_id} in Project22 created"

project5_id=$(${CCURL} "${GITLAB_URL}/api/v4/projects/user/${user11_id}" -d "access_token=${GITLAB_PRIVATE_TOKEN}&name=Project23" | jq '.id')
print::result "${project5_id}" "${user11_id} in Project23 created"

echo ""
# create hooks for projects
info "Creating hooks for projects"
${CCURL} "${GITLAB_URL}/api/v4/projects/${project1_id}/hooks" -d "access_token=${GITLAB_PRIVATE_TOKEN}&url=http%3A%2F%2Fwww.sample1.com%2F"; echo
success "hooks in Project1 created"

${CCURL} "${GITLAB_URL}/api/v4/projects/${project2_id}/hooks" -d "access_token=${GITLAB_PRIVATE_TOKEN}&url=http%3A%2F%2Fwww.sample2.com%2F"; echo
success "hooks in Project2 created"

${CCURL} "${GITLAB_URL}/api/v4/projects/${project3_id}/hooks" -d "access_token=${GITLAB_PRIVATE_TOKEN}&url=http%3A%2F%2Fwww.sample3.com%2F"; echo
success "hooks in Project3 created"

${CCURL} "${GITLAB_URL}/api/v4/projects/${project4_id}/hooks" -d "access_token=${GITLAB_PRIVATE_TOKEN}&url=http%3A%2F%2Fwww.sample4.com%2F"; echo
success "hooks in Project4 created"

${CCURL} "${GITLAB_URL}/api/v4/projects/${project5_id}/hooks" -d "access_token=${GITLAB_PRIVATE_TOKEN}&url=http%3A%2F%2Fwww.sample5.com%2F"; echo
success "hooks in Project5 created"

echo ""
# push some commits in preparation for creation of git tags
info "Push commits for projects"
${UCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/repository/files/README.md" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=Test&commit_message=First+commit"; echo
success "commit main branch to Group1/Project1"

${UCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/repository/files/README.md" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=Test&commit_message=First+commit"; echo
success "commit main branch to Group1/Project2"

${UCURL} "${GITLAB_URL}/api/v4/projects/${group_project3_id}/repository/files/README.md" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=Test&commit_message=First+commit"; echo
success "commit main branch to Group1/Project3"


${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project13_id}/repository/files/README.md" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=Test&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project13_id}/repository/files/test%2Ftest.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 1&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project13_id}/repository/files/test%2Ftest2.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 2&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project13_id}/repository/files/test%2Ftest3.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 3&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project13_id}/repository/files/test%2Ftest4.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 4&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project13_id}/repository/files/test%2Ftest5.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 5&commit_message=First+commit"; echo
success "commit main branch to Group2/SubGroup3/Project13"

${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project14_id}/repository/files/README.md" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=Test&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project14_id}/repository/files/test%2Ftest.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 1&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project14_id}/repository/files/test%2Ftest2.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 2&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project14_id}/repository/files/test%2Ftest3.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 3&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project14_id}/repository/files/test%2Ftest4.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 4&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project14_id}/repository/files/test%2Ftest5.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 5&commit_message=First+commit"; echo
success "commit main branch to Group2/SubGroup3/Project14"

${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project15_id}/repository/files/README.md" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=Test&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project15_id}/repository/files/test%2Ftest.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 1&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project15_id}/repository/files/test%2Ftest2.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 2&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project15_id}/repository/files/test%2Ftest3.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 3&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project15_id}/repository/files/test%2Ftest4.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 4&commit_message=First+commit"; echo
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project15_id}/repository/files/test%2Ftest5.yaml" -d "access_token=${GITLAB_PRIVATE_TOKEN}&branch=main&content=test: 5&commit_message=First+commit"; echo
success "commit main branch to Group2/SubGroup3/Project15"

echo ""
# create tags
info "Creating tags for projects"
${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project1_id}/repository/tags" -d "access_token=${GITLAB_PRIVATE_TOKEN}&tag_name=sample_1.0&ref=main"; echo
success "Group1/Project1 sample_1.0 main branch tag created"

${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project2_id}/repository/tags" -d "access_token=${GITLAB_PRIVATE_TOKEN}&tag_name=sample_1.0&ref=main"; echo
success "Group1/Project2 sample_1.0 main branch tag created"

${CCURL} "${GITLAB_URL}/api/v4/projects/${group_project3_id}/repository/tags" -d "access_token=${GITLAB_PRIVATE_TOKEN}&tag_name=sample_1.0&ref=main"; echo
success "Group1/Project3 sample_1.0 main branch tag created"; echo

success "Add test data successful"