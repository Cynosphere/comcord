package state

import (
  "time"

  "github.com/bwmarrin/discordgo"
)

type ComcordState struct {
  Connected bool
  RPCConnected bool
  StartTime int64
  CurrentGuild string
  CurrentChannel string
  NameLength int32
  InPrompt bool
  AFK bool
  MessageQueue []discordgo.Message
  LastChannel map[string]string
}

var state ComcordState

func Setup() {
  state = ComcordState{}
  state.Connected = true
  state.RPCConnected = false
  state.StartTime = time.Now().Unix()
  state.NameLength = 2
  state.InPrompt = false
  state.AFK = false
  state.MessageQueue = make([]discordgo.Message, 0)
  state.LastChannel = make(map[string]string)
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

func GetNameLength() int32 {
  return state.NameLength
}

func SetNameLength(value int32) {
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
