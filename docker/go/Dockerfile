FROM golang:1.16.5-buster

RUN apt update && apt install -y git 

ENTRYPOINT mkdir src && git clone ${GIT_URL} src && cd src && \
  go run main.go