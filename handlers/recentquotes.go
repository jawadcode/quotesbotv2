package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jawadcode/quotesbotv2/db"

	"strconv"
)

// RecentQuotes handles the command "<PREFIX> recent <optional limit>" by sending an embed containing previews of the latest quotes (with the limit being either the one specified or 15)
func RecentQuotes(s *discordgo.Session, m *discordgo.MessageCreate, limitStr string) {
	var quotes []db.QuotePreview
	chID := m.ChannelID
	// Convert limit to integer
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit > 30 {
		limit = 15
	}
	// Get most recent quotes with limit and then assign it to `quotes`
	err = db.DB.Select(&quotes, "SELECT content, author FROM quotes ORDER BY id DESC LIMIT $1", limit)
	if err != nil {
		s.ChannelMessageSend(chID, "An Error Occurred while getting the quotes")
		return
	}

	limit = len(quotes)
	quotesEmbed := Embed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Recent Quotes: ",
		Description: limitStr + " most recent quotes",
		Color:       COLOUR,
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
