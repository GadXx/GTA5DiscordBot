package app

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken       string
	GuildID        string
	DBPath         string
	LogChannelID   string
	RoleVacID      string
	RoleSanctionID string
	RoleLeaderID   string
	DefaultRoleID  string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		BotToken:       os.Getenv("BOT_TOKEN"),
		GuildID:        os.Getenv("GUILD_ID"),
		DBPath:         os.Getenv("DB_PATH"),
		LogChannelID:   os.Getenv("LOG_CHANNEL_ID"),
		RoleVacID:      os.Getenv("ROLE_VAC_ID"),
		RoleSanctionID: os.Getenv("ROLE_SANCTION_ID"),
		RoleLeaderID:   os.Getenv("ROLE_LEADER_ID"),
		DefaultRoleID:  os.Getenv("DEFAULT_ROLE_ID"),
	}, nil
}
