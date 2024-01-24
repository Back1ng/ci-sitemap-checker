package warmer

import (
	"fmt"
	"github.com/gosuri/uilive"
	"net/http"
	"sync"
	"time"
)

type warmer struct {
	Urls map[string]int64
	mu   *sync.Mutex
}

type FailedCheck struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
}

var client http.Client

func New() Warmer {
	client = http.Client{}

	return &warmer{
		Urls: make(map[string]int64),
		mu:   &sync.Mutex{},
	}
}

// Process Perform check on low latency
func (w *warmer) Process(url string) *FailedCheck {
	req := prepareUrl(url)

	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return &FailedCheck{URL: url, StatusCode: resp.StatusCode}
	}
	resp.Body.Close()

	return nil
}

func (w *warmer) Add(url string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, ok := w.Urls[url]
	if !ok {
		w.Urls[url] = time.Now().Unix()
	}
}

func (w *warmer) Refresh() []FailedCheck {
	writer := uilive.New()
	writer.Start()

	counter := 1

	var failed []FailedCheck

	for url, _ := range w.Urls {
		fmt.Fprintf(writer, "Checking [%d/%d]\n", counter, len(w.Urls))
		err := w.Process(url)
		if err != nil {
			failed = append(failed, *err)
		}
		counter++
	}

	return failed
}

func prepareUrl(url string) *http.Request {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")

	return req
}
