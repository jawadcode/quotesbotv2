package handlers

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/jawadcode/quotesbotv2/db"
)

// GetQuote handles the command "<PREFIX> get <number>" by getting the quote with an id of <number>
func GetQuote(s *discordgo.Session, m *discordgo.MessageCreate, IDStr string) {
	var quote db.Quote
	// Convert ID to integer (and potentially prevent SQL injection)
	chID := m.ChannelID
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		s.ChannelMessageSend(chID, IDStr+" is not a valid number :(")
		Help(s, m)
		return
	}
	// Get quote with id of ID
	err = db.DB.Get(&quote, "SELECT * FROM quotes WHERE id=$1", ID)
	if err != nil {
		s.ChannelMessageSend(chID, "An Error Occurred while getting the quote")
		fmt.Println(err.Error())
		return
	}

	s.ChannelMessageSendEmbed(chID, embedQuote(s, chID, &quote))
}
