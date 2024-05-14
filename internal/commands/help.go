package commands

import (
	"github.com/yimincai/health-checker/internal/bot"
	"github.com/yimincai/health-checker/internal/config"
	"github.com/yimincai/health-checker/internal/errs"
	"github.com/yimincai/health-checker/internal/utils"
	"github.com/yimincai/health-checker/pkg/logger"
)

type CommandHelp struct {
	Cfg *config.Config
}

func (c *CommandHelp) IsAdminRequired() bool {
	return false
}

func (c *CommandHelp) Invokes() []string {
	return []string{"Help", "h", "help"}
}

func (c *CommandHelp) Description() string {
	return "Show help message, list of available commands, use `img` as argument to get image"
}

func (c *CommandHelp) Exec(ctx *bot.Context) (err error) {
	var response string
	var respArr []string

	if len(ctx.Args) > 0 && ctx.Args[0] == "img" {
		respArr = append(respArr, "Commands:")
		for _, command := range ctx.Commands {
			respArr = append(respArr, c.Cfg.Prefix+command.Invokes()[0]+": "+command.Description())
		}
	} else {
		for _, command := range ctx.Commands {
			response += c.Cfg.Prefix + command.Invokes()[0] + ": "
			response += command.Description() + "\n"
		}
	}

	if len(ctx.Args) > 0 && ctx.Args[0] == "img" {
		br, err := utils.GenerateImage(respArr)
		if err != nil {
			logger.Errorf("Error generating image: %v", err)
			return err
		}

		_, err = ctx.Session.ChannelFileSend(ctx.Message.ChannelID, "help.png", br)
		if err != nil {
			logger.Errorf("Error sending image to channel %s: %v", ctx.Message.ChannelID, err)
		}
		return nil
	}

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
	return
}
