package crawler

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/purell"
	radix "github.com/armon/go-radix"
)

const (
	// defaultHTTPTimeout is the default timeout on the http client
	defaultHTTPTimeout = 60 * time.Second

	// User agent
	userAgent = "mz-crawl/2018-07-04"
)

// ErrVisitedRetry is the error thrown if url has already been visited
var ErrVisitedRetry = errors.New("URL already visited")

type responseCh chan *Page

// Page describe a crawled page
type Page struct {
	URL    *url.URL
	Title  string
	URLs   []*url.URL
	Indent int
}

// Fetcher is the interface that wraps the crawling mechanism
// Implements Fetch and IsProcessed methods
type Fetcher interface {
	// Fetch returns the url and URLs found on that page.
	Fetch(url *url.URL) (page *Page, err error)

	// IsProcessed returns true if URL has already been visited
	IsProcessed(url string) bool
}

// Crawler does the actual crawling of the page
type Crawler struct {
	sync.Mutex
	log        *log.Logger
	host       string
	responseCh responseCh
	client     *http.Client
	startURL   *url.URL
	seenURLs   *radix.Tree
	verbose    bool
}

// NewCrawler creates a new crawler with the given base url, crawl depth,
// and HTTP client, allowing overriding of the HTTP client to use
func NewCrawler(url *url.URL, responseCh responseCh, verbose bool) *Crawler {
	normalized := normalizeURL(url)
	c := &Crawler{
		startURL:   normalized,
		log:        log.New(os.Stderr, "", log.LstdFlags),
		host:       normalized.Hostname(),
		responseCh: responseCh,
		client: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
		seenURLs: radix.New(),
		verbose:  verbose,
	}
	return c
}

// Fetch returns scraped page. Page contains title, URLs found
func (c *Crawler) Fetch(url *url.URL) (*Page, error) {
	if c.IsProcessed(url.Path) {
		return nil, ErrVisitedRetry
	}
	// mark as seen
	c.Lock()
	c.seenURLs.Insert(url.Path, struct{}{})
	c.Unlock()

	start := time.Now()
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		c.debugf("Error creating http request: %v\n", err)
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	c.debugf("Completed in %v\n", time.Since(start))
	defer resp.Body.Close()
	page := c.processPage(resp)
	return page, err
}

// IsProcessed returns True if url has already been visited
func (c *Crawler) IsProcessed(url string) bool {
	c.Lock()
	defer c.Unlock()
	_, exists := c.seenURLs.Get(url)
	return exists
}

func (c *Crawler) processPage(resp *http.Response) *Page {
	doc, err := goquery.NewDocumentFromResponse(resp)
	title := doc.Find("title").Contents().Text()

	if err != nil {
		log.Fatalf("couldn't parse page %s", err)
	}

	// base URL to use for all relative URLs contained within the page
	var base *url.URL
	if baseURL, _ := doc.Find("base[href]").Attr("href"); baseURL != "" {
		base, _ = url.Parse(baseURL)
	} else {
		base = doc.Url
	}

	// gather all links on page
	urls := doc.Find("a[href]").Map(func(_ int, item *goquery.Selection) string {
		link, _ := item.Attr("href")
		return link
	})

	// resolve urls to absolute uri from the base url
	var result []*url.URL
	seen := radix.New()

	for _, u := range urls {
		if uu, err := url.Parse(u); err == nil {
			nu := normalizeURL(base.ResolveReference(uu))

			if nu.Hostname() == c.host {
				// first time seeing this
				if !c.IsProcessed(nu.Path) {
					// collect only unique links on page
					if _, ok := seen.Get(nu.Path); !ok {
						seen.Insert(nu.Path, struct{}{})
						result = append(result, nu)
					}
				}
			} else {
				c.debugf("ignored: external link, %s", nu)
			}
		}
	}

	return &Page{Title: title, URL: doc.Url, URLs: result}
}

// Crawl recursively scans the url and extract the children links
func (c *Crawler) Crawl(uri *url.URL, depth int) {
	var wg sync.WaitGroup

	if c.IsProcessed(uri.Path) || depth <= 0 {
		return
	}
	page, err := c.Fetch(uri)
	if err == ErrVisitedRetry {
		return
	}

	if err != nil {
		log.Fatal(err)
	}
	// mark indentation for printing
	page.Indent++
	c.responseCh <- page

	for _, u := range page.URLs {
		wg.Add(1)
		go func(u *url.URL) {
			defer wg.Done()
			c.Crawl(u, depth-1)
		}(u)
	}
	wg.Wait()
}

func (c *Crawler) debugf(format string, args ...interface{}) {
	if c.verbose {
		c.log.Printf(format, args...)
	}
}

// utils
func normalizeURL(u *url.URL) *url.URL {
	flags := purell.FlagsUsuallySafeGreedy | purell.FlagRemoveFragment | purell.FlagRemoveDuplicateSlashes
	s := purell.NormalizeURL(u, flags)
	u, _ = url.Parse(s)
	return u
}
