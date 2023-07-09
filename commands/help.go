package commands

import (
	"fmt"
	"strings"

	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
	"github.com/mgutz/ansi"
)

const format string = "  %s - %s%s"

func HelpCommand(session *discordgo.Session) {
  noColor := state.HasNoColor()

  fmt.Println("\r\nCOMcord (c)left 2023\n\r")

  index := 0
  for key, cmd := range GetAllCommands() {
    str := fmt.Sprintf(format, key, cmd.Description, "")
    length := len(str)
    padding := strings.Repeat(" ", 25 - length)

    if noColor {
      fmt.Printf(format, key, cmd.Description, padding)
    } else {
      coloredKey := ansi.Color(key, "yellow+b")
      fmt.Printf(format, coloredKey, cmd.Description, padding)
    }

    index++
    if index % 3 == 0 {
      fmt.Print("\n\r")
    }
  }
  if index % 3 != 0 {
    fmt.Print("\n\r")
  }

  fmt.Println("\r\nTo begin TALK MODE, press [SPACE]\n\r")
}
