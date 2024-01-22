.PHONY: all build run gotool install clean help

BINARY_NAME=init_data
BIN_DIR=../../bin/


all: gotool build

build-task-user-head:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o task_user_head -buildvcs=false && mv task_user_head ./bin

build-task-follow-twitter:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o task_user_follow -buildvcs=false && mv task_user_follow ./bin

build-task-member:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o task_group_member_join -buildvcs=false && mv task_group_member_join ./bin

build-task-create-group:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o task_group_create -buildvcs=false && mv task_group_create ./bin

build-create-task:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o create_task -buildvcs=false && mv create_task ./bin
build-group:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o init_group -buildvcs=false && mv init_group ./bin


run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

install:
	make build
	mv ${BINARY_NAME} ${BIN_DIR}

clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi

