from flask import Flask
from flask import request
from flask import jsonify
from email_sender import EmailSender

app = Flask(__name__)


@app.route("/", methods=["POST"])
def send_email():
    data = request.get_json()
    print(data)
    to = data["to"]
    subject = data["subject"]
    text = data["text"]
    email_sender = EmailSender()
    email_sender.send(to, subject, text)
    return jsonify({"code": 1, "message": "succeed"})


if __name__ == "__main__":
    # email_sender.send("18435155427@163.com", "test", "test")
    app.run(host="0.0.0.0", port=8080, debug=False)
    pass
