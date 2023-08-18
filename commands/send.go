package commands

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mgutz/ansi"
)

var REGEX_MENTION = regexp.MustCompile("@([a-z0-9._]{1,32})")

func SendMode() {
  client := state.GetClient()

  currentGuild := state.GetCurrentGuild()

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

  client.Typing(channel.ID)

  lib.MakePrompt(prompt, true, func(input string, interrupt bool) {
    if input == "" {
      if interrupt {
        fmt.Print("^C<no message sent>\n\r")
      } else {
        fmt.Print(prompt, "<no message sent>\n\r")
      }
    } else {
      fmt.Print(prompt, input, "\n\r")

      input := REGEX_MENTION.ReplaceAllStringFunc(input, func(match string) string {
        matches := REGEX_MENTION.FindStringSubmatch(match)
        username := matches[1]

        parsedGuildId, err := discord.ParseSnowflake(currentGuild)
        if err != nil {
          return match
        }

        members, err := client.MemberStore.Members(discord.GuildID(parsedGuildId))
        if err != nil {
          return match
        }

        for _, member := range members {
          if member.User.Username == username {
            return member.User.ID.Mention()
          }
        }

        return match
      })

      _, err := client.SendMessage(channel.ID, input)

      if err != nil {
        fmt.Print("<failed to send message: ", err, ">\n\r")
      }

      // TODO: update afk state
    }
  })
}
