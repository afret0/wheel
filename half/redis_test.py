'''
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
    Description : 
    Usage: 
-------------------------------------------------
'''

import redis


r = redis.Redis(host='192.168.3.3', decode_responses=True)
# test1 = {'afreto': 4}
test1 = {'afreto1': 5}
# r.mset(test1)
# resp = r.msetnx(test1)
resp = r.zadd(name='test_zdd', mapping=test1)
# resp = r.zcount()
# r.zincrby()

print(resp)
