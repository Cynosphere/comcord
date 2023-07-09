package commands

import "github.com/bwmarrin/discordgo"

var commandMap map[string]Command

type Command struct {
  Run func(*discordgo.Session)
  Description string
}

func Setup() {
  commandMap = make(map[string]Command)

  commandMap["q"] = Command{
    Run: QuitCommand,
    Description: "quit comcord",
  }

  commandMap["h"] = Command{
    Run: HelpCommand,
    Description: "command help",
  }

  commandMap["c"] = Command{
    Run: ClearCommand,
    Description: "clear",
  }
}

func GetCommand(key string) (Command, bool) {
  command, has := commandMap[key]
  return command, has
}

func GetAllCommands() map[string]Command {
  return commandMap
}
