FROM afreto/debian
COPY . /root
RUN ["pipenv","install"]
RUN ["pipenv","install","flask"]
CMD ["pipenv","run","python","/root/main.py"]


#docker build . --build-arg "HTTP_PROXY=http://127.0.0.1:1087/" --build-arg "HTTPS_PROXY=http://127.0.0.1:1087/"  --network host -t smtp:1