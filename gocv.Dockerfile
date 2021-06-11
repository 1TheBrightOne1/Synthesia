FROM ubuntu

RUN apt update && apt install -y wget && apt install -y make

RUN wget https://golang.org/dl/go1.16.5.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz
RUN mkdir /go
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
RUN go version

RUN apt install -y sudo

RUN go get -u -d gocv.io/x/gocv 
WORKDIR $GOPATH/pkg/mod/gocv.io/x/gocv@v0.27.0
