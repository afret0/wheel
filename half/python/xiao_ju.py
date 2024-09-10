import hashlib
import json
from collections import OrderedDict
from Crypto.Cipher import AES
import base64
import urllib.parse


class PiggyOpenBasicsVerifyUtil:

    @staticmethod
    def get_sig_str(secret, params):
        # 获取map转换为String 字符串
        params_string = PiggyOpenBasicsVerifyUtil.build_map_to_sign(params)
        # 获取密文串
        base_string = secret + params_string

        # 使用MD5对待签名串求签
        md5 = hashlib.md5()
        md5.update(base_string.encode('utf-8'))
        bytes = md5.digest()

        # 将MD5输出的二进制结果转换为小写的十六进制
        sign = ''.join(['{:02x}'.format(byte) for byte in bytes])
        return sign

    @staticmethod
    def build_map_to_sign(params):
        sort_map = OrderedDict(sorted(params.items()))
        return json.dumps(sort_map, separators=(',', ':'))


class OpenBasicsAesUtil:

    @staticmethod
    def produce_aes_data(secret, iv_str, params):
        try:
            cipher = AES.new(secret.encode('utf-8'), AES.MODE_CBC, iv_str.encode('utf-8'))
            json_param = PiggyOpenBasicsVerifyUtil.build_map_to_sign(params)
            pad = lambda s: s + (AES.block_size - len(s) % AES.block_size) * chr(
                AES.block_size - len(s) % AES.block_size)
            encrypted = cipher.encrypt(pad(json_param).encode('utf-8'))
            en_string = base64.b64encode(encrypted).decode('utf-8')
            return urllib.parse.quote(en_string)
        except Exception as e:
            print("AES加密异常", e)
        return ""

    @staticmethod
    def analysis_aes_data_dev(secret, iv_str, aes_data):
        try:
            raw = secret.encode('utf-8')
            cipher = AES.new(raw, AES.MODE_CBC, iv_str.encode('utf-8'))
            temp_str = urllib.parse.unquote(aes_data)
            encrypted1 = base64.b64decode(temp_str)
            pad = lambda s: s[:-ord(s[len(s) - 1:])]
            original = pad(cipher.decrypt(encrypted1)).decode('utf-8')
            return original
        except Exception as e:
            print("AES解密异常", e)
        return ""


# 示例用法
if __name__ == "__main__":
    params = {
        "wechatAppId": "wx6bc35656563fc27e6",
        "payAccount": "ooKio5yG5888888AX8U3tO7ppI",
        "positionName": "新手提现",
        "licenseType": "ID_CARD",
        "empPhone": "15811112222",
        "payAmount": 0.01,
        "month": "2021-05",
        "empName": "李飞",
        "notifyUrl": "https://www.baidu.com",
        "outerTradeNo": "91b3f16bafcc440b976b32b7262684d3",
        "taxFundId": "a22f81a74ee98ff962ff723dc85401e8",
        "licenseId": "110010200010102323",
        "settleType": "wechatpay"
    }
    secret = "18a9722d3bc84812830b2306b72e4605"
    iv_str = "0000000000000000"

    # 签名生成
    p = PiggyOpenBasicsVerifyUtil()
    sign = p.get_sig_str(secret, params)
    print("Signature:", sign)

    # AES加密
    o = OpenBasicsAesUtil()
    encrypted = o.produce_aes_data(secret, iv_str, params)
    print("Encrypted:", encrypted)

    # AES解密
    decrypted = o.analysis_aes_data_dev(secret, iv_str, encrypted)
    print("Decrypted:", decrypted)
