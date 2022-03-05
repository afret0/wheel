import json
shelf = {"house":{},"stock":{},"assets":{}}
class House:
    pass

class Stock:
    id = ""
    price = 0
    min = 0
    max = 0 



class User:
    name = ""

    def __init__(self, name: str) -> None:
        self.name = name

class Shelf:
    house = {}

class DB:
    def __init__(self) -> None:
        self.db = "DB.json"
    def save(self,data:str):
        with open(self.db,"w") as f:
            f.write(json.dumps(data))
    
    def reload(self):
        with open(self.db) as f :
            return json.loads(f.read)


class Stock_Operator:
    def __init__(self) -> None:
        self.stock = shelf["stock"]
        self.db = DB()

    def update(self,id:str,price: int,min:int=0,max:int=0):
        self.stock["id"]=id
        self.stock["price"]=price
        self.stock['min']=min
        self.stock['max'] = max
        shelf["stock"] = self.stock
        self.db.save(shelf)
        pass

    def sell(self,name)

def show_shelf():
    print(shelf)

so = Stock_Operator()
so.update("test",3)
show_shelf()