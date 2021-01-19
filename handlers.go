package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Type Aliases because typing out the full path and typename everytime is annoying
type (
	// Embed is a discordgo embed
	Embed = discordgo.MessageEmbed
	// EmbedField is a single field used in an embed
	EmbedField = discordgo.MessageEmbedField
)

// Help handles the command "<PREFIX> help" and it displays a description of the bot andlists all of the commands
func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Array of `MessageEmbedField`s to be used as as `Fields` in help embed
	helpEmbedFields := []*EmbedField{
		{
			Name:  "`" + PREFIX + " help`",
			Value: "Pretty self explanatory",
		},
		{
			Name:  "`" + PREFIX + " save <@user>`",
			Value: "Save message that the command is replying to or save the most recent message by `@user`",
		},
		{
			Name:  "`" + PREFIX + " get <number>`",
			Value: "Get quote `#<number>`",
		},
		{
			Name:  "`" + PREFIX + " search <query>`",
			Value: "Returns all quotes that roughly match `<query>`",
		},
	}

	helpEmbed := Embed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "QuotesBot V2 Help:",
		Description: "QuotesBot is a bot that allows you to save and recall quotes ~~to blackmail your friends~~ for future use.",
		Color:       COLOUR,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "We take no responsibility for the stuff you say",
		},
		Fields: helpEmbedFields,
	}

	// Send embed
	s.ChannelMessageSendEmbed(m.ChannelID, &helpEmbed)
}

// SaveQuote handles the command "<PREFIX> save <optional @user>" by saving the referenced message or the most recent message from the mentioned user
func SaveQuote(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error
	// Message Content to be saved to database
	var message *discordgo.Message
	// This will prevent unnecessary dereferencing and it looks slightly cleaner
	chID := m.ChannelID
	// Check if message referenced another message
	referencedMessage := m.MessageReference
	// If not then check if the message contains a mention
	if referencedMessage == nil {
		mentions := m.Mentions
		if len(mentions) < 1 {
			s.ChannelMessageSend(chID, "You must mention at least one user or reply to a message :(")
			Help(s, m)
			return
		}
		mention := mentions[0]
		// Get the 30 most recent messages and find one by the mentioned user
		recent, err := s.ChannelMessages(chID, 30, m.Message.ID, "", "")
		if err != nil {
			s.ChannelMessageSend(chID, "An Error Occurred while getting recent messages :(")
			return
		}

		// Loop through those messages and find the most recent one by
		for _, msg := range recent {
			if msg.Author.ID == mention.ID && !msg.Author.Bot {
				message = msg
				break
			}
		}
	} else {
		// If there is a referenced message then grab that message
		message, err = s.ChannelMessage(chID, referencedMessage.MessageID)
		if err != nil {
			s.ChannelMessageSend(chID, "An Error Occurred while getting the referenced message :(")
			return
		}
	}

	if message == nil {
		s.ChannelMessageSend(chID, "Message Not Found :(")
		return
	}

	// Create Quote struct to insert into the database
	quote := QuoteSave{
		Content: message.Content,
		Author:  message.Author.ID,
		AddedBy: m.Message.Author.ID,
		AddedAt: uint64(time.Now().UnixNano()),
	}

	// Insert into DB and catch any errors
	rows, err := DB.NamedQuery(
		`INSERT INTO quotes (content, author, added_by, added_at)
			 VALUES (:content, :author, :added_by, :added_at)
		RETURNING id`,
		&quote,
	)

	if rows.Err() != nil || err != nil {
		s.ChannelMessageSend(m.ChannelID, "An Error Occurred while saving the quote :(")
		return
	}

	var ID string

	// Using if instead of for because there is only one row to be read
	if rows.Next() {
		rows.Scan(&ID)
	} else {
		s.ChannelMessageSend(m.ChannelID, "An Error Occurred while saving the quote :(")
		return
	}

	s.ChannelMessageSend(
		chID,
		fmt.Sprintf("<@%s>, saved quote number %s, by <@%s>", m.Author.ID, ID, message.Author.ID),
	)
}

// RecentQuotes handles the command "<PREFIX> recent <optional limit>" by sending an embed containing previews of the latest quotes (with the limit being either the one specified or 15)
func RecentQuotes(s *discordgo.Session, m *discordgo.MessageCreate, limitStr string) {
	var quotes []QuotePreview
	chID := m.ChannelID
	// Convert limit to integer
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit > 30 {
		limit = 15
	}
	// Get most recent quotes with limit and then assign it to `quotes`
	err = DB.Select(&quotes, "SELECT content, author FROM quotes ORDER BY id DESC LIMIT $1", limit)
	if err != nil {
		s.ChannelMessageSend(chID, "An Error Occurred while getting the quotes")
		return
	}

	limit = len(quotes)
	quotesEmbed := Embed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Recent Quotes: ",
		Description: limitStr + " most recent quotes",
		Fields:      make([]*discordgo.MessageEmbedField, limit),
	}
	// Premature optimisation at its finest :)
	// Create a string to string map (to get usernames of authors as efficiently as possible) and make it at least as big as the limit to avoid resizing
	authorsMap := make(map[string]string, limit)
	// Loop through quotes and add values one by one to Embed Fields
	for i, quote := range quotes {
		var username string
		// If user ID is in map then get username and assign it to username
		if val, ok := authorsMap[quote.Author]; ok {
			username = val
		} else {
			// User ID is not in map so find it using API
			user, err := s.User(quote.Author)
			if err != nil {
				username = "*unknown*"
			} else {
				// Assign the username and then add the ID and username to the map for future use
				username = user.Username
			}
			authorsMap[quote.Author] = username
		}
		// Add Field to Embed Fields
		quotesEmbed.Fields[i] = &discordgo.MessageEmbedField{
			Name:   "@" + username,
			Value:  quote.Content,
			Inline: true,
		}
	}
	s.ChannelMessageSendEmbed(chID, &quotesEmbed)
}
