package commands

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
	"github.com/mgutz/ansi"
)

func SendMode(session *discordgo.Session) {
  channelId := state.GetCurrentChannel()
  if channelId == "" {
    fmt.Print("<not in a channel>\n\r")
  return
  }

  state.SetInPrompt(true)

  length := utf8.RuneCountInString(session.State.User.Username) + 2
  curLength := state.GetNameLength()

  prompt := fmt.Sprintf("[%s]%s", session.State.User.Username, strings.Repeat(" ", (curLength - length) + 1))
  if !state.HasNoColor() {
    prompt = ansi.Color(prompt, "cyan+b")
  }

  lib.MakePrompt(session, prompt, true, func(session *discordgo.Session, input string, interrupt bool) {
    if input == "" {
      if interrupt {
        fmt.Print("^C<no message sent>\n\r")
      } else {
        fmt.Print(prompt, "<no message sent>\n\r")
      }
    } else {
      fmt.Print(prompt, input, "\n\r")
      _, err := session.ChannelMessageSend(channelId, input)

      if err != nil {
        fmt.Print("<failed to send message: ", err, ">\n\r")
      }

      // TODO: update afk state
    }
  })
}
