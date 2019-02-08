package kit

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// BotSession creates a new DiscordGo Client for the given BotID
// requires DISCORD_TOKEN_<BOT ID> to be set
func BotSession(botID string) (*discordgo.Session, error) {
	token := os.Getenv("DISCORD_TOKEN_" + botID)

	if token == "" {
		return nil, errors.New("token for bot ID is not configured")
	}

	newSession, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, errors.New("error creating bot session")
	}

	return newSession, nil
}
