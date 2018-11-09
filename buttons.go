package buttons

import (
	"bytes"
	"crypto/sha1"

	"github.com/bwmarrin/discordgo"
)

type ButtonFunction func(s *discordgo.Session, r *discordgo.MessageReactionAdd)

type Button struct {
	Name     string
	Reaction string
}

func NewButton(s *discordgo.Session, name string, reaction string, callback ButtonFunction) Button {
	h := sha1.New()
	h.Write([]byte(name + reaction))
	buttonSig := h.Sum(nil)

	f := func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		// Reactions added by the bot do not count
		if r.UserID == s.State.User.ID {
			return
		}

		m, err := s.ChannelMessage(r.ChannelID, r.MessageID)
		if err != nil {
			return
		}

		// Reactions to messages not by the bot are irrelevant
		if m.Author.ID != s.State.User.ID {
			return
		}

		// The last embed in a given message is reserved for buttons
		// TODO: Think of alternate methods
		buttonEmbed := m.Embeds[len(m.Embeds)-1]
		for _, button := range buttonEmbed.Fields {
			h := sha1.New()
			h.Write([]byte(button.Name + r.Emoji.Name))
			signature := h.Sum(nil)
			if bytes.Equal(signature, buttonSig) {
				callback(s, r)
				return
			}
		}
	}

	s.AddHandler(f)
	return Button{
		Name:     name,
		Reaction: reaction,
	}
}

// AddButtons Creates an embed in the given message containing the buttons.
// This should be done last, after all other preperation for the button is complete.
func AddButton(message *discordgo.MessageEmbed, button Button) error {
	buttonField := discordgo.MessageEmbedField{
		Name:   button.Name,
		Value:  button.Reaction,
		Inline: true,
	}

	message.Fields = append(message.Fields, &buttonField)
	return nil
}
