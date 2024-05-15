package commands

import (
	"errors"
	"fmt"

	"github.com/yimincai/health-checker/internal/bot"
	"github.com/yimincai/health-checker/internal/errs"
	"github.com/yimincai/health-checker/internal/service"
	"github.com/yimincai/health-checker/pkg/logger"
)

type CommandDeleteWatcher struct {
	Svc service.Service
}

func (c *CommandDeleteWatcher) IsAdminRequired() bool {
	return false
}

func (c *CommandDeleteWatcher) Invokes() []string {
	return []string{"DeleteWatcher", "deletewatcher", "dw", "Deletewatcher", "DW"}
}

func (c *CommandDeleteWatcher) Description() string {
	return "Delete a watcher"
}

func (c *CommandDeleteWatcher) Exec(ctx *bot.Context) (err error) {
	if args := ctx.Args; len(args) < 1 {
		usage := fmt.Sprintf("Usage: %sDeleteWatcher <id>", c.Svc.Cfg.Prefix)
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, usage)
		if err != nil {
			return errs.ErrSendingMessage
		}
		return nil
	}

	id := ctx.Args[0]

	watcher, err := c.Svc.Repo.FindWatcherByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrWatcherNotFound) {
			response := fmt.Sprintf("❌ Watcher ID %s not found", id)
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
			if err != nil {
				return errs.ErrSendingMessage
			}
			return nil
		}
		return errs.ErrInternalError
	}

	err = c.Svc.Repo.DeleteWatcher(watcher.ID)
	if err != nil {
		return errs.ErrInternalError
	}

	c.Svc.Cron.Remove(watcher.GetCronID())
	watcher.RemoveCronID()
	watcher.RemoveContinueErrorTimes()
	watcher.RemoveLastStatus()

	response := fmt.Sprintf("✅ Watcher %s deleted", watcher.Name)
	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, response)
	if err != nil {
		return errs.ErrSendingMessage
	}

	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)

	return err
}
