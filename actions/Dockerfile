FROM debian:bookworm

# docker build -f actions/Dockerfile -t googlecontainertools/container-diff .

LABEL "com.github.actions.name"="container-diff GitHub Action"
LABEL "com.github.actions.description"="use Container-Diff in Github Actions Workflows"
LABEL "com.github.actions.icon"="cloud"
LABEL "com.github.actions.color"="blue"

LABEL "repository"="https://www.github.com/GoogleContainerTools/container-diff"
LABEL "homepage"="https://www.github.com/GoogleContainerTools/container-diff"
LABEL "maintainer"="Google Inc."

# Install container-diff latest release
RUN apt-get update && apt-get install -y curl && \
    curl -LO https://storage.googleapis.com/container-diff/latest/container-diff-linux-amd64 && \
    install container-diff-linux-amd64 /usr/local/bin/container-diff

ADD entrypoint.sh /entrypoint.sh

RUN mkdir -p /root/.docker && \
    echo {} > /root/.docker/config.json && \
    chmod u+x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
