#! /bin/bash
# 判断目录是否存在
if [ ! -d  ~/mine ] 
then
    echo "~/mine not exists"
else
    echo -e "\033[32mexists\033[1m"
fi