'''
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
    Description : 
    Usage: 
-------------------------------------------------
'''

import dramatiq
import sys
import time
from dramatiq.brokers.rabbitmq import RabbitmqBroker

rabbitmq_broker = RabbitmqBroker(host='127.0.0.1')
dramatiq.set_broker(rabbitmq_broker)

@dramatiq.actor()
def count_to(n):
    for i in range(n):
        print('Current value', i)
        time.sleep(1)
    print('all done')


if __name__ == '__main__':
    count_to.send(10)
