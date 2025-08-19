package service

import (
	"discordbot/internal/repository"
	"log/slog"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type TrackingService struct {
	wg             sync.WaitGroup
	repo           *repository.Repository
	session        *discordgo.Session
	roleVacID      string
	roleSanctionID string
	guildID        string
}

func NewTrackingService(r *repository.Repository, s *discordgo.Session, guildID, roleVacID, roleSanctionID string) *TrackingService {
	return &TrackingService{
		wg:             sync.WaitGroup{},
		repo:           r,
		session:        s,
		roleVacID:      roleVacID,
		roleSanctionID: roleSanctionID,
		guildID:        guildID,
	}
}

func (s *TrackingService) StartTracking() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.TrackIvent()
	}()
}

func (s *TrackingService) TrackIvent() {
	ticker := time.NewTicker(time.Second * 10) // Проверка каждые 10 секунд
	defer ticker.Stop()
	slog.Info("Starting tracking service")
	for range ticker.C {
		users, err := s.repo.ListWithEndAt(time.Now())
		if err != nil {
			slog.Error("failed to list users", slog.String("error", err.Error()))
			continue
		}

		for _, user := range users {
			slog.Info("Processing user with ended vacation", slog.String("userID", user.UserID), slog.String("userName", user.UserName))

			channel, err := s.session.UserChannelCreate(user.UserID)
			if err != nil {
				slog.Error("failed to create DM channel", slog.String("userID", user.UserID), slog.String("error", err.Error()))
				continue
			}

			_, err = s.session.ChannelMessageSend(channel.ID, "Ваш отпуск закончился. И теперь вы находитесь в санкционированном состоянии. Для снятия ......")
			if err != nil {
				slog.Error("failed to send DM", slog.String("userID", user.UserID), slog.String("error", err.Error()))
				continue
			}

			err = s.session.GuildMemberRoleRemove(s.guildID, user.UserID, s.roleVacID)
			if err != nil {
				slog.Error("failed to remove vacation role", slog.String("userID", user.UserID), slog.String("error", err.Error()))
				continue
			}

			err = s.session.GuildMemberRoleAdd(s.guildID, user.UserID, s.roleSanctionID)
			if err != nil {
				slog.Error("failed to add sanction role", slog.String("userID", user.UserID), slog.String("error", err.Error()))
				continue
			}

			err = s.repo.Delete(user.ID)
			if err != nil {
				slog.Error("failed to delete request", slog.String("userID", user.UserID), slog.String("error", err.Error()))
				continue
			}

			slog.Info("Successfully processed user", slog.String("userID", user.UserID), slog.String("userName", user.UserName))
		}
	}
}

func (s *TrackingService) StopTracking() {
	s.wg.Wait()
}
