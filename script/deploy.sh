#!/usr/bin/env bash

source ../.env
# # step1: check the environment
deployment=${DEPLOYMENT}
branch=dev

if [ $deployment == "dev" ]; then
    branch=dev
elif [ $deployment == "stag" ]; then
    branch=dev
elif [ $deployment == "prod" ]; then
    branch=main
else
    echo "deployment is invalid"
    exit 1
fi

# step2: pull the latest code
git fetch origin $branch
git checkout $branch
git reset --hard origin/$branch

# printf "Deployment: %s\n" $deployment

# step3: build the docker image

./build_all_service.sh

find . -type f -exec dos2unix {} \;;
# mv ../config/config.${DEPLOYMENT}.yaml ../config/config.yaml
cd ../
docker-compose build open_im_server
docker-compose stop open_im_server
docker-compose up -d open_im_server

printf "deploy success\n"
