package events

import (
	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
)

func MessageCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
  if msg.Author.ID == session.State.User.ID {
    return
  }

  channel, err := session.State.Channel(msg.ChannelID)
  if err != nil {
    return
  }

  isDM := channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM

  if state.IsInPrompt() {
    state.AddMessageToQueue(msg.Message)
  } else {
    lib.ProcessMessage(session, msg.Message, lib.MessageOptions{NoColor: state.HasNoColor()})
  }

  if isDM {
    state.SetLastDM(msg.ChannelID)
  }
}

func MessageUpdate(session *discordgo.Session, msg *discordgo.MessageUpdate) {

}
