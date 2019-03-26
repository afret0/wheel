'''
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
'''

import os
from config import LocConfig, DevConfig, ProConfig


def __res(code: int = 0, message: str = 'success') -> dict:
    return {'code': code, 'message': message}


def get_config() -> object:
    """
    根据环境变量 ENV 切换配置
    :return: 配置
    """
    env = os.environ.get('ENV', 'loc')
    rules = {'loc': LocConfig, 'dev': DevConfig, 'pro': ProConfig}
    return rules.get(env)


config = get_config()


def phone_check(phone: str) -> tuple:
    """
    检查电话号码
    :param phone:
    :return: ([phone 校验成功则返回 phone, 校验失败则返回错误信息],错误标志)
    """
    if not phone:
        return __res(-1, 'unable to get phone'), False
    phone = phone.strip()
    phone_err_res = __res(-1, 'Abnormal phone number'), False
    if len(phone) != 11:
        return phone_err_res
    return phone, True
    # try:
    #     phone = int(phone)
    # except Exception as e:
    #     return phone_err_res
    # else:
    #     # 可在此修饰 phone
    #     return phone, True
