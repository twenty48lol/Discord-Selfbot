package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/lmittmann/tint"
	"github.com/twenty48lol/selfbot/commands"
	"github.com/twenty48lol/selfbot/selfbot"
)

var (
	realConfig string
)

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
		NoColor:    false,
	})))

	flag.StringVar(&realConfig, "config", "config.json", "Specify the config file to use.")
	flag.Parse()

	var config selfbot.Config

	if _, err := os.Stat(realConfig); os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		slog.Info("Could not find config.")
		slog.Info("Creating new config.")
		slog.Info("Enter your discord token here: ")
		token, _ := reader.ReadString('\n')
		token = strings.TrimSpace(token)

		config = selfbot.NewDefaultConfig(token)

		file, err := os.Create(realConfig)
		if err != nil {
			panic(err)
		}
		encoder := json.NewEncoder(file)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(config)
	} else {
		file, err := os.ReadFile(realConfig)
		if err != nil {
			slog.Error(err.Error())
		}
		config, _ = selfbot.LoadConfig(file)
	}
	slog.Info("Loaded Config.")

	commands := commands.InitCommands()
	bot, err := selfbot.NewBot(config, commands)
	if err != nil {
		slog.Error(err.Error())
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	go func() {
		bot.Close()
		waitgroup.Done()
	}()
	waitgroup.Wait()

}
