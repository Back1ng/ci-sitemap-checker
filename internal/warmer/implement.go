package warmer

import (
	"errors"
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

var client http.Client

func New() Warmer {
	client = http.Client{}

	return &warmer{
		Urls: make(map[string]int64),
		mu:   &sync.Mutex{},
	}
}

// Process Perform check on low latency
func (w *warmer) Process(url string) error {
	req := prepareUrl(url)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return errors.New(fmt.Sprintf("Page: %v. Status code: %v", url, resp.StatusCode))
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

func (w *warmer) Refresh() error {
	writer := uilive.New()
	writer.Start()

	counter := 1

	for url, _ := range w.Urls {
		fmt.Fprintf(writer, "Checking [%d/%d]\n", counter, len(w.Urls))
		err := w.Process(url)
		if err != nil {
			return err
		}
		counter++
	}

	return nil
}

func prepareUrl(url string) *http.Request {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		fmt.Println(err)
	}

	return req
}
