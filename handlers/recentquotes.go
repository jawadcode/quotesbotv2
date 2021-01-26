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
	// Rare case but has to be handled
	if len(quotes) == 0 {
		s.ChannelMessageSend(chID, "0 quotes")
	}
	quotesEmbed := embedQuotePreviews(s, "Recent Quotes", "most recent quotes", &quotes)
	s.ChannelMessageSendEmbed(chID, quotesEmbed)
}
