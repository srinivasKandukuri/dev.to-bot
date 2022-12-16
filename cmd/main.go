package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"math/rand"
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

func waitToEnd(num int64) {
	time.Sleep(time.Duration(num) * time.Second)
}

func assertEquals(a string, b string, msg string) {
	if a != b {
		fmt.Printf("Failed on -> %s: \"%v\" != \"%v\"\n", msg, a, b)
		panic("Assertion Error")
	}
}

func GenerateRandomTitle(titles []string) string {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(titles)
	return titles[n]
}

func UserLogin(page *rod.Page, email string, pass string) {
	page.MustElement("[type='email']").MustInput(email).MustType(input.Tab)
	page.MustElement("[type='password']").MustInput(pass).MustType(input.Enter)
	assertEquals(page.MustElement("nav header .crayons-subtitle-3").MustText(), "My Tags", "User not able to login")
	fmt.Println("User Logged IN")
}

func OpenTagPage(browser *rod.Browser, tag string, url string) *rod.Page {
	page := browser.MustPage(url)
	fmt.Printf("Navigated to tag : %s , url : %s\n", tag, url)
	return page
}

func GetTopPosts(browser *rod.Browser, page *rod.Page, topPostsCount int) (listOfPosts map[int]Post) {
	waitToEnd(2)
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

func addNewPost(browser *rod.Browser, posts map[int]Post, tag string, title string) {
	fmt.Println("Navigating to NewPost")
	url := "https://dev.to/new"
	newPage := browser.MustPage(url).MustWaitLoad()
	waitToEnd(2)
	fmt.Println("Page opened")
	pageTitle := newPage.Timeout(5 * time.Second).MustElement("#article-form-title")
	assertEquals(pageTitle.MustText(), "New post title here...", "The title is not present. Perhaps the homepage is not opening")

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
	waitToEnd(2)
	fmt.Println("Preparing Body")
	body := ""
	for _, v := range posts {
		body += "##" + v.Title + "\n"
		body += v.Description + "\n"
		body += v.PostLink + "\n"
	}
	newPage.Timeout(5 * time.Second).MustElement("#article_body_markdown").MustInput(body)
	fmt.Println("Body entered")
	waitToEnd(2)
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
}

/*
func AddHeaderImage(browser *rod.Browser) {
	url := "https://dev.to/new"
	page := browser.MustPage(url).MustWaitLoad()
	time.Sleep(2 * time.Second)
	fmt.Println("Page opened")
	storjLogo, _ := filepath.Abs("./go.jpeg")
	page.MustElement("input[type=file]").MustSetFiles(storjLogo)

	page.SetDocumentContent()
	textarea := page.Timeout(5 * time.Second).MustElements("textarea")
	for _, b := range textarea {
		t, err := b.Text()
		if err != nil {
			// error handling
		}
		if t == "New post title here..." {
			b.MustInput("Top 5 Featured DEV Tag(#go) Posts from the Past Week")
			break
		}
	}

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
			fmt.Printf("============================END tag post : %s=========================\n", "go")
			break
		}
	}

}
*/
func main() {

	titles := make([]string, 0)
	titles = append(titles,
		"Top 5 Featured DEV Tag(#%s) Posts from the Past Week",
		"Last week top 5 posts tagged(#%s)",
		"Top 5 Posts tagged(#%s) last week",
		"Awesome Posts from last week tagged(#%s)",
		"Checkout Last week top 5 posts tagged(#%s)",
		"Popular tag(#%s) last week top 5",
	)

	tags := make(map[string]string)
	tags["go"] = "https://dev.to/t/go/top/week"
	tags["javascript"] = "https://dev.to/t/javascript/top/week"
	tags["programming"] = "https://dev.to/t/programming/top/week"
	tags["devops"] = "https://dev.to/t/devops/top/week"
	tags["node"] = "https://dev.to/t/node/top/week"
	tags["python"] = "https://dev.to/t/python/top/week"
	tags["opensource"] = "https://dev.to/t/opensource/top/week"
	tags["ai"] = "https://dev.to/t/ai/top/week"

	userEmail := os.Getenv("USER_EMAIL")
	userPassword := os.Getenv("USER_PASSWORD")
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page := browser.MustPage("https://dev.to/enter")
	UserLogin(page, userEmail, userPassword)
	//AddHeaderImage(browser)

	for k, v := range tags {
		tagPage := OpenTagPage(browser, k, v)
		posts := GetTopPosts(browser, tagPage, 4)
		addNewPost(browser, posts, k, GenerateRandomTitle(titles))
	}
}
