#!/usr/bin/env python
# -*- coding: utf-8 -*-
# Created on 2018/11/7
import subprocess
import requests
from selenium import webdriver

driver = webdriver.PhantomJS(executable_path=r'C:\Program Files (x86)\phantomjs-2.1.1-windows\bin\phantomjs')
driver.get('http://zhyq.hqiic.com/cockpit/front/pages/index.html#/app/zsxx')
tar = driver.page_source
tar = tar
with open('index.html', 'w+',encoding='utf-8') as  f:
    f.writelines(tar)
driver.quit()
