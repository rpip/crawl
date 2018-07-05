package crawler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const noLinksPage = `
<!DOCTYPE html>
<html>
<head>
<title>No links</title>
</head>
<body>

<h1>This is a Heading</h1>
<p>This is a paragraph.</p>

</body>
</html>
`

const homePage = `
<!DOCTYPE html>
<html>
<head>
<title>Welcome</title>
</head>
<body>

<h1>This is a Heading</h1>
<p>This is a paragraph.</p>
<a href="/about">About</a>
<a href="/help">Help</a>
<a href="https://www.w3schools.com/html/default.asp">Learn HTML</a>
<a href="http://google.com">Google</a>
<a href="http://github.com">Github</a>
<a href="https://monzo.com">Monzo</a>
</body>
</html>
`

const aboutPage = `
<!DOCTYPE html>
<html>
<head>
<title>About</title>
</head>
<body>

<h1>This is a Heading</h1>
<p>This is a paragraph.</p>
<a href="/about">About</a>

</body>
</html>
`

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(homePage))
	})
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(aboutPage))
	})
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(noLinksPage))
	})

	return httptest.NewServer(mux)
}

func TestFetch(t *testing.T) {
	testServer := newTestServer()

	var responseCh = make(chan *Page)
	defer close(responseCh)

	url, _ := url.Parse(testServer.URL)
	// TODO: make responseCh optional
	c := NewCrawler(url, responseCh, true)

	// TODO: support string type url
	page, err := c.Fetch(url)
	if err != nil {
		t.Fatal(err)
	}
	if page == nil {
		t.Fatal("expected page, got nil")
	}
	if len(page.URLs) != 2 {
		t.Errorf("expected 2 page links got %d", len(page.URLs))
	}

	// no page links
	url, _ = url.Parse("/test")
	page, err = c.Fetch(url)
	if len(page.URLs) != 0 {
		t.Errorf("expected zero page links, got %d", len(page.URLs))
	}

	// already visisted url
	page, err = c.Fetch(url)
	if err != ErrVisitedRetry {
		t.Errorf("Expected error, already seen URL. got %v", page)
	}

	// check is processed
	if !c.IsProcessed("/test") {
		t.Errorf("Expected /test to be marked as visited. Not marked")
	}
	testServer.Close()
}
