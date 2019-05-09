import dramatiq
import requests
import asyncio
from pysnooper import snoop
from tenacity import retry
from tenacity.stop import stop_after_attempt


@retry(stop=stop_after_attempt(3))
@snoop()
def ts():
    tag = [1]
    data = 'asdf' if isinstance(tag, list) else tag
    print(data)
    print(isinstance(tag, list))
    r = requests.get('https://www.google.com', timeout=1)
    print(r.text)


ts()






# import sentry_sdk
#
# sentry_sdk.init("https://e2926eae717a496c947432aaa80be374@sentry.io/1414241")
#
# class Person:
#     @property
#     def eat(self):
#         return 1
#
# tom = Person()
# print(tom.eat)


from raven import Client

# DNS = 'https://e2926eae717a496c947432aaa80be374@sentry.io/1414241'
# client = Client(DNS)
#
#
# try:
#     1/0
# except:
#     client.captureException()
#


#
# @dramatiq.actor
# def count_words(url):
#     r = requests.get(url)
#     count = len(r.text.split(' '))
#     print('there are {} words at {}'.format(count, url))


#
# async def test(v):
#     print('start test')
#     r = await asyncio.sleep(3)
#     print('input {}'.format(v))
#
# loop = asyncio.get_event_loop()
# c = test(u'xxx')
# loop.run_until_complete(c)
# loop.close()
if __name__ == '__main__':
    # count_words.send('http://example.com')
    # count_words.send('http://baidu.com')
    pass
