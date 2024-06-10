package main

import (
	"github.com/uzzalhcse/crawl-test/config"
	"github.com/uzzalhcse/crawl-test/handlers"
	"github.com/uzzalhcse/crawl-test/pkg/ninjacrawler"
)

// const (
//
//	name = "markt"
//	url  = "https://markt-mall.jp"
//
// )
const (
	name = "kyocera"
	url  = "https://www.kyocera.co.jp/prdct/tool/category/product"
)

func main() {
	crawler := ninjacrawler.NewCrawler(name, url, ninjacrawler.Engine{
		//BoostCrawling:  true,
		BlockResources: true,
		DevCrawlLimit:  2,
		BlockedURLs:    []string{"syncsearch.jp"},
	})
	crawler.Start()
	defer crawler.Stop()
	config.Register(crawler)
	handlers.UrlHandler(crawler)
	handlers.ProductHandler(crawler)
}
