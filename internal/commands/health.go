package commands

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/yimincai/health-checker/internal/bot"
	"github.com/yimincai/health-checker/internal/service"
	"github.com/yimincai/health-checker/internal/utils"
	"github.com/yimincai/health-checker/pkg/logger"
)

type CommandHealth struct {
	Svc service.Service
}

func (c *CommandHealth) IsAdminRequired() bool {
	return false
}

func (c *CommandHealth) Invokes() []string {
	return []string{"Health", "health"}
}

func (c *CommandHealth) Description() string {
	return "Checking if services is healthy use `img` as argument to get image"
}

func (c *CommandHealth) Exec(ctx *bot.Context) (err error) {
	result, err := c.Svc.CheckHealth()
	if err != nil {
		return err
	}

	if len(result) == 0 {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Result is empty, maybe the watchlist is empty")
		if err != nil {
			logger.Errorf("Error sending message to channel %s: %v", ctx.Message.ChannelID, err)
		}
		return nil
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Service", "Status", "Response Time"})
	for _, r := range result {
		status := ""
		if len(ctx.Args) > 0 && ctx.Args[0] == "img" {
			status = "Down"
			if r.Status {
				status = "Up"
			}
		} else {
			status = "❌"
			if r.Status {
				status = "✅"
			}
		}
		t.AppendRow(table.Row{r.Name, status, fmt.Sprintf("%d %s", r.ResponseTime, "ms")})
	}
	t.SetStyle(table.StyleLight)

	if len(ctx.Args) > 0 && ctx.Args[0] == "img" {
		text := strings.Split(t.Render(), "\n")
		br, err := utils.GenerateImage(text)
		if err != nil {
			logger.Errorf("Error generating image: %v", err)
			return err
		}

		_, err = ctx.Session.ChannelFileSend(ctx.Message.ChannelID, "health.png", br)
		if err != nil {
			logger.Errorf("Error sending image to channel %s: %v", ctx.Message.ChannelID, err)
		}
		return nil
	}

	response := "```" + t.Render() + "```"

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		logger.Errorf("Error sending message to channel %s: %v", ctx.Message.ChannelID, err)
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)

	return err
}
