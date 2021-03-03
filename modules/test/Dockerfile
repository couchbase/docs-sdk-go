ARG BASE_IMAGE=ubuntu:20.04

ARG CB_EDITION=enterprise
ARG CB_VERSION=7.0.0-beta
ARG CB_IMAGE=couchbase:$CB_EDITION-$CB_VERSION

ARG CB_CLIENT_OS=ubuntu2004
ARG CB_CLIENT_OS_TYPE=focal

# SDK related images...

FROM golang AS golang

FROM $CB_IMAGE

ARG CB_EDITION
ARG CB_VERSION
ARG CB_IMAGE

ARG CB_CLIENT_OS
ARG CB_CLIENT_OS_TYPE

ENV TZ=America/Los_Angeles
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && \
    apt-get install -y \
	curl wget jq unzip zip \
	build-essential cmake libssl-dev \
	atop htop psmisc strace time \
	vim npm

# ------------------------------------------------------

COPY --from=golang /usr/local/go /usr/local/go
COPY --from=golang /go /go

ENV GOPATH /go
ENV PATH="$PATH:/usr/local/go/bin:$GOPATH/bin"
RUN go version && \
    go mod init test && \
    go get github.com/couchbase/gocb/v2

# Install golang couchbase SDK:

RUN go get github.com/couchbase/gocb

# Copy init-couchbase files into image.

RUN mkdir -p /init-couchbase
COPY test/scripts/init-couchbase /init-couchbase
RUN chmod +x /init-couchbase/*.sh

# Append to /opt/couchbase/etc/couchbase/static_config...

RUN if [ ! -d /opt/couchbase/etc/couchbase ]; then mkdir -p /opt/couchbase/etc/couchbase; fi \
    && cat /init-couchbase/init-static-config.txt >> \
        /opt/couchbase/etc/couchbase/static_config
