package config

import "github.com/dev.to-bot/pkg/types"

var Urls = types.Urls{
	"target":     "https://dev.to",
	"login":      "https://dev.to/enter",
	"createPost": "https://dev.to/new",
	"linkPrefix": "https://dev.to",
	"top":        "https://dev.to/t/%s/top/year",
	"dashboard":  "https://dev.to/dashboard",
}

var Tags = types.Tags{
	"go",
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
	"css",
}

var Titles = types.Titles{
	"Top 5 Featured DEV Tag(#%s) Posts from the year 2022",
}

var comments = []string{
	"Shoutout to all the awesome authors featured in this years's Top 5 in 2022: %s 🙌.",
}
