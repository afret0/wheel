'''
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
'''

import os
from icecream import ic
from pysnooper import snoop
import time


def cur_time():
    return time.strftime('%Y-%m-%d %H:%M:%S:  ', time.localtime(time.time()))


ic.configureOutput(prefix=cur_time(), includeContext=True)


def __res(code: int = 0, message: str = 'success', data='') -> dict:
    if data == '':
        return {'code': code, 'message': message}
    return {'code': code, 'message': message, 'data': data}


if __name__ == '__main__':
    pass
