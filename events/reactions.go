package events

import (
	"fmt"
	"time"

	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func ReactionAdd(event *gateway.MessageReactionAddEvent) {
  client := state.GetClient()
  currentChannel := state.GetCurrentChannel()

  if event.ChannelID.String() != currentChannel {
    return
  }

  emote := event.Emoji.Name
  if event.Emoji.IsCustom() {
    emote = ":" + emote + ":"
  }

  now := time.Now()
  nowSnowflake := discord.NewSnowflake(now)

  message, err := client.MessageStore.Message(event.ChannelID, event.MessageID)
  if err != nil {
    message, err = client.Message(event.ChannelID, event.MessageID)
    if err != nil {
      return
    }
  }

  msg := discord.Message{
    Content: fmt.Sprintf("*reacted with %s*", emote),
    Author: event.Member.User,
    ChannelID: event.ChannelID,
    GuildID: event.GuildID,
    ID: discord.MessageID(nowSnowflake),
    ReferencedMessage: message,
    Type: discord.InlinedReplyMessage,
    Timestamp: discord.Timestamp(now),
  }

  if state.IsInPrompt() {
    state.AddMessageToQueue(msg)
  } else {
    lines := lib.ProcessMessage(msg, lib.MessageOptions{NoColor: state.HasNoColor()})
    for _, line := range lines {
      fmt.Print(line)
    }
  }
}
