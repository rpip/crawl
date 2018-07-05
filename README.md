[![Build Status](https://travis-ci.org/rpip/crawl.svg?branch=master)](https://travis-ci.org/rpip/crawl)

## Installation

``` shell
$ go get github.com/rpip/crawl
$ crawl --help
Usage of crawl:
  -d int
    	Depth of lookup within page (default 4)
  -u string
    	URL to start crawl on
  -v	Verbose mode (default true)
```


``` shell
$ crawl -u=https://golang.org -d=2 -v=false > sitemap.txt
```

## Docker

```shell
$ make docker
$ docker run -i -t mz/crawl
```
