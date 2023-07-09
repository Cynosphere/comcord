package lib

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
	"github.com/mgutz/ansi"
)

var /*const*/ REGEX_CODEBLOCK = regexp.MustCompile(`(?i)\x60\x60\x60(?:([a-z0-9_+\-\.]+?)\n)?\n*([^\n].*?)\n*\x60\x60\x60`)

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

  // TODO: timestamps for pin and join

  // TODO: history lines

  nameLength := len(options.Name) + 2
  stateNameLength := state.GetNameLength()
  if nameLength > stateNameLength {
    state.SetNameLength(nameLength)
  }

  // TODO: replies

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

      if options.NoColor {
        fmt.Print(str)
      }
    }
  } else {
    // TODO: markdown

    if options.IsDM {
      name := fmt.Sprintf("*%s*", options.Name)
      if !options.NoColor {
        name = ansi.Color(name, "red+b")
      }

      fmt.Printf("%s %s\x07\n\r", name, options.Content)
    } else if len(options.Content) > 1 &&
    (strings.HasPrefix(options.Content, "*") && strings.HasSuffix(options.Content, "*") && !strings.HasPrefix(options.Content, "**") && !strings.HasSuffix(options.Content, "**")) ||
    (strings.HasPrefix(options.Content, "_") && strings.HasSuffix(options.Content, "_") && !strings.HasPrefix(options.Content, "__") && !strings.HasSuffix(options.Content, "__")) {
      str := fmt.Sprintf("<%s %s>", options.Name, options.Content[1:len(options.Content)-1])

      if options.NoColor {
        fmt.Print(str + "\n\r")
      } else {
        fmt.Print(ansi.Color(str, "green+b") + "\n\r")
      }
    } else if options.IsJoin {
      // TODO
    } else if options.IsPin {
      // TODO
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

      // FIXME: where is this off by 4 actually from
      padding := strings.Repeat(" ", int(math.Abs(float64(stateNameLength) - float64(nameLength) - 4)))
      str := fmt.Sprintf("%s%s %s", name, padding, options.Content)
      if options.IsMention {
        str = str + "\x07"
      }
      fmt.Print(str + "\n\r")
    }
  }

  // TODO: attachments

  // TODO: stickers

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
  options.IsDump = REGEX_CODEBLOCK.MatchString(content)

  FormatMessage(session, options)
}

func ProcessQueue(session *discordgo.Session) {
  queue := state.GetMessageQueue()

  for _, msg := range queue {
    ProcessMessage(session, msg, MessageOptions{NoColor: state.HasNoColor()})
  }

  state.EmptyMessageQueue()
}
