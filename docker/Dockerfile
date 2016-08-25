FROM ubuntu:wily
MAINTAINER caktux

ENV DEBIAN_FRONTEND noninteractive

# Usual update / upgrade
RUN apt-get update
RUN apt-get upgrade -q -y
RUN apt-get dist-upgrade -q -y

# Install Krypton
RUN apt-get install -q -y software-properties-common
RUN add-apt-repository ppa:krypton/krypton
RUN add-apt-repository ppa:krypton/krypton-dev
RUN apt-get update
RUN apt-get install -q -y gkr

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/usr/bin/gkr"]
