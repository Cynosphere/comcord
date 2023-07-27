package events

import (
	"fmt"
	"unicode/utf8"

	"github.com/Cynosphere/comcord/commands"
	"github.com/Cynosphere/comcord/state"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/mgutz/ansi"
)

func Ready(event *gateway.ReadyEvent) {
  client := state.GetClient()
  self, err := client.Me()
  if err != nil {
    fmt.Print("\r% Failed to get self: ", err.Error(), "\n\r")
    return
  }

  fmt.Printf("\rLogged in as: %s\n\r", ansi.Color(fmt.Sprintf("%s (%s)", self.Username, self.ID), "yellow"))

  state.SetNameLength(utf8.RuneCountInString(self.Username) + 2)

  commands.ListGuildsCommand()

  defaultGuild := state.GetConfigValue("defaultGuild")
  defaultChannel := state.GetConfigValue("defaultChannel")
  if defaultGuild != "" {
    parsedGuildId, err := discord.ParseSnowflake(defaultGuild)
    if err != nil {
      fmt.Print("\r% Failed to parse guild ID: ", err.Error(), "\n\r")
      return
    }

    guild, err := client.Guild(discord.GuildID(parsedGuildId))
    if err == nil {
      if defaultChannel != "" {
        state.SetCurrentChannel(defaultChannel)
        state.SetLastChannel(defaultGuild, defaultChannel)
      }
      commands.SwitchGuild(guild.Name)
    } else {
      fmt.Println("\r% This account is not in the defined default guild.")
    }
  } else {
    if defaultChannel != "" {
      fmt.Println("\r% Default channel defined without defining default guild.")
    }
  }
}
