package handlers

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/jawadcode/quotesbotv2/db"
)

// Get Quote handles the command "<PREFIX> get <number>" by getting the quote with an id of <number>
func GetQuote(s *discordgo.Session, m *discordgo.MessageCreate, IDStr string) {
	var quote db.Quote

	chID := m.ChannelID
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		s.ChannelMessageSend(chID, IDStr+" is not a valid number :(")
		Help(s, m)
		return
	}

	err = db.DB.Get(&quote, "SELECT * FROM quotes WHERE id=$1", ID)

	quoteEmbed := Embed{
		Title: "",
	}
}
