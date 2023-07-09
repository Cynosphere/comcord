package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ClearCommand(session *discordgo.Session) {
  fmt.Print("\n\r\033[H\033[2J")
}
