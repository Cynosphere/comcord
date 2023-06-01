package main

import (
  "os"
  "path/filepath"
  "strings"
)

func GetRCPath() string {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    panic(err)
  }

  return filepath.Join(homeDir, ".comcordrc")
}

func LoadRCFile() map[string]string {
  config := make(map[string]string)
  file, err := os.ReadFile(GetRCPath())
  if err != nil {
    panic(err)
  }

  lines := strings.Split(string(file), "\n")
  for _, line := range lines {
    kvs := strings.Split(line, "=")
    if len(kvs) == 2 {
      config[kvs[0]] = kvs[1]
    }
  }

  return config
}

func SaveRCFile(config map[string]string) {
  out := ""

  for key, value := range config {
    out = out + key + "=" + value + "\n"
  }

  err := os.WriteFile(GetRCPath(), []byte(out), 0644)
  if err != nil {
    panic(err)
  }
}
