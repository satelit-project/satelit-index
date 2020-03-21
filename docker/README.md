# Docker

This directory contains necessary files to build new Docker image for
running it in production.

## Building

To build new version of the container run the following command:

``` sh
VERSION="<version>"
docker build -t satelit/satelit-index:"$VERSION" -f docker/Dockerfile .
docker push satelit/satelit-index:"$VERSION"
```
