FROM debian:unstable-slim
MAINTAINER afreto kongandmarx@163.com
RUN ["apt-get","install","python3","-y"]
RUN ["apt","install","python3-pip","-y"]
RUN ["pipe","install","pipenv"]


#docker build . --build-arg "HTTP_PROXY=http://127.0.0.1:1087/" --build-arg "HTTPS_PROXY=http://127.0.0.1:1087/"  --network host -t debian:1