#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
shimroot="$HOME/.goenv"
bindir="$shimroot/bin"
shimdir="$shimroot/shims"

suggest_export_to_profile() {
  shell="$(basename $SHELL)"

  case "$shell" in
  bash )
    if [ -f "${HOME}/.bashrc" ] && [ ! -f "${HOME}/.bash_profile" ]; then
      profile="$HOME/.bashrc"
    else
      profile="$HOME/.bash_profile"
    fi
    ;;
  zsh )
    profile="$HOME/.zshrc"
    ;;
  ksh )
    profile="$HOME/.profile"
    ;;
  fish )
    profile="$HOME/.config/fish/config.fish"
    ;;
  * )
    profile=''
    ;;
  esac

  echo "To use goenv, modify your profile file (probably $profile) to include"
  echo '"~/.goenv/bin" in the PATH. Important note: for security, ensure that'
  echo "both goenv\'s bin directory, and GOPATH are at the end of your PATH,"
  echo 'not the beginning. For goenv to work, it will need to be earlier in'
  echo 'your PATH than GOPATH. Your regular Go installation, and goenv should'
  echo 'look something like this in your profile:'
  echo ''
  echo '  export GOPATH=/usr/local/go/bin'
  echo '  export PATH=$PATH:$HOME/.goenv/bin:$GOPATH'
  echo ''
  echo 'or:'
  echo ''
  echo '  export GOPATH=/usr/local/go/bin'
  echo '  export PATH=$PATH:$HOME/.goenv/bin'
  echo '  export PATH=$PATH:$GOPATH'
  echo ''
  echo 'Once you have modified your profile, either restart your terminal, or,'
  echo "presuming you are using $profile, do:"
  echo ''
  echo "  source $profile"
}

copy_package() {
  if [ -e $bindir ]; then
    rm -rf $bindir
  fi
  mkdir -p $shimdir

  system_go="$(which go)"
  if [ -z $system_go ]; then
    echo "No system installation found"
  else
    echo "Found system installation at $system_go -- saving as fallback"
    echo "$(dirname $system_go)" > $shimroot/system_go
  fi

  echo "Installing goenv into $shimroot"
  cp -R $DIR/bin $shimroot

}

copy_package
if [[ $PATH =~ $bindir ]]; then
  echo "goenv already installed in PATH"
else
  suggest_export_to_profile
fi
