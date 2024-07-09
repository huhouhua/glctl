#!/bin/bash

set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
export TERM=xterm
source /root/.bashrc 
source $DIR/../common.sh
set +o noglob

REGISTRY_PREFIX=registry.cn-beijing.aliyuncs.com
REGISTRY_USERNAME="huhouhua@1369200651364267"
REGISTRY_PASSWORD="huhouhua123456"
DOCKER_COMPOSE_FILE=$DIR/docker-compose-dev.yaml
IS_DOWN="false"
IMAGE_TAG=dev
TEMPLATE_TAG=dev
IMAGE=""
IMAGES=(registry.cn-beijing.aliyuncs.com\\/ruijie-ruizhi\\/riil-bridge-frontend)

usage=$'提供镜像仓库,镜像版本,dockerfile地址 ! 
如果需要设置部署的docker-compose.yaml,请设置--file=docker-compose.yaml文件地址,默认为:'$DOCKER_COMPOSE_FILE'
如果需要部署指定镜像版本号，请设置--image-tag=xxx,默认为:'$IMAGE_TAG'
如果需要关闭,请设置--down,默认为:'不关闭'
查看帮助信息，请设置--help'
item=0
while [ $# -gt 0 ]; do
  case $1 in
  --help)
    note "$usage"
    exit 0
    ;;
  --file=*)
    DOCKER_COMPOSE_FILE=$DIR/${1#*--file=}
    ;;
  --image-tag=*)
    IMAGE_TAG=${1#*--image-tag=}
    ;;
  --image=*)
    IMAGE=${1#*--image=}
    ;;
  --env=*)
    ENV=${1#*--env=}
    ;;
  --down)
    IS_DOWN="true"
    ;;
  *)
    note "$usage"
    exit 1
    ;;
  esac
  shift || true
done

h2 "[Step $item]: checking if docker is installed ..."; let item+=1
check_docker

h2 "[Step $item]: checking docker-compose is installed ..."; let item+=1
check_dockercompose

h2 "[Step $item]: try create docker-compose.yaml ..."
let item+=1
if [ -n "${IMAGE}" ] && [ "${IS_DOWN}"!="true" ]; then
 NEW_FILE=${DIR}/docker-compose-dev.yaml
 tag=${IMAGE##*:}
if [ -z "${tag}" ] || [ "${IMAGE}" == "${tag}" ]; then
  IMAGE="${IMAGE%:}:${IMAGE_TAG}"
fi
 for image in ${IMAGES[@]}
 do
 # todo 
 # currently, only one mirror can be replaced. You need to replace it with multiple mirrors
  sed -i "s#${image}:${TEMPLATE_TAG}#${IMAGE}#g" $NEW_FILE
 done
else
 NEW_FILE=${DIR}/docker-compose-${IMAGE_TAG}.yaml
if [ "${IMAGE_TAG}" != "dev" ] && [ "${IS_DOWN}"!="true" ]; then
   cp $DOCKER_COMPOSE_FILE $NEW_FILE
 for image in ${IMAGES[@]} 
 do
 sed -i "s/${image}:${TEMPLATE_TAG}/${image}:${IMAGE_TAG}/g" $NEW_FILE
 done
 fi
fi
echo ""

h2 "[Step $item]: checking docker-compose.yaml ..."
let item+=1
DOCKER_COMPOSE_FILE=$NEW_FILE
if [ ! -f $DOCKER_COMPOSE_FILE ]; then
    error "${DOCKER_COMPOSE_FILE} docker-compose file not exist!"
    exit 1
fi
echo ""

if [ "$IS_DOWN" == "true" ]; then
  docker-compose -f $DOCKER_COMPOSE_FILE down
  success $"---- riil-bridge-frontend down successfully.----"
  exit 0
fi

h2 "[Step $item]: try create network ..."
let item+=1
DOCKER_COMPOSE_FILE_LAST=$(echo $DOCKER_COMPOSE_FILE | awk -F '[/]' '{print $NF}')
if [ $DOCKER_COMPOSE_FILE_LAST == "docker-compose-dev.yaml" ]; then
ret=$(docker network ls | awk '$2=="cloud_server" && $3=="bridge" {print "exist"}')
if [ -z "${ret}" ]; then
    docker network create -d bridge cloud_server
    success "create bridge is cloud_server network successfully.----"
else
  info "cloud_server network already existed skipping ...."
fi
fi
echo ""

if [ -n "docker-compose -f $DOCKER_COMPOSE_FILE ps -q"  ]; then
        note "stopping existing riil-bridge-frontend instance ..." 
        docker-compose -f $DOCKER_COMPOSE_FILE down -v --remove-orphans
fi
echo ""

h2 "[Step $item]: starting run riil-bridge-frontend ..."
let item+=1
info "login ${REGISTRY_PREFIX} ...."
docker login ${REGISTRY_PREFIX} --username=${REGISTRY_USERNAME}  --password=${REGISTRY_PASSWORD}

docker-compose --project-name riil-bridge-frontend -f $DOCKER_COMPOSE_FILE up -d --remove-orphans --pull always || docker-compose --project-name riil-bridge-frontend -f $DOCKER_COMPOSE_FILE up -d --remove-orphans

echo ""

h2 "[Step $item]: service view"
let item+=1
docker-compose -f $DOCKER_COMPOSE_FILE ps
echo ""

success $"---- riil-bridge-frontend started successfully.----"
