package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yimincai/health-checker/pkg/logger"
)

type CommandHandler struct {
	prefix       string
	cmdInstances []Command
	cmdMap       map[string]Command
	middlewares  []Middleware
	OnError      func(ctx *Context, err error)
}

// NewCommandHandler creates a new CommandHandler.
func NewCommandHandler(prefix string) *CommandHandler {
	return &CommandHandler{
		prefix:       prefix,
		cmdInstances: make([]Command, 0),
		cmdMap:       make(map[string]Command),
		middlewares:  make([]Middleware, 0),
		OnError:      func(ctx *Context, err error) {},
	}
}

func (c *CommandHandler) RegisterCommand(cmd Command) {
	c.cmdInstances = append(c.cmdInstances, cmd)
	for _, invoke := range cmd.Invokes() {
		c.cmdMap[invoke] = cmd
	}

	logger.Infof("Command Registered: %v, Description: %s", cmd.Invokes(), cmd.Description())
}

func (c *CommandHandler) RegisterMiddleware(mw Middleware) {
	c.middlewares = append(c.middlewares, mw)
}

func (c *CommandHandler) HandleMessage(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.ID == s.State.User.ID || e.Author.Bot || !strings.HasPrefix(e.Content, c.prefix) {
		return
	}

	logger.Infof("Command Received, Content: %s, UserID: %s", e.Content, e.Author.ID)

	split := strings.Split(e.Content[len(c.prefix):], " ")
	if len(split) < 1 {
		return
	}

	invoke := split[0]
	args := split[1:]

	cmd, ok := c.cmdMap[invoke]
	if !ok || cmd == nil {
		embedMsg := &discordgo.MessageEmbed{
			Title:       "Command not found",
			Description: "Try `!help`",
			Color:       0xff0000,
		}
		_, err := s.ChannelMessageSendEmbed(e.ChannelID, embedMsg)
		if err != nil {
			logger.Errorf("Error sending message: %s", err)
		}
		return
	}

	ctx := &Context{
		Session:  s,
		Args:     args,
		Handler:  c,
		Message:  e.Message,
		Commands: c.GetCommands(),
	}

	for _, mw := range c.middlewares {
		next, err := mw.Exec(ctx, cmd)
		if err != nil {
			c.OnError(ctx, err)
			return
		}
		if !next {
			return
		}
	}

	if err := cmd.Exec(ctx); err != nil {
		c.OnError(ctx, err)
	}
}

func (c *CommandHandler) GetCommands() []Command {
	return c.cmdInstances
}
