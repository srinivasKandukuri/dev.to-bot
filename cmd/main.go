package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"os"
	"strings"
	"time"
	"unicode"
)

type Post struct {
	Title       string
	Link        string
	PostLink    string
	Description string
}

var posts map[string]Post

func UserLogin(page *rod.Page, email string, pass string) {
	page.MustElement("[type='email']").MustInput(email).MustType(input.Tab)
	page.MustElement("[type='password']").MustInput(pass).MustType(input.Enter)
	fmt.Println("User Logged IN")
}

func OpenTagPage(browser *rod.Browser, tag string, url string) *rod.Page {
	page := browser.MustPage(url)
	fmt.Printf("Navigated to tag : %s , url : %s\n", tag, url)
	return page
}

func GetTopPosts(browser *rod.Browser, page *rod.Page, topPostsCount int) (listOfPosts map[int]Post) {
	time.Sleep(2 * time.Second)
	elements := page.MustElements(".crayons-story")
	myMap := make(map[int]Post)
	for i, ele := range elements {
		fmt.Printf("Collecting data from element : %s\n", i)
		title := ele.MustElement("a").MustText()
		link, _ := ele.MustElement("a").Attribute("href")
		fullLink := fmt.Sprintf("%s%s", "https://dev.to", *link)
		description := browser.MustPage(fullLink).MustWaitLoad().MustElement("#article-body").MustText()
		description = strings.TrimRight(description, "\r\n")
		description = EllipticalTruncate(description, 200)
		time.Sleep(2 * time.Second)
		p := Post{
			Title:       title,
			Link:        fullLink,
			PostLink:    fmt.Sprintf("%s%s%s", "{% link https://dev.to", *link, " %}"),
			Description: description,
		}
		myMap[i] = p
		if i == topPostsCount {
			break
		}
	}
	fmt.Println("Scraping data from tag page is done.")
	return myMap
}

func EllipticalTruncate(text string, maxLen int) string {
	lastSpaceIx := maxLen
	len := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len > maxLen {
			return text[:lastSpaceIx] + "..."
		}
	}
	// If here, string is shorter or equal to maxLen
	return text
}

func addNewPost(browser *rod.Browser, posts map[int]Post, tag string) {
	fmt.Println("Navigating to NewPost")
	url := "https://dev.to/new"
	page := browser.MustPage(url)
	time.Sleep(2 * time.Second)
	fmt.Println("Page opened")
	textarea := page.Timeout(5 * time.Second).MustElements("textarea")
	for _, b := range textarea {
		t, err := b.Text()
		if err != nil {
			// error handling
		}
		if t == "New post title here..." {
			b.MustInput("Top 5 Featured DEV Tag(#" + tag + ") Posts from the Past Week")
			break
		}
	}
	time.Sleep(2 * time.Second)
	fmt.Println("Title Entered")
	fmt.Println("Preparing Body")
	body := ""
	for _, v := range posts {
		body += "##" + v.Title + "\n"
		body += v.Description + "\n"
		body += v.PostLink + "\n"
	}
	page.Timeout(5 * time.Second).MustElement("#article_body_markdown").MustInput(body)
	time.Sleep(2 * time.Second)
	fmt.Println("Body entered")
	time.Sleep(2 * time.Second)
	fmt.Println("Looking for Publish button")
	page.MustElement("#tag-input").MustInput(tag).MustType(input.Tab).MustInput("c4r4x35").MustType(input.Tab)
	buttons := page.MustElements("button")
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
}

func main() {
	tags := make(map[string]string)
	tags["go"] = "https://dev.to/t/go/top/week"
	tags["javascript"] = "https://dev.to/t/javascript/top/week"
	tags["programming"] = "https://dev.to/t/programming/top/week"
	tags["devops"] = "https://dev.to/t/devops/top/week"
	tags["node"] = "https://dev.to/t/node/top/week"
	tags["python"] = "https://dev.to/t/python/top/week"

	userEmail := os.Getenv("USER_EMAIL")
	userPassword := os.Getenv("USER_PASSWORD")
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page := browser.MustPage("https://dev.to/enter")
	UserLogin(page, userEmail, userPassword)

	for k, v := range tags {
		tagPage := OpenTagPage(browser, k, v)
		posts := GetTopPosts(browser, tagPage, 4)
		addNewPost(browser, posts, k)
	}

}
