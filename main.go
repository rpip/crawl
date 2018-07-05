package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/rpip/crawl/crawler"
)

func main() {

	depth := flag.Int("d", 4, "Depth of lookup within page")
	verbose := flag.Bool("v", true, "Verbose mode")
	baseURL := flag.String("u", "", "URL to start crawl on")
	flag.Parse()

	startURL, err := url.Parse(*baseURL)
	if err != nil {
		log.Fatalf("failed to crawl %s: %v", startURL, err)
	}
	var responseCh = make(chan *crawler.Page)
	defer close(responseCh)

	c := crawler.NewCrawler(startURL, responseCh, *verbose)
	go c.Crawl(startURL, *depth)

	for page := range responseCh {
		print(page)
	}
}

func print(page *crawler.Page) {
	fmt.Print(strings.Repeat(" ", page.Indent))
	fmt.Printf("%s \"%s\"\n", page.URL.Path, page.Title)
	for _, u := range page.URLs {
		fmt.Print(strings.Repeat(" ", page.Indent+1))
		fmt.Printf("%s \n", u.Path)
	}
}
