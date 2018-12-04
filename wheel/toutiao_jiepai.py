#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''
-------------------------------------------------
    File Name：     toutiao_jiepai.py
    Author :        Afreto
    E-mail:         kongandmarx@163.com
    Date:           2018/6/22
-------------------------------------------------
    Description :   下载头条街拍图片
    Usage:  python3 toutiao_jiepai.py
-------------------------------------------------
'''

import requests, os, time


def get_page(offset: int):
    url = r'https://www.toutiao.com/search_content/?offset={}&format=json&keyword=%E8%A1%97%E6%8B%8D&autoload=true&count=20&cur_tab=1&from=search_tab'.format(
        offset)
    try:
        r = requests.get(url=url)
    except Exception as e:
        raise e
    else:
        if r.status_code == 200:
            return r.json()


def get_images(json):
    res = []
    if json.get('data'):
        for item in json.get('data'):
            title = item.get('title')
            images = item.get('image_list')
            if images:
                for image in images:
                    # print(image)
                    res.append({'title': title, 'image': image.get('url')}.copy())
    return res


def save_image(item: dict):
    usr = os.environ.get('USERPROFILE')
    tag = os.path.join(os.path.join(usr, 'Downloads'), 'demo')
    file_path = tag
    if not os.path.exists(file_path):
        os.mkdir(file_path)
    large_img = 'http:' + str(item.get('image')).replace('list', 'large')
    try:
        r = requests.get(large_img)
    except Exception:
        raise Exception
    else:
        if r.status_code == 200:
            file_path = r'{}\{}.{}'.format(file_path, large_img.split('/')[-1], 'jpg')
            with open(file_path, 'wb') as f:
                f.write(r.content)
        else:
            print('already download')


def begin(group):
    page = get_page(group)
    if page:
        for img in get_images(page):
            print('img --> ', img)
            save_image(img)
            time.sleep(0.1)


if __name__ == '__main__':
    from multiprocessing.pool import Pool

    GROUP_START = 1
    GROUP_END = 5
    pool = Pool()
    group = ([x * 20 for x in range(GROUP_START, GROUP_END + 1)])
    pool.map(begin, group)
    pool.close()
    pool.join()
