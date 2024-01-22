#!/usr/bin/env bash
source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
check=$(ps aux | grep -w ./${enssearch_server_name} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  oldPid=$(ps aux | grep -w ./${enssearch_server_name} | grep -v grep | awk '{print $2}')
  kill -9 $oldPid
fi
sleep 1
cd ${enssearch_server_binary_root}
allPorts=10030
nohup ./${enssearch_server_name} -port ${allPorts} >> ../logs/openIM.log 2>> ../logs/error.log &
sleep 3
check=$(ps aux | grep -w ./${enssearch_server_name} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  newPid=$(ps aux | grep -w ./${enssearch_server_name} | grep -v grep | awk '{print $2}')
  ports=$(netstat -netulp | grep -w ${newPid} | awk '{print $4}' | awk -F '[:]' '{print $NF}')
  echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS "${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${enssearch_server_name}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${newPid}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allPorts}${COLOR_SUFFIX}
else
  echo -e ${YELLOW_PREFIX}${enssearch_server_name}${COLOR_SUFFIX}${RED_PREFIX}"SERVICE START ERROR, PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
