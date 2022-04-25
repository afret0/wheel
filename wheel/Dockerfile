FROM debian
RUN ["apt-get","update"]
RUN ["apt-get","install","wget","-y"]
RUN ["apt-get","install","python3"]
RUN ["apt","install","python3-pip"]
RUN ["wget","https://raw.githubusercontent.com/kong5664546498/half_a_wheel/master/wheel/onekey-omzsh/onekey-omzsh.sh"]
RUN ["chmod", "+x", "onekey-omzsh.sh"]
RUN ["./onekey-omzsh.sh"]

# RUN ["apt-get","install","npm","-y"]
# RUN ["npm","install","-g","nodeppt"]

# docker run  -it -d -p 10119:8080 --name nodeppt -v /Users/afreto/mine/buffer:/data/ nodeppt /bin/bash