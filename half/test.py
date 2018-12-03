import time


class LazyDB():
    def __init__(self):
        self.exists = 5

    def __getattr__(self, name):
        value = 'value for {}'.format(name)
        setattr(self, name, value)
        return value


if __name__ == '__main__':
    data = LazyDB()
    print('before: {}'.format(data.__dict__))
    print('foo: {}'.format(data.foo))
    print('after: {}'.format(data.__dict__))
