package main

import (
	"fmt"
	"github.com/dev.to-bot/pkg/config"
	"github.com/dev.to-bot/pkg/types"
	"github.com/dev.to-bot/pkg/utils"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"os"
	"strings"
	"time"
)

func UserLogin(page *rod.Page, email string, pass string) {
	page.MustElement("[type='email']").MustInput(email).MustType(input.Tab)
	page.MustElement("[type='password']").MustInput(pass).MustType(input.Enter)
	utils.AssertEquals(page.MustElement("nav header .crayons-subtitle-3").MustText(), "My Tags", "User not able to login")
	fmt.Println("User Logged IN")
}

func OpenTagPage(browser *rod.Browser, tag string, url string) *rod.Page {
	page := browser.MustPage(url)
	fmt.Printf("Navigated to tag : %s , url : %s\n", tag, url)
	return page
}

func GetTopPosts(browser *rod.Browser, page *rod.Page, topPostsCount int) types.Posts {
	utils.WaitToEnd(2)
	elements := page.MustElements(".crayons-story")
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
		description := browser.MustPage(fullLink).MustWaitLoad().MustElement("#article-body").MustText()
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
	fmt.Println("Scraping data from tag page is done.")
	return myMap
}

func addNewPost(browser *rod.Browser, posts types.Posts, tag string, title string) string {
	fmt.Println("Navigating to NewPost")
	newPage := browser.MustPage(config.Urls["createPost"]).MustWaitLoad()
	utils.WaitToEnd(2)
	fmt.Println("Page opened")
	pageTitle := newPage.Timeout(5 * time.Second).MustElement("#article-form-title")
	utils.AssertEquals(pageTitle.MustText(), "New post title here...", "The title is not present. Perhaps the homepage is not opening")

	textarea := newPage.Timeout(5 * time.Second).MustElements("textarea")
	for _, b := range textarea {
		t, err := b.Text()
		if err != nil {
			// error handling
		}
		if t == "New post title here..." {
			title = fmt.Sprintf(title, tag)
			b.MustInput(title)
			fmt.Println("Title Entered")
			break
		}
	}
	utils.WaitToEnd(2)
	fmt.Println("Preparing Body")
	body := ""
	ids := ""
	for _, v := range posts {
		ids += "@" + v.UserId + " , "
		body += "##" + v.Title + "\n"
		body += v.Description + "\n"
		body += v.PostLink + "\n"
	}
	ids += "& @c4r4x35 ."
	newPage.Timeout(5 * time.Second).MustElement("#article_body_markdown").MustInput(body)
	fmt.Println("Body entered")
	utils.WaitToEnd(2)
	fmt.Println("Looking for Publish button")
	newPage.MustElement("#tag-input").MustInput(tag).MustType(input.Tab).MustInput("c4r4x35").MustType(input.Tab)
	buttons := newPage.MustElements("button")
	for _, b := range buttons {
		t, err := b.Text()
		if err != nil {
			// error handling
		}
		if t == "Publish" {
			fmt.Println("Found Publish button")
			b.MustClick().MustWaitLoad()
			fmt.Println("Clicking on Publish")
			fmt.Printf("============================END tag post : %s=========================\n", tag)
			break
		}
	}
	return ids
}

func AddComments(browser *rod.Browser, userIds string) {
	fmt.Println("Adding comment to the published post.")
	fmt.Println("Navigating to Dashboard.")
	dashboardPage := browser.MustPage(config.Urls["dashboard"]).MustWaitLoad()
	title := dashboardPage.MustElement(".crayons-title").MustText()
	utils.AssertEquals(title, "Dashboard", "Dashboard title not matched, navigation into dashboard page went wrong.")
	elem := dashboardPage.MustElements(".crayons-layout__content .crayons-card .dashboard-story.js-dashboard-story.spec__dashboard-story.single-article")[0]
	elem.MustElement(".dashboard-story__title a").MustClick()
	utils.WaitToEnd(3)
	comment := fmt.Sprintf("Awesome articles by authors: %s  🙌", userIds)
	fmt.Printf("adding comment Message : %s", comment)
	dashboardPage.MustElement("#comments #comments-container #text-area").MustInput(comment)
	//dashboardPage.MustSearch("Submit").MustClick()
	elm := dashboardPage.MustElement("#comments #comments-container #new_comment .comment-form__buttons button.crayons-btn")
	el := elm.MustAttribute("disabled")
	if el == nil {
		elm.MustClick()
		utils.WaitToEnd(2)
		fmt.Println("comment submitted successfully.")
	}
}

func main() {
	fmt.Printf("Target website %s\n", config.Urls["target"])
	userEmail := os.Getenv("USER_EMAIL")
	userPassword := os.Getenv("USER_PASSWORD")
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page := browser.MustPage(config.Urls["login"])
	UserLogin(page, userEmail, userPassword)

	for _, tag := range config.Tags {
		fmt.Printf("******Started %s post ******", tag)
		url := fmt.Sprintf(config.Urls["top"], tag)
		tagPage := OpenTagPage(browser, tag, url)
		posts := GetTopPosts(browser, tagPage, 4)
		title := utils.GenerateRandomTitle(config.Titles)
		userIds := addNewPost(browser, posts, tag, title)
		utils.WaitToEnd(5)
		AddComments(browser, userIds)
		fmt.Printf("******Completed %s post successfully******", tag)
	}

}
