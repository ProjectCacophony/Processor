package cache

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	discord      *discordgo.Session
	discordMutex sync.RWMutex
)

func SetDiscord(s *discordgo.Session) {
	discordMutex.Lock()
	discord = s
	discordMutex.Unlock()
}

func GetDiscord() *discordgo.Session {
	discordMutex.RLock()
	defer discordMutex.RUnlock()

	return discord
}
