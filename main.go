package main

import (
  "fmt"
  "os"
  "strings"

  "github.com/bwmarrin/discordgo"
)

func main() {
  var config map[string]string = make(map[string]string)
  var token string

  homeDir, homeErr := os.UserHomeDir()
  if homeErr != nil {
    panic(homeErr)
  }
  RCPATH := GetRCPath()

  _, rcErr := os.Stat(RCPATH)
  if !os.IsNotExist(rcErr) {
    fmt.Printf("%% Reading %s ...\n", strings.Replace(RCPATH, homeDir, "~", 1))
    config = LoadRCFile()
  }

  if len(os.Args) > 1 {
    token = os.Args[1]
    if os.IsNotExist(rcErr) {
      fmt.Println("% Writing token to ~/.comcordrc")
      config["token"] = token
      SaveRCFile(config)
    }
  } else {
    configToken, tokenInConfig := config["token"]
    if tokenInConfig {
      token = configToken
    } else {
      fmt.Println("No token provided.")
      os.Exit(1)
    }
  }
  fmt.Println(token)
}
