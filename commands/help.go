package commands

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Cynosphere/comcord/state"
	"github.com/mgutz/ansi"
)

const format string = "  %s - %s%s"

func lessLower(sa, sb string) bool {
  for {
    rb, nb := utf8.DecodeRuneInString(sb)
    if nb == 0 {
      // The number of runes in sa is greater than or
      // equal to the number of runes in sb. It follows
      // that sa is not less than sb.
      return false
    }

    ra, na := utf8.DecodeRuneInString(sa)
    if na == 0 {
      // The number of runes in sa is less than the
      // number of runes in sb. It follows that sa
      // is less than sb.
      return true
    }

    rbl := unicode.ToLower(rb)
    ral := unicode.ToLower(ra)

    if ral != rbl {
      return ral < rbl
    } else {
      return ra > rb
    }
  }
}

func HelpCommand() {
  noColor := state.HasNoColor()

  fmt.Println("\r\nCOMcord (c)left 2023\n\r")

  commands := GetAllCommands()
  keys := make([]string, 0, len(commands))
  for key := range commands {
    keys = append(keys, key)
  }
  sort.Slice(keys, func(i, j int) bool {
    return lessLower(keys[i], keys[j])
  })

  index := 0
  for _, key := range keys {
    cmd := commands[key]
    str := fmt.Sprintf(format, key, cmd.Description, "")
    length := len(str)
    padding := strings.Repeat(" ", 25 - length)

    if noColor {
      fmt.Printf(format, key, cmd.Description, padding)
    } else {
      coloredKey := ansi.Color(key, "yellow+b")
      fmt.Printf(format, coloredKey, cmd.Description, padding)
    }

    index++
    if index % 3 == 0 {
      fmt.Print("\n\r")
    }
  }
  if index % 3 != 0 {
    fmt.Print("\n\r")
  }

  fmt.Println("\r\nTo begin TALK MODE, press [SPACE]\n\r")
}
