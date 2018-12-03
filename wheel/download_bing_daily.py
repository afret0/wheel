#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''
-------------------------------------------------
    File Name：     download_bing_daily.py
    Author :        Afreto
    Date：          2018/12/3
    E-mail:         kongandmarx@163.com
-------------------------------------------------
    Description :   下载 bing 当日壁纸
-------------------------------------------------
'''
import requests, re, time, os
# import  win32api, win32gui, win32con
from datetime import timedelta, datetime


def download_wallpaper():
    ''':return bool'''
    now = time.strftime('%m_%d', time.localtime())
    png = now + '.png'
    if os.path.exists(png):
        print(png + ' exited')
    else:
        try:
            res = requests.get('http://area.sinaapp.com/bingImg/')
            with open(png, 'wb')as f:
                f.write(res.content)
            return 1
        except:
            return 0


def wallpaper_name():
    now = time.strftime('%m_%d', time.localtime())
    name = now + '.png'
    if os.path.exists(name):
        paper = os.path.abspath(name)
        return paper
    else:
        return 0


def set_wallpaper(png):
    key = win32api.RegOpenKeyEx(win32con.HKEY_CURRENT_USER, "Control Panel\\Desktop", 0, win32con.KEY_SET_VALUE)
    win32api.RegSetValueEx(key, "WallpaperStyle", 0, win32con.REG_SZ, "2")
    # 2拉伸适应桌面,0桌面居中
    # win32api.RegSetValueEx(key, "TileWallpaper", 0, win32con.REG_SZ, "0")
    win32gui.SystemParametersInfo(win32con.SPI_SETDESKWALLPAPER, png, 1 + 2)


if __name__ == '__main__':
    while 1:
        if wallpaper_name():
            # set_wallpaper(png=wallpaper_name())
            break
        else:
            download_wallpaper()

    pass
