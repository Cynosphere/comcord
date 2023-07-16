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

  commandMap["e"] = Command{
    Run: EmoteCommand,
    Description: "emote",
  }

  commandMap["L"] = Command{
    Run: ListGuildsCommand,
    Description: "list guilds",
  }

  commandMap["l"] = Command{
    Run: ListChannelsCommand,
    Description: "list channels",
  }

  commandMap["G"] = Command{
    Run: SwitchGuildsCommand,
    Description: "goto guild",
  }

  commandMap["g"] = Command{
    Run: SwitchChannelsCommand,
    Description: "goto channel",
  }

  commandMap["w"] = Command{
    Run: ListUsersCommand,
    Description: "who is in channel",
  }

  commandMap["r"] = Command{
    Run: HistoryCommand,
    Description: "channel history",
  }

  commandMap["R"] = Command{
    Run: ExtendedHistoryCommand,
    Description: "extended history",
  }

  commandMap["p"] = Command{
    Run: PeekCommand,
    Description: "peek at channel",
  }

  commandMap["P"] = Command{
    Run: CrossPeekCommand,
    Description: "cross-guild peek",
  }
}

func GetCommand(key string) (Command, bool) {
  command, has := commandMap[key]
  return command, has
}

func GetAllCommands() map[string]Command {
  return commandMap
}
