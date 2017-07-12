FROM ubuntu:17.04

RUN apt-get update \
	&& apt-get install --no-install-recommends --no-install-suggests -y golang git npm ca-certificates tzdata \
	&& rm -rf /var/lib/apt/lists/*

# forward nginx logs to docker logs
# RUN ln -sf /dev/stdout /var/log/nginx/access.log \
#	&& ln -sf /dev/stderr /var/log/nginx/error.log

RUN ln -s /usr/bin/nodejs /usr/bin/node

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src/github.com/wtg/shuttletracker" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/github.com/wtg/shuttletracker
RUN go get -u github.com/kardianos/govendor

# ADD ./package.json /app
RUN npm install -g bower

COPY ./bower.json .
RUN bower install --allow-root

COPY . .

RUN govendor sync
RUN go build -o shuttletracker

CMD ["./shuttletracker"]
