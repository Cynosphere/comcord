package lib

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Cynosphere/comcord/state"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mergestat/timediff"
	"github.com/mgutz/ansi"
)

var REGEX_CODEBLOCK = regexp.MustCompile(`(?i)\x60\x60\x60(?:([a-z0-9_+\-\.]+?)\n)?\n*([^\n](?:.|\n)*?)\n*\x60\x60\x60`)
var REGEX_MENTION = regexp.MustCompile(`<@!?(\d+)>`)
var REGEX_ROLE_MENTION = regexp.MustCompile(`<@&(\d+)>`)
var REGEX_CHANNEL = regexp.MustCompile(`<#(\d+)>`)
var REGEX_EMOTE = regexp.MustCompile(`<(?:\x{200b}|&)?a?:(\w+):(\d+)>`)
var REGEX_COMMAND = regexp.MustCompile(`</([^\s]+?):(\d+)>`)
var REGEX_BLOCKQUOTE = regexp.MustCompile(`^ *>>?>? +`)
var REGEX_GREENTEXT = regexp.MustCompile(`^(>.+?)(?:\n|$)`)
var REGEX_SPOILER = regexp.MustCompile(`\|\|(.+?)\|\|`)
var REGEX_BOLD = regexp.MustCompile(`\*\*(.+?)\*\*`)
var REGEX_UNDERLINE = regexp.MustCompile(`__(.+?)__`)
var REGEX_ITALIC_1 = regexp.MustCompile(`\*(.+?)\*`)
var REGEX_ITALIC_2 = regexp.MustCompile(`_(.+?)_`)
var REGEX_STRIKE = regexp.MustCompile(`~~(.+?)~~`)
var REGEX_3Y3 = regexp.MustCompile(`[\x{e0020}-\x{e007e}]{1,}`)
var REGEX_TIMESTAMP = regexp.MustCompile(`<t:(-?\d{1,17})(?::(t|T|d|D|f|F|R))?>`)

type MessageOptions struct {
  Content string
  Name string
  Channel discord.ChannelID
  Bot bool
  Webhook bool
  Attachments []discord.Attachment
  Stickers []discord.StickerItem
  Reply *discord.Message
  Timestamp time.Time
  IsMention bool
  IsDM bool
  IsJoin bool
  IsPin bool
  IsDump bool
  NoColor bool
  InHistory bool
}

func Parse3y3(content string) string {
  out := []rune{}
  for i, w := 0, 0; i < len(content); i += w {
    runeValue, width := utf8.DecodeRuneInString(content[i:])
    w = width

    out = append(out, rune(int(runeValue) - 0xe0000))
  }

  return string(out)
}

func ReplaceStyledMarkdown(content string) string {
  content = REGEX_BLOCKQUOTE.ReplaceAllString(content, ansi.Color("\u258e", "black+h"))
  content = REGEX_GREENTEXT.ReplaceAllStringFunc(content, func(match string) string {
    return ansi.Color(match, "green")
  })

  if state.GetConfigValue("enable3y3") == "true" {
    parsed := REGEX_3Y3.FindString(content)
    parsed = Parse3y3(parsed)
    parsed = "\033[3m" + parsed + "\033[23m"
    content = REGEX_3Y3.ReplaceAllString(content, ansi.Color(parsed, "magenta"))
  }

  content = REGEX_SPOILER.ReplaceAllString(content, "\033[30m\033[40m$1\033[39m\033[49m")
  content = REGEX_STRIKE.ReplaceAllString(content, "\033[9m$1\033[29m")
  content = REGEX_BOLD.ReplaceAllString(content, "\033[1m$1\033[22m")
  content = REGEX_UNDERLINE.ReplaceAllString(content, "\033[4m$1\033[24m")
  content = REGEX_ITALIC_1.ReplaceAllString(content, "\033[3m$1\033[23m")
  content = REGEX_ITALIC_2.ReplaceAllString(content, "\033[3m$1\033[23m")

  return content
}

func replaceAllWithCallback(re regexp.Regexp, content string, callback func(matches []string) string) string {
  return re.ReplaceAllStringFunc(content, func(match string) string {
    matches := re.FindStringSubmatch(match)
    return callback(matches)
  })
}

func ReplaceMarkdown(content string, noColor bool) string {
  client := state.GetClient()

  content = replaceAllWithCallback(*REGEX_MENTION, content, func(matches []string) string {
    id := matches[1]
    parsedId, err := discord.ParseSnowflake(id)
    if err != nil {
      return "@Unknown User"
    }

    currentGuild := state.GetCurrentGuild()
    if currentGuild == "" {
      user, err := client.User(discord.UserID(parsedId))
      if err != nil {
        return "@Unknown User"
      }

      return "@" + user.Username
    } else {
      parsedGuildId, err := discord.ParseSnowflake(currentGuild)
      if err != nil {
        return "@Unknown User"
      }

      member, err := client.MemberStore.Member(discord.GuildID(parsedGuildId), discord.UserID(parsedId))
      if err != nil {
        return "@Unknown User"
      }

      return "@" + member.User.Username
    }
  })

  content = replaceAllWithCallback(*REGEX_ROLE_MENTION, content, func(matches []string) string {
    id := matches[1]
    parsedId, err := discord.ParseSnowflake(id)
    if err != nil {
      return "[@Unknown Role]"
    }

    currentGuild := state.GetCurrentGuild()
    if currentGuild == "" {
      return "[@Unknown Role]"
    }
    parsedGuildId, err := discord.ParseSnowflake(currentGuild)
    if err != nil {
      return "[@Unknown Role]"
    }

    role, err := client.RoleStore.Role(discord.GuildID(parsedGuildId), discord.RoleID(parsedId))
    if err != nil {
      return "[@Unknown Role]"
    }

    return fmt.Sprintf("[@%s]", role.Name)
  })

  content = replaceAllWithCallback(*REGEX_CHANNEL, content, func(matches []string) string {
    id := matches[1]
    parsedId, err := discord.ParseSnowflake(id)
    if err != nil {
      return "#Unknown"
    }

    channel, err := client.ChannelStore.Channel(discord.ChannelID(parsedId))
    if err != nil {
      return "#Unknown"
    }

    return "#" + channel.Name
  })

  content = REGEX_EMOTE.ReplaceAllString(content, ":$1:")
  content = REGEX_COMMAND.ReplaceAllString(content, "/$1")

  content = replaceAllWithCallback(*REGEX_TIMESTAMP, content, func (matches []string) string {
    timestamp, err := strconv.Atoi(matches[1])
    if err != nil {
      return "Invalid Date"
    }
    timeObj := time.Unix(int64(timestamp), 0).UTC()

    format := matches[2]

    switch format {
      case "t":
        return timeObj.Format("15:04")
      case "T":
        return timeObj.Format("15:04:05")
      case "d":
        return timeObj.Format("2006/01/02")
      case "D":
        return timeObj.Format("2 January 2006")
      case "f":
      default:
        return timeObj.Format("2 January 2006 15:04")
      case "F":
        return timeObj.Format("Monday, 2 January 2006 15:04")
      case "R":
        return timediff.TimeDiff(timeObj)
    }

    return "Invalid Date"
  })

  if !noColor {
    content = ReplaceStyledMarkdown(content)
  } else {
    if state.GetConfigValue("enable3y3") == "true" {
      parsed := REGEX_3Y3.FindString(content)
      parsed = Parse3y3(parsed)
      content = REGEX_3Y3.ReplaceAllString(content, "<3y3:" + parsed + ">")
    }
  }

  return content
}

func FormatMessage(options MessageOptions) []string {
  client := state.GetClient()

  lines := make([]string, 0)

  timestamp := options.Timestamp.UTC().Format("[15:04:05]")

  nameLength := utf8.RuneCountInString(options.Name) + 2
  stateNameLength := state.GetNameLength()
  if nameLength > stateNameLength {
    state.SetNameLength(nameLength)
    stateNameLength = nameLength
  }

  if options.Reply != nil {
    nameColor := "cyan+b"
    if options.Reply.Author.Bot {
      nameColor = "yellow+b"
    }

    headerLength := 6 + utf8.RuneCountInString(options.Reply.Author.Username)

    content := options.Reply.Content
    replyContent := strings.ReplaceAll(content, "\n", " ")

    replyContent = ReplaceMarkdown(replyContent, options.NoColor)

    attachmentCount := len(options.Reply.Attachments)
    if attachmentCount > 0 {
      attachmentPlural := ""
      if attachmentCount > 1 {
        attachmentPlural = "s"
      }

      replyContent = strings.TrimSpace(replyContent + fmt.Sprintf(" <%d attachment%s>", attachmentCount, attachmentPlural))
    }

    stickerCount := len(options.Reply.Stickers)
    if stickerCount > 0 {
      stickerPlural := ""
      if stickerCount > 0 {
        stickerPlural = "s"
      }

      replyContent = strings.TrimSpace(replyContent + fmt.Sprintf(" <%d sticker%s>", stickerCount, stickerPlural))
    }

    length := headerLength + utf8.RuneCountInString(replyContent)

    replySymbol := " \u00bb "
    if !options.NoColor {
      replySymbol = ansi.Color(replySymbol, "white+b")
    }

    name := fmt.Sprintf("[%s] ", options.Reply.Author.Username)
    if !options.NoColor {
      name = ansi.Color(name, nameColor)
    }

    moreContent := "\u2026"
    if !options.NoColor {
      moreContent = ansi.Color(moreContent, "reset")
    }

    if length > 79 {
      replyContent = replyContent[:79 - headerLength] + moreContent
    }

    lines = append(lines, replySymbol, name, replyContent, "\n\r")
  }

  if options.IsDump {
    if options.InHistory {
      headerLength := 80 - (utf8.RuneCountInString(options.Name) + 5)
      dumpHeader := fmt.Sprintf("--- %s %s\n\r", options.Name, strings.Repeat("-", headerLength))

      contentLines := strings.Split(options.Content, "\n")

      lines = append(lines, dumpHeader)
      for _, line := range contentLines {
        lines = append(lines, line + "\n\r")
      }
      lines = append(lines, dumpHeader)
    } else {
      wordCount := len(strings.Split(options.Content, " "))
      lineCount := len(strings.Split(options.Content, "\n"))
      wordsPlural := ""
      linesPlural := ""

      if wordCount > 1 {
        wordsPlural = "s"
      }
      if lineCount > 1 {
        linesPlural = "s"
      }

      str := fmt.Sprintf("<%s DUMPs in %d characters of %d word%s in %d line%s>", options.Name, len(options.Content), wordCount, wordsPlural, lineCount, linesPlural)

      if !options.NoColor {
        str = ansi.Color(str, "yellow+b")
      }

      lines = append(lines, str + "\n\r")
    }
  } else {
    content := options.Content

    if options.IsDM {
      name := fmt.Sprintf("*%s*", options.Name)
      if !options.NoColor {
        name = ansi.Color(name, "red+b")
      }

      content = ReplaceMarkdown(content, options.NoColor)

      lines = append(lines, fmt.Sprintf("%s %s\x07\n\r", name, content))
    } else if utf8.RuneCountInString(content) > 1 &&
    (strings.HasPrefix(content, "*") && strings.HasSuffix(content, "*") && !strings.HasPrefix(content, "**") && !strings.HasSuffix(content, "**")) ||
    (strings.HasPrefix(content, "_") && strings.HasSuffix(content, "_") && !strings.HasPrefix(content, "__") && !strings.HasSuffix(content, "__")) {
      str := fmt.Sprintf("<%s %s>", options.Name, content[1:len(content)-1])

      if !options.NoColor {
        str = ansi.Color(str, "green+b")
      }

      lines = append(lines, str + "\n\r")
    } else if options.IsJoin {
      channel, err := client.ChannelStore.Channel(options.Channel)
      if err != nil {
        return lines
      }
      guild, err := client.GuildStore.Guild(channel.GuildID)
      if err != nil {
        return lines
      }

      str := fmt.Sprintf("%s %s has joined %s", timestamp, options.Name, guild.Name)
      if !options.NoColor {
        str = ansi.Color(str, "yellow+b")
      }

      lines = append(lines, str + "\n\r")
    } else if options.IsPin {
      str := fmt.Sprintf("%s %s pinned a message to this channel", timestamp, options.Name)
      if !options.NoColor {
        str = ansi.Color(str, "yellow+b")
      }

      lines = append(lines, str + "\n\r")
    } else {
      nameColor := "cyan+b"
      if options.IsMention {
        nameColor = "red+b"
      } else if options.Webhook {
        nameColor = "magenta+b"
      } else if options.Bot {
        nameColor = "yellow+b"
      }

      content = ReplaceMarkdown(content, options.NoColor)

      name := fmt.Sprintf("[%s]", options.Name)
      if !options.NoColor {
        name = ansi.Color(name, nameColor)
      }

      padding := strings.Repeat(" ", int(math.Abs(float64(stateNameLength) - float64(nameLength))) + 1)
      str := name + padding + content
      if options.IsMention {
        str = str + "\x07"
      }
      lines = append(lines, str + "\n\r")
    }
  }

  if len(options.Attachments) > 0 {
    for _, attachment := range options.Attachments {
      str := fmt.Sprintf("<attachment: %s >", attachment.URL)
      if !options.NoColor {
        str = ansi.Color(str, "yellow+b")
      }

      lines = append(lines, str + "\n\r")
    }
  }

  if len(options.Stickers) > 0 {
    for _, sticker := range options.Stickers {
      str := fmt.Sprintf("<sticker: \"%s\" https://cdn.discordapp.com/stickers/%s.png >", sticker.Name, sticker.ID)
      if !options.NoColor {
        str = ansi.Color(str, "yellow+b")
      }

      lines = append(lines, str + "\n\r")
    }
  }

  // TODO: links

  // TODO: embeds

  // TODO: lines output for history
  return lines
}

func ProcessMessage(msg discord.Message, options MessageOptions) []string {
  client := state.GetClient()
  lines := make([]string, 0)

  channel, err := client.ChannelStore.Channel(msg.ChannelID)
  if err != nil {
    return lines
  }

  guild, err := client.GuildStore.Guild(channel.GuildID)
  if err != nil {
    return lines
  }

  self, err := client.MeStore.Me()
  if err != nil {
    return lines
  }

  selfMember, err := client.MemberStore.Member(guild.ID, self.ID)
  if err != nil {
    return lines
  }

  hasMentionedRole := false
  for _, role := range msg.MentionRoleIDs {
    for _, selfRole := range selfMember.RoleIDs {
      if role == selfRole {
        hasMentionedRole = true
        break;
      }
    }
  }

  isDirectlyMentioned := false
  for _, user := range msg.Mentions {
    if user.ID == self.ID {
      isDirectlyMentioned = true
      break;
    }
  }

  isPing := msg.MentionEveryone || hasMentionedRole || isDirectlyMentioned
  isDM := channel.Type == discord.DirectMessage || channel.Type == discord.GroupDM
  isEdit := msg.EditedTimestamp.IsValid()

  currentChannel := state.GetCurrentChannel()
  isCurrentChannel := currentChannel == msg.ChannelID.String()

  if !isCurrentChannel && !isDM && !isPing && !options.InHistory {
    return lines
  }

  if isPing && !isCurrentChannel && !isDM && !options.InHistory {
    str := fmt.Sprintf("**mentioned by %s in #%s in %s**", msg.Author.Username, channel.Name, guild.Name)
    if !options.NoColor {
      str = ansi.Color(str, "red+b")
    }
    str = str + "\x07\n\r"
    lines = append(lines, str)
  } else {
    content := msg.Content
    if isEdit {
      content = content + " (edited)"
    }

    isDump := REGEX_CODEBLOCK.MatchString(content)

    if strings.Index(content, "\n") > -1 && !isDump {
      for i, line := range strings.Split(content, "\n") {
        options.Content = line
        options.Name = msg.Author.Username
        options.Channel = msg.ChannelID
        options.Bot = msg.Author.Bot
        options.Webhook = msg.WebhookID.IsValid()
        options.Attachments = msg.Attachments
        options.Stickers = msg.Stickers
        if i == 0 {
          options.Reply = msg.ReferencedMessage
        } else {
          options.Reply = nil
        }
        options.Timestamp = time.Time(msg.Timestamp)
        options.IsMention = isPing
        options.IsDM = isDM
        options.IsJoin = msg.Type == discord.GuildMemberJoinMessage
        options.IsPin = msg.Type == discord.ChannelPinnedMessage
        options.IsDump = false

        msgLines := FormatMessage(options)
        for _, line := range msgLines {
          lines = append(lines, line)
        }
      }
    } else {
      options.Content = content
      options.Name = msg.Author.Username
      options.Channel = msg.ChannelID
      options.Bot = msg.Author.Bot
      options.Webhook = msg.WebhookID.IsValid()
      options.Attachments = msg.Attachments
      options.Stickers = msg.Stickers
      options.Reply = msg.ReferencedMessage
      options.Timestamp = time.Time(msg.Timestamp)
      options.IsMention = isPing
      options.IsDM = isDM
      options.IsJoin = msg.Type == discord.GuildMemberJoinMessage
      options.IsPin = msg.Type == discord.ChannelPinnedMessage
      options.IsDump = isDump

      lines = FormatMessage(options)
    }
  }

  return lines
}

func ProcessQueue() {
  queue := state.GetMessageQueue()

  for _, msg := range queue {
    lines := ProcessMessage(msg, MessageOptions{NoColor: state.HasNoColor()})
    for _, line := range lines {
      fmt.Print(line)
    }
  }

  state.EmptyMessageQueue()
}
