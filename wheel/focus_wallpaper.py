#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''
-------------------------------------------------
    File Name：     focus_wallpaper.py
    Author :        Afreto
    E-mail:         kongandmarx@163.com
    Date:           2018/6/22
-------------------------------------------------
    Description :   提取 windows 聚焦壁纸
    Usage:  python3 focus_wallpaper.py
-------------------------------------------------
'''
import os, os.path
import shutil, filecmp
import stat
from PIL import Image


def move_file(srcpath, dstpath):  # 复制文件
    filelist = os.listdir(srcpath)
    for files in filelist:
        Olddir = os.path.join(srcpath, files)  # 原来的文件路径
        shutil.copy(Olddir, dstpath)
    print('.............move_file............ done')


def rename_file(dstpath):  # 重命名文件
    filelist = os.listdir(dstpath)
    for files in filelist:
        Olddir = os.path.join(dstpath, files)  # 原来的文件路径
        if os.path.isdir(Olddir):  # 如果是文件夹则跳过
            continue
        filename = os.path.splitext(files)[0]  # 文件名
        # filetype=os.path.splitext(files)[1];#文件扩展名
        Newdir = os.path.join(dstpath, filename + ".jpg")  # 新的文件路径
        try:
            os.rename(Olddir, Newdir)  # 重命名
        except:
            pass
    print('.............rename_file............ done')


def modify_file_attribute(dstpath):  # 修改文件属性
    filelist = os.listdir(dstpath)
    for files in filelist:
        Olddir = os.path.join(dstpath, files);  # 原来的文件路径
        os.chmod(Olddir, stat.S_IWRITE)  # 修改文件只读属性
    print('.............modify_file_attribute............ done')


def delete_file_less_200kb(dstpath):  # 删除小于200kb的文件
    filelist = os.listdir(dstpath)
    filelist.remove('midpath')
    for files in filelist:
        Olddir = os.path.join(dstpath, files);  # 原来的文件路径
        if os.path.getsize(Olddir) < 200000:  # 小于200kb且不是JPEG文件 and (im.mode != 'RGB' or im.format != 'JPEG')
            os.remove(Olddir)
    print('.............delete_less_200kb_file............ done')


def delete_file_gif(dstpath):  # 删除gif的文件
    filelist = os.listdir(dstpath)
    filelist.remove('midpath')
    for files in filelist:
        Olddir = os.path.join(dstpath, files);  # 原来的文件路径
        im = Image.open(Olddir)
        if im.format != 'JPEG':  # 小于200kb且不是JPEG文件 and (im.mode != 'RGB' or im.format != 'JPEG')
            im.close()
            os.remove(Olddir)
        else:
            im.close()

    print('.............delete_gif_file............ done')


def picture_length_width(path, midpath):  # 判断图片长宽将电脑壁纸复制到midpath
    filelist = os.listdir(path)
    filelist.remove('midpath')
    for files in filelist:
        Olddir = os.path.join(path, files);  # 原来的文件路径
        img = Image.open(Olddir)  # img.size 输出 (长, 宽)
        if img.size[0] > img.size[1]:  # 比较文件长宽
            shutil.copy(Olddir, midpath)
            pass
    print('.............picture_length_width............ done')


def delete_dstpath_file(dstpath):
    '''删除dstpath路径下图片文件'''
    file_list = os.listdir(dstpath)
    for file in file_list:
        file = os.path.join(dstpath, file)
        try:
            os.remove(file)
        except:
            pass


def main():
    srcpath = os.path.join(os.environ['LOCALAPPDATA'],
                           r'Packages\Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy\LocalState\Assets')
    # dstpath = r'C:\Users\kong5\OneDrive\图片\focus_paper'
    dstpath = os.path.join(os.environ['HOMEPATH'], r'OneDrive\图片\focus_paper')
    # midpath = r'C:\Users\kong5\OneDrive\图片\focus_paper\midpath'
    midpath = os.path.join(os.environ['HOMEPATH'], r'OneDrive\图片\focus_paper\midpath')
    for files in os.listdir(dstpath):
        files_dir = os.path.join(dstpath, files)
        os.remove(files_dir)
    try:
        os.mkdir(midpath)
        move_file(srcpath, dstpath)
        rename_file(dstpath)
        modify_file_attribute(dstpath)
        delete_file_less_200kb(dstpath)
        delete_file_gif(dstpath)
        picture_length_width(dstpath, midpath)
        delete_dstpath_file(dstpath)
        move_file(midpath, dstpath)
        shutil.rmtree(midpath)
    except Exception as e:
        with open('log.txt', 'a+') as f:
            f.writelines(e)


if __name__ == '__main__':
    main()
