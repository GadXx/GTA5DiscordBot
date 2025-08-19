package handler

import (
	"discordbot/internal/service"
	"log/slog"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type InactiveCommand struct {
	service      *service.Service
	logChannelID string
	roleVacID    string
}

func NewInactiveCommand(s *service.Service, logChannelID, roleVacID string) *InactiveCommand {
	return &InactiveCommand{
		service:      s,
		logChannelID: logChannelID,
		roleVacID:    roleVacID,
	}
}

func (c *InactiveCommand) Name() string {
	return "inactive"
}

func (c *InactiveCommand) Register(s *discordgo.Session, guildID string) error {
	cmd := &discordgo.ApplicationCommand{
		Name:        "inactive",
		Description: "Заявка на инакт",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "days",
				Description: "Число дней отсутствия",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "Причина отсутствия",
				Required:    true,
			},
		},
	}
	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
	return err
}

func (c *InactiveCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		slog.Error("invalid interaction type", slog.String("type", i.Type.String()))
		return
	}
	if i.ApplicationCommandData().Name != c.Name() {
		return
	}

	days := i.ApplicationCommandData().Options[0].IntValue()
	reason := i.ApplicationCommandData().Options[1].StringValue()
	slog.Info("Handling interaction", slog.Int64("days", days), slog.String("reason", reason))

	user := i.Member.User
	guildID := i.GuildID
	roleID := c.roleVacID

	now := time.Now()
	end := now.Add(time.Hour * 24 * time.Duration(days))

	err := c.service.SaveRequest(user.ID, user.Username, reason, now, end)
	if err != nil {
		slog.Error("failed to save request", slog.String("error", err.Error()))
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Ошибка при сохранении заявки"},
		})
		return
	}

	err = s.GuildMemberRoleAdd(guildID, user.ID, roleID)
	if err != nil {
		slog.Error("failed to add role", slog.String("error", err.Error()))
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Ошибка при добавлении роли"},
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "Заявка на отпуск успешно обработана!"},
	})

	_, err = s.ChannelMessageSend(c.logChannelID, "Пользователь <@"+user.ID+"> заявил отпуск, вернется через "+strconv.Itoa(int(days))+" дней.")
	if err != nil {
		slog.Error("failed to send DM", slog.String("userID", user.ID), slog.String("error", err.Error()))
	}
}
