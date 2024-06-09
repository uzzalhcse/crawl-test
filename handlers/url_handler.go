package handlers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/uzzalhcse/crawl-test/constant"
	"github.com/uzzalhcse/crawl-test/pkg/ninjacrawler"
	"strings"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	categorySelector := ninjacrawler.UrlSelector{
		Selector:     "div.index.clearfix ul.clearfix li",
		SingleResult: false,
		FindSelector: "a",
		Attr:         "href",
		Handler:      customHandler,
	}
	categoryProductSelector := ninjacrawler.UrlSelector{
		Selector:     "div.index.clearfix ul.clearfix li",
		SingleResult: false,
		FindSelector: "a",
		Attr:         "href",
		Handler: func(collection ninjacrawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
			if strings.Contains(fullUrl, "/tool/product/") {
				return fullUrl, nil
			}
			return "", nil
		},
	}

	crawler.Collection(constant.Categories).CrawlUrls(crawler.GetBaseCollection(), categorySelector)
	crawler.Collection(constant.Products).CrawlUrls(crawler.GetBaseCollection(), categoryProductSelector)
}
func customHandler(collection ninjacrawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
	if strings.Contains(fullUrl, "/tool/category/") {
		return fullUrl, nil
	}
	return "", nil
}
