package handlers

import (
	"github.com/bwmarrin/discordgo"

	"github.com/jawadcode/quotesbotv2/db"

	"fmt"
	"time"
)

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
	quote := db.QuoteSave{
		GuildID:   m.GuildID,
		ChannelID: chID,
		MessageID: message.ID,
		Content:   message.Content,
		Author:    message.Author.ID,
		AddedBy:   m.Message.Author.ID,
		AddedAt:   uint64(time.Now().UnixNano()),
	}
	// Insert into DB and catch any errors
	rows, err := db.DB.NamedQuery(
		`INSERT INTO quotes (guild_id, channel_id, message_id, content, author, added_by, added_at)
			 VALUES (:guild_id, :channel_id, :message_id, :content, :author, :added_by, :added_at)
		RETURNING id`,
		&quote,
	)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "An Error Occurred while saving the quote :(")
		fmt.Println(err.Error())
		return
	}

	var ID string
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
