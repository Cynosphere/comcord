package commands

import (
	"fmt"
)

func ClearCommand() {
  fmt.Print("\n\r\033[H\033[2J")
}
