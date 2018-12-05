#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''
-------------------------------------------------
    File Name：     test
    Author :        Afreto
    Date：          2018/12/4
    E-mail:         kongandmarx@163.com
-------------------------------------------------
    Description :   test
    Usage:  
-------------------------------------------------
'''


class LazyDB():
    def __init__(self):
        self.exists = 5

    def __getattr__(self, name):
        value = 'value for {}'.format(name)
        setattr(self, name, value)
        return value


class LoggingLazyDB(LazyDB):
    def __getattr__(self, name):
        print('Called __getattr__{}'.format(name))
        return super().__getattr__(name)


class ValidatingDB():
    def __init__(self):
        self.exists = 5

    def __getattribute__(self, name):
        print('Called __getattrbute__ {}'.format(name))
        try:
            return super().__getattribute__(name)
        except AttributeError:
            value = 'value for {}'.format(name)
            setattr(self, name, value)
            return value


class MissingPropertyDB():
    def __getattr__(self, name):
        if name == 'bad_name':
            raise AttributeError('{} is missing'.format(name))


if __name__ == '__main__':


    # print(t1.pop())
    # print(t1)

    # data = MissingPropertyDB()
    # data.bad_name

    # data = ValidatingDB()
    # print('exists:  {}'.format(data.exists))
    # print('foo: {}'.format(data.foo))
    # print('foo: {}'.format(data.foo))

    # data = LoggingLazyDB()
    # print('before:  {}'.format(data.__dict__))
    # print('foo exists:  {}'.format(hasattr(data, 'foo')))
    # print('after:   {}'.format(data.__dict__))
    # print('foo exists:  {}'.format(hasattr(data, 'foo')))

    # print('exists:  {}'.format(data.exists))
    # print('foo: {}'.format(data.foo))
    # print('foo: {}'.format(data.foo))

    # data = LazyDB()
    # print('before: {}'.format(data.__dict__))
    # print('foo: {}'.format(data.foo))
    # print('after: {}'.format(data.__dict__))
    pass
