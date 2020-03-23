FROM golang:alpine3.11 AS builder
WORKDIR /satelit-index/
COPY . .
RUN docker/scripts/build-project.sh

FROM alpine:latest
WORKDIR /satelit-index/
EXPOSE 8080
COPY --from=builder /satelit-index/satelit-index.tar.gz /satelit-index/
RUN ["tar", "-xvf", "satelit-index.tar.gz"]
RUN ["apk", "add", "curl"]
ENTRYPOINT ["/satelit-index/entry.sh"]
HEALTHCHECK --interval=1m --start-period=10s CMD ["curl", "-s", "-I", "http://127.0.0.1:8080/anidb/"]
