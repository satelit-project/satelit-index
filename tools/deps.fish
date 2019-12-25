#!/usr/bin/env fish

# Prints help message
function print_usage
  echo "Install or update project dependencies

Usage: deps.fish <COMMAND>

COMMAND:
  install installs project dependencies
  update  updates project dependencies" >&2
end

function fatal -a msg -d "Prints error message and exits"
  set_color red
  printf "$msg\n" >&2
  exit 1
end

# Checks if current system has all required dependencies
function check_system
  if ! type -q go
    fatal "Go is not installed"
  end

  if ! type -q tar
    fatal "Tar is not installed"
  end
  
  if ! type -q curl
    fatal "Curl is not installed"
  end

  set gobin (go env GOBIN)
  if test -z "$gobin"
    fatal "GOBIN is not set"
  end
end

# Installs project dependencies
# Args:
#   $forced - wherever to re-install/update dependencies
function install_deps -a forced
  set -x GO111MODULE on

  pushd (go env GOBIN)
  install_gopls $forced
  install_goose $forced
  install_delve $forced
  popd

  install_sqlc $forced
end

# Installs gopls LSP
# Args:
#   $forced - wherever to update gopls if already installed
function install_gopls -a forced
  if ! type -q gopls || test -n "$forced"
    echo "Installing 'gopls'" >&2
    go get golang.org/x/tools/gopls@latest
  end
end

# Installs goose SQL migration tool
# Args:
#   $forced - wherever to update goose if already installed
function install_goose -a forced
  if ! type -q goose || test -n "$forced"
    echo "Installing 'goose'" >&2
    go get -u github.com/pressly/goose/cmd/goose
  end
end

# Installs delve debugger
# Args:
#   $forced - wherever to update delve if already installed
function install_delve -a forced
  if ! type -q dlv || test -n "$forced"
    echo "Installing 'delve'" >&2
    go get -u github.com/go-delve/delve/cmd/dlv
  end
end

# Installs sqlc SQL/Go generation tool
# Args:
#   $forced - wherever to update sqlc if already installed
function install_sqlc -a forced
  if ! type -q sqlc || test -n "$forced"
    set tmpdir (mktemp -d)
    trap "rm -rf $tmpdir" EXIT
    
    echo "Installing 'sqlc'" >&2
    set system (uname -s | string lower)
    set arch (uname -m)

    if test $system = linux -a $arch = x86_64
      download_sqlc_linux $tmpdir
    else if test $system = darwin -a $arch = x86_64
      download_sqlc_darwin $tmpdir
    else
      fatal "'sqlc' binary is not available for your platform"
    end

    mv $tmpdir/sqlc $GOBIN/sqlc
  end
end

# Downloads sqlc for linux platform
# Args:
#   $dir - path to download directory
function download_sqlc_linux -a dir
  set path $dir/sqlc.tgz
  set url "https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-linux-amd64.tgz"

  curl -o $path -OL $url

  pushd $dir
  tar zxf $path
  popd
end

# Downloads sqlc for macOS platform
# Args:
#   $dir - path to download directory
function download_sqlc_darwin -a dir
  set path $dir/sqlc.zip
  set url "https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-darwin-amd64.zip"

  curl -o $path -OL $url

  pushd $dir
  unzip $path >/dev/null
  popd
end

function main
  check_system
  set cmd $argv[1]

  switch "$cmd"
    case install
      install_deps
    case update
      install_deps "forced"
    case --help -h
      print_usage
    case "*"
      fatal "Unknown command $cmd\nTry '--help' for help"
  end
end

main $argv
