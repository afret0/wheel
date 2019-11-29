from xml.dom.minidom import parse
import xml.dom.minidom
import xml.etree.ElementTree as ET
import re

# from yattag import indent


def get_tag(s: str, tag: str) -> str:
    reg = f"<{tag}>.*</{tag}>"
    # print(reg)
    v = re.findall(reg, s)
    # print(v)
    v1 = v[0].split(">")[1]
    # print(v1)
    v2 = v1.split("<")[0]
    return v2
    # print(v2)


def main():
    s = '<appid>wx50aa92aff9dece89</appid><mch_id>1501874481</mch_id><nonce_str>Rnd58CRxhbfax6Zs</nonce_str><prepay_id>wx1911402597053289bf4ba8741034009600</prepay_id><result_code>SUCCESS</result_code><return_code>SUCCESS</return_code><return_msg>OK</return_msg><sign>appid=wx50aa92aff9dece89&noncestr=Rnd58CRxhbfax6Zs&package=Sign=WXPay&partnerid=1501874481&prepayid=wx1911402597053289bf4ba8741034009600&timestamp=1574163625&sign=128C78725118195EA5096B7D73302D51</sign><trade_type>APP</trade_type>'

    # tree = ET.fromstring(s)
    # root = tree.getroot()
    # appid = root.find("appid")
    # s1 = indent(s.encode('utf-8'))

    # dom_tree = xml.dom.minidom.parseString(s.encode("utf-8"))
    # collection = dom_tree.documentElement
    # root = collection.getElementsByTagName("xml")
    # appid = collection.getElementsByTagName("appid")
    # parser = xml.sax.make_parser()
    # parser.setFeature(xml.sax.handler.feature_namespaces, 0)

    print(get_tag(s,"sign"))
    raw_sign = get_tag(s,"sign")
    sign = raw_sign.split("=")[-1]

    pay_timestamp = raw_sign.split("timestamp")[-1]
    pay_timestamp = pay_timestamp.split("=")[1].split("&")[0]
    
    print(pay_timestamp)
    # print(collection)

    print("stop")


main()
