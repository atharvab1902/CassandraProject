#!/bin/bash

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get update
apt-get install -y docker-ce docker-compose-plugin
usermod -aG docker ubuntu

curl -OL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
tar -C /usr/local -xvf go1.21.0.linux-amd64.tar.gz
rm go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/ubuntu/.bashrc

su - ubuntu
