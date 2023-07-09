package events

import (
	"fmt"

	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
	"github.com/mgutz/ansi"
)

func Ready(session *discordgo.Session, event *discordgo.Ready) {
  fmt.Printf("\rLogged in as: %s%s (%s)%s\n\r", ansi.ColorCode("yellow"), session.State.User.Username, session.State.User.ID, ansi.ColorCode("reset"))

  state.SetNameLength(len(session.State.User.Username) + 2)

  defaultGuild := state.GetConfigValue("defaultGuild")
  defaultChannel := state.GetConfigValue("defaultChannel")
  if defaultGuild != "" {
    //var guild discordgo.Guild
    hasGuild := false
    for _, g := range session.State.Guilds {
      if g.ID == defaultGuild {
        //guild = *g
        hasGuild = true
        break
      }
    }
    if hasGuild {
      if defaultChannel != "" {
        state.SetCurrentChannel(defaultChannel)
        state.SetLastChannel(defaultGuild, defaultChannel)
      }
    } else {
      fmt.Println("\r% This account is not in the defined default guild.")
    }
  } else {
    if defaultChannel != "" {
      fmt.Println("\r% Default channel defined without defining default guild.")
    }
  }
}
