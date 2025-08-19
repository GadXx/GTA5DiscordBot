package handler

import (
	"discordbot/internal/service"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type SanctionCommand struct {
	service        *service.Service
	roleSanctionID string
	roleLeaderID   string
}

func NewSanctionCommand(s *service.Service, roleSanctionID, roleLeaderID string) *SanctionCommand {
	return &SanctionCommand{
		service:        s,
		roleSanctionID: roleSanctionID,
		roleLeaderID:   roleLeaderID,
	}
}

func (c *SanctionCommand) Name() string {
	return "sanction"
}

func (c *SanctionCommand) Register(s *discordgo.Session, guildID string) error {
	cmd := &discordgo.ApplicationCommand{
		Name:        "sanction",
		Description: "Наложить санкцию на пользователя",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "Причина санкции",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Пользователь, на которого накладывается санкция",
				Required:    true,
			},
		},
	}
	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
	return err
}

func (c *SanctionCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		slog.Error("invalid interaction type", slog.String("type", i.Type.String()))
		return
	}
	if i.ApplicationCommandData().Name != c.Name() {
		return
	}

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

	reason := i.ApplicationCommandData().Options[0].StringValue()
	user := i.ApplicationCommandData().Options[1].UserValue(s)
	slog.Info("Handling sanction command", slog.String("reason", reason), slog.String("userID", user.ID))

	guildID := i.GuildID
	roleID := c.roleSanctionID

	err := s.GuildMemberRoleAdd(guildID, user.ID, roleID)
	if err != nil {
		slog.Error("failed to add sanction role", slog.String("error", err.Error()))
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Ошибка при добавлении роли санкции"},
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("Санкция наложена на %s, роль добавлена", user.Username)},
	})
}
