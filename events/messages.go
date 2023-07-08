package events

import (
	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
)

func MessageCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
  if (msg.Author.ID == session.State.User.ID) {
    return
  }

  channel, err := session.State.Channel(msg.ChannelID)
  if err != nil {
    return
  }

  if state.IsInPrompt() {
    state.AddMessageToQueue(*msg.Message)
  } else {
    // TODO
  }

  if channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM {
    state.SetLastDM(msg.ChannelID)
  }
}

func MessageUpdate(session *discordgo.Session, msg *discordgo.MessageUpdate) {

}
