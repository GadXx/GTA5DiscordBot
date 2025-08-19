package handler

import (
	"discordbot/internal/service"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// InfoCommand реализует команду /info
type InfoCommand struct {
	service      *service.Service
	roleLeaderID string
}

func NewInfoCommand(s *service.Service, roleLeaderID string) *InfoCommand {
	return &InfoCommand{
		service:      s,
		roleLeaderID: roleLeaderID,
	}
}

func (c *InfoCommand) Name() string {
	return "info"
}

func (c *InfoCommand) Register(s *discordgo.Session, guildID string) error {
	cmd := &discordgo.ApplicationCommand{
		Name:        "info",
		Description: "Получить информацию о пользователях",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "leave",
				Description: "Список пользователей в отпуске",
			},
		},
	}
	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
	if err != nil {
		slog.Error("failed to register info command", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (c *InfoCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		slog.Error("invalid interaction type", slog.String("type", i.Type.String()))
		return
	}
	if i.ApplicationCommandData().Name != c.Name() {
		return
	}

	// Проверка прав
	hasPermission := false
	for _, roleID := range i.Member.Roles {
		slog.Info("Role ID", slog.String("roleID", roleID))
		if roleID == c.roleLeaderID {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Недостаточно прав для выполнения команды"},
		})
		return
	}

	// Проверка подкоманды
	if len(i.ApplicationCommandData().Options) == 0 || i.ApplicationCommandData().Options[0].Name != "leave" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Неизвестная подкоманда"},
		})
		return
	}

	// Получение списка пользователей в отпуске
	users, err := c.service.ListVacations()
	if err != nil {
		slog.Error("failed to list vacations", slog.String("error", err.Error()))
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Ошибка при получении списка отпусков"},
		})
		return
	}

	// Формирование ответа
	if len(users) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Нет пользователей в отпуске"},
		})
		return
	}

	var response strings.Builder
	response.WriteString("**Пользователи в отпуске:**\n")
	for _, user := range users {
		duration := user.EndAt.Sub(user.CreatedAt).Hours() / 24 // Продолжительность в днях
		response.WriteString(fmt.Sprintf(
			"- **%s**\n  Причина: %s\n  Начало: %s\n  Конец: %s\n  Длительность: %.0f дней\n",
			user.UserName, user.Reason,
			user.CreatedAt.Format("2006-01-02 15:04"),
			user.EndAt.Format("2006-01-02 15:04"),
			duration,
		))
	}

	// Отправка ответа
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: response.String()},
	})
}
