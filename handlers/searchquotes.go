package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jawadcode/quotesbotv2/db"
)

// SearchQuotes handles the command "<PREFIX> search <query>" by using full-text search to find quotes relevant to the query
func SearchQuotes(s *discordgo.Session, m *discordgo.MessageCreate, q []string) {
	var quotes []db.QuotePreview
	chID := m.ChannelID
	// Convert query from form "this is a query" to "this & is & a & query"
	query := q[0]
	if len(q) > 1 {
		query = strings.Join(q, " & ")
	}
	// Fetch results that match full-text query
	err := db.DB.Select(
		&quotes,
		"SELECT content, author FROM quotes WHERE tsv @@ to_tsquery($1)",
		query,
	)
	if err != nil {
		s.ChannelMessageSend(chID, "An Error Occurred while searching for quotes")
	}
	quotesEmbed := embedQuotePreviews(
		s,
		"Search results for \""+strings.Join(q, " ")+"\":",
		" result(s)", &quotes,
	)
	s.ChannelMessageSendEmbed(chID, quotesEmbed)
}
