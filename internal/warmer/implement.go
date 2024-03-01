package warmer

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	"github.com/gosuri/uilive"
)

type warmer struct {
	client http.Client
	writer *uilive.Writer
	mu     *sync.Mutex

	url <-chan string
	err chan FailedCheck
}

type FailedCheck struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
}

func New(url <-chan string, errs chan FailedCheck) Warmer {
	writer := uilive.New()
	writer.Start()

	return &warmer{
		client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		mu:     &sync.Mutex{},
		writer: writer,
		url:    url,
		err:    errs,
	}
}

// Process Perform check on low latency
func (w *warmer) Process(url string) *FailedCheck {
	req := prepareUrl(url)

	resp, _ := w.client.Do(req)
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return &FailedCheck{URL: url, StatusCode: resp.StatusCode}
	}
	resp.Body.Close()

	return nil
}

func (w *warmer) Refresh() {
	for url := range w.url {
		//fmt.Fprintf(w.writer, "Urls left: %d\n", len(w.url))

		fmt.Printf("Process %s...\n", url)

		if err := w.Process(url); err != nil {
			w.err <- *err
		}
	}

	w.writer.Flush()
}

func prepareUrl(url string) *http.Request {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")

	return req
}
