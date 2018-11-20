# DATE :  2018/7/23 
# AUTHOR : kongandmarx@163.com


import prometheus_client
from prometheus_client import Gauge
from flask import Response, Flask

test = {r'monitor': 1}
app = Flask(__name__)
gauges = test
# Gauge列表
gauge_buffer = []
# 位置列表
index_buffer = []
for ele in gauges:
    # 记录Gauge 列表位置
    index_buffer.append(ele)
    # 循环生成Gauge实例,添加到列表
    ele = Gauge('{}'.format(ele), '', ['instance', ])
    gauge_buffer.append(ele)


@app.route('/metrics')
def metrics():
    # 每次访问调用 test
    data = test
    res_list = []
    for t in gauge_buffer:
        t.labels(instance='{}:{}'.format('ip', 'port')).set(data[index_buffer[gauge_buffer.index(t)]])
        # generate_latest更新列表
        res_list.append(prometheus_client.generate_latest(t))
    return Response(res_list, mimetype='text/plain')
    pass


if __name__ == '__main__':
    app.debug = True
    app.run(host='localhost', port=10110)
    pass