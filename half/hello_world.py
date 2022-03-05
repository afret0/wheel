"""
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
    Description : 
    Usage: 
-------------------------------------------------
"""
import time
import pretty_errors


def hello(world: str) -> str:
    print("1111111111")
    time.sleep(5)

    return world if world else print("kitty")


class Bird:
    def __init__(self, name: str):
        self.name = name

    def fly(self):
        print(f"{self.name} is flying")

    def jump(self):
        print(f"{self.name} is jumping")

    def eat(self):
        print(f"eating")


# import requests
# url = "http://www.baidu.com"
# res = requests.get(url=url, params={})
# print(res.status_code)
if __name__ == "__main__":
    tag = "all"
    # tag = ["a","b"]
    if isinstance(tag, str) and tag != "all":
        tag = [tag]
    a = {"tag": tag} if isinstance(tag, list) else tag
    print(a)
    # s = "a"

    # a = int(s)

    # for i in range(1, 100):
    # print(1)
