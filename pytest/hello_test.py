'''
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
    Description : 
    Usage: 
-------------------------------------------------
'''
import os
import sys

sys.path.append('..')
from half.hello_world import hello


def hello_test():
    assert hello('world') == 'world'
    assert hello('') == 'kitty'
    assert hello('world') == 'kitty'
