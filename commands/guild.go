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
	"github.com/diamondburned/arikawa/v3/discord"
	tsize "github.com/kopoli/go-terminal-size"
	"github.com/mgutz/ansi"
)

var REGEX_EMOTE = regexp.MustCompile(`<(?:\x{200b}|&)?a?:(\w+):(\d+)>`)

type GuildListing struct {
  Name string
  Members int
  Online int
}

func ListGuildsCommand() {
  client := state.GetClient()

  longest := 0
  guilds := make([]GuildListing, 0)

  clientGuilds, err := client.Guilds()
  if err != nil {
    fmt.Print("<failed to get guilds: ", err.Error(), ">\n\r")
    return
  }

  for _, guild := range clientGuilds {
    length := utf8.RuneCountInString(guild.Name)
    if length > longest {
      longest = length
    }

    withCount, err := client.GuildWithCount(guild.ID)
    if err == nil {
      guilds = append(guilds, GuildListing{
        Name: guild.Name,
        Members: int(withCount.ApproximateMembers),
        Online: int(withCount.ApproximatePresences),
      })
    } else {
      guilds = append(guilds, GuildListing{
        Name: guild.Name,
        Members: int(guild.ApproximateMembers),
        Online: int(guild.ApproximatePresences),
      })
    }
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

func GetSortedChannels(guildId string, withCategories bool, withPrivate bool) []discord.Channel {
  client := state.GetClient()
  channels := make([]discord.Channel, 0)

  guildSnowflake, err := discord.ParseSnowflake(guildId)
  if err != nil {
    return channels
  }

  parsedGuildId := discord.GuildID(guildSnowflake)

  guild, err := client.GuildStore.Guild(parsedGuildId)
  if err != nil {
    return channels
  }

  guildChannels, err := client.ChannelStore.Channels(guild.ID)
  if err != nil {
    return channels
  }

  self, err := client.MeStore.Me()
  if err != nil {
    return channels
  }

  selfMember, err := client.MemberStore.Member(guild.ID, self.ID)
  if err != nil {
    return channels
  }

  if withCategories {
    categories := make(map[string][]discord.Channel)

    for _, channel := range guildChannels {
      if channel.Type != discord.GuildText && channel.Type != discord.GuildAnnouncement {
        continue
      }

      perms := discord.CalcOverwrites(*guild, channel, *selfMember)

      private := !perms.Has(discord.PermissionViewChannel)
      if private && !withPrivate {
        continue
      }

      categoryID := "0"
      if channel.ParentID.IsValid() {
        categoryID = channel.ParentID.String()
      }

      _, has := categories[categoryID]
      if !has {
        categories[categoryID] = make([]discord.Channel, 0)
      }

      categories[categoryID] = append(categories[categoryID], channel)
    }

    for id, category := range categories {
      // sort channels by position
      sort.Slice(category, func(i, j int) bool {
        return category[i].Position < category[j].Position
      })
      categoryChannels := make([]discord.Channel, 0)

      // append category channel to top
      if id != "0" {
        parsedCategoryId, err := discord.ParseSnowflake(id)
        if err != nil {
          continue
        }

        for _, channel := range guildChannels {
          if channel.ID == discord.ChannelID(parsedCategoryId) {
            categoryChannels = append(categoryChannels, channel)
            break
          }
        }
      }

      // append channels
      for _, channel := range category {
        categoryChannels = append(categoryChannels, channel)
      }
      categories[id] = categoryChannels
    }

    keys := make([]string, 0)
    for id := range categories {
      if id == "0" {
        continue
      }
      keys = append(keys, id)
    }
    sort.Slice(keys, func(i, j int) bool {
      pa, _ := discord.ParseSnowflake(keys[i])
      pb, _ := discord.ParseSnowflake(keys[i])

      ca, _ := client.ChannelStore.Channel(discord.ChannelID(pa))
      cb, _ := client.ChannelStore.Channel(discord.ChannelID(pb))

      return ca.Position < cb.Position
    })
    sortedCategories := make(map[string][]discord.Channel)
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
    for _, channel := range guildChannels {
      if channel.Type != discord.GuildText && channel.Type != discord.GuildAnnouncement {
        continue
      }

      perms := discord.CalcOverwrites(*guild, channel, *selfMember)

      private := !perms.Has(discord.PermissionViewChannel)
      if private && !withPrivate {
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

func ListChannelsCommand() {
  client := state.GetClient()
  self, err := client.MeStore.Me()
  if err != nil {
    fmt.Print("<failed to get self: ", err.Error(), ">\n\r")
    return
  }

  currentGuild := state.GetCurrentGuild()
  if currentGuild == "" {
    fmt.Print("<not in a guild>\n\r")
    return
  }

  guildSnowflake, err := discord.ParseSnowflake(currentGuild)
  if err != nil {
    fmt.Print("<failed to parse current guild id: ", err.Error(), ">\n\r")
    return
  }

  parsedGuildId := discord.GuildID(guildSnowflake)
  guild, err := client.GuildStore.Guild(parsedGuildId)
  if err != nil {
    fmt.Print("<failed to get current guild: ", err.Error(), ">\n\r")
    return
  }

  selfMember, err := client.MemberStore.Member(parsedGuildId, self.ID)
  if err != nil {
    fmt.Print("<failed to get self member: ", err.Error(), ">\n\r")
    return
  }

  longest := 0
  channels := GetSortedChannels(currentGuild, true, false)

  for _, channel := range channels {
    perms := discord.CalcOverwrites(*guild, channel, *selfMember)

    private := !perms.Has(discord.PermissionViewChannel)
    category := channel.Type == discord.GuildCategory

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
    perms := discord.CalcOverwrites(*guild, channel, *selfMember)

    private := !perms.Has(discord.PermissionViewChannel)
    category := channel.Type == discord.GuildCategory
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
    timestamp := channel.CreatedAt()
    created = timestamp.UTC().Format("02-Jan-06")

    fmt.Printf("  %*s  %s  %s\n\r", longest, name, created, topic)
  }
  fmt.Print(strings.Repeat("-", 80) + "\n\r")
  fmt.Print("\n\r")
}

type ListedMember struct {
  Name string
  Bot bool
  Status discord.Status
  Position int
}

func ListUsersCommand() {
  client := state.GetClient()

  currentGuild := state.GetCurrentGuild()
  currentChannel := state.GetCurrentChannel()

  if currentGuild == "" {
    fmt.Print("<not in a guild>\n\r")
    return
  }
  if currentChannel == "" {
    fmt.Print("<not in a channel>\n\r")
    return
  }

  parsedGuildId, err := discord.ParseSnowflake(currentGuild)
  if err != nil {
    fmt.Print("<failed to parse guild id: ", err.Error(), ">\n\r")
    return
  }
  parsedChannelId, err := discord.ParseSnowflake(currentChannel)
  if err != nil {
    fmt.Print("<failed to parse channel id: ", err.Error(), ">\n\r")
    return
  }

  guild, err := client.GuildStore.Guild(discord.GuildID(parsedGuildId))
  if err != nil {
    fmt.Print("<failed to get guild: ", err.Error(), ">\n\r")
    return
  }
  channel, err := client.ChannelStore.Channel(discord.ChannelID(parsedChannelId))
  if err != nil {
    fmt.Print("<failed to get channel: ", err.Error(), ">\n\r")
    return
  }

  longest := 0

  sortedMembers := make([]ListedMember, 0)

  presences, err := client.Presences(guild.ID)
  if err != nil {
    fmt.Print("<failed to get presences: ", err.Error(), ">\n\r")
    return
  }

  for _, presence := range presences {
    if presence.Status == discord.OfflineStatus {
      continue
    }
    member, err := client.MemberStore.Member(guild.ID, presence.User.ID)
    if err != nil {
      continue
    }

    perms := discord.CalcOverwrites(*guild, *channel, *member)
    if !perms.Has(discord.PermissionViewChannel) {
      continue
    }

    length := utf8.RuneCountInString(member.User.Username) + 3
    if length > longest {
      longest = length
    }

    position := 0
    for _, id := range member.RoleIDs {
      role, err := client.RoleStore.Role(guild.ID, id)
      if err != nil {
        continue
      }

      if role.Hoist && role.Position > position {
        position = role.Position
      }
    }

    sortedMembers = append(sortedMembers, ListedMember{
      Name: member.User.Username,
      Bot: member.User.Bot,
      Status: presence.Status,
      Position: position,
    })
  }

  fmt.Print("\n\r")
  fmt.Printf("[you are in '%s' in '#%s' among %d]\n\r", guild.Name, channel.Name, len(sortedMembers))
  fmt.Print("\n\r")

  membersByPosition := make(map[int][]ListedMember)
  for _, member := range sortedMembers {
    _, has := membersByPosition[member.Position]
    if !has {
      membersByPosition[member.Position] = make([]ListedMember, 0)
    }

    membersByPosition[member.Position] = append(membersByPosition[member.Position], member)
  }
  for _, members := range membersByPosition {
    sort.Slice(members, func(i, j int) bool {
      return members[i].Name < members[j].Name
    })
  }

  positions := make([]int, 0, len(membersByPosition))
  for k := range membersByPosition {
    positions = append(positions, k)
  }
  sort.Slice(positions, func(i, j int) bool {
    return positions[i] > positions[j]
  })

  size, err := tsize.GetSize()
  if err != nil {
    return
  }
  columns := int(math.Floor(float64(size.Width) / float64(longest)))

  index := 0
  for _, position := range positions {
    members := membersByPosition[position]
    if len(members) > 150 {
      continue
    }
    for _, member := range members {

      statusColor := "reset"
      if member.Status == discord.OnlineStatus {
        statusColor = "green+b"
      } else if member.Status == discord.IdleStatus {
        statusColor = "yellow+b"
      } else if member.Status == discord.DoNotDisturbStatus {
        statusColor = "red+b"
      }

      nameColor := "reset"
      if member.Bot {
        nameColor = "yellow"
      }

      nameAndStatus := ansi.Color(" \u2022 ", statusColor) + ansi.Color(member.Name, nameColor)
      nameLength := utf8.RuneCountInString(member.Name) + 3

      index++

      pad := 0
      if index % columns != 0 {
        pad = longest - nameLength
      }
      if pad < 0 {
        pad = 0
      }

      fmt.Printf(nameAndStatus + strings.Repeat(" ", pad))

      if index % columns == 0 {
        fmt.Print("\n\r")
      }
    }
  }
  if index % columns != 0 {
    fmt.Print("\n\r")
  }
  fmt.Print("\n\r")

  if channel.Topic != "" {
    fmt.Print("--Topic" + strings.Repeat("-", 73) + "\n\r")
    for _, line := range strings.Split(channel.Topic, "\n") {
      fmt.Print(line + "\n\r")
    }
    fmt.Print(strings.Repeat("-", 80) + "\n\r")
    fmt.Print("\n\r")
  }
}

func SwitchGuild(input string) {
  client := state.GetClient()

  if input == "" {
    ListChannelsCommand()
    ListUsersCommand()
  } else {
    target := ""

    guilds, err := client.GuildStore.Guilds()
    if err != nil {
      fmt.Print("<failed to get guilds: ", err.Error(), ">\n\r")
      return
    }

    for _, guild := range guilds {
      if strings.Index(strings.ToLower(guild.Name), strings.ToLower(input)) > -1 {
        target = guild.ID.String()
        break;
      }
    }

    if target == "" {
      fmt.Print("<guild not found>\n\r")
    } else {
      state.SetCurrentGuild(target)
      last := state.GetLastChannel(target)
      if last == "" {
        channels := GetSortedChannels(target, false, false)
        if len(channels) > 0 {
          topChannel := channels[0]

          state.SetCurrentChannel(topChannel.ID.String())
          state.SetLastChannel(target, topChannel.ID.String())
        }
      } else {
        state.SetCurrentChannel(last)
      }

      ListChannelsCommand()
      ListUsersCommand()

      lib.UpdatePresence()
    }
  }
}

func SwitchGuildsCommand() {
  lib.MakePrompt(":guild> ", false, func(input string, interrupt bool) {
    fmt.Print("\r")
    SwitchGuild(input)
  })
}

func SwitchChannelsCommand() {
  currentGuild := state.GetCurrentGuild()

  if currentGuild == "" {
    fmt.Print("<not in a guild>\n\r")
    return
  }

  lib.MakePrompt(":channel> ", false, func(input string, interrupt bool) {
    fmt.Print("\r")
    if input == "" {
      ListUsersCommand()
    } else {
      target := ""

      channels := GetSortedChannels(currentGuild, false, false)

      for _, channel := range channels {
        if strings.Index(strings.ToLower(channel.Name), strings.ToLower(input)) > -1 {
          target = channel.ID.String()
          break
        }
      }

      if target == "" {
        fmt.Print("<channel not found>\n\r")
      } else {
        state.SetCurrentChannel(target)
        state.SetLastChannel(currentGuild, target)

        ListUsersCommand()

        lib.UpdatePresence()
      }
    }
  })
}
