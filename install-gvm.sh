#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
shimroot="$HOME/.gvm"
bindir="$shimroot/bin"
shimdir="$shimroot/shims"
system_go="$(which go)"

add_export_to_profile() {
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
  echo "Adding gvm to PATH in $profile"

  if [ -z $profile ]; then
    echo "To use gvm, add the following to your profile file ($profile)"
    echo '  export PATH=~/.gvm/bin:$PATH'
    echo ''
    echo 'Then, either restart your terminal or enter:'
    echo "  source $profile"
  else
    echo 'export PATH=~/.gvm/bin:$PATH' >> $profile
    echo 'gvm has been installed--either restart your terminal or enter:'
    echo "  source $profile"
  fi
}

copy_package() {
  if [ -e $bindir ]; then
    rm -rf $bindir
  fi
  mkdir -p $shimdir
  echo "Installing gvm into $shimroot"
  cp -R $DIR/bin $shimroot
  if [ -z $system_go ]; then
    echo "Found system installation at $system_go -- saving as fallback"
    echo "$system_go\n" > $shimroot/system_go
  fi
}

copy_package
if [[ $PATH =~ $bindir ]]; then
  echo "gvm already installed in PATH"
else
  add_export_to_profile
fi
