# comcord (`rewrite-go`)
A CLI-based client for Discord inspired by [SDF](https://sdf.org)'s [commode](https://sdf.org/?tutorials/comnotirc).

## Why?
1. All CLI/TUI Discord clients are outdated/unmaintained or have flaws.
2. I've been spending more time in commode on SDF and have been accustomed to the experience.

## Usage
TODO

## Rewrite Design Decisions
Go is more portable than Node.js

## TODO
- [x] Send mode
- [x] Commands
  - [x] Quit (q)
  - [ ] Switch guilds (G)
  - [ ] Switch channels (g)
  - [ ] List online users in guild (w)
  - [ ] Emote (e)
    - Just sends message surrounded in `*`'s
  - [ ] Finger (f)
    - Shows presence data if available
    - Creation date, join date, ID, etc
  - [ ] Room history (r)
  - [ ] Extended room history (R)
  - [ ] Peek (p)
  - [ ] Cross-guild peek (P)
  - [ ] List channels (l)
  - [ ] List guilds (L)
  - [x] Clear (c)
  - [ ] Surf channels forwards (>)
  - [ ] Surf channels backwards (<)
  - [ ] AFK toggle (A)
  - [ ] Send DM (s)
  - [ ] Answer DM (a)
- [x] Message Receiving
  - Markdown styling
    - [x] Emotes
    - [ ] Timestamp parsing
    - [x] Mentions parsing
  - [ ] Embeds in the style of commode's posted links
  - [x] Messages wrapped in `*`'s or `_`'s parsed as emotes
  - [x] Inline DMs to replicate commode's private messages
  - [x] Replies
  - [ ] Group DMs
- [x] Message sending
  - [x] Puts incoming messages into queue whilst in send mode
- [x] Configuration
  - [x] Write token from argv into rc file if rc file doesn't exist
  - [x] Default guild/channel
- [ ] Threads
- [ ] External rich presence when using bot accounts
