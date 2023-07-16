package commands

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
)

func GetHistory(session *discordgo.Session, limit int, channel string) {
  messages, err := session.ChannelMessages(channel, 100, "", "", "")
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
    msgLines := lib.ProcessMessage(session, msg, lib.MessageOptions{NoColor: true, InHistory: true})
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

func HistoryCommand(session *discordgo.Session) {
  currentChannel := state.GetCurrentChannel()
  if currentChannel == "" {
    fmt.Print("<not in a channel>\n\r")
    return
  }

  GetHistory(session, 20, currentChannel)
}

func ExtendedHistoryCommand(session *discordgo.Session) {
  currentChannel := state.GetCurrentChannel()
  if currentChannel == "" {
    fmt.Print("<not in a channel>\n\r")
    return
  }

  lib.MakePrompt(session, ":lines> ", false, func(session *discordgo.Session, input string, interrupt bool) {
    fmt.Print("\r")
    limit, err := strconv.Atoi(input)

    if err != nil {
      fmt.Print("<not a number>\n\r")
    } else {
      GetHistory(session, limit, currentChannel)
    }
  })
}

func PeekHistory(session *discordgo.Session, guild, channel string) {
  target := ""

  channels := GetSortedChannels(session, guild, false, false)

  for _, c := range channels {
    if strings.Index(strings.ToLower(c.Name), strings.ToLower(channel)) > -1 {
      target = c.ID
      break
    }
  }

  if target == "" {
    fmt.Print("<channel not found>\n\r")
  } else {
    GetHistory(session, 20, target)
  }
}

func PeekCommand(session *discordgo.Session) {
  currentGuild := state.GetCurrentGuild()
  if currentGuild == "" {
    fmt.Print("<not in a guild>\n\r")
    return
  }

  lib.MakePrompt(session, ":peek> ", false, func(session *discordgo.Session, input string, interrupt bool) {
    fmt.Print("\r")

    if input != "" {
      PeekHistory(session, currentGuild, input)
    }
  })
}

func CrossPeekCommand(session *discordgo.Session) {
  lib.MakePrompt(session, ":guild> ", false, func(session *discordgo.Session, input string, interrupt bool) {
    fmt.Print("\r")

    if input != "" {
      targetGuild := ""

      for _, guild := range session.State.Guilds {
        if strings.Index(strings.ToLower(guild.Name), strings.ToLower(input)) > -1 {
          targetGuild = guild.ID
          break;
        }
      }

      if targetGuild == "" {
        fmt.Print("<guild not found>\n\r")
      } else {
        lib.MakePrompt(session, ":peek> ", false, func(session *discordgo.Session, input string, interrupt bool) {
          fmt.Print("\r")

          if input != "" {
            PeekHistory(session, targetGuild, input)
          }
        })
      }
    }
  })
}
