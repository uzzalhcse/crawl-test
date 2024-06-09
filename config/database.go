package config

import "github.com/uzzalhcse/crawl-test/pkg/ninjacrawler"

func loadDatabaseConfig(crawler *ninjacrawler.Crawler) {
	crawler.Config.Add("database", map[string]any{
		"username": crawler.Config.Env("DB_USERNAME", "ninjacrawler"),
		"password": crawler.Config.Env("DB_PASSWORD", "password"),
		"host":     crawler.Config.Env("DB_HOST", "localhost"),
		"port":     crawler.Config.Env("DB_PORT", "27017"),
	})
}
