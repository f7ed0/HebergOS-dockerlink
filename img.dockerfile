FROM ubuntu:latest

RUN useradd -rm -d /var/www -s /bin/bash -g root -G sudo -u 1000 admin

RUN  echo 'admin:changeitpls' | chpasswd

RUN apt update && apt install nano git pip mysql-client screen openssh-server gnupg2 sudo -y

RUN sudo sh -c 'echo "deb https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'

RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -

RUN sudo apt-get update

RUN sudo apt-get -y install postgresql-client

RUN echo "#!bin/bash\nexit 0;" > /var/www/starter.sh

RUN service ssh start