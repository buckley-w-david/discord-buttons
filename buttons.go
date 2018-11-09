package discordbuttons

import (
	"github.com/bwmarrin/discordgo"
)

type Button struct {
	Data     interface{}
	Reaction string
	Callback func(s *discordgo.Session, r *discordgo.MessageReactionAdd, mID string, cID string, data interface{})
}

// AddButton adds a listener for a specific reaction to a given message
func AddButton(s *discordgo.Session, messageID string, channelID string, button Button, once bool) error {
	var remove func()
	f := func(sess *discordgo.Session, r *discordgo.MessageReactionAdd) {
		// Reactions added by the bot do not count
		// Reactions to other messages are irrelevant
		// Reaction recieved needs to match reaction of the button
		if r.UserID == s.State.User.ID || messageID != r.MessageID || channelID != r.ChannelID || r.Emoji.Name != button.Reaction {
			return
		}

		button.Callback(s, r, messageID, channelID, button.Data)
		if remove != nil && once {
			remove()
		}
		return
	}
	remove = s.AddHandler(f)
	s.MessageReactionAdd(channelID, messageID, button.Reaction)

	return nil
}
