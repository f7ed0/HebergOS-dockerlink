FROM ubuntu:latest

RUN useradd -rm -d /var/www -s /bin/bash -g root -G sudo -u 1000 admin

RUN  echo 'admin:changeitpls' | chpasswd

RUN apt update && apt install  openssh-server nginx sudo -y

RUN service ssh start