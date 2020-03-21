FROM golang:alpine3.11 AS builder
WORKDIR /satelit-index/
COPY . .
RUN docker/scripts/build-project.sh

FROM alpine:latest
WORKDIR /satelit-index/
COPY --from=builder /satelit-index/satelit-index.tar.gz /satelit-index/
RUN ["tar", "-xvf", "satelit-index.tar.gz"]
ENTRYPOINT ["/satelit-index/entry.sh"]
