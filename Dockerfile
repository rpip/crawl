FROM alpine:latest
MAINTAINER Yao Adzaku <yao.adzaku@gmail.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

COPY . /go/src/github.com/rpip/crawl

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		go \
		git \
		gcc \
		libc-dev \
		libgcc \
	&& cd /go/src/github.com/rpip/crawl \
	&& go build -o /usr/bin/crawl . \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

ENTRYPOINT ["crawl"]

CMD ["--help"]
