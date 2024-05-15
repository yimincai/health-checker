package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yimincai/health-checker/internal/bot"
	"github.com/yimincai/health-checker/internal/commands"
	"github.com/yimincai/health-checker/internal/events"
	"github.com/yimincai/health-checker/internal/middlewares"
	"github.com/yimincai/health-checker/pkg/logger"
	"github.com/yimincai/health-checker/pkg/snowflake"
)

var (
	goVersion = runtime.Version()
	oSArch    = fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
)

func init() {
	loc := time.FixedZone("UTC+8", 8*60*60)
	time.Local = loc
	logger.New()
	snowflake.New()
}

func main() {
	server := bot.New()

	// Register events
	registerEvents(server)

	// Register commands
	registerCommands(server)

	logger.Infof("Go Version: %s", goVersion)
	logger.Infof("OS/Arch: %s", oSArch)
	logger.Infof("Bot timezone is %s", time.Local.String())
	server.Run()
	defer func() {
		server.Close()
		logger.Info("Bot closed")
	}()

	// ============================== Graceful shutdown ==============================
	// Wait for interrupt signal to gracefully shut down the server with a timeout.
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catching, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// ========================= Server block and serve here =========================
	q := <-quit
	// ============================== Graceful shutdown ==============================

	logger.Infof("Got signal: %s", q.String())
	logger.Info("Shutdown server ...")
}

func registerEvents(s *bot.Bot) {
	s.Session.AddHandler(events.NewMessageHandler(s.Svc).Handler)
}

func registerCommands(b *bot.Bot) {
	// Register commands here
	cmdHandler := bot.NewCommandHandler(b.Cfg.Prefix)
	cmdHandler.OnError = func(ctx *bot.Context, err error) {
		logger.Errorf("Error executing command: %v", err)
		embedMsg := &discordgo.MessageEmbed{
			Title:       "Error",
			Description: "error: " + err.Error(),
			Color:       0xff0000,
		}
		_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embedMsg)
		if err != nil {
			logger.Errorf("Error sending message: %v", err)
		}
	}

	cmdHandler.RegisterCommand(&commands.CommandHelp{Cfg: b.Cfg})
	cmdHandler.RegisterCommand(&commands.CommandHealth{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandAddWatcher{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandPrintWatchers{Svc: b.Svc})
	cmdHandler.RegisterCommand(&commands.CommandDeleteWatcher{Svc: b.Svc})

	cmdHandler.RegisterMiddleware(middlewares.NewRequiredAdminPermission(b.Repo))

	b.Session.AddHandler(cmdHandler.HandleMessage)
}
