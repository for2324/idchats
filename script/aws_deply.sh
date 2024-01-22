#!/usr/bin/env bash

dateDir=date +"%Y-%m-%d"

echo $dateDir

cd /home/project/idchats_services_v2/script

git checkout main

git pull origin main

./build_all_service.sh

ssh aws mkdir /home/ubuntu/updatelog/$dateDir

tar -zcvf bin bin.tar

scp -r ./bin.tar aws:/home/ubuntu/updatelog/$dateDir/

ssh aws tar -zxvf /home/ubuntu/updatelog/$dateDir/bin.tar

ssh aws docker cp /home/ubuntu/updatelog/$dateDir/bin open_im_server:/Open-IM-Server/

ssh aws docker restart open_im_server

rm bin.tar
