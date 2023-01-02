package pkg

import (
	devArtifact "github.com/dev.to-bot/pkg/artifact"
	"github.com/dev.to-bot/pkg/config"
	baseScraper "github.com/dev.to-bot/pkg/scraper"
	"go.uber.org/zap"
	"os"
)

func NewApp(version string) {
	logger := zap.NewExample()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	var scraper = baseScraper.NewScraper(logger)
	userEmail := os.Getenv("USER_EMAIL")
	userPassword := os.Getenv("USER_PASSWORD")
	artifact := devArtifact.NewArtifact(scraper, userEmail, userPassword, config.Tags, config.Titles, config.Urls)
	artifact.RunArtifact()
}
