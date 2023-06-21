# comcord
A CLI-based client for Discord inspired by [SDF](https://sdf.org)'s [commode](https://sdf.org/?tutorials/comnotirc).

## Why?
1. All CLI/TUI Discord clients are outdated/unmaintained or have flaws.
2. I've been spending more time in commode on SDF and have been accustomed to the experience.

## Usage
1. `pnpm i`
2. `node src/index.js <token>`

Your token will be then stored in `.comcordrc` after the first launch.

### User Accounts
User accounts are _partially_ supported via `allowUserAccounts=true` in your `.comcordrc`.
This is use at your own risk, despite spoofing the official client. I am not responsible for any banned accounts.

#### Guild members not populating
This is due to most libraries not implementing Lazy Guilds, as bots do not need lazy guilds to function.

If you are willing to implement Lazy Guilds based off of [unofficial documentation](https://luna.gitlab.io/discord-unofficial-docs/lazy_guilds.html)
and my already existing horrible hacks to make user accounts work in the first place, feel free to send a PR (on GitLab, GitHub repo is a read only mirror).

### Bot Accounts (prefered)
You **MUST** grant your bot all Privileged Gateway Intents.

## Design Decisions
- Node.js was chosen currently due to familiarity.
- Dysnomia was chosen due to familiarity and the nature of everything not being abstracted out to 200 different classes unlike discord.js.
- "Jank" by design. While I don't expect anyone to actually use comcord on serial terminals or teletypes other than for meme factor, the option is still there.

## TODO
- [x] Commands
  - [x] Quit (q)
  - [x] Switch guilds (G)
  - [x] Switch channels (g)
  - [x] List online users in guild (w)
  - [x] Emote (e)
    - Just sends message surrounded in `*`'s
  - [ ] Finger (f)
    - [ ] Shows presence data if available
    - [ ] Creation date, join date, ID, etc
  - [x] Room history (r)
  - [x] Extended room history (R)
  - [x] List channels (l)
  - [x] List guilds (L)
  - [x] Clear (c)
  - [ ] Surf channels forwards (>)
  - [ ] Surf channels backwards (<)
  - [x] AFK toggle (A)
  - [x] Send DM (s)
  - [x] Answer DM (a)
  - [x] Peek (p)
- [x] Message Receiving
  - [x] Markdown styling
    - [x] Common markdown (bold, italic, etc)
    - [x] Figure out how spoilers would work
    - [x] Emotes?????
    - [x] Timestamp parsing
    - [x] Mentions parsing
  - [ ] Embeds in the style of commode's posted links
  - [x] Messages wrapped in `*`'s or `_`'s parsed as emotes
  - [x] Inline DMs to replicate commode's private messages
  - [x] Replies
- [x] Message sending
  - [x] Puts incoming messages into queue whilst in send mode
  - [ ] Mentions
  - [ ] Replies
- [x] Configuration
  - [x] Default guild/channel
    - No way to set in client (yet?), `defaultChannel=` and `defaultGuild=` in your `.comcordrc`.
- [ ] Threads
- [x] Not have the token just be in argv
- [x] Not have everything in one file

## Repository
If you're viewing this on GitHub or GitLab, you are viewing a read only mirror.
The main repository is located on [Gitdab](https://gitdab.com/Cynosphere/comcord) and is push mirrored to the other two.
