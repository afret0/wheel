from flask import Flask

app = Flask(__name__)


@app.route('/', methods=['GET'])
def hello():
    print('hello world')
    return 'Hello World'


if __name__ == '__main__':
    app.run(host='127.0.0.1', port=10112)
