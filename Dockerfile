FROM golang:bullseye

COPY . /go-docs-repo
WORKDIR /go-docs-repo/tests

RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    libssl-dev \
    jq \
    curl \
    npm

RUN npm install -g bats

ENTRYPOINT [ "./wait-for-couchbase.sh", "1" ]