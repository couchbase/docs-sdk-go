FROM golang:bullseye

WORKDIR /app/tests

RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    libssl-dev \
    jq curl \
    npm

RUN npm install -g bats

WORKDIR /app/tests

ENTRYPOINT [ "./wait-for-couchbase.sh" ]