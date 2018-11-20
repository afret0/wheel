from sqlalchemy import create_engine
from sqlalchemy import Column, Integer, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

# 创建连接
engine = create_engine('mysql+pymysql://bird:qpzm2745@127.0.0.1/demo?charset=utf8', echo=True)
# 创建元类
Base = declarative_base()
# 连接池
Session = sessionmaker(bind=engine)
# 实例化
session = Session()


class Phone(Base):
    # 表名
    __tablename__ = 'Phone_test'
    id = Column(Integer(), autoincrement=True, primary_key=True, unique=True)
    item_id = Column(String(10))
    name = Column(String(120))
    crawl_time = Column(String(100))

# 查
test1 = session.query(Phone).all()
print('test1 --> {}'.format(test1))
# 增
session.add(Phone(item_id='item_id_2',name='name_2',crawl_time='crawl_time_2'))
session.commit()
session.close()

