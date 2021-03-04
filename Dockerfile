FROM thebrightone/gocv

WORKDIR /go/src/app

ENV GO111MODULE=on

ENV skipframes=240

ENV row=50

RUN go get -u gocv.io/x/gocv@v0.26.0

RUN mkdir /outputs

COPY . .

COPY videoplayback.mp4 /out.mp4

RUN go build -o /app .

ENTRYPOINT ["/app"]