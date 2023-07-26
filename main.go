package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/Cynosphere/comcord/commands"
	"github.com/Cynosphere/comcord/events"
	"github.com/Cynosphere/comcord/rcfile"
	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/term"
)

func main() {
  oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
  if err != nil {
    panic(err)
  }
  defer term.Restore(int(os.Stdin.Fd()), oldState)

  var config map[string]string = make(map[string]string)
  var token string

  homeDir, homeErr := os.UserHomeDir()
  if homeErr != nil {
    panic(homeErr)
  }
  RCPATH := rcfile.GetPath()

  _, rcErr := os.Stat(RCPATH)
  if !os.IsNotExist(rcErr) {
    fmt.Printf("%% Reading %s ...\n", strings.Replace(RCPATH, homeDir, "~", 1))
    config = rcfile.Load()
  }

  if len(os.Args) > 1 {
    token = os.Args[1]
    if os.IsNotExist(rcErr) {
      fmt.Println("% Writing token to ~/.comcordrc")
      config["token"] = token
      rcfile.Save(config)
    }
  } else {
    configToken, tokenInConfig := config["token"]
    if tokenInConfig {
      token = configToken
    } else {
      fmt.Println("No token provided.")
      os.Exit(1)
      return
    }
  }

  fmt.Println("\rCOMcord (c)left 2023")
  fmt.Println("\rType 'h' for Commands")
  fmt.Print("\r")

  state.Setup(config)
  commands.Setup()

  allowUserAccounts := config["allowUserAccounts"] == "true"
  tokenPrefix := "Bot "
  if allowUserAccounts {
    tokenPrefix = ""
  }

  fullToken := tokenPrefix + token

  client, err := discordgo.New(fullToken)
  if err != nil {
    fmt.Println("% Failed to create client:", err)
    fmt.Print("\r")
    os.Exit(1)
    return
  }

  //client.LogLevel = -1
  client.LogLevel = discordgo.LogDebug

  client.Identify.Intents = discordgo.IntentsAll

  client.Identify.Properties = discordgo.IdentifyProperties{
    OS: runtime.GOOS,
  }
  statusType := config["statusType"]
  if statusType == "mobile" {
    client.Identify.Properties.Browser = "Discord Android"
  } else if statusType == "embedded" {
    client.Identify.Properties.Browser = "Discord Embedded"
  } else if statusType == "desktop" {
    client.Identify.Properties.Browser = "Discord Client"
  } else {
    client.Identify.Properties.Browser = "comcord"
  }

  status := "online"
  defaultStatus := config["defaultStatus"]
  if defaultStatus != "" {
    status = defaultStatus
  }
  startTime := state.GetStartTime()

  client.Identify.Presence = discordgo.GatewayStatusUpdate{
    Since: 0,
    Status: status,
    AFK: false,
    Game: discordgo.Activity{
      Type: 0,
      Name: "comcord",
      ApplicationID: "1026163285877325874",
      CreatedAt: startTime,
    },
  }

  events.Setup(client)

  err = client.Open()
  if err != nil {
    fmt.Println("% Failed to connect to Discord:", err)
    fmt.Print("\r")
    os.Exit(1)
    return
  }

  keyboard.Listen(func(key keys.Key) (stop bool, err error) {
    if !state.IsInPrompt() {
      if key.Code == keys.CtrlC {
        commands.QuitCommand(client)
        return true, nil
      } else {
        command, has := commands.GetCommand(key.String())
        if has {
          command.Run(client)
        } else {
          commands.SendMode(client)
        }
      }
    }

    return false, nil
  })

	/*sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

  client.Close()*/
}
