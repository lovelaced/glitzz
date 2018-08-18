FROM golang:alpine
RUN apk add --no-cache git make
WORKDIR /go/src/glitzz
RUN mkdir /data
RUN git clone https://github.com/lovelaced/glitzz.git .
RUN make
ENTRYPOINT /bin/sh -c 'glitzz default_config > /data/config.json && glitzz run /data/config.json'