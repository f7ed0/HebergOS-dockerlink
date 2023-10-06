FROM ubuntu:latest

RUN useradd -rm -d /var/www -s /bin/bash -g root -G sudo -u 1000 admin

RUN  echo 'admin:changeitpls' | chpasswd

RUN apt update && apt install nano git screen openssh-server sudo -y

RUN echo "#!bin/bash\nexit 0;" > /var/www/starter.sh

RUN service ssh start