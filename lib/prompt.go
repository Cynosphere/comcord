package lib

import (
	"strings"

	"github.com/Cynosphere/comcord/state"
	"github.com/ergochat/readline"
)

func MakePrompt(prompt string, uniqueLine bool, callback func(input string, interrupt bool)) {
  state.SetInPrompt(true)
  state.SetPromptText(prompt)

  rl, _ := readline.NewFromConfig(&readline.Config{
    Prompt: prompt,
    UniqueEditLine: uniqueLine,
  })
  defer rl.Close()

  input, err := rl.Readline()
  input = strings.TrimSpace(input)
  rl.Close()

  interrupt := err == readline.ErrInterrupt

  callback(input, interrupt)

  state.SetInPrompt(false)
  state.SetPromptText("")

  ProcessQueue()
}
