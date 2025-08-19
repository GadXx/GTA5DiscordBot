package discord

import (
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler interface {
	Register(s *discordgo.Session, guildID string) error
	HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type Bot struct {
	Session       *discordgo.Session
	handlers      []CommandHandler
	guildID       string
	DefaultRoleID string
}

func NewBot(token, guildID, defaultRoleID string, handlers []CommandHandler) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("failed to create Discord session", slog.String("error", err.Error()))
		return nil, err
	}

	s.Identify.Intents = discordgo.IntentsGuildMembers | discordgo.IntentsGuilds

	bot := &Bot{
		Session:       s,
		handlers:      handlers,
		guildID:       guildID,
		DefaultRoleID: defaultRoleID,
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.GuildID != bot.guildID {
			slog.Debug("Ignoring interaction from another guild", slog.String("guildID", i.GuildID))
			return
		}
		for _, h := range bot.handlers {
			h.HandleInteraction(s, i)
		}
	})

	if defaultRoleID != "" {
		s.AddHandler(bot.handleGuildMemberAdd)
		slog.Info("GuildMemberAdd handler registered", slog.String("defaultRoleID", defaultRoleID))
	} else {
		slog.Warn("DefaultRoleID is empty, GuildMemberAdd handler not registered")
	}

	return bot, nil
}

func (b *Bot) handleGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	slog.Info("GuildMemberAdd event received",
		slog.String("userID", m.User.ID),
		slog.String("username", m.User.Username),
		slog.String("guildID", m.GuildID))

	if m.GuildID != b.guildID {
		slog.Debug("Ignoring member add event from another guild", slog.String("guildID", m.GuildID))
		return
	}

	if m.User.Bot {
		slog.Info("Ignoring bot user", slog.String("userID", m.User.ID))
		return
	}

	roles, err := s.GuildRoles(b.guildID)
	if err != nil {
		slog.Error("failed to fetch guild roles", slog.String("error", err.Error()))
		return
	}
	roleExists := false
	for _, role := range roles {
		if role.ID == b.DefaultRoleID {
			roleExists = true
			break
		}
	}
	if !roleExists {
		slog.Error("default role does not exist", slog.String("roleID", b.DefaultRoleID))
		return
	}

	for retries := 3; retries > 0; retries-- {
		err = s.GuildMemberRoleAdd(b.guildID, m.User.ID, b.DefaultRoleID)
		if err == nil {
			slog.Info("Default role added to new member",
				slog.String("userID", m.User.ID),
				slog.String("username", m.User.Username),
				slog.String("roleID", b.DefaultRoleID))
			return
		}
		if discordErr, ok := err.(*discordgo.RESTError); ok && discordErr.Message != nil {
			slog.Warn("rate limited, retrying",
				slog.String("userID", m.User.ID),
				slog.Int("retries_left", retries-1))
			time.Sleep(time.Second * 2)
			continue
		}
		slog.Error("failed to add default role",
			slog.String("userID", m.User.ID),
			slog.String("roleID", b.DefaultRoleID),
			slog.String("error", err.Error()))
		return
	}
	slog.Error("failed to add default role after retries",
		slog.String("userID", m.User.ID),
		slog.String("roleID", b.DefaultRoleID))
}

func (b *Bot) Open() error {
	if err := b.Session.Open(); err != nil {
		slog.Error("failed to open Discord connection", slog.String("error", err.Error()))
		return err
	}

	for _, h := range b.handlers {
		if err := h.Register(b.Session, b.guildID); err != nil {
			slog.Error("failed to register handler", slog.String("error", err.Error()))
			return err
		}
	}

	slog.Info("Bot successfully connected and commands registered", slog.String("guildID", b.guildID))
	return nil
}

func (b *Bot) Close() {
	slog.Info("Closing bot connection")
	if err := b.Session.Close(); err != nil {
		slog.Error("failed to close Discord connection", slog.String("error", err.Error()))
	}
}
