package types

type Posts map[int]Post

type Tags []Tag

type Tag string

type Urls map[string]string

type Title string

type Titles []Title

type Post struct {
	Title       string
	Link        string
	PostLink    string
	Description string
	UserId      string
}
