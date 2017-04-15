FROM ubuntu:17.04

RUN apt-get update \
	&& apt-get install --no-install-recommends --no-install-suggests -y golang git npm ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

# forward nginx logs to docker logs
# RUN ln -sf /dev/stdout /var/log/nginx/access.log \
#	&& ln -sf /dev/stderr /var/log/nginx/error.log

RUN ln -s /usr/bin/nodejs /usr/bin/node

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/shuttle_tracking_2

# ADD ./package.json /app
RUN npm install -g bower


COPY ./bower.json $GOPATH/src/shuttle_tracking_2
RUN bower install --allow-root

COPY . $GOPATH/src/shuttle_tracking_2

RUN go get

EXPOSE 8080
CMD ["go", "run", "main.go"]
