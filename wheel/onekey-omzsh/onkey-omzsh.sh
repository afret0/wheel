#!/usr/bin/env bash
if ! [ -x "$(command -v git)" ]; then
  echo 'Error: git is not installed.' >&2
  exit 1
fi
wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | sh
git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting
# 下载 .zshrc
wget https://raw.githubusercontent.com/kong5664546498/half_a_wheel/master/wheel/onekey-omzsh/.zshrc

# mv /root/.zshrc /root/.zshrc.bk
chsh -s /bin/zsh
zsh


