package commands

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

type GuildListing struct {
  Name string
  Members int
  Online int
}

func ListGuildsCommand(session *discordgo.Session) {
  longest := 0
  guilds := make([]GuildListing, 0)

  for _, guild := range session.State.Guilds {
    length := utf8.RuneCountInString(guild.Name)
    if length > longest {
      longest = length
    }

    guildWithCounts, err := session.GuildWithCounts(guild.ID)
    if err != nil {
      guilds = append(guilds, GuildListing{
        Name: guild.Name,
        Members: guild.MemberCount,
        Online: 0,
      })
      return
    }

    guilds = append(guilds, GuildListing{
      Name: guild.Name,
      Members: guildWithCounts.ApproximateMemberCount,
      Online: guildWithCounts.ApproximatePresenceCount,
    })
  }

  fmt.Print("\n\r")
  fmt.Printf("  %*s  online  total\n\r", longest, "guild-name")
  fmt.Print(strings.Repeat("-", 80) + "\n\r")
  for _, guild := range guilds {
    fmt.Printf("  %*s  %6d  %5d\n\r", longest, guild.Name, guild.Online, guild.Members)
  }
  fmt.Print(strings.Repeat("-", 80) + "\n\r")
  fmt.Print("\n\r")
}
