# comcord
A CLI-based client for Discord inspired by [SDF](https://sdf.org)'s [commode](https://sdf.org/?tutorials/comnotirc).

## Why?
1. All CLI/TUI Discord clients are outdated/unmaintained or have flaws.
2. I've been spending more time in commode on SDF and have been accustomed to the experience.

## Usage
1. `pnpm i`
2. `node src/index.js <token>`

Currently only bot accounts are supported, and that is unlikely to change anytime soon.
Eris has a lot of user-only endpoints implemented, but it would require hacking apart Eris to do the things nessicary to spoof being the actual client.
I also don't want to give skids an easy point of reference of how to spoof the client. :^)

You **MUST** grant your bot all Privileged Gateway Intents.

## Design Decisions
* Node.js was chosen currently due to familiarity.
* Eris was chosen due to familiarity and the nature of everything not being abstracted out to 200 different classes unlike discord.js.
* "Jank" by design. While I don't expect anyone to actually use comcord on serial terminals or teletypes other than for meme factor, the option is still there.

## TODO
- [ ] Commands
  - [ ] Quit (q)
  - [ ] Switch guilds (G)
  - [ ] Switch channels (g)
  - [ ] List online users in guild (w)
  - [ ] Emote (e)
    - Just sends message surrounded in `*`'s
  - [ ] Finger (f)
    - [ ] Shows presence data if available
    - [ ] Creation date, join date, ID, etc
  - [ ] Room history (r)
  - [ ] Extended room history (R)
- [ ] Message Receiving
  - [ ] Markdown styling
    - [ ] Common markdown (bold, italic, etc)
    - [ ] Figure out how spoilers would work
    - [ ] Emotes?????
    - [ ] Timestamp parsing
  - [ ] Embeds in the style of commode's posted links
  - [ ] Messages wrapped in `*`'s or `_`'s parsed as emotes
  - [ ] Inline DMs to replicate commode's private messages
  - [ ] Replies
- [ ] Message sending
  - [ ] Puts incoming messages into queue whilst in send mode
  - [ ] Mentions
  - [ ] Replies
- [ ] Configuration
  - [ ] Default guild/channel
- [ ] Threads
