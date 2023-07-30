package commands

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mgutz/ansi"
)

func SendMode() {
  client := state.GetClient()

  channelId := state.GetCurrentChannel()
  if channelId == "" {
    fmt.Print("<not in a channel>\n\r")
    return
  }
  parsedChannelId, err := discord.ParseSnowflake(channelId)
  if err != nil {
    fmt.Print("<failed to parse channel id: ", err.Error(), ">\n\r")
    return
  }

  channel, err := client.ChannelStore.Channel(discord.ChannelID(parsedChannelId))
  if err != nil {
    fmt.Print("<error getting channel: ", err.Error(), ">\n\r")
    return
  }

  guild, err := client.GuildStore.Guild(channel.GuildID)
  if err != nil {
    fmt.Print("<failed to get current guild: ", err.Error(), ">\n\r")
    return
  }

  self, err := client.MeStore.Me()
  if err != nil {
    fmt.Print("<failed to get self: ", err.Error(), ">\n\r")
    return
  }

  selfMember, err := client.MemberStore.Member(guild.ID, self.ID)
  if err != nil {
    fmt.Print("<failed to get self as member: ", err.Error(), ">\n\r")
    return
  }

  perms := discord.CalcOverwrites(*guild, *channel, *selfMember)
  cannotSend := !perms.Has(discord.PermissionSendMessages)
  if perms == 0 {
    cannotSend = false
  }

  if cannotSend {
    fmt.Print("<you do not have permission to send messages here>\n\r")
    return
  }

  length := utf8.RuneCountInString(self.Username) + 2
  curLength := state.GetNameLength()

  prompt := fmt.Sprintf("[%s]%s", self.Username, strings.Repeat(" ", (curLength - length) + 1))
  if !state.HasNoColor() {
    prompt = ansi.Color(prompt, "cyan+b")
  }

  lib.MakePrompt(prompt, true, func(input string, interrupt bool) {
    if input == "" {
      if interrupt {
        fmt.Print("^C<no message sent>\n\r")
      } else {
        fmt.Print(prompt, "<no message sent>\n\r")
      }
    } else {
      fmt.Print(prompt, input, "\n\r")
      _, err := client.SendMessage(channel.ID, input)

      if err != nil {
        fmt.Print("<failed to send message: ", err, ">\n\r")
      }

      // TODO: update afk state
    }
  })
}
