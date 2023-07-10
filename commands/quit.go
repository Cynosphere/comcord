package commands

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func QuitCommand(session *discordgo.Session) {
  fmt.Print("Unlinking TTY...\n\r")
  session.Close()
  os.Exit(0)
}
