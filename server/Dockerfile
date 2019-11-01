FROM ubuntu:14.04

RUN sudo apt-get update
RUN sudo apt-get -y upgrade
RUN sudo apt-get -y install golang

export GOPATH="/home/ubuntu/go" 
export GOROOT="/opt/go"

echo export GOPATH=/home/ubuntu/go >> ~/.bashrc
echo export GOROOT=/opt/go >> ~/.bashrc
echo export PATH=$PATH:$GOROOT/bin:$GOPATH/bin >> ~/.bashrc

RUN source ~/.bashrc

RUN go get github.com/boltdb/bolt
RUN go get github.com/nu7hatch/gouuid



ADD ./D7024E /image

WORKDIR /image

CMD sudo go run main.go



