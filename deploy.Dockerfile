# FROM golang as build

# # go mod Installation source, container environment variable addition will override the default variable value
# ENV GO111MODULE=on
# ENV GOPROXY=https://goproxy.cn,direct

# # Set up the working directory
# WORKDIR /Open-IM-Server

# COPY go.mod go.sum ./

# RUN go mod download
# RUN cd sdk && go mod download
# # add all files to the container
# WORKDIR /Open-IM-Server
# COPY . .

# WORKDIR /Open-IM-Server/script

# RUN chmod +x *.sh

# RUN ./build_all_service.sh

#Blank image Multi-Stage Build
FROM ubuntu

WORKDIR /Open-IM-Server
COPY . /Open-IM-Server/

RUN rm -rf /var/lib/apt/lists/*
RUN apt-get update && apt-get install apt-transport-https && apt-get install procps\
&&apt-get install net-tools
#Non-interactive operation
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get install -y vim curl tzdata gawk dos2unix
#Time zone adjusted to East eighth District
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure -f noninteractive tzdata

#set directory to map logs,config file,script file.
VOLUME ["/Open-IM-Server/logs","/Open-IM-Server/config", "/Open-IM-Server/db/sdk"]

WORKDIR /Open-IM-Server/script

RUN find . -type f -exec dos2unix {} \;;

CMD ["./docker_start_all.sh"]

# CMD ["tail","-f","/dev/null"]
