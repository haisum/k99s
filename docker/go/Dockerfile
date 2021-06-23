FROM golang:1.16.5-buster

RUN apt update && apt install -y git 

ENTRYPOINT mkdir tmp && git clone ${GIT_URL} tmp && \
  cp -r tmp/${SUBPATH}/* src/ && rm -rf tmp && cd src && \
  go run main.go