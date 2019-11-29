from threading import Thread
from hello_world import hello

t = Thread(target=hello, args=("kitty",))
c = Thread(target=hello, args=("kitty",))
t.start()
c.start()
