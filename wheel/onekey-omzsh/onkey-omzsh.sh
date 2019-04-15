#!/usr/bin/env bash
if ! [ -x "$(command -v git)" ]; then
  echo 'Error: git is not installed.' >&2
  exit 1
fi
wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | sh
git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting
# 下载 .zshrc
#git clone https://gist.github.com/63281965e83dfbea69a8601a5da10476.git

# mv /root/.zshrc /root/.zshrc.bk
#mv 63281965e83dfbea69a8601a5da10476/.zshrc /root/.zshrc

# source /root/.zshrc

chsh -s /bin/zsh
rm -rf 63281965e83dfbea69a8601a5da10476
zsh


