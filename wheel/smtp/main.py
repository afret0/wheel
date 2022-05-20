import smtplib
import time
from email.header import Header
from email.mime.text import MIMEText
from flask import Flask
from flask import request
from flask import jsonify

app = Flask(__name__)


class EmailSender:
    def __init__(self) -> None:
        self.receiver = "kongandmarx@163.com"
        self.sender = "kongandmarx@163.com"
        self.smtp_obj = smtplib.SMTP_SSL("smtp.163.com",port=994)
        # self.smtp_obj.connect("smtp.163.com", 25)
        self.smtp_obj.login("kongandmarx@163.com", "YVLZXZWJBYAHLCAJ")

    def send(self, to, subject, text):
        message = MIMEText(f"{text}")
        message["Subject"] = Header(subject)
        message["From"] = Header(f"{self.sender}")
        message["To"] = Header(f"{to}")
        self.smtp_obj.sendmail(self.sender, self.receiver, message.as_string())


email_sender = EmailSender()


@app.route("/", methods=["POST"])
def send_email():
    data = request.get_json()
    print(data)
    to = data["to"]
    subject = data["subject"]
    text = data["text"]
    email_sender.send(to, subject, text)
    return jsonify({"code": 1, "msg": "succeed"})


if __name__ == "__main__":
    # email_sender.send("18435155427@163.com", "test", "test")
    app.run(host="0.0.0.0", port=8080, debug=False)
    pass
