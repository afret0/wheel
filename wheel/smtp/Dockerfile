FROM afreto/debian
MAINTAINER afreto kongandmarx@163.com
COPY main.py /root
RUN ["pipenv","install"]
CMD ["pipenv","run","python","/root/main.py"]
