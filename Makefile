BINARY_NAME=helm-wrapper

GOPATH = $(shell go env GOPATH)

LDFLAGS="-s -w"

build:
	go build -ldflags ${LDFLAGS} -o ${BINARY_NAME} 

# cross compilation
build-linux:
# 构建程序 -s -w缩小构建文件大小
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.buildTime=`date +%Y-%m-%d,%H:%M:%S` -X main.buildVersion=v1.0 -X main.gitCommitID=`git rev-parse HEAD`" -o ${BINARY_NAME}

# build docker image
build-docker:
# docker镜像采用scratch最小化镜像构建,所以采取静态编译方法构建程序,不依赖任何动态链接库
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o ${BINARY_NAME}
	docker build -t helm-wrapper:v1 .

.PHONY: golangci-lint
golangci-lint: $(GOLANGCILINT)
	@echo
	$(GOPATH)/bin/golangci-lint run

$(GOLANGCILINT):
	(cd /; GO111MODULE=on GOPROXY="direct" GOSUMDB=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0)