package handlers

import "github.com/bwmarrin/discordgo"

// Type Aliases because typing out the full path and typename everytime is annoying
type (
	// Embed is a discordgo embed
	Embed = discordgo.MessageEmbed
	// EmbedField is a single field used in an embed
	EmbedField = discordgo.MessageEmbedField
)

// This comment is useless, you know what these are
const (
	// PREFIX is the prefix for every command
	PREFIX = "Q!"
	// COLOUR is a nice bright discord orange colour int
	COLOUR = 0xFAA61A
	// URL is the URL that all message links start with
	URL = "https://discord.com/channels/"
)
