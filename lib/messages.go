package lib

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
	"github.com/mgutz/ansi"
)

var REGEX_CODEBLOCK = regexp.MustCompile(`(?i)\x60\x60\x60(?:([a-z0-9_+\-\.]+?)\n)?\n*([^\n].*?)\n*\x60\x60\x60`)
var REGEX_EMOTE = regexp.MustCompile(`<(?:\x{200b}|&)?a?:(\w+):(\d+)>`)

type MessageOptions struct {
  Content string
  Name string
  Channel string
  Bot bool
  Attachments []*discordgo.MessageAttachment
  Stickers []*discordgo.Sticker
  Reply *discordgo.Message
  Timestamp time.Time
  IsMention bool
  IsDM bool
  IsJoin bool
  IsPin bool
  IsDump bool
  NoColor bool
  InHistory bool
}

func FormatMessage(session *discordgo.Session, options MessageOptions) {
  timestamp := options.Timestamp.Format("[15:04:05]")

  // TODO: history lines

  nameLength := utf8.RuneCountInString(options.Name) + 2
  stateNameLength := state.GetNameLength()
  if nameLength > stateNameLength {
    state.SetNameLength(nameLength)
    stateNameLength = nameLength
  }

  if options.Reply != nil {
    nameColor := "cyan+b"
    if options.Bot {
      nameColor = "yellow+b"
    }

    headerLength := 6 + utf8.RuneCountInString(options.Reply.Author.Username)

    content, _ := options.Reply.ContentWithMoreMentionsReplaced(session)
    replyContent := strings.ReplaceAll(content, "\n", " ")

    // TODO: markdown
    replyContent = REGEX_EMOTE.ReplaceAllString(replyContent, ":$1:")

    attachmentCount := len(options.Reply.Attachments)
    if attachmentCount > 0 {
      attachmentPlural := ""
      if attachmentCount > 1 {
        attachmentPlural = "s"
      }

      replyContent = strings.TrimSpace(replyContent + fmt.Sprintf(" <%d attachment%s>", attachmentCount, attachmentPlural))
    }

    stickerCount := len(options.Reply.StickerItems)
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

    fmt.Print(replySymbol, name, replyContent, "\n\r")
  }

  if options.IsDump {
    if options.InHistory {
      // TODO
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

      fmt.Print(str + "\n\r")
    }
  } else {
    // TODO: markdown
    content := options.Content
    content = REGEX_EMOTE.ReplaceAllString(content, ":$1:")

    if options.IsDM {
      name := fmt.Sprintf("*%s*", options.Name)
      if !options.NoColor {
        name = ansi.Color(name, "red+b")
      }

      fmt.Printf("%s %s\x07\n\r", name, content)
    } else if utf8.RuneCountInString(content) > 1 &&
    (strings.HasPrefix(content, "*") && strings.HasSuffix(content, "*") && !strings.HasPrefix(content, "**") && !strings.HasSuffix(content, "**")) ||
    (strings.HasPrefix(content, "_") && strings.HasSuffix(content, "_") && !strings.HasPrefix(content, "__") && !strings.HasSuffix(content, "__")) {
      str := fmt.Sprintf("<%s %s>", options.Name, content[1:len(content)-1])

      if !options.NoColor {
        str = ansi.Color(str, "green+b")
      }

      fmt.Print(str + "\n\r")
    } else if options.IsJoin {
      channel, err := session.State.Channel(options.Channel)
      if err != nil {
        return
      }
      guild, err := session.State.Guild(channel.GuildID)
      if err != nil {
        return
      }

      str := fmt.Sprintf("%s %s has joined %s", timestamp, options.Name, guild.Name)
      if !options.NoColor {
        str = ansi.Color(str, "yellow+b")
      }

      fmt.Print(str + "\n\r")
    } else if options.IsPin {
      str := fmt.Sprintf("%s %s pinned a message to this channel", timestamp, options.Name)
      if !options.NoColor {
        str = ansi.Color(str, "yellow+b")
      }

      fmt.Print(str + "\n\r")
    } else {
      nameColor := "cyan+b"
      if options.IsMention {
        nameColor = "red+b"
      } else if options.Bot {
        nameColor = "yellow+b"
      }

      name := fmt.Sprintf("[%s]", options.Name)
      if !options.NoColor {
        name = ansi.Color(name, nameColor)
      }

      padding := strings.Repeat(" ", int(math.Abs(float64(stateNameLength) - float64(nameLength))) + 1)
      str := name + padding + content
      if options.IsMention {
        str = str + "\x07"
      }
      fmt.Print(str + "\n\r")
    }
  }

  if len(options.Attachments) > 0 {
    for _, attachment := range options.Attachments {
      str := fmt.Sprintf("<attachment: %s >", attachment.URL)
      if !options.NoColor {
        str = ansi.Color(str, "yellow+b")
      }

      fmt.Print(str + "\n\r")
    }
  }

  if len(options.Stickers) > 0 {
    for _, sticker := range options.Stickers {
      str := fmt.Sprintf("<sticker: \"%s\" https://cdn.discordapp.com/stickers/%s.png >", sticker.Name, sticker.ID)
      if !options.NoColor {
        str = ansi.Color(str, "yellow+b")
      }

      fmt.Print(str + "\n\r")
    }
  }

  // TODO: links

  // TODO: embeds

  // TODO: lines output for history
}

func ProcessMessage(session *discordgo.Session, msg *discordgo.Message, options MessageOptions) {
  channel, err := session.State.Channel(msg.ChannelID)
  if err != nil {
    return
  }

  guild, err := session.State.Guild(channel.GuildID)
  if err != nil {
    return
  }

  selfMember, err := session.State.Member(guild.ID, session.State.User.ID)
  if err != nil {
    return
  }

  hasMentionedRole := false
  for _, role := range msg.MentionRoles {
    for _, selfRole := range selfMember.Roles {
      if role == selfRole {
        hasMentionedRole = true
        break;
      }
    }
  }

  isDirectlyMentioned := false
  for _, user := range msg.Mentions {
    if user.ID == session.State.User.ID {
      isDirectlyMentioned = true
      break;
    }
  }

  isPing := msg.MentionEveryone || hasMentionedRole || isDirectlyMentioned
  isDM := channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM
  isEdit := msg.EditedTimestamp != nil

  currentChannel := state.GetCurrentChannel()
  isCurrentChannel := currentChannel == msg.ChannelID

  if !isCurrentChannel && !isDM && !isPing {
    return
  }

  if isPing && !isCurrentChannel && !isDM {
    str := fmt.Sprintf("**mentioned by %s in #%s in %s**", msg.Author.Username, channel.Name, guild.Name)
    if options.NoColor {
      fmt.Print(ansi.Color(str, "red+b"))
    } else {
      fmt.Print(str)
    }
    fmt.Print("\x07\n\r")
    return
  }

  content, _ := msg.ContentWithMoreMentionsReplaced(session)
  if isEdit {
    content = content + " (edited)"
  }

  isDump := REGEX_CODEBLOCK.MatchString(content)

  if strings.Index(content, "\n") > -1 && !isDump {
    for _, line := range strings.Split(content, "\n") {
      options.Content = line
      options.Name = msg.Author.Username
      options.Channel = msg.ChannelID
      options.Bot = msg.Author.Bot
      options.Attachments = msg.Attachments
      options.Stickers = msg.StickerItems
      options.Reply = msg.ReferencedMessage
      options.IsMention = isPing
      options.IsDM = isDM
      options.IsJoin = msg.Type == discordgo.MessageTypeGuildMemberJoin
      options.IsPin = msg.Type == discordgo.MessageTypeChannelPinnedMessage
      options.IsDump = false

      FormatMessage(session, options)
    }
  } else {
    options.Content = content
    options.Name = msg.Author.Username
    options.Channel = msg.ChannelID
    options.Bot = msg.Author.Bot
    options.Attachments = msg.Attachments
    options.Stickers = msg.StickerItems
    options.Reply = msg.ReferencedMessage
    options.IsMention = isPing
    options.IsDM = isDM
    options.IsJoin = msg.Type == discordgo.MessageTypeGuildMemberJoin
    options.IsPin = msg.Type == discordgo.MessageTypeChannelPinnedMessage
    options.IsDump = isDump

    FormatMessage(session, options)
  }
}

func ProcessQueue(session *discordgo.Session) {
  queue := state.GetMessageQueue()

  for _, msg := range queue {
    ProcessMessage(session, msg, MessageOptions{NoColor: state.HasNoColor()})
  }

  state.EmptyMessageQueue()
}
