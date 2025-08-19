package app

import (
	"discordbot/internal/database"
	"discordbot/internal/discord"
	"discordbot/internal/handler"
	"discordbot/internal/repository"
	"discordbot/internal/service"
	"log/slog"
)

type App struct {
	Bot *discord.Bot
}

func NewApp() (*App, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	db, err := database.Connect(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	repo := repository.NewRepository(db)
	inactiveService := service.NewInactiveService(repo)
	handler := handler.NewHandler(inactiveService, cfg.LogChannelID, cfg.RoleVacID, cfg.RoleSanctionID, cfg.RoleLeaderID)
	bot, err := discord.NewBot(cfg.BotToken, cfg.GuildID, cfg.DefaultRoleID, []discord.CommandHandler{
		handler,
	})
	if err != nil {
		return nil, err
	}
	trackingService := service.NewTrackingService(repo, bot.Session, cfg.GuildID, cfg.RoleVacID, cfg.RoleSanctionID)
	trackingService.StartTracking() // Запускаем трекинг
	return &App{Bot: bot}, nil
}

func (a *App) Run() error {
	slog.Info("Bot is running...")
	return a.Bot.Open()
}

func (a *App) Close() {
	a.Bot.Close()
}
