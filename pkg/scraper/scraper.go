package scraper

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"go.uber.org/zap"
)

type Scraper struct {
	Browser *rod.Browser
	Page    *rod.Page
	Logger  *zap.Logger
}

func NewScraper(logger *zap.Logger) *Scraper {
	return &Scraper{
		Logger: logger,
	}
}

func (S *Scraper) NewBrowser() *rod.Browser {
	S.Logger.Debug("Initializing scraper browser.")
	S.Browser = rod.New().MustConnect()
	return S.Browser
}

func (S *Scraper) NewPage(URL string) *rod.Page {
	S.Logger.Debug("Browsing To URL: " + URL)
	S.Page = S.Browser.MustPage(URL)
	S.Page.MustEmulate(devices.Clear)
	return S.Page
}
