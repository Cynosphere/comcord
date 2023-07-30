package events

import (
	"github.com/diamondburned/ningen/v3"
)

func Setup(session *ningen.State) {
  session.PreHandler.AddHandler(Ready)
  session.PreHandler.AddHandler(MessageCreate)
  session.PreHandler.AddHandler(MessageUpdate)
  session.PreHandler.AddHandler(ReactionAdd)
  SetupClock()
}
