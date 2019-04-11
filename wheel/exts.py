'''
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
'''

import os
from config import LocConfig, DevConfig, ProConfig


def __res(code: int = 0, message: str = 'success', data='') -> dict:
    if data == '':
        return {'code': code, 'message': message}
    return {'code': code, 'message': message, 'data': data}


def get_config() -> object:
    """
    根据环境变量 ENV 切换配置
    :return: 配置
    """
    env = os.environ.get('ENV', 'loc')
    rules = {'loc': LocConfig, 'dev': DevConfig, 'pro': ProConfig}
    return rules.get(env)


config = get_config()



