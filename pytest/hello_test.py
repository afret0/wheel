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


def test_hello():
    assert hello('world') == 'world'
    # assert hello('world') == 'kitty'
    # assert hello('') == 'kitty'
