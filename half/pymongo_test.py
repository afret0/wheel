'''
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
    Description : 
    Usage: 
-------------------------------------------------
'''


import sanic
from sanic import Sanic
from sanic.response import json
app = Sanic('test')

@app.route('/v1/relationship/test')
def nasha(request):
    return json([111])


app.run()
