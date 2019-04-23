#!/usr/bin/env bash
cd /root
if [ -x "$(command -v yum)" ]; then
    echo ' use yum'>&2
    yum install git -y
    yum install zsh -y
    yum install wget -y
fi

if [ -x "$(command -v apt-get)" ]; then
    echo 'use apt-get'>&2
    apt-get install git -y
    apt-get install zsh -y
    apt-get install wget -y
fi


#if ! [ -x "$(command -v git)" ]; then
#  echo 'Error: git is not installed.' >&2
#  exit 1
#fi
wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | sh
git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting

if [ -f "/root/.zshrc" ]; then
     mv /root/.zshrc /root/.zshrc.bk
fi

# 下载 .zshrc
wget https://raw.githubusercontent.com/kong5664546498/half_a_wheel/master/wheel/onekey-omzsh/.zshrc

chsh -s /bin/zsh
zsh
#

