package commands

import (
	"fmt"

	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/diamondburned/arikawa/v3/discord"
)

func EmoteCommand() {
  channelId := state.GetCurrentChannel()
  if channelId == "" {
    fmt.Print("<not in a channel>\n\r")
    return
  }

  prompt := ":emote> "
  lib.MakePrompt(prompt, true, func(input string, interrupt bool) {
    if input == "" {
      if interrupt {
        fmt.Print("^C<no message sent>\n\r")
      } else {
        fmt.Print(prompt, "<no message sent>\n\r")
      }
    } else {
      fmt.Print(prompt, input, "\n\r")
      client := state.GetClient()

      snowflake, err := discord.ParseSnowflake(channelId)
      if err != nil {
        fmt.Print("<failed to parse channel id: ", err.Error(), ">\n\r")
        return
      }

      _, err = client.SendMessage(discord.ChannelID(snowflake), "*" + input + "*")

      if err != nil {
        fmt.Print("<failed to send message: ", err.Error(), ">\n\r")
      }

      // TODO: update afk state
    }
  })
}
