package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func HelpCommand(session *discordgo.Session) {
  fmt.Println("\r\nCOMcord (c)left 2023\n\r")

  index := 0
  for key, cmd := range GetAllCommands() {
    str := fmt.Sprintf("  %s - %s", key, cmd.Description)
    length := len(str)

    fmt.Print("  ")
    color.Set(color.FgYellow, color.Bold)
    fmt.Print(key)
    color.Unset()
    fmt.Printf(" - %s", cmd.Description)
    fmt.Print(strings.Repeat(" ", 25 - length))

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
