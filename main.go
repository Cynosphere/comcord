package main

import (
	"context"
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
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	arikawa_state "github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/state/store/defaultstore"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/diamondburned/ningen/v3"
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

  commands.Setup()

  allowUserAccounts := config["allowUserAccounts"] == "true"
  tokenPrefix := "Bot "
  if allowUserAccounts {
    tokenPrefix = ""
  }

  fullToken := tokenPrefix + token

  props := gateway.IdentifyProperties{
    OS: runtime.GOOS,
  }
  statusType := config["statusType"]
  if statusType == "mobile" {
    props.Browser = "Discord Android"
  } else if statusType == "embedded" {
    props.Browser = "Discord Embedded"
  } else if statusType == "desktop" {
    props.Browser = "Discord Client"
  } else {
    props.Browser = "comcord"
  }

  ident := gateway.IdentifyCommand{
    Token: fullToken,
    Properties: props,

    Compress: true,
    LargeThreshold: 50,
  }

  status := "online"
  defaultStatus := config["defaultStatus"]
  if defaultStatus != "" {
    status = defaultStatus
  }
  startTime := state.GetStartTime()

  activity := discord.Activity{
    Name: "comcord",
    Type: discord.GameActivity,
    CreatedAt: discord.UnixTimestamp(startTime.Unix()),
    Timestamps: &discord.ActivityTimestamps{
      Start: discord.UnixMsTimestamp(startTime.Unix()),
    },
  }

  presence := gateway.UpdatePresenceCommand{
    Since: 0,
    Activities: make([]discord.Activity, 0),
    Status: discord.Status(status),
    AFK: false,
  }
  presence.Activities = append(presence.Activities, activity)
  ident.Presence = &presence

  gwURL, err := gateway.URL(context.Background())
  if err != nil {
    fmt.Print("% Failed to get gateway URL: ", err, "\n\r")
    os.Exit(1)
  }
  gw := gateway.NewCustomWithIdentifier(gateway.AddGatewayParams(gwURL), gateway.NewIdentifier(ident), nil)
  ses := session.NewWithGateway(gw, handler.New())
  st := arikawa_state.NewFromSession(ses, defaultstore.New())
  client := ningen.FromState(st)
  client.PreHandler = handler.New()

  client.AddIntents(gateway.IntentGuilds)
  client.AddIntents(gateway.IntentGuildPresences)
  client.AddIntents(gateway.IntentGuildMembers)
  client.AddIntents(gateway.IntentGuildMessages)
  client.AddIntents(gateway.IntentGuildMessageReactions)
  client.AddIntents(gateway.IntentDirectMessages)
  client.AddIntents(gateway.IntentDirectMessageReactions)
  client.AddIntents(gateway.IntentMessageContent)

  state.Setup(config, client)
  events.Setup(client)

  err = client.Open(context.Background())
  if err != nil {
    fmt.Print("% Failed to connect to Discord: ", err, "\n\r")
    os.Exit(1)
    return
  }
  defer client.Close()

  keyboard.Listen(func(key keys.Key) (stop bool, err error) {
    if !state.IsInPrompt() {
      if key.Code == keys.CtrlC {
        term.Restore(int(os.Stdin.Fd()), oldState)
        commands.QuitCommand()
        return true, nil
      } else {
        command, has := commands.GetCommand(key.String())
        if has {
          if key.String() == "q" {
            term.Restore(int(os.Stdin.Fd()), oldState)
          }
          command.Run()
        } else {
          commands.SendMode()
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
