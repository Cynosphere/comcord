package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func TimeCommand(session *discordgo.Session) {
  now := time.Now().UTC()

  fmt.Printf("%s\n\r", now.Format("[Mon 02-Jan-06 15:04:05]"))
}
