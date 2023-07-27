package commands

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/diamondburned/arikawa/v3/discord"
)

func GetHistory(limit int, channel string) {
  client := state.GetClient()

  parsedChannelId, err := discord.ParseSnowflake(channel)
  if err != nil {
    fmt.Print("<failed to parse channel id: ", err.Error(), ">\n\r")
    return
  }

  messages, err := client.Messages(discord.ChannelID(parsedChannelId), 50)
  if err != nil {
    fmt.Print("<failed to get messages: ", err.Error(), ">\n\r")
    return
  }

  for i, j := 0, len(messages) - 1; i < j; i, j = i + 1, j - 1 {
    messages[i], messages[j] = messages[j], messages[i]
  }

  fmt.Print("--Beginning-Review", strings.Repeat("-", 62), "\n\r")

  lines := make([]string, 0)
  for _, msg := range messages {
    msgLines := lib.ProcessMessage(msg, lib.MessageOptions{NoColor: true, InHistory: true})
    for _, line := range msgLines {
      lines = append(lines, line)
    }
  }

  length := len(lines)
  startIndex := int(math.Max(float64(0), float64(length - limit)))
  for _, line := range lines[startIndex:] {
    fmt.Print(line)
  }

  fmt.Print("--Review-Complete", strings.Repeat("-", 63), "\n\r")
}

func HistoryCommand() {
  currentChannel := state.GetCurrentChannel()
  if currentChannel == "" {
    fmt.Print("<not in a channel>\n\r")
    return
  }

  GetHistory(20, currentChannel)
}

func ExtendedHistoryCommand() {
  currentChannel := state.GetCurrentChannel()
  if currentChannel == "" {
    fmt.Print("<not in a channel>\n\r")
    return
  }

  lib.MakePrompt(":lines> ", false, func(input string, interrupt bool) {
    fmt.Print("\r")
    limit, err := strconv.Atoi(input)

    if err != nil {
      fmt.Print("<not a number>\n\r")
    } else {
      GetHistory(limit, currentChannel)
    }
  })
}

func PeekHistory(guild, channel string) {
  target := ""

  channels := GetSortedChannels(guild, false, false)

  for _, c := range channels {
    if strings.Index(strings.ToLower(c.Name), strings.ToLower(channel)) > -1 {
      target = c.ID.String()
      break
    }
  }

  if target == "" {
    fmt.Print("<channel not found>\n\r")
  } else {
    GetHistory(20, target)
  }
}

func PeekCommand() {
  currentGuild := state.GetCurrentGuild()
  if currentGuild == "" {
    fmt.Print("<not in a guild>\n\r")
    return
  }

  lib.MakePrompt(":peek> ", false, func(input string, interrupt bool) {
    fmt.Print("\r")

    if input != "" {
      PeekHistory(currentGuild, input)
    }
  })
}

func CrossPeekCommand() {
  client := state.GetClient()

  lib.MakePrompt(":guild> ", false, func(input string, interrupt bool) {
    fmt.Print("\r")

    if input != "" {
      targetGuild := ""

      guilds, err := client.GuildStore.Guilds()
      if err != nil {
        fmt.Print("<failed to get guilds: ", err.Error(), ">\n\r")
        return
      }

      for _, guild := range guilds {
        if strings.Index(strings.ToLower(guild.Name), strings.ToLower(input)) > -1 {
          targetGuild = guild.ID.String()
          break;
        }
      }

      if targetGuild == "" {
        fmt.Print("<guild not found>\n\r")
      } else {
        lib.MakePrompt(":peek> ", false, func(input string, interrupt bool) {
          fmt.Print("\r")

          if input != "" {
            PeekHistory(targetGuild, input)
          }
        })
      }
    }
  })
}
