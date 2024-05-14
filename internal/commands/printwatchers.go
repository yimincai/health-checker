package commands

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/yimincai/health-checker/internal/bot"
	"github.com/yimincai/health-checker/internal/errs"
	"github.com/yimincai/health-checker/internal/service"
	"github.com/yimincai/health-checker/pkg/logger"
	"github.com/yimincai/health-checker/pkg/utils"
)

type CommandPrintWatchers struct {
	Svc service.Service
}

func (c *CommandPrintWatchers) IsAdminRequired() bool {
	return false
}

func (c *CommandPrintWatchers) Invokes() []string {
	return []string{"PrintWatchers", "printwatchers", "pw", "Printwatchers", "PW"}
}

func (c *CommandPrintWatchers) Description() string {
	return "Print watchers, use `img` as argument to get image"
}

func (c *CommandPrintWatchers) Exec(ctx *bot.Context) (err error) {
	watchers, err := c.Svc.Repo.FindWatchers()
	if err != nil {
		return err
	}

	if len(watchers) == 0 {
		response := "❌ No watchers found"
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return nil
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"ID", "Name", "Url", "Interval"})
	for _, w := range watchers {
		t.AppendRow(table.Row{w.ID, w.Name, w.Location, fmt.Sprintf("%d %s", w.Interval, "s")})
	}
	t.SetStyle(table.StyleLight)

	text := strings.Split(t.Render(), "\n")

	br, err := utils.GenerateImage(text)
	if err != nil {
		logger.Errorf("Error generating image: %v", err)
		return err
	}

	if len(ctx.Args) > 0 && ctx.Args[0] == "img" {
		_, err = ctx.Session.ChannelFileSend(ctx.Message.ChannelID, "watchlist.png", br)
		if err != nil {
			logger.Errorf("Error sending image to channel %s: %v", ctx.Message.ChannelID, err)
		}
		return nil
	} else {
		response := "```" + t.Render() + "```"
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
		if err != nil {
			logger.Errorf("Error sending message to channel %s: %v", ctx.Message.ChannelID, err)
		}

		logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)

		return err
	}
}
