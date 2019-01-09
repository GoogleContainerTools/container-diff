FROM golang:1.11.3-stretch

# docker build -f actions/Dockerfile -t googlecontainertools/container-diff .

RUN apt-get update && \
    apt-get install -y automake \
                       libffi-dev \ 
                       libxml2 \
                       libxml2-dev \
                       libxslt-dev \
                       libxslt1-dev \
                       git \
                       gcc g++ \
                       wget \
                       locales

RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    locale-gen
ENV LANG en_US.UTF-8  
ENV LANGUAGE en_US:en  
ENV LC_ALL en_US.UTF-8

LABEL "com.github.actions.name"="container-diff GitHub Action"
LABEL "com.github.actions.description"="use Container-Diff in Github Actions Workflows"
LABEL "com.github.actions.icon"="cloud"
LABEL "com.github.actions.color"="blue"

LABEL "repository"="https://www.github.com/GoogleContainerTools/container-diff"
LABEL "homepage"="https://www.github.com/GoogleContainerTools/container-diff"
LABEL "maintainer"="Google Inc."

# Install container-diff from master
RUN go get github.com/GoogleContainerTools/container-diff && \
    cd ${GOPATH}/src/github.com/GoogleContainerTools/container-diff && \
    go get && \
    make && \
    go install && \
    mkdir -p /code && \
    apt-get autoremove

ADD entrypoint.sh /entrypoint.sh

RUN mkdir -p /root/.docker && \
    echo {} > /root/.docker/config.json && \
    chmod u+x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
