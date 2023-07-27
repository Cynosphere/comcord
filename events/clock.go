package events

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/Cynosphere/comcord/state"
)

var sentTime bool = false

func SetupClock() {
  clock := time.NewTicker(500 * time.Millisecond)
  go func() {
    for {
      select {
        case <- clock.C: {
          now := time.Now().UTC()
          if now.Minute() % 15 == 0 && now.Second() < 2 && !sentTime {
            if state.IsInPrompt() {
              // TODO
            } else {
              fmt.Printf("%s\n\r", now.Format("[Mon 02-Jan-06 15:04:05]"))
            }

            client := state.GetClient()
            self, err := client.MeStore.Me()
            if err != nil {
              return
            }

            state.SetNameLength(utf8.RuneCountInString(self.Username) + 2)
            sentTime = true
          } else if now.Second() > 2 && sentTime {
            sentTime = false
          }
        }
      }
    }
  }()
}
