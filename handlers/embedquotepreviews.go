package handlers

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/jawadcode/quotesbotv2/db"
)

func embedQuotePreviews(
	s *discordgo.Session,
	title string,
	description string,
	quotes *[]db.QuotePreview,
) *Embed {
	limit := len(*quotes)
	quotesEmbed := Embed{
		Type:        discordgo.EmbedTypeRich,
		Title:       title,
		Description: strconv.Itoa(limit) + " " + description,
		Color:       COLOUR,
		Fields:      make([]*discordgo.MessageEmbedField, limit),
	}
	// Premature optimisation at its finest :)
	// Create a string to string map (to get usernames of authors as efficiently as possible) and make it at least as big as the limit to avoid resizing
	authorsMap := make(map[string]string, limit)
	// Loop through quotes and add values one by one to Embed Fields
	for i, quote := range *quotes {
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
	return &quotesEmbed
}
