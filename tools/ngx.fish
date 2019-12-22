#!/usr/bin/env fish

set REPO_DIR (git rev-parse --show-toplevel)
set SERVE_PATH $REPO_DIR/static
set CONTAINER_NAME ngx-satelit-index
set SERVE_PORT 8081

source  $REPO_DIR/tools/docker.fish

function print_usage
  echo "Manipulates Nginx Docker container.

Usage:
  ngx COMMAND

  COMMAND:
    start   starts new Nginx container and serves files from 'static' dir.
    stop    stops running Nginx container."
end

function start_ngx
  if test ! -d $SERVE_PATH
    mkdir -p $SERVE_PATH
  end

  docker run --rm -d \
    -p $SERVE_PORT:80 \
    --mount type=bind,source=$SERVE_PATH,target=/static \
    --name $CONTAINER_NAME \
    flashspys/nginx-static; or exit $status
end

function stop_ngx
  docker stop $CONTAINER_NAME
end

function main
  set subcmd $argv[1]
  switch $subcmd
    case start
      start_ngx
    case stop
      stop_ngx
    case -h --help
      print_usage
    case '*'
      echo "Unknown command $subcmd"
      echo "Try '--help' for help"
      exit 1
  end
end

main $argv
