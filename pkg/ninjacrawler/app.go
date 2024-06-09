package ninjacrawler

import (
	"github.com/playwright-community/playwright-go"
	"sync"
	"time"
)

var Once sync.Once
var app *Crawler
var startTime time.Time

const (
	baseCollection = "sites"
)

type Crawler struct {
	*client
	Config                *configService
	Name                  string
	Url                   string
	BaseUrl               string
	pw                    *playwright.Playwright
	collection            string
	UrlSelectors          []UrlSelector
	ProductDetailSelector ProductDetailSelector
	engine                *Engine
	Logger                *defaultLogger
}

func NewCrawler(name, url string, engines ...Engine) *Crawler {
	Once.Do(func() {
		startTime = time.Now()
		logger := newDefaultLogger(name)
		logger.Info("Program started! 🚀")

		// Create default engine configuration
		defaultEngine := Engine{
			BrowserType:     "chromium",
			ConcurrentLimit: 10,
			IsDynamic:       false,
			DevCrawlLimit:   10,
			BlockResources:  false,
			BlockedURLs: []string{
				"www.googletagmanager.com",
				"google.com",
				"googleapis.com",
				"gstatic.com",
			},
			BoostCrawling: false,
			ProxyServers:  []Proxy{},
		}

		// Override defaults with provided engine configuration if available
		if len(engines) > 0 {
			eng := engines[0]
			if eng.BrowserType != "" {
				defaultEngine.BrowserType = eng.BrowserType
			}
			if eng.ConcurrentLimit > 0 {
				defaultEngine.ConcurrentLimit = eng.ConcurrentLimit
			}
			if eng.IsDynamic {
				defaultEngine.IsDynamic = eng.IsDynamic
			}
			if eng.DevCrawlLimit > 0 {
				defaultEngine.DevCrawlLimit = eng.DevCrawlLimit
			}
			if eng.BlockResources {
				defaultEngine.BlockResources = eng.BlockResources
			}
			if eng.BoostCrawling {
				defaultEngine.BoostCrawling = eng.BoostCrawling
				defaultEngine.ProxyServers = eng.getProxyList()
			}
			if len(eng.ProxyServers) > 0 {
				defaultEngine.ProxyServers = eng.ProxyServers
			}
			defaultEngine.BlockedURLs = append(defaultEngine.BlockedURLs, eng.BlockedURLs...)
		}

		app = &Crawler{
			Name:       name,
			Url:        url,
			BaseUrl:    getBaseUrl(url),
			collection: app.GetBaseCollection(),
			engine:     &defaultEngine,
			Logger:     logger,
			Config:     newConfig(),
		}
	})

	return app
}

func (a *Crawler) Start() {
	client := connectDB()
	client.newSite()
	pw, err := GetPlaywright()
	if err != nil {
		a.Logger.Fatal("failed to initialize playwright: %v\n", err)
	}

	a.client = client
	a.pw = pw
}

func (a *Crawler) Stop() {
	if a.pw != nil {
		a.pw.Stop()
	}
	if a.client != nil {
		a.client.close()
	}
	duration := time.Since(startTime)
	a.Logger.Info("Program stopped in ⚡ %v", duration)
}

func (a *Crawler) Collection(collection string) *Engine {
	a.collection = collection
	return a.engine
}

func (a *Crawler) GetCollection() string {
	return a.collection
}

func (a *Crawler) GetBaseCollection() string {
	return baseCollection
}
