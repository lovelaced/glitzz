FROM golang:alpine
RUN apk add --no-cache git make
WORKDIR /go/src/glitzz
RUN mkdir /data
RUN git clone https://github.com/lovelaced/glitzz.git .
RUN make
VOLUME ./_data:data
ONBUILD RUN /bin/sh -c 'glitzz default_config > /data/config.json'
ENTRYPOINT /bin/sh -c 'glitzz run /data/config.json'
