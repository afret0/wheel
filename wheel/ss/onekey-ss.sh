#!/usr/bin/env bash
wget -qO- get.docker.com | bash 
systemctl start docker 
systemctl enable docker 
docker pull teddysun/shadowsocks-r
mkdir -p  /etc/shadowsocks-r 
wget -O /etc/shadowsocks-r/config.json https://raw.githubusercontent.com/kong5664546498/half_a_wheel/master/wheel/ss/config.json
docker run -d -p 80:80 -p 80:80/udp --name ssr -v /etc/shadowsocks-r:/etc/shadowsocks-r teddysun/shadowsocks-r 
# bbr
wget --no-check-certificate https://github.com/teddysun/across/raw/master/bbr.sh && chmod +x bbr.sh && ./bbr.sh