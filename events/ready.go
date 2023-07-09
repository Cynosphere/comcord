package events

import (
	"fmt"
	"unicode/utf8"

	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
	"github.com/mgutz/ansi"
)

func Ready(session *discordgo.Session, event *discordgo.Ready) {
  fmt.Printf("\rLogged in as: %s\n\r", ansi.Color(fmt.Sprintf("%s (%s)", session.State.User.Username, session.State.User.ID), "yellow"))

  state.SetNameLength(utf8.RuneCountInString(session.State.User.Username) + 2)

  defaultGuild := state.GetConfigValue("defaultGuild")
  defaultChannel := state.GetConfigValue("defaultChannel")
  if defaultGuild != "" {
    _, err := session.State.Guild(defaultGuild)
    if err == nil {
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
