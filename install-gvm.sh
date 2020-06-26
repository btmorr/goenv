#!/usr/bin/env bash

shimroot="$HOME/.gvm"
bindir="$shimroot/bin"
shimdir="$shimroot/shims"

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
  mkdir -p $bindir
  mkdir -p $shimdir
  cp -R bin $shimroot
}

# todo: save file with path to previous installation of go

if [[ $PATH =~ $bindir ]]; then
  echo "gvm already installed in PATH"
else
  copy_package
  add_export_to_profile
fi
