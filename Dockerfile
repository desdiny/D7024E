FROM google/golang

WORKDIR /gopath/src/D7024E.git/branches/Objective-2
ADD . /gopath/src/D7024E.git/branches/Objective-2

RUN go get github.com/nu7hatch/gouuid

RUN go get github.com/desdiny/D7024E

CMD []
ENTRYPOINT ["/gopath/"]