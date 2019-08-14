"""
-------------------------------------------------
    Author :        Afreto
    E-mail:         kongandmarx@163.com
-------------------------------------------------
    Description : 
    Usage: 
-------------------------------------------------
"""


def hello(world: str) -> str:
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


import requests

url = "http://www.baidu.com"
res = requests.get(url=url, params={})
print(res.status_code)

if __name__ == "__main__":
    for i in range(1, 100):
        print(i, end=">>>\r", flush=False)
