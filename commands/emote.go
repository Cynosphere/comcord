package commands

import (
	"fmt"

	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
)

func EmoteCommand(session *discordgo.Session) {
  channelId := state.GetCurrentChannel()
  if channelId == "" {
    fmt.Print("<not in a channel>\n\r")
    return
  }

  prompt := ":emote> "
  lib.MakePrompt(session, prompt, true, func(session *discordgo.Session, input string, interrupt bool) {
    if input == "" {
      if interrupt {
        fmt.Print("^C<no message sent>\n\r")
      } else {
        fmt.Print(prompt, "<no message sent>\n\r")
      }
    } else {
      fmt.Print(prompt, input, "\n\r")
      _, err := session.ChannelMessageSend(channelId, "*" + input + "*")

      if err != nil {
        fmt.Print("<failed to send message: ", err, ">\n\r")
      }

      // TODO: update afk state
    }
  })
}
