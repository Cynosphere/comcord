package events

import "github.com/bwmarrin/discordgo"

func Setup(session *discordgo.Session) {
  session.AddHandlerOnce(Ready)
  session.AddHandler(MessageCreate)
  session.AddHandler(MessageUpdate)
}
