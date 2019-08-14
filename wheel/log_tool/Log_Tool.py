# -*- coding: utf-8 -*-
# DATE :  2018/6/8
import logging.handlers
import configparser
import os
import sys


class Log_Tool:
    def __init__(self, logger_name=""):
        # 获取log配置
        config = configparser.ConfigParser()
        curdir = os.path.abspath(os.curdir)
        CONF = os.path.join(curdir, "log_tool")

        CONF = os.path.join(CONF, "logger.conf")
        config.read(CONF)
        self.log_dir = config.get("LOG", "dir")
        # self.fmt = config.get('FMT', 'fmt')
        self.fmt = "%(asctime)s - %(levelname)s - %(message)s"
        if logger_name:
            self.looger_name = logger_name
        else:
            self.looger_name = config.get("LOGGER", "name")
        self.level = config.get("LEVEL", "level")
        self.maxBytes = eval(config.get("LOG", "maxBytes"))
        self.backupCount = config.getint("LOG", "backupCount")

    def get_logger(self) -> object:
        # 获取handler
        handler = logging.handlers.RotatingFileHandler(
            self.log_dir,
            maxBytes=self.maxBytes,
            backupCount=self.backupCount,
            encoding="utf-8",
        )
        console_handler = logging.StreamHandler(sys.stdout)
        # 实例化formatter
        formatter = logging.Formatter(self.fmt)
        # 给handler添加formatter
        handler.setFormatter(formatter)
        # 获取logger
        logger = logging.getLogger(self.looger_name)
        # 给logger添加handler
        logger.addHandler(handler)
        logger.addHandler(console_handler)
        # 设置过滤
        logger.setLevel(self.level)
        return logger


logger = Log_Tool().get_logger()

if __name__ == "__main__":
    logger = Log_Tool()
    logger = logger.get_logger()
    logger.info("test_log")
