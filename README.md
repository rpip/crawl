[![Build Status](https://travis-ci.org/rpip/crawl.svg?branch=master)](https://travis-ci.org/rpip/crawl)

## Installation

``` shell
λ git clone github.com/rpip/crawl
λ cd crawl && make deps && make
λ ./crawl --help
Usage of ./crawl:
  -d int
    	Depth of lookup within page (default 4)
  -u string
    	URL to start crawl on (default "https://www.stearsng.com")
```


```shell
λ go run *.go --u=https://jeiwan.cc --d=3
```

## Docker

```shell
λ make docker
λ docker run -i -t mz/crawl
```
