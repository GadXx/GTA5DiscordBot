package handler

import (
	"discordbot/internal/service"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	service        *service.Service
	logChannelID   string
	roleVacID      string
	roleSanctionID string
	roleLeaderID   string
	commands       []Command
}

func NewHandler(s *service.Service, logChannelID, roleVacID, roleSanctionID, roleLeaderID string) *Handler {

	commands := []Command{
		NewInactiveCommand(s, logChannelID, roleVacID),
		NewSanctionCommand(s, roleSanctionID, roleLeaderID),
		NewInfoCommand(s, roleLeaderID),
	}

	return &Handler{
		service:        s,
		logChannelID:   logChannelID,
		roleVacID:      roleVacID,
		roleSanctionID: roleSanctionID,
		roleLeaderID:   roleLeaderID,
		commands:       commands,
	}
}

func (h *Handler) Register(s *discordgo.Session, guildID string) error {
	for _, cmd := range h.commands {
		if err := cmd.Register(s, guildID); err != nil {
			slog.Error("failed to register command", slog.String("command", cmd.Name()), slog.String("error", err.Error()))
			return err
		}
		slog.Info("Registered command", slog.String("command", cmd.Name()))
	}
	return nil
}

func (h *Handler) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, cmd := range h.commands {
		if i.ApplicationCommandData().Name == cmd.Name() {
			cmd.Handle(s, i)
			return
		}
	}
	slog.Error("unknown command", slog.String("name", i.ApplicationCommandData().Name))
}
