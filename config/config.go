package config

import "github.com/uzzalhcse/crawl-test/pkg/ninjacrawler"

func Register(crawler *ninjacrawler.Crawler) {
	loadAppConfig(crawler)
}
