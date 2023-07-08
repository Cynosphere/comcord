package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
  "github.com/fatih/color"
)

func Ready(session *discordgo.Session, event *discordgo.Ready) {
  fmt.Print("Logged in as: ")
  color.Set(color.FgYellow)
  fmt.Printf("%s (%s)\n", session.State.User.Username, session.State.User.ID)
  color.Unset()
}
