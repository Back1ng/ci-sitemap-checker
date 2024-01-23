package internal

import (
	"flag"
	"fmt"
	"gitlab.com/back1ng1/prerender-warming/internal/sitemapper"
	"gitlab.com/back1ng1/prerender-warming/internal/warmer"
	"log"
	"os"
	"time"
)

var sleeping time.Duration

func Run() {
	warm := warmer.New()

	flag.String("sitemap", "", "Parsable sitemap.xml")
	flag.Parse()

	sitemap := flag.Lookup("sitemap").Value.String()

	if len(sitemap) == 0 {
		log.Fatal("Sitemap not defined. Add arg: -sitemap https://example.com/sitemap.xml")
		os.Exit(1)
	}

	for {
		sitemapParser := sitemapper.New()
		sitemap := sitemapParser.Get(sitemap)

		for _, url := range sitemap.URL {
			warm.Add(url.Loc)
		}

		err := warm.Refresh()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("All pages are success checked.")
		os.Exit(0)
	}
}
