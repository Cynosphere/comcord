package state

import (
  "time"

  "github.com/bwmarrin/discordgo"
)

type ComcordState struct {
  Config map[string]string
  Connected bool
  RPCConnected bool
  StartTime int64
  CurrentGuild string
  CurrentChannel string
  NameLength int
  InPrompt bool
  AFK bool
  MessageQueue []discordgo.Message
  LastChannel map[string]string
  LastDM string
}

var state ComcordState

func Setup(config map[string]string) {
  state = ComcordState{}
  state.Config = config
  state.Connected = true
  state.RPCConnected = false
  state.StartTime = time.Now().Unix()
  state.CurrentGuild = ""
  state.CurrentChannel = ""
  state.NameLength = 2
  state.InPrompt = false
  state.AFK = false
  state.MessageQueue = make([]discordgo.Message, 0)
  state.LastChannel = make(map[string]string)
  state.LastDM = ""
}

func IsConnected() bool {
  return state.Connected
}

func SetConnected(value bool) {
  state.Connected = value
}

func IsRPCConnected() bool {
  return state.RPCConnected
}

func SetRPCConnected(value bool) {
  state.RPCConnected = value
}

func GetStartTime() int64 {
  return state.StartTime
}

func GetCurrentGuild() string {
  return state.CurrentGuild
}

func SetCurrentGuild(value string) {
  state.CurrentGuild = value
}

func GetCurrentChannel() string {
  return state.CurrentChannel
}

func SetCurrentChannel(value string) {
  state.CurrentChannel = value
}

func GetNameLength() int {
  return state.NameLength
}

func SetNameLength(value int) {
  state.NameLength = value
}

func IsInPrompt() bool {
  return state.InPrompt
}

func SetInPrompt(value bool) {
  state.InPrompt = value
}

func IsAFK() bool {
  return state.AFK
}

func SetAFK(value bool) {
  state.AFK = value
}

func GetMessageQueue() []discordgo.Message {
  return state.MessageQueue
}

func AddMessageToQueue(msg discordgo.Message) {
  state.MessageQueue = append(state.MessageQueue, msg)
}

func SetLastChannel(guild string, channel string) {
  state.LastChannel[guild] = channel
}

func GetLastChannel(guild string) string {
  channel, has := state.LastChannel[guild]

  if has {
    return channel
  } else {
    return ""
  }
}

func GetLastDM() string {
  return state.LastDM
}

func SetLastDM(value string) {
  state.LastDM = value
}

func GetConfigValue(key string) string {
  value, has := state.Config[key]

  if has {
    return value
  } else {
    return ""
  }
}
