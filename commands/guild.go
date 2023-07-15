package commands

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/Cynosphere/comcord/lib"
	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
)

var REGEX_EMOTE = regexp.MustCompile(`<(?:\x{200b}|&)?a?:(\w+):(\d+)>`)

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

func GetSortedChannels(session *discordgo.Session, guildId string, withCategories bool) []*discordgo.Channel {
  channels := make([]*discordgo.Channel, 0)
  guild, err := session.State.Guild(guildId)
  if err != nil {
    return channels
  }

  if withCategories {
    categories := make(map[string][]*discordgo.Channel)

    for _, channel := range guild.Channels {
      categoryID := "0"
      if channel.ParentID != "" {
        categoryID = channel.ParentID
      }

      _, has := categories[categoryID]
      if !has {
        categories[categoryID] = make([]*discordgo.Channel, 0)
      }

      if channel.Type != discordgo.ChannelTypeGuildText && channel.Type != discordgo.ChannelTypeGuildNews {
        continue
      }
      categories[categoryID] = append(categories[categoryID], channel)
    }

    for id, channels := range categories {
      sort.Slice(channels, func(i, j int) bool {
        return channels[i].Position < channels[j].Position
      })
      categoryChannels := make([]*discordgo.Channel, 0)
      if id != "0" {
        for _, channel := range guild.Channels {
          if channel.ID == id {
            categoryChannels = append(categoryChannels, channel)
            break
          }
        }
      }
      for _, channel := range channels {
        categoryChannels = append(categoryChannels, channel)
      }
      categories[id] = categoryChannels
    }

    keys := make([]string, 0, len(categories) - 1)
    for id := range categories {
      if id == "0" {
        continue
      }
      keys = append(keys, id)
    }
    sort.Slice(keys, func(i, j int) bool {
      ca, _ := session.State.Channel(keys[i])
      cb, _ := session.State.Channel(keys[j])

      return ca.Position < cb.Position
    })
    sortedCategories := make(map[string][]*discordgo.Channel)
    sortedCategories["0"] = categories["0"]

    for _, id := range keys {
      sortedCategories[id] = categories[id]
    }

    for _, categoryChannels := range sortedCategories {
      for _, channel := range categoryChannels {
        channels = append(channels, channel)
      }
    }
  } else {
    for _, channel := range guild.Channels {
      if channel.Type != discordgo.ChannelTypeGuildText && channel.Type != discordgo.ChannelTypeGuildNews {
        continue
      }
      channels = append(channels, channel)
    }

    sort.Slice(channels, func(i, j int) bool {
      return channels[i].Position < channels[j].Position
    })
  }

  return channels
}

func ListChannelsCommand(session *discordgo.Session) {
  currentGuild := state.GetCurrentGuild()
  if currentGuild == "" {
    fmt.Print("<not in a guild>\n\r")
    return
  }

  longest := 0
  channels := GetSortedChannels(session, currentGuild, true)

  for _, channel := range channels {
    perms, err := session.State.UserChannelPermissions(session.State.User.ID, channel.ID)
    if err != nil {
      continue
    }

    private := perms & discordgo.PermissionViewChannel == 0
    category := channel.Type == discordgo.ChannelTypeGuildCategory

    catLen := 0
    if category {
      catLen = 6
    }

    privLen := 0
    if private {
      privLen = 1
    }
    length := utf8.RuneCountInString(channel.Name) + privLen + catLen

    if length > longest {
      longest = int(math.Min(25, float64(length)))
    }
  }

  fmt.Print("\n\r")
  fmt.Printf("  %*s    created  topic\n\r", longest, "channel-name")
  fmt.Print(strings.Repeat("-", 80) + "\n\r")
  for _, channel := range channels {
    perms, err := session.State.UserChannelPermissions(session.State.User.ID, channel.ID)
    if err != nil {
      continue
    }

    private := perms & discordgo.PermissionViewChannel == 0
    category := channel.Type == discordgo.ChannelTypeGuildCategory
    topic := REGEX_EMOTE.ReplaceAllString(channel.Topic, ":$1:")
    topic = strings.ReplaceAll(topic, "\n", " ")
    name := channel.Name
    if category {
      name = "-- " + name + " --"
    }
    if private {
      name = "*" + name
    }

    nameLength := utf8.RuneCountInString(name)
    if nameLength > 25 {
      name = name[:24] + "\u2026"
    }

    topicLength := utf8.RuneCountInString(topic)
    longestTopic := 80 - (longest + 5) - 11
    if topicLength > longestTopic {
      topic = topic[:(longestTopic - 1)] + "\u2026"
    }

    created := "??-???-??"
    timestamp, err := discordgo.SnowflakeTimestamp(channel.ID)
    if err == nil {
      created = timestamp.Format("02-Jan-06")
    }

    fmt.Printf("  %*s  %s  %s\n\r", longest, name, created, topic)
  }
  fmt.Print(strings.Repeat("-", 80) + "\n\r")
  fmt.Print("\n\r")
}

func ListUsersCommand(session *discordgo.Session) {

}

func SwitchGuild(session *discordgo.Session, input string) {
  if input == "" {
    ListChannelsCommand(session)
    ListUsersCommand(session)
  } else {
    target := ""

    for _, guild := range session.State.Guilds {
      if strings.Index(strings.ToLower(guild.Name), strings.ToLower(input)) > -1 {
        target = guild.ID
        break;
      }
    }

    if target == "" {
      fmt.Print("<guild not found>\n\r")
    } else {
      state.SetCurrentGuild(target)
      last := state.GetLastChannel(target)
      if last == "" {
        channels := GetSortedChannels(session, target, false)
        topChannel := channels[0]

        state.SetCurrentChannel(topChannel.ID)
        state.SetLastChannel(target, topChannel.ID)
      } else {
        state.SetCurrentChannel(last)
      }

      ListChannelsCommand(session)
      ListUsersCommand(session)

      // TODO: update presence
    }
  }
}

func SwitchGuildsCommand(session *discordgo.Session) {
  lib.MakePrompt(session, ":guild> ", false, func(session *discordgo.Session, input string, interrupt bool) {
    fmt.Print("\r")
    SwitchGuild(session, input)
  })
}
