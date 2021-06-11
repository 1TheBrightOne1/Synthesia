FROM thebrightone/gocv:go1.16

WORKDIR /go/src/app

ENV GO111MODULE=on

ENV skipframes=240

ENV row=50

RUN go get -u gocv.io/x/gocv@v0.26.0

RUN mkdir /outputs && mkdir /outputs/sessions

COPY . .

RUN go build -o /app .

ENTRYPOINT ["/app"]