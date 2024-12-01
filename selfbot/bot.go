package selfbot

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lmittmann/tint"
	"github.com/twenty48lol/selfbot/commands"
)

type Bot struct {
	Session  *discordgo.Session
	User     *discordgo.User
	Config   Config
	Commands commands.CommandList
	Logger   *slog.Logger
}

func NewBot(config Config, commands commands.CommandList) (Bot, error) {
	var err error
	t := time.Now()
	session, err := discordgo.New(config.Token)
	if err != nil {
		return Bot{}, err
	}
	if err = session.Open(); err != nil {
		return Bot{}, err
	}
	user, err := session.User("@me")
	if err != nil {
		return Bot{}, err
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.Kitchen,
	}))
	logger.Info(fmt.Sprintf("Found token for user: %s", user.Username))
	logger.Info(fmt.Sprintf("Started in %dms", time.Since(t).Milliseconds()))

	bot := Bot{
		Session:  session,
		User:     user,
		Config:   config,
		Commands: commands,
		Logger:   logger,
	}

	err = bot.RegisterEvents()

	return bot, err
}

func (bot *Bot) Close() {
	t := time.Now()
	_ = bot.Session.Close()
	bot.Logger.Info(fmt.Sprintf("Shutdown in %dms", time.Since(t).Milliseconds()))
	bot.Logger.Info("Shutdown")
}
