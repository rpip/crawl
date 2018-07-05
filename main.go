package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/PuerkitoBio/purell"
)

func main() {

	depth := flag.Int("d", 4, "Depth of lookup within page")
	// also try https://stratechery.com
	baseURL := flag.String("u", "https://www.stearsng.com", "URL to start crawl on")
	flag.Parse()

	startURL, err := url.Parse(*baseURL)
	if err != nil {
		log.Fatal("Invalid url. Please check and try again", err)
		os.Exit(1)
	}

	var responseCh = make(chan *Page)
	c := NewCrawler(startURL, responseCh)
	go c.crawl(startURL, *depth)
	for page := range responseCh {
		fmt.Printf("found: %s \n %s \n %q\n", page.title, page.url, page.urls)
	}
	close(responseCh)
}

func (c *Crawler) crawl(link *url.URL, depth int) {
	var wg sync.WaitGroup

	if c.IsProcessed(link.String()) || depth <= 0 {
		return
	}
	page, err := c.Fetch(link)
	if err != nil {
		log.Fatal(err)
	}
	c.responseCh <- page

	for _, u := range page.urls {
		wg.Add(1)
		go func(u *url.URL) {
			defer wg.Done()
			c.crawl(u, depth-1)
		}(u)
	}
	wg.Wait()
	return
}

// utils
func normalizeURL(u *url.URL) *url.URL {
	flags := purell.FlagsUsuallySafeGreedy | purell.FlagRemoveFragment | purell.FlagRemoveDuplicateSlashes
	s := purell.NormalizeURL(u, flags)
	u, _ = url.Parse(s)
	return u
}
