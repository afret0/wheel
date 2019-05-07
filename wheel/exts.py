'''
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
'''

import os
from icecream import ic


def __res(code: int = 0, message: str = 'success', data='') -> dict:
    if data == '':
        return {'code': code, 'message': message}
    return {'code': code, 'message': message, 'data': data}
