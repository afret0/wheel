#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
    Description : 脚本启动管理器
    Usage: python3 ctl.py -h
-------------------------------------------------
'''

import os, subprocess, getopt, sys, time, argparse


class Ctl():
    def __init__(self, script: str):
        self.pwd = os.getcwd()
        self.script = os.path.join(self.pwd, script)

    def is_run(self):
        # 检查是否启动
        p = subprocess.Popen('ps -ax | grep {} | grep -v grep'.format(self.script), stdout=subprocess.PIPE, shell=True)
        p.wait()
        # 捕获stdout
        pids = []
        for readline in p.stdout.readlines():
            readline = str(readline, encoding='utf-8').split(' ')[1]
            pids.append(readline)
        if pids:
            return pids
        else:
            return False

    def start(self, options: str = None, dev: bool = False):
        if self.is_run():
            print('{} is running ... '.format(self.script))
        else:
            print('{} to start running ...'.format(self.script))
            if dev:
                command = '{python3} {pro} {op}'.format(python3=sys.executable, pro=self.script, op=options)
            else:
                command = 'nohup {python3} {pro} {op} & '.format(python3=sys.executable, pro=self.script, op=options)
            p = subprocess.Popen(command, shell=True)
            p.wait()
            time.sleep(1)
            # 再次检查进程确认是否存活
            if self.is_run():
                print('{} is running ... '.format(self.script))
            else:
                print('{} startup failed ...'.format(self.script))

    def stop(self):
        pids = self.is_run()
        if pids:
            print('{} is stopping ...'.format(self.script))
            for pid in pids:
                p = subprocess.Popen('kill -9 {}'.format(pid), shell=True)
                p.wait()
            print('{} stopped ...'.format(self.script))
            return True
        else:
            print('{} stopped ...'.format(self.script))
            print('Make sure the {} is started ...'.format(self.script))
            return False

    def status(self):
        if self.is_run():
            print('{} is running ...'.format(self.script))
        else:
            print('{} stopped ...'.format(self.script))


def manage(script: str = None):
    """
    处理命令行参数
    :param script:
    """
    parser = argparse.ArgumentParser(description='A Python script management tool')
    parser.add_argument('-r', '-run', dest='run', action='store_true', help='Startup script ...')
    parser.add_argument('-d', '-down', dest='down', action='store_true', help='Stop script ...')
    parser.add_argument('-restart', dest='restart', action='store_true', help='Restart the script ...')
    parser.add_argument('-status', dest='status', action='store_true', help='Startup script ...')
    parser.add_argument('-o', '-options', dest='options', action='store', help='Parameters at startup ...')
    parser.add_argument('-s', '-script', dest='script', action='store', help='Script to run ...')
    parser.add_argument('-dev', dest='dev', action='store_true', help='Output All ...')
    args = parser.parse_args()
    if args.script:
        ctr = Ctl(args.script)
    else:
        ctr = Ctl(script)
    if args.run:
        if args.dev:
            ctr.start(options=args.options, dev=True)
        else:
            ctr.start(options=args.options)
    elif args.down:
        ctr.stop()
    elif args.restart:
        ctr.stop()
        if args.dev:
            ctr.start(options=args.options, dev=True)
        else:
            ctr.start(options=args.options)
    elif args.status:
        ctr.status()
    else:
        print('usage: ctl.py [-h] for HELP ...')


if __name__ == '__main__':
    manage('app.py')
