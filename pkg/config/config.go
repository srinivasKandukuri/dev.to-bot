package config

var Urls = map[string]string{
	"target":     "https://dev.to",
	"login":      "https://dev.to/enter",
	"createPost": "https://dev.to/new",
	"linkPrefix": "https://dev.to",
}

var Tags = map[string]string{
	"go":          "https://dev.to/t/go/top/week",
	"javascript":  "https://dev.to/t/javascript/top/week",
	"programming": "https://dev.to/t/programming/top/week",
	"devops":      "https://dev.to/t/devops/top/week",
	"node":        "https://dev.to/t/node/top/week",
	"python":      "https://dev.to/t/python/top/week",
	"opensource":  "https://dev.to/t/opensource/top/week",
	"ai":          "https://dev.to/t/ai/top/week",
	"docker":      "https://dev.to/t/docker/top/week",
}

var Titles = []string{
	"Top 5 Featured DEV Tag(#%s) Posts from the Past Week",
	"Last week top 5 posts tagged(#%s)",
	"Top 5 Posts tagged(#%s) last week",
	"Awesome Posts from last week tagged(#%s)",
	"Checkout Last week top 5 posts tagged(#%s)",
	"Popular tag(#%s) last week top 5",
}
