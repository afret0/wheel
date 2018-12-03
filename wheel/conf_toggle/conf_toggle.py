#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''
-------------------------------------------------
    File Name：     conf_toggle.py
    Author :        Afreto
    E-mail:         kongandmarx@163.com
    Date:           2018/12/3
-------------------------------------------------
    Description :
        一键切换配置文件, 需切换的配置文件在 conf 文件中读取
-------------------------------------------------
'''

import os
import json
import shutil
import sys
import argparse


class ConfToggle():
    def __init__(self, conf: str):
        self.conf = conf
        self.files = {}

    @staticmethod
    def merge_path(path_list: list) -> str:
        """
        合并路径
        :param path_list: 路径分级列表
        :return: 合并后的路径
        """
        path = ''
        for arg in path_list:
            path = os.path.join(path, arg)
        return path

    @staticmethod
    def priority(file: dict):
        """
        处理路径优先级
        """
        file_basename = file.get('file', None)
        if 'abspath' in file:
            return file.get('abspath', None)
        elif 'dir' in file:
            file_dirpath = ConfToggle.merge_path(file.get('dir', None))
            return os.path.join(file_dirpath, file_basename)

    def read_conf(self) -> list:
        """
        读取 conf 文件内容
        """
        files = []
        # 读取内容 解析 json
        with open(self.conf, 'r') as f:
            content = f.readlines()
        content = ''.join(content)
        content = content.replace('\n', '').replace(' ', '')
        content = json.loads(content)
        # 处理内容 判断路径优先级
        for data in content:
            dev_files = data.get('dev', False)
            dist_files = data.get('dist', False)
            dev = ConfToggle.priority(dev_files)
            dist = ConfToggle.priority(dist_files)
            files.append({'dev': dev, 'dist': dist})
        self.files = files
        return self.files

    def re_backup(self, tar: str = 'dev'):
        """
        从备份恢复
        :param tar: dev: 操作dev dist: 操作dist
        """
        for file in self.files:
            dst = file[tar]
            src = '{}_{}_bk'.format(file[tar], tar)
            if os.path.exists(src):
                try:
                    shutil.move(src, dst)
                except Exception as e:
                    raise e
                else:
                    print('已从备份恢复 {}'.format(src))
                    print('删除备份')
            else:
                print(r'备份 {} 不存在'.format(src))

    def mk_backup(self, tar: str = 'dev'):
        """
        创建备份
        :param tar: dev: 操作dev dist: 操作dist
        """
        for file in self.files:
            src = file[tar]
            dst = '{}_{}_bk'.format(file[tar], tar)
            if os.path.exists(dst):
                print(r'''备份已存在,可能存在错误覆盖,程序已主动退出.
                请确认覆盖内容, 可使用 re_back 方法恢复备份后重新运行''')
                sys.exit(0)
            else:
                print('创建备份 {} '.format(dst))
                shutil.copy(src, dst)

    def toggle_conf_from_dist_to_dev(self):
        """
        替换 conf 文件
        """
        for file in self.files:
            src = file['dist']
            dst = file['dev']
            print('正在替换配置文件 {} '.format(src))
            shutil.copy(src, dst)
            print('配置文件 {} 替换完成'.format(src))


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Toggle configuration Files')
    parser.add_argument('-f', '--conf', dest='conf', action='store',
                        help='Specifies that the configuration file is sought by default./conf')
    parser.add_argument('-t', '--toggle', dest='toggle', action='store_true', help='Toggle configuration Files')
    parser.add_argument('-r', '--re_backup', dest='re_backup', action='store_true',
                        help='Restore configuration from backup files')
    args = parser.parse_args()
    print(args.re_backup)
    if args.conf:
        confer = ConfToggle(args.conf)
    else:
        confer = ConfToggle('./conf')
    confer.read_conf()
    if args.toggle:
        confer.mk_backup()
        confer.toggle_conf_from_dist_to_dev()
    elif args.re_backup:
        confer.re_backup()
    else:
        print('''usage: conf_toggle.py [-h] for help''')
