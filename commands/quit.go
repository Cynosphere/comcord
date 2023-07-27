package commands

import (
	"fmt"
	"os"

	"github.com/Cynosphere/comcord/state"
)

func QuitCommand() {
  client := state.GetClient()

  fmt.Print("Unlinking TTY...\n\r")
  client.Close()
  os.Exit(0)
}
