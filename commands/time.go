package commands

import (
	"fmt"
	"time"
)

func TimeCommand() {
  now := time.Now().UTC()

  fmt.Printf("%s\n\r", now.Format("[Mon 02-Jan-06 15:04:05]"))
}
