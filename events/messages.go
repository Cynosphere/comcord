package events

import (
	"fmt"

	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/discord"
)

func MessageCreate(msg *gateway.MessageCreateEvent) {
  client := state.GetClient()
  self, err := client.MeStore.Me()
  if err != nil {
    return
  }

  if msg.Author.ID == self.ID {
    return
  }

  channel, err := client.ChannelStore.Channel(msg.ChannelID)
  if err != nil {
    return
  }

  isDM := channel.Type == discord.DirectMessage || channel.Type == discord.GroupDM

  if state.IsInPrompt() {
    state.AddMessageToQueue(msg.Message)
  } else {
    lines := lib.ProcessMessage(msg.Message, lib.MessageOptions{NoColor: state.HasNoColor()})
    for _, line := range lines {
      fmt.Print(line)
    }
  }

  if isDM {
    state.SetLastDM(msg.ChannelID.String())
  }
}

func MessageUpdate(msg *gateway.MessageUpdateEvent) {
  client := state.GetClient()
  self, err := client.MeStore.Me()
  if err != nil {
    return
  }

  if msg.Author.ID == self.ID {
    return
  }

  old, err := client.MessageStore.Message(msg.ChannelID, msg.ID)
  if err != nil {
    return
  }

  if msg.Content == old.Content {
    return
  }

  // dont process embed updates as messages
  if !msg.EditedTimestamp.IsValid() {
    return
  }

  channel, err := client.ChannelStore.Channel(msg.ChannelID)
  if err != nil {
    return
  }

  isDM := channel.Type == discord.DirectMessage || channel.Type == discord.GroupDM

  if state.IsInPrompt() {
    state.AddMessageToQueue(msg.Message)
  } else {
    lines := lib.ProcessMessage(msg.Message, lib.MessageOptions{NoColor: state.HasNoColor()})
    for _, line := range lines {
      fmt.Print(line)
    }
  }

  if isDM {
    state.SetLastDM(msg.ChannelID.String())
  }
}
