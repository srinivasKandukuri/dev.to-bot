package config

var Urls = map[string]string{
	"target":     "https://dev.to",
	"login":      "https://dev.to/enter",
	"createPost": "https://dev.to/new",
	"linkPrefix": "https://dev.to",
	"top":        "https://dev.to/t/%s/top/week",
	"dashboard":  "https://dev.to/dashboard",
}

var Tags = []string{
	/*	"go",
		"javascript",
		"programming",
		"devops",
		"node",
		"opensource",
		"ai",
		"docker",
		"github",
		"openai",
		"python",
		"blockchain",
		"security",
		"react",

	*/
	"css",
}

var Titles = []string{
	"Top 5 Featured DEV Tag(#%s) Posts from the Past Week",
	"Last week top 5 posts tagged(#%s)",
	"Top 5 Posts tagged(#%s) last week",
	"Awesome top 5 Posts from last week tagged(#%s)",
	"Checkout Last week top 5 posts tagged(#%s)",
	"Popular tag(#%s) last week top 5",
}

var comments = []string{
	"Shoutout to all the awesome authors featured in this years's Top 12 in 2022: %s 🙌.",
	"",
}
