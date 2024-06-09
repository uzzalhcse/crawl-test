package config

import "github.com/uzzalhcse/crawl-test/pkg/ninjacrawler"

func loadAppConfig(crawler *ninjacrawler.Crawler) {
	crawler.Config.Add("app", map[string]any{
		"app_env":    crawler.Config.Env("APP_ENV", "local"),
		"user_agent": crawler.Config.Env("USER_AGENT", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"),
	})
}
