import smtplib
from email.header import Header
from email.mime.text import MIMEText


class EmailSender:
    def __init__(self) -> None:
        self.receiver = "kongandmarx@163.com"
        self.sender = "kongandmarx@163.com"
        self.smtp_obj = smtplib.SMTP_SSL("smtp.163.com", port=994)
        # self.smtp_obj.connect("smtp.163.com", 25)
        self.smtp_obj.login("kongandmarx@163.com", "YVLZXZWJBYAHLCAJ")

    def send(self, to, subject, text):
        t = f"""
        <h1> {text} </h1>
        """
        message = MIMEText(t, "html")
        message["Subject"] = Header(subject)
        message["From"] = Header(f"{self.sender}")
        message["To"] = Header(f"{to}")
        self.smtp_obj.sendmail(self.sender, self.receiver, message.as_string())