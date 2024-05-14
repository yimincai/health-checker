package events

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yimincai/health-checker/internal/service"
	"github.com/yimincai/health-checker/pkg/logger"
)

type MessageHandler struct {
	Svc service.Service
}

func NewMessageHandler(svc service.Service) *MessageHandler {
	return &MessageHandler{
		Svc: svc,
	}
}

func (h *MessageHandler) Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, h.Svc.Cfg.Prefix) {
		return
	}

	if m.Author.ID == s.State.User.ID || m.Author.Bot || m.GuildID != "" {
		return
	}

	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		logger.Errorf("Error creating channel: %s", err)
	}

	response := fmt.Sprintf("Hello! I'm a bot.\n I'm not able to respond to messages right now.\n If you need help, please using `%shelp` command or contact the server administrator.", h.Svc.Cfg.Prefix)

	_, err = s.ChannelMessageSend(channel.ID, response)
	if err != nil {
		logger.Errorf("Error sending message: %s", err)
	}
}
