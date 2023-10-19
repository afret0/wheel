import json
import hashlib

if __name__=="__main__":
    with open("./2269cc780b355851f56da3b6ffc581d1.png", "rb") as f:
        data = f.read()
    with open("./35678fc36007bba831a1177f6deb73bd.png", "rb") as fq:
        data1 = fq.read()
    print(data == data1)
    # version = "1.0.8"
    # url = "https://v1.kekeyuyin.com/user/getUserInfo?uid=123123"
    # body = {"param": "123134"}
    # param = url.split("?")[1]
    # body_str = json.dumps(body)
    # salt = "X3y7v9T2m1N8J6k0P4o5W7u2E6H9F7g2L3b4A1Z8D9Q5R6S7C0x5V3nXr3m0r9s2D65B8V7"
    # s1 = version + param + body_str + salt
    # s2 = sorted(s1)
    # s3 = "".join(s2)
    # hash_object = hashlib.md5(s3.encode())
    # md5_hash = hash_object.hexdigest()
    # print(md5_hash)