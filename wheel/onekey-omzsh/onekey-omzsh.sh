#!/usr/bin/env bash
# ============== 安装 git zsh wget ==============
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
# ===================================


# ============== 安装配置 oh-my-zsh ==============
wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | sh
git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting

if [ -f "/root/.zshrc" ]; then
     mv /root/.zshrc /root/.zshrc.bk
fi

if [ -f "./onekey-omzsh.sh" ]; then
    rm ./onekey-omzsh.sh
fi

# 下载 .zshrc
wget -P /root https://raw.githubusercontent.com/kong5664546498/half_a_wheel/master/wheel/onekey-omzsh/.zshrc
# ===================================

# ============== gotop ==============
# git clone --depth 1 https://github.com/cjbassi/gotop /tmp/gotop

# /tmp/gotop/scripts/download.sh

# mv ./gotop /usr/bin
# ===================================

# autojump
#git clone git://github.com/wting/autojump.git
#cd autojump
#./install.py

# the fuck
pip3 install thefuck

# 修改默认 bash
chsh -s /bin/zsh
zsh


