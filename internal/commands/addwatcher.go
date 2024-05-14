package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/yimincai/health-checker/internal/bot"
	"github.com/yimincai/health-checker/internal/enums"
	"github.com/yimincai/health-checker/internal/errs"
	"github.com/yimincai/health-checker/internal/service"
	"github.com/yimincai/health-checker/models"
	"github.com/yimincai/health-checker/pkg/logger"
)

type CommandAddWatcher struct {
	Svc service.Service
}

func (c *CommandAddWatcher) IsAdminRequired() bool {
	return false
}

func (c *CommandAddWatcher) Invokes() []string {
	return []string{"AddWatcher", "addwatcher", "aw", "Addwatcher", "AW"}
}

func (c *CommandAddWatcher) Description() string {
	return "Add watcher"
}

func (c *CommandAddWatcher) Exec(ctx *bot.Context) (err error) {
	if args := ctx.Args; len(args) < 3 {
		usage := fmt.Sprintf("Usage: %sAddWatcher <name> <url> <interval in second>", c.Svc.Cfg.Prefix)
		example := fmt.Sprintf("Example: %sAddWatcher Google https://google.com 5", c.Svc.Cfg.Prefix)
		response := fmt.Sprintf("%s\n%s", usage, example)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return nil
	}

	name := ctx.Args[0]
	location := ctx.Args[1]
	interval, err := strconv.Atoi(ctx.Args[2])
	if err != nil {
		response := fmt.Sprintf("❌ Invalid interval: %s", ctx.Args[3])
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return nil
	}

	// check type, location, interval
	watcher := &models.Watcher{
		Name:     name,
		Type:     enums.Watcher_HTTP,
		Location: location,
		Interval: interval,
		IsEnable: true,
	}

	err = c.Svc.Validator.Struct(watcher)
	if err != nil {
		response := fmt.Sprintf("❌ Invalid input: %s", err)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return nil
	}

	_, err = c.Svc.Repo.CreateWatcher(watcher)
	if err != nil {
		if errors.Is(err, errs.ErrDuplicateWatcher) {
			return errs.ErrDuplicateWatcher
		}
		return errs.ErrInternalError
	}

	err = c.Svc.AddWatcher(watcher)
	if err != nil {
		return errs.ErrInternalError
	}

	response := fmt.Sprintf("✅ Watcher %s added", name)
	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)

	return err
}
