package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/jawadcode/quotesbotv2/db"
	"github.com/jawadcode/quotesbotv2/handlers"
)

const (
	// PREFIX is the prefix for every command
	PREFIX = "Q!"
	// COLOUR is a nice bright discord orange colour int
	COLOUR = 0xFAA61A
)

// All of this setup code is stolen straight from discordgo's examples :)
func main() {
	// Initialise DB connection
	db.DBInit()
	// Get discord token from environment variable
	token := os.Getenv("DISCORD_TOKEN")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	// We only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// To keep dereferencing to a minimum and also to make code easy to read
	content := m.Content
	// Ignore all messages created by the bot itself and also check if the length of the command is less than 7 (minimum possible length of command)
	if m.Author.ID == s.State.User.ID || len(content) < 7 {
		return
	}
	// Ignore messages that don't start with the prefix
	if content[:2] != PREFIX {
		return
	}
	// Split everything after prefix into individual arguments (by whitespace)
	args := strings.Fields(content[3:])
	// Command handlers
	switch strings.ToLower(args[0]) {
	case "help":
		handlers.Help(s, m)
	case "save":
		handlers.SaveQuote(s, m)
	case "recent":
		if len(args) < 2 {
			args = append(args, "15")
		}
		handlers.RecentQuotes(s, m, args[1])
	}
}
