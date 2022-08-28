const Eris = require("eris");
const chalk = require("chalk");

const token = process.argv[2];
const stdin = process.stdin;
const stdout = process.stdout;

stdin.setRawMode(true);
stdin.resume();
stdin.setEncoding("utf8");

let currentGuild,
  currentChannel,
  inSendMode = false,
  inEmoteMode = false,
  guildSwitch = false,
  channelSwitch = false,
  extendedHistory = false,
  nameLength = 2;

const messageQueue = [];

const commands = {
  q: "quit comcord",
  e: "emote",
  g: "goto a channel",
  G: "goto a guild",
  l: "list channels",
  L: "list guilds",
  w: "who is in guild",
  f: "finger",
  r: "channel history",
  R: "extended history",
  h: "command help",
  c: "clear",
  "<": "surf backwards",
  ">": "surf forwards",
};

const client = new Eris("Bot " + token, {
  defaultImageFormat: "png",
  defaultImageSize: 1024,
  intents: Eris.Constants.Intents.all,
});

client.once("ready", function () {
  console.log(
    "Logged in as: " +
      chalk.yellow(
        `${client.user.username}#${client.user.discriminator} (${client.user.id})`
      )
  );
  nameLength = client.user.username.length + 2;

  listGuilds();
});

function processMessage({
  name,
  content,
  bot,
  attachments,
  reply,
  isHistory = false,
}) {
  if (name.length + 2 > nameLength) nameLength = name.length + 2;

  if (reply) {
    const nameColor = reply.author.bot ? chalk.bold.yellow : chalk.bold.cyan;

    const headerLength = 5 + reply.author.username.length;
    const length = headerLength + reply.content.length;

    if (isHistory) {
      console.log(
        ` \u250d [${reply.author.username}] ${
          length > 79
            ? reply.content.substring(0, length - headerLength) + "\u2026"
            : reply.content
        }`
      );
    } else {
      console.log(
        chalk.bold.white(" \u250d ") +
          nameColor(`[${reply.author.username}] `) +
          chalk.reset(
            `${
              length > 79
                ? reply.content.substring(0, length - headerLength) + "\u2026"
                : reply.content
            }`
          )
      );
    }
  }

  if (
    (content.startsWith("*") && content.endsWith("*")) ||
    (content.startsWith("_") && content.endsWith("_"))
  ) {
    if (isHistory) {
      console.log(`<${name} ${content.subString(1, content.length - 1)}>`);
    } else {
      console.log(
        chalk.bold.green(
          `<${name} ${content.substring(1, content.length - 1)}>`
        )
      );
    }
  } else {
    if (isHistory) {
      console.log(
        `[${name}]${" ".repeat(nameLength - (name.length + 2))} ${content}`
      );
    } else {
      const nameColor = bot ? chalk.bold.yellow : chalk.bold.cyan;

      // TODO: markdown
      console.log(
        nameColor(`[${name}]`) +
          " ".repeat(nameLength - (name.length + 2)) +
          chalk.reset(" " + content)
      );
    }
  }

  for (const attachment of attachments) {
    if (isHistory) {
      console.log(`<attachment: ${attachment.url} >`);
    } else {
      console.log(chalk.bold.yellow(`<attachment: ${attachment.url} >`));
    }
  }
}

function processQueue() {
  for (const msg of messageQueue) {
    if (msg.content.indexOf("\n") > -1) {
      const lines = msg.content.split("\n");
      for (const index in lines) {
        const line = lines[index];
        processMessage({
          name: msg.author.username,
          bot: msg.author.bot,
          content: line,
          attachments: index == lines.length - 1 ? msg.attachments : null,
          reply: index == 0 ? msg.referencedMessage : null,
        });
      }
    } else {
      processMessage({
        name: msg.author.username,
        bot: msg.author.bot,
        content: msg.content,
        attachments: msg.attachments,
        reply: msg.referencedMessage,
      });
    }
  }

  messageQueue.splice(0, messageQueue.length);
}

client.on("messageCreate", function (msg) {
  if (msg.author.id === client.user.id) return;

  if (msg.channel.id == currentChannel) {
    if (inSendMode || inEmoteMode) {
      messageQueue.push(msg);
    } else {
      if (msg.content.indexOf("\n") > -1) {
        const lines = msg.content.split("\n");
        for (const index in lines) {
          const line = lines[index];
          processMessage({
            name: msg.author.username,
            bot: msg.author.bot,
            content: line,
            attachments: index == lines.length - 1 ? msg.attachments : null,
            reply: index == 0 ? msg.referencedMessage : null,
          });
        }
      } else {
        processMessage({
          name: msg.author.username,
          bot: msg.author.bot,
          content: msg.content,
          attachments: msg.attachments,
          reply: msg.referencedMessage,
        });
      }
    }
  }
});
client.on("messageUpdate", function (msg, old) {
  if (msg.author.id === client.user.id) return;

  if (msg.channel.id == currentChannel) {
    if (msg.content == old.content) return;

    if (inSendMode || inEmoteMode) {
      messageQueue.push(msg);
    } else {
      if (msg.content.indexOf("\n") > -1) {
        const lines = msg.content.split("\n");
        for (const index in lines) {
          const line = lines[index];
          processMessage({
            name: msg.author.username,
            bot: msg.author.bot,
            content: line + index == lines.length - 1 ? " (edited)" : null,
            attachments: index == lines.length - 1 ? msg.attachments : null,
            reply: index == 0 ? msg.referencedMessage : null,
          });
        }
      } else {
        processMessage({
          name: msg.author.username,
          bot: msg.author.bot,
          content: msg.content + " (edited)",
          attachments: msg.attachments,
          reply: msg.referencedMessage,
        });
      }
    }
  }
});

let toSend = "";
async function setupSendMode() {
  inSendMode = true;
  toSend = "";
  stdout.write(
    chalk.bold.cyan(`[${client.user.username}]`) +
      " ".repeat(nameLength - (client.user.username.length + 2)) +
      chalk.reset(" ")
  );
  try {
    await client.guilds
      .get(currentGuild)
      .channels.get(currentChannel)
      .sendTyping();
  } catch (err) {
    //
  }
}
async function sendMessage() {
  toSend = toSend.trim();
  if (toSend === "") {
    stdout.write("<no message sent>\n");
  } else {
    try {
      stdout.write("\n");
      await client.createMessage(currentChannel, toSend);
    } catch (err) {
      console.log("<failed to send message: " + err.message + ">");
    }
  }
  inSendMode = false;
  processQueue();
}

function showHelp() {
  console.log("\nCOMcord (c)left 2022\n");

  const keys = Object.keys(commands);
  keys.sort((a, b) => a.localeCompare(b));

  let index = 0;
  for (const key of keys) {
    const desc = commands[key];
    const length = `  ${key} - ${desc}`.length;

    stdout.write(
      "  " +
        chalk.bold.yellow(key) +
        chalk.reset(" - " + desc) +
        " ".repeat(Math.abs(25 - length))
    );

    index++;
    if (index % 3 == 0) stdout.write("\n");
  }
  if (index % 3 != 0) stdout.write("\n");

  console.log("\nTo begin TALK MODE, press [SPACE]\n");
}

function listGuilds() {
  let longest = 0;
  const guilds = [];

  for (const guild of client.guilds.values()) {
    if (guild.name.length > longest) longest = guild.name.length;

    const online = [...guild.members.values()].filter((m) => m.status).length;
    guilds.push({name: guild.name, members: guild.memberCount, online});
  }

  console.log("");
  console.log("  " + "guild-name".padStart(longest, " ") + "  online  total");
  console.log("-".repeat(80));
  for (const guild of guilds) {
    console.log(
      "  " +
        guild.name.padStart(longest, " ") +
        "  " +
        guild.online.toString().padStart(6, " ") +
        "  " +
        guild.members.toString().padStart(5, " ")
    );
  }
  console.log("");
}

let targetGuild = "";
function gotoGuild() {
  targetGuild = "";
  guildSwitch = true;

  stdout.write(":guild> ");
}

function findTopChannel(guildId) {
  const guild = client.guilds.get(guildId);
  const channels = [...guild.channels.values()].filter((c) => c.type == 0);
  channels.sort((a, b) => a.position - b.position);

  return channels[0];
}

function getStatus(status) {
  let color;
  switch (status) {
    case "online":
      color = chalk.bold.green;
      break;
    case "idle":
      color = chalk.bold.yellow;
      break;
    case "dnd":
      color = chalk.bold.red;
      break;
    default:
      color = chalk.bold;
      break;
  }

  return color(" \u2022 ");
}

function listUsers() {
  const guild = client.guilds.get(currentGuild);
  const channel = guild.channels.get(currentChannel);

  console.log(
    `\n[you are in '${guild.name}' in '${channel.name}' among ${guild.memberCount}]\n`
  );

  const online = [...guild.members.values()].filter((m) => m.status);
  online.sort((a, b) => a.name - b.name);

  let longest = 0;
  for (const member of online) {
    const name = member.user.username + "#" + member.user.discriminator;
    if (name.length + 3 > longest) longest = name.length + 3;
  }

  const columns = Math.ceil(stdout.columns / longest);

  let index = 0;
  for (const member of online) {
    const name = member.user.username + "#" + member.user.discriminator;
    const status = getStatus(member.status);
    const nameAndStatus = chalk.reset(name) + status;

    index++;
    stdout.write(
      nameAndStatus +
        " ".repeat(
          index % columns == 0 ? 0 : Math.abs(longest - (name.length + 3))
        )
    );

    if (index % columns == 0) stdout.write("\n");
  }
  if (index % columns != 0) stdout.write("\n");
  console.log("");

  if (channel.topic != null) {
    console.log("--Topic".padEnd(80, "-"));
    console.log(channel.topic);
    console.log("-".repeat(80));
    console.log("");
  }
}

function switchGuild() {
  targetGuild = targetGuild.trim();
  if (targetGuild == "") {
    listUsers();
    guildSwitch = false;
    return;
  }

  let target;

  for (const guild of client.guilds.values()) {
    if (guild.name.toLowerCase().indexOf(targetGuild.toLowerCase()) > -1) {
      target = guild.id;
      break;
    }
  }

  if (target == null) {
    console.log("<guild not found>");
  } else {
    currentGuild = target;
    // TODO: store last visited channel and switch to it if we've been to this guild before
    const topChannel = findTopChannel(target);
    currentChannel = topChannel.id;

    listUsers();
  }

  guildSwitch = false;
}

let targetChannel = "";
function gotoChannel() {
  targetChannel = "";
  channelSwitch = true;

  stdout.write(":channel> ");
}

function switchChannel() {
  targetChannel = targetChannel.trim();
  if (targetChannel == "") {
    listUsers();
    channelSwitch = false;
    return;
  }
  let target;

  const guild = client.guilds.get(currentGuild);
  const channels = [...guild.channels.values()].filter((c) => c.type == 0);
  channels.sort((a, b) => a.position - b.position);

  for (const channel of channels) {
    if (channel.name.toLowerCase().indexOf(targetChannel.toLowerCase()) > -1) {
      target = channel.id;
      break;
    }
  }

  if (target == null) {
    console.log("<channel not found>");
  } else {
    currentChannel = target;

    listUsers();
  }

  channelSwitch = false;
}

function startEmote() {
  toSend = "";
  inEmoteMode = true;

  stdout.write(":emote> ");
}

async function sendEmote() {
  toSend = toSend.trim();
  if (toSend === "") {
    console.log("<no message sent>");
  } else {
    try {
      await client.createMessage(currentChannel, "*" + toSend + "*");
      console.log(`<${client.user.username} ${toSend}>`);
    } catch (err) {
      console.log("<failed to send message: " + err.message + ">");
    }
  }
  inEmoteMode = false;
  processQueue();
}

async function getHistory(limit = 20) {
  const messages = await client.getMessages(currentChannel, {limit});
  messages.reverse();

  console.log("--Beginning-Review".padEnd(72, "-"));

  for (const msg of messages) {
    if (msg.content.indexOf("\n") > -1) {
      const lines = msg.content.split("\n");
      for (const index in lines) {
        const line = lines[index];
        processMessage({
          name: msg.author.username,
          bot: msg.author.bot,
          content:
            line +
            (msg.editedTimestamp != null && index == lines.length - 1
              ? " (edited)"
              : ""),
          attachments: index == lines.length - 1 ? msg.attachments : null,
          reply: index == 0 ? msg.referencedMessage : null,
          isHistory: true,
        });
      }
    } else {
      processMessage({
        name: msg.author.username,
        bot: msg.author.bot,
        content: msg.content + (msg.editedTimestamp != null ? " (edited)" : ""),
        attachments: msg.attachments,
        reply: msg.referencedMessage,
        isHistory: true,
      });
    }
  }

  console.log("--Review-Complete".padEnd(73, "-"));
}

let numLines = "";
function startExtendedHistory() {
  numLines = "";
  extendedHistory = true;

  stdout.write(":lines> ");
}

async function getExtendedHistory() {
  numLines = numLines.trim();
  numLines = parseInt(numLines);
  if (isNaN(numLines)) {
    console.log("<not a number>");
    extendedHistory = false;
    return;
  }

  try {
    await getHistory(numLines);
  } catch (err) {
    console.log("<failed to get history: " + err.message + ">");
  }

  extendedHistory = false;
}

stdin.on("data", function (key) {
  if (guildSwitch) {
    if (key === "\r") {
      console.log("");
      switchGuild();
    } else {
      if (key === "\b") {
        if (targetGuild.length > 0) {
          stdout.moveCursor(-1);
          stdout.write(" ");
          stdout.moveCursor(-1);
          targetGuild = targetGuild.substring(0, targetGuild.length - 1);
        }
      } else {
        stdout.write(key);
        targetGuild += key;
      }
    }
  } else if (channelSwitch) {
    if (key === "\r") {
      console.log("");
      switchChannel();
    } else {
      if (key === "\b") {
        if (targetChannel.length > 0) {
          stdout.moveCursor(-1);
          stdout.write(" ");
          stdout.moveCursor(-1);
          targetChannel = targetChannel.substring(0, targetChannel.length - 1);
        }
      } else {
        stdout.write(key);
        targetChannel += key;
      }
    }
  } else if (inSendMode) {
    if (key === "\r") {
      sendMessage();
    } else {
      if (key === "\b") {
        if (toSend.length > 0) {
          stdout.moveCursor(-1);
          stdout.write(" ");
          stdout.moveCursor(-1);
          toSend = toSend.substring(0, toSend.length - 1);
        }
      } else {
        stdout.write(key);
        toSend += key;
      }
    }
  } else if (inEmoteMode) {
    if (key === "\r") {
      console.log("");
      sendEmote();
    } else {
      if (key === "\b") {
        if (toSend.length > 0) {
          stdout.moveCursor(-1);
          stdout.write(" ");
          stdout.moveCursor(-1);
          toSend = toSend.substring(0, toSend.length - 1);
        }
      } else {
        stdout.write(key);
        toSend += key;
      }
    }
  } else if (extendedHistory) {
    if (key === "\r") {
      console.log("");
      getExtendedHistory();
    } else {
      if (key === "\b") {
        if (numLines.length > 0) {
          stdout.moveCursor(-1);
          stdout.write(" ");
          stdout.moveCursor(-1);
          numLines = numLines.substring(0, numLines.length - 1);
        }
      } else {
        stdout.write(key);
        numLines += key;
      }
    }
  } else {
    switch (key) {
      case "\u0003":
      case "q": {
        client.disconnect(false);
        process.exit(0);
        break;
      }
      case "h": {
        showHelp();
        break;
      }
      case "g": {
        if (currentGuild == null) {
          console.log("<not in a guild>");
          break;
        }
        gotoChannel();
        break;
      }
      case "G": {
        gotoGuild();
        break;
      }
      case "l": {
        if (currentGuild == null) {
          console.log("<not in a guild>");
          break;
        }
        break;
      }
      case "L": {
        listGuilds();
        break;
      }
      case "w": {
        if (currentGuild == null) {
          console.log("<not in a guild>");
          break;
        }
        listUsers();
        break;
      }
      case "e": {
        if (currentChannel == null) {
          console.log("<not in a channel>");
          break;
        }
        startEmote();
        break;
      }
      case "r": {
        if (currentChannel == null) {
          console.log("<not in a channel>");
          break;
        }
        getHistory();
        break;
      }
      case "R": {
        if (currentChannel == null) {
          console.log("<not in a channel>");
          break;
        }
        startExtendedHistory();
        break;
      }
      case "c": {
        console.clear();
        break;
      }
      case "<": {
        break;
      }
      case ">": {
        break;
      }
      case " ":
      case "\r":
      default: {
        if (currentChannel == null) {
          console.log("<not in a channel>");
          break;
        }
        setupSendMode();
        break;
      }
    }
  }
});

client.connect();

console.log("COMcord (c)left 2022");
console.log("Type 'h' for Commands");
