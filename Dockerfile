FROM golang:1.11

RUN go get -d github.com/spf13/cobra/cobra && \
    cd /go/src/github.com/spf13/cobra && \
    git checkout d2d81d9a96e23f0255397222bb0b4e3165e492dc && \
    go install github.com/spf13/cobra/cobra

WORKDIR /go/src/github.com/hiroara/bunshin-no-jutsu
