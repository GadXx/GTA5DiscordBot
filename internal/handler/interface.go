package handler

import "github.com/bwmarrin/discordgo"

type Command interface {
	Register(s *discordgo.Session, guildID string) error
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
	Name() string
}
