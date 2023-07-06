package events

import "github.com/bwmarrin/discordgo"

func MessageCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
  if (msg.Author.ID == session.State.User.ID) {
    return
  }
}
