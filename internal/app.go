package internal

import (
	"encoding/json"
	"flag"
	"fmt"
	"gitlab.com/back1ng1/prerender-warming/internal/sitemapper"
	"gitlab.com/back1ng1/prerender-warming/internal/warmer"
	"log"
	"os"
)

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
		filename := "e2e-tests-result.json"

		for _, url := range sitemap.URL {
			warm.Add(url.Loc)
		}

		err := warm.Refresh()
		output, _ := json.Marshal(err)
		if len(err) > 0 {
			os.WriteFile(filename, output, 0777)
			os.Exit(1)
		} else {
			os.WriteFile(filename, []byte("[]"), 0777)
		}

		fmt.Println("All pages are success checked.")
		os.Exit(0)
	}
}
