package ninjacrawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

// GetPlaywright initializes and runs the Playwright framework.
// It returns a Playwright instance if successful, otherwise returns an error.
func GetPlaywright() (*playwright.Playwright, error) {
	err := playwright.Install()
	if err != nil {
		return nil, err
	}
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	return pw, nil
}

// GetBrowserPage launches a browser instance and creates a new page using the Playwright framework.
// It supports Chromium, Firefox, and WebKit browsers, and can configure Proxy settings if provided.
// It returns the browser and page instances, or an error if the operation fails.
func GetBrowserPage(pw *playwright.Playwright, browserType string, proxy Proxy) (playwright.Browser, playwright.Page, error) {
	var browser playwright.Browser
	var err error

	var browserTypeLaunchOptions playwright.BrowserTypeLaunchOptions
	browserTypeLaunchOptions.Headless = playwright.Bool(!isLocalEnv())
	browserTypeLaunchOptions.Devtools = playwright.Bool(isLocalEnv())

	if len(app.engine.ProxyServers) > 0 {
		// Set Proxy options
		browserTypeLaunchOptions.Proxy = &playwright.Proxy{
			Server:   proxy.Server,
			Username: playwright.String(proxy.Username),
			Password: playwright.String(proxy.Password),
		}
	}
	switch browserType {
	case "chromium":
		browser, err = pw.Chromium.Launch(browserTypeLaunchOptions)
	case "firefox":
		browser, err = pw.Firefox.Launch(browserTypeLaunchOptions)
	case "webkit":
		browser, err = pw.WebKit.Launch(browserTypeLaunchOptions)
	default:
		return nil, nil, fmt.Errorf("unsupported browser type: %s", browserType)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: playwright.String(app.Config.GetString("USER_AGENT")),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create page: %w", err)
	}

	// Conditionally intercept and block resources based on configuration
	if app.engine.BlockResources {
		err := page.Route("**/*", func(route playwright.Route) {
			req := route.Request()
			resourceType := req.ResourceType()
			url := req.URL()

			// Check if the resource should be blocked based on resource type or URL
			if shouldBlockResource(resourceType, url) {
				route.Abort()
			} else {
				route.Continue()
			}
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to set up request interception: %w", err)
		}
	}

	return browser, page, nil
}

// NavigateToURL navigates to a specified URL using the given Playwright page.
// It waits until the page is fully loaded and returns a goquery document representing the DOM.
// If navigation fails, it logs the page content to a file and returns an error.
func NavigateToURL(page playwright.Page, url string) (*goquery.Document, error) {
	waitUntil := playwright.WaitUntilStateDomcontentloaded
	if app.engine.IsDynamic {
		waitUntil = playwright.WaitUntilStateNetworkidle
	}

	_, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: waitUntil,
	})
	if err != nil {
		app.Logger.Html(page, err.Error())
		return nil, err
	}
	return GetPageDom(page)
}
