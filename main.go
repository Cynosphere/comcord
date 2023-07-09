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

  // TODO: user account support
  client, err := discordgo.New("Bot " + token)
  if err != nil {
    fmt.Println("\r% Failed to create client:", err)
    os.Exit(1)
    return
  }

  // TODO: dont set for user accounts(? never really tested if it matters)
  client.Identify.Intents = discordgo.IntentsAll

  if config["useMobile"] == "true" {
    client.Identify.Properties = discordgo.IdentifyProperties{
      OS: "Android",
      Browser: "Discord Android",
      Device: "Pixel, raven",
    }
  } else {
    // TODO: user account support
    client.Identify.Properties = discordgo.IdentifyProperties{
      OS: runtime.GOOS,
      Browser: "comcord",
      Device: "comcord",
    }
  }

  events.Setup(client)

  err = client.Open()
  if err != nil {
    fmt.Println("\r% Failed to connect to Discord:", err)
    os.Exit(1)
    return
  }

  keyboard.Listen(func(key keys.Key) (stop bool, err error) {
    if !state.IsInPrompt() {
      if key.Code == keys.CtrlC {
        client.Close()
        os.Exit(0)
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
