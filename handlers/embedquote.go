package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jawadcode/quotesbotv2/db"
)

func embedQuote(s *discordgo.Session, chID string, quote *db.Quote) *discordgo.MessageEmbed {
	author, err := s.User(quote.Author)
	addedBy, err := s.User(quote.AddedBy)
	channel, err := s.Channel(quote.ChannelID)

	if err != nil {
		return nil
	}

	link := URL + quote.GuildID + "/" + quote.ChannelID + "/" + quote.MessageID

	return &discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    author.Username,
			IconURL: author.AvatarURL("128"),
		},
		Title:       fmt.Sprintf("Quote #%d:", quote.ID),
		Description: quote.Content,
		Fields: []*EmbedField{
			{
				Name:  "Info:",
				Value: "added by: <@" + addedBy.ID + ">, in <#" + channel.ID + ">",
			},
			{
				Name:  "Link:",
				Value: link,
			},
		},
		Color: COLOUR,
	}
}
