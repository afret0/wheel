FROM afreto/debian
MAINTAINER afreto kongandmarx@163.com
COPY . /root
RUN ["pipenv","install"]
RUN ["pipenv","install","flask"]
CMD ["pipenv","run","python","/root/main.py"]
