package artifact

import (
	"github.com/dev.to-bot/pkg/scraper"
	"github.com/dev.to-bot/pkg/types"
	"github.com/dev.to-bot/pkg/utils"
	"github.com/go-rod/rod/lib/input"
)

type Artifact struct {
	Scraper  *scraper.Scraper
	UserName string
	Password string
	Tags     types.Tags
	Titles   types.Titles
	Urls     types.Urls
}

type Option func(*Artifact)

func NewArtifact(scraper *scraper.Scraper, username string, password string, tags types.Tags, titles types.Titles,
	urls types.Urls, opts ...Option) *Artifact {
	artifact := &Artifact{
		Scraper:  scraper,
		UserName: username,
		Password: password,
		Tags:     tags,
		Titles:   titles,
		Urls:     urls,
	}
	for _, o := range opts {
		o(artifact)
	}
	return artifact
}

func (a Artifact) RunArtifact() {
	a.Scraper.NewBrowser()
	a.Scraper.Browser.MustPage(a.Urls["login"])
	a.LoginUser()
}

func (a Artifact) LoginUser() {
	a.Scraper.NewPage(a.Urls["login"])
	a.Scraper.Page.MustElement("[type='email']").MustInput(a.UserName).MustType(input.Tab)
	a.Scraper.Page.MustElement("[type='password']").MustInput(a.Password).MustType(input.Enter)
	utils.AssertEquals(a.Scraper.Page.MustElement("nav header .crayons-subtitle-3").MustText(), "My Tags", "User not able to login")
	a.Scraper.Logger.Debug("User Logged IN")
}
