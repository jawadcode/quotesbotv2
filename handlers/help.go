package handlers

import "github.com/bwmarrin/discordgo"

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
