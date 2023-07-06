package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
  "github.com/fatih/color"
)

func Ready(session *discordgo.Session, event *discordgo.Ready) {
  fmt.Print("Logged in as: ")
  color.Yellow("%s (%s)", session.State.User.Username, session.State.User.ID)
}
