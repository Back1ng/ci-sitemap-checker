package internal

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gitlab.com/back1ng1/prerender-warming/internal/sitemapper"
	"gitlab.com/back1ng1/prerender-warming/internal/warmer"
)

func Run() {
	filename := "e2e-tests-result.json"

	var url = flag.String("url", "https://example.com/sitemap.xml", "Sitemap that will be parsed.")
	var threads = flag.Int("threads", 2, "Count of threads to warm prerender.")
	flag.Parse()

	if *threads < 1 {
		log.Fatal("Count of threads cannot be less then 1.")
	}

	sitemapParser := sitemapper.New()
	sitemapLinksStream := make(chan string, 10000)
	errorsStream := make(chan warmer.FailedCheck, 1)

	var wg sync.WaitGroup

	sitemap := sitemapParser.Get(*url)
	wg.Add(len(sitemap.URL))

	go func() {
		for _, url := range sitemap.URL {
			sitemapLinksStream <- url.Loc
			wg.Done()
		}

		close(sitemapLinksStream)
	}()

	var errs []warmer.FailedCheck

	go func() {
		for err := range errorsStream {
			errs = append(errs, err)
		}
	}()

	warm := warmer.New(sitemapLinksStream, errorsStream)

	for i := 0; i < *threads; i++ {
		go warm.Refresh()
	}

	wg.Wait()
	for len(sitemapLinksStream) > 0 {
		<-time.After(time.Second)
	}

	close(errorsStream)

	for _, err := range errs {
		fmt.Printf("Page: %v. Status: %d\n", err.URL, err.StatusCode)
	}

	if len(errs) > 0 {
		output, _ := json.Marshal(errs)
		os.WriteFile(filename, output, 0777)
		os.Exit(1)
	}

	fmt.Println("All pages are success checked.")
	os.Exit(0)
}
