# Dockerfile used to build a build step that builds container-diff in CI.
FROM golang:1.9
RUN apt-get update && apt-get install make
RUN mkdir -p /go/src/github.com/GoogleContainerTools/
RUN ln -s /workspace /go/src/github.com/GoogleContainerTools/container-diff
WORKDIR /go/src/github.com/GoogleContainerTools/container-diff
