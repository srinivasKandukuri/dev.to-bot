package artifact

import (
	"fmt"
	"github.com/dev.to-bot/pkg/config"
	"github.com/dev.to-bot/pkg/scraper"
	"github.com/dev.to-bot/pkg/types"
	"github.com/dev.to-bot/pkg/utils"
	"github.com/go-rod/rod/lib/input"
	"strings"
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

	for _, tag := range config.Tags {
		a.Scraper.Logger.Debug(string("Started Tag : " + tag))
		url := fmt.Sprintf(config.Urls["top"], tag)
		a.Scraper.NewPage(url)
		posts := a.GetTopPosts(4)
		title := utils.GenerateRandomTitle(config.Titles)
		userIds := a.addNewPost(posts, tag, title)
		utils.WaitToEnd(5)
		a.AddComments(userIds)
		a.Scraper.Logger.Debug("****** Completed post successfully ******")
	}
}

func (a Artifact) LoginUser() {
	a.Scraper.NewPage(a.Urls["login"])
	a.Scraper.Page.MustElement("[type='email']").MustInput(a.UserName).MustType(input.Tab)
	a.Scraper.Page.MustElement("[type='password']").MustInput(a.Password).MustType(input.Enter)
	utils.AssertEquals(a.Scraper.Page.MustElement("nav header .crayons-subtitle-3").MustText(), "My Tags", "User not able to login")
	a.Scraper.Logger.Debug("User Logged IN")
}

func (a Artifact) GetTopPosts(topPostsCount int) types.Posts {
	utils.WaitToEnd(2)
	elements := a.Scraper.Page.MustElements(".crayons-story")
	myMap := make(types.Posts)
	for i, ele := range elements {
		fmt.Printf("Collecting data from element : %v\n", i)
		e1 := ele.MustElements(".crayons-story__top .crayons-story__author-pic a")
		user, _ := e1[0].Attribute("href")
		if len(e1) > 1 {
			user, _ = e1[1].Attribute("href")
		}
		userId := fmt.Sprintf("%s", strings.TrimPrefix(*user, `/`))
		title := ele.MustElement("a").MustText()
		link, _ := ele.MustElement("a").Attribute("href")
		fullLink := fmt.Sprintf("%s%s", config.Urls["linkPrefix"], *link)
		description := a.Scraper.Browser.MustPage(fullLink).MustWaitLoad().MustElement("#article-body").MustText()
		description = strings.TrimRight(description, "\r\n")
		description = utils.EllipticalTruncate(description, 200)
		p := types.Post{
			Title:       title,
			Link:        fullLink,
			PostLink:    fmt.Sprintf("%s%s%s", "{% link https://dev.to", *link, " %}"),
			Description: description,
			UserId:      userId,
		}
		myMap[i] = p
		if i == topPostsCount {
			break
		}

	}
	a.Scraper.Logger.Debug("Scraping data from tag page is done.")
	return myMap
}

func (a Artifact) addNewPost(posts types.Posts, tag types.Tag, title string) string {
	a.Scraper.Logger.Debug("Navigating to NewPost")
	a.Scraper.NewPage(config.Urls["createPost"])
	utils.WaitToEnd(1)
	pageTitle := a.Scraper.Page.MustElement("#article-form-title")
	utils.AssertEquals(pageTitle.MustText(), "New post title here...", "The title is not present. Perhaps the homepage is not opening")

	textarea := a.Scraper.Page.MustElements("textarea")
	for _, b := range textarea {
		t, err := b.Text()
		if err != nil {
			// error handling
		}
		if t == "New post title here..." {
			title = fmt.Sprintf(title, tag)
			b.MustInput(title)
			a.Scraper.Logger.Debug("Title entered")
			break
		}
	}
	utils.WaitToEnd(1)
	a.Scraper.Logger.Debug("Preparing Body")
	body := ""
	ids := ""
	postsLen := len(posts)
	counter := 1
	for _, v := range posts {
		if postsLen == counter {
			ids += "and @" + v.UserId
		} else {
			ids += "@" + v.UserId + " , "
		}
		body += "##" + v.Title + "\n"
		body += v.Description + "\n"
		body += v.PostLink + "\n"
		counter++
	}
	a.Scraper.Page.MustElement("#article_body_markdown").MustInput(body)
	a.Scraper.Logger.Debug("Body entered")
	utils.WaitToEnd(1)
	a.Scraper.Logger.Debug("Looking for Publish button")

	a.Scraper.Page.MustElement("#tag-input").MustInput(string(tag)).MustType(input.Tab).MustInput("c4r4x35").MustType(input.Tab)
	buttons := a.Scraper.Page.MustElements("button")
	for _, b := range buttons {
		t, err := b.Text()
		if err != nil {
			// error handling
		}
		if t == "Publish" {
			a.Scraper.Logger.Debug("Found Publish button")
			b.MustClick().MustWaitLoad()
			a.Scraper.Logger.Debug("Post Published.")
			break
		}
	}
	return ids
}

func (a Artifact) AddComments(userIds string) {
	a.Scraper.Logger.Debug("Adding comment to the published post.")
	a.Scraper.Logger.Debug("Navigating to Dashboard.")
	dashboardPage := a.Scraper.Browser.MustPage(config.Urls["dashboard"]).MustWaitLoad()
	title := dashboardPage.MustElement(".crayons-title").MustText()
	utils.AssertEquals(title, "Dashboard", "Dashboard title not matched, navigation into dashboard page went wrong.")
	elem := dashboardPage.MustElements(".crayons-layout__content .crayons-card .dashboard-story.js-dashboard-story.spec__dashboard-story.single-article")[0]
	elem.MustElement(".dashboard-story__title a").MustClick()
	utils.WaitToEnd(3)
	comment := fmt.Sprintf("Shoutout to all the awesome authors featured in this years's Top 5 in 2022: %s 🙌.", userIds)
	fmt.Printf("adding comment Message : %s", comment)
	dashboardPage.MustElement("#comments #comments-container #text-area").MustInput(comment)
	//dashboardPage.MustSearch("Submit").MustClick()
	elm := dashboardPage.MustElement("#comments #comments-container #new_comment .comment-form__buttons button.crayons-btn")
	el := elm.MustAttribute("disabled")
	if el == nil {
		elm.MustClick()
		utils.WaitToEnd(2)
		a.Scraper.Logger.Debug("comment submitted successfully.")
	}
}
