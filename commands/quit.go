package commands

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func QuitCommand(session *discordgo.Session) {
  session.Close()
  os.Exit(0)
}
